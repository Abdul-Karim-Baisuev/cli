package view

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/internal/prompter"
	"github.com/cli/cli/v2/pkg/cmd/run/shared"
	workflowShared "github.com/cli/cli/v2/pkg/cmd/workflow/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdView(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wants    ViewOptions
		wantsErr bool
		tty      bool
	}{
		{
			name: "blank tty",
			tty:  true,
			wants: ViewOptions{
				Prompt: true,
			},
		},
		{
			name:     "blank nontty",
			wantsErr: true,
		},
		{
			name: "nontty jobID",
			cli:  "1234",
			wants: ViewOptions{
				JobID: "1234",
			},
		},
		{
			name: "log tty",
			tty:  true,
			cli:  "--log",
			wants: ViewOptions{
				Prompt: true,
				Log:    true,
			},
		},
		{
			name: "log nontty",
			cli:  "--log 1234",
			wants: ViewOptions{
				JobID: "1234",
				Log:   true,
			},
		},
		{
			name: "exit status",
			cli:  "--exit-status 1234",
			wants: ViewOptions{
				JobID:      "1234",
				ExitStatus: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			io.SetStdinTTY(tt.tty)
			io.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *ViewOptions
			cmd := NewCmdView(f, func(opts *ViewOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tt.wants.JobID, gotOpts.JobID)
			assert.Equal(t, tt.wants.Log, gotOpts.Log)
			assert.Equal(t, tt.wants.Prompt, gotOpts.Prompt)
			assert.Equal(t, tt.wants.ExitStatus, gotOpts.ExitStatus)
		})
	}
}

func TestRunView(t *testing.T) {
	tests := []struct {
		name      string
		opts      *ViewOptions
		httpStubs func(*httpmock.Registry)
		askStubs  func(*prompter.MockPrompter)
		tty       bool
		wantErr   bool
		wantOut   string
	}{
		{
			name: "exit status respected with --log",
			opts: &ViewOptions{
				JobID:      "20",
				ExitStatus: true,
				Log:        true,
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/20"),
					httpmock.JSONResponse(shared.FailedJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/20/logs"),
					httpmock.StringResponse("it's a log\nfor this job\nbeautiful log\n"))
			},
			wantErr: true,
			wantOut: "it's a log\nfor this job\nbeautiful log\n",
		},
		{
			name: "exit status respected",
			opts: &ViewOptions{
				JobID:      "20",
				ExitStatus: true,
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/20"),
					httpmock.JSONResponse(shared.FailedJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/check-runs/20/annotations"),
					httpmock.JSONResponse(shared.FailedJobAnnotations))
			},
			wantErr: true,
			wantOut: "sad job (ID 20)\nX 59m ago in 4m34s\n\n✓ barf the quux\nX quux the barf\n\nANNOTATIONS\nX the job is sad\nblaze.py#420\n\n\nTo see the full logs for this job, try: gh job view 20 --log\nView this job on GitHub: https://github.com/jobs/20\n",
		},
		{
			name: "interactive flow, multi-job",
			tty:  true,
			opts: &ViewOptions{
				Prompt: true,
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.TestRuns,
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs/3"),
					httpmock.JSONResponse(shared.SuccessfulRun))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows/123"),
					httpmock.JSONResponse(shared.TestWorkflow))
				reg.Register(
					httpmock.REST("GET", "runs/3/jobs"),
					httpmock.JSONResponse(shared.JobsPayload{
						Jobs: []shared.Job{
							shared.SuccessfulJob,
							shared.FailedJob,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10"),
					httpmock.JSONResponse(shared.SuccessfulJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/check-runs/10/annotations"),
					httpmock.JSONResponse([]shared.Annotation{}))
			},
			askStubs: func(pm *prompter.MockPrompter) {
				pm.RegisterSelect("Select a workflow run", []string{
					"X cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"✓ cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
					"- cool commit, CI [trunk] Feb 23, 2021",
					"- cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
				},
					func(_, _ string, opts []string) (int, error) { return 2, nil })
				pm.RegisterSelect("Select a job to view", []string{"✓ cool job", "X sad job"},
					func(_, _ string, opts []string) (int, error) { return 0, nil })
			},
			wantOut: "\n\ncool job (ID 10)\n✓ 59m ago in 4m34s\n\n✓ fob the barz\n✓ barz the fob\n\nTo see the full logs for this job, try: gh job view 10 --log\nView this job on GitHub: https://github.com/jobs/10\n",
		},
		{
			name: "interactive, run has only one job",
			tty:  true,
			opts: &ViewOptions{
				Prompt: true,
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.TestRuns,
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs/3"),
					httpmock.JSONResponse(shared.SuccessfulRun))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows/123"),
					httpmock.JSONResponse(shared.TestWorkflow))
				reg.Register(
					httpmock.REST("GET", "runs/3/jobs"),
					httpmock.JSONResponse(shared.JobsPayload{
						Jobs: []shared.Job{
							shared.SuccessfulJob,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10"),
					httpmock.JSONResponse(shared.SuccessfulJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/check-runs/10/annotations"),
					httpmock.JSONResponse([]shared.Annotation{}))
			},
			askStubs: func(pm *prompter.MockPrompter) {
				pm.RegisterSelect("Select a workflow run", []string{
					"X cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"✓ cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
					"- cool commit, CI [trunk] Feb 23, 2021",
					"- cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
				},
					func(_, _ string, opts []string) (int, error) { return 2, nil })
			},
			wantOut: "\n\ncool job (ID 10)\n✓ 59m ago in 4m34s\n\n✓ fob the barz\n✓ barz the fob\n\nTo see the full logs for this job, try: gh job view 10 --log\nView this job on GitHub: https://github.com/jobs/10\n",
		},
		{
			name: "interactive with log",
			tty:  true,
			opts: &ViewOptions{
				Prompt: true,
				Log:    true,
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.TestRuns,
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs/3"),
					httpmock.JSONResponse(shared.SuccessfulRun))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows/123"),
					httpmock.JSONResponse(shared.TestWorkflow))
				reg.Register(
					httpmock.REST("GET", "runs/3/jobs"),
					httpmock.JSONResponse(shared.JobsPayload{
						Jobs: []shared.Job{
							shared.SuccessfulJob,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10"),
					httpmock.JSONResponse(shared.SuccessfulJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10/logs"),
					httpmock.StringResponse("it's a log\nfor this job\nbeautiful log\n"))
			},
			askStubs: func(pm *prompter.MockPrompter) {
				pm.RegisterSelect("Select a workflow run", []string{
					"X cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"✓ cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
					"- cool commit, CI [trunk] Feb 23, 2021",
					"- cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"* cool commit, CI [trunk] Feb 23, 2021",
					"X cool commit, CI [trunk] Feb 23, 2021",
				},
					func(_, _ string, opts []string) (int, error) { return 2, nil })
			},
			wantOut: "\n\nit's a log\nfor this job\nbeautiful log\n",
		},
		{
			name: "noninteractive with log",
			opts: &ViewOptions{
				JobID:  "10",
				Prompt: false,
				Log:    true,
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10"),
					httpmock.JSONResponse(shared.SuccessfulJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10/logs"),
					httpmock.StringResponse("it's a log\nfor this job\nbeautiful log\n"))
			},
			wantOut: "it's a log\nfor this job\nbeautiful log\n",
		},
		{
			name: "noninteractive",
			opts: &ViewOptions{
				JobID: "10",
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/10"),
					httpmock.JSONResponse(shared.SuccessfulJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/check-runs/10/annotations"),
					httpmock.JSONResponse([]shared.Annotation{}))
			},
			wantOut: "cool job (ID 10)\n✓ 59m ago in 4m34s\n\n✓ fob the barz\n✓ barz the fob\n\nTo see the full logs for this job, try: gh job view 10 --log\nView this job on GitHub: https://github.com/jobs/10\n",
		},
		{
			name: "shows annotations for failed job",
			opts: &ViewOptions{
				JobID: "20",
			},
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/jobs/20"),
					httpmock.JSONResponse(shared.FailedJob))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/check-runs/20/annotations"),
					httpmock.JSONResponse(shared.FailedJobAnnotations))
			},
			wantOut: "sad job (ID 20)\nX 59m ago in 4m34s\n\n✓ barf the quux\nX quux the barf\n\nANNOTATIONS\nX the job is sad\nblaze.py#420\n\n\nTo see the full logs for this job, try: gh job view 20 --log\nView this job on GitHub: https://github.com/jobs/20\n",
		},
	}

	for _, tt := range tests {
		reg := &httpmock.Registry{}
		tt.httpStubs(reg)
		tt.opts.HttpClient = func() (*http.Client, error) {
			return &http.Client{Transport: reg}, nil
		}

		tt.opts.Now = func() time.Time {
			notnow, _ := time.Parse("2006-01-02 15:04:05", "2021-02-23 05:50:00")
			return notnow
		}

		io, _, stdout, _ := iostreams.Test()
		io.SetStdoutTTY(tt.tty)
		tt.opts.IO = io
		tt.opts.BaseRepo = func() (ghrepo.Interface, error) {
			return ghrepo.FromFullName("OWNER/REPO")
		}

		pm := prompter.NewMockPrompter(t)
		tt.opts.Prompter = pm
		if tt.askStubs != nil {
			tt.askStubs(pm)
		}
		t.Run(tt.name, func(t *testing.T) {
			err := runView(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				if !tt.opts.ExitStatus {
					return
				}
			}
			if !tt.opts.ExitStatus {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantOut, stdout.String())
			reg.Verify(t)
		})
	}

}
