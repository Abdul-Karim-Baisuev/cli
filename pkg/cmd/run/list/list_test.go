package list

import (
	"bytes"
	"fmt"
<<<<<<< HEAD
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmd/run/shared"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/httpmock"
	"github.com/cli/cli/pkg/iostreams"
=======
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/cmd/run/shared"
	workflowShared "github.com/cli/cli/v2/pkg/cmd/workflow/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/cli/cli/v2/pkg/iostreams"
>>>>>>> origin/trunk
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdList(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		tty      bool
		wants    ListOptions
		wantsErr bool
	}{
		{
			name: "blank",
			wants: ListOptions{
				Limit: defaultLimit,
			},
		},
		{
			name: "limit",
			cli:  "--limit 100",
			wants: ListOptions{
				Limit: 100,
			},
		},
		{
			name:     "bad limit",
			cli:      "--limit hi",
			wantsErr: true,
		},
<<<<<<< HEAD
=======
		{
			name: "workflow",
			cli:  "--workflow foo.yml",
			wants: ListOptions{
				Limit:            defaultLimit,
				WorkflowSelector: "foo.yml",
			},
		},
		{
			name: "branch",
			cli:  "--branch new-cool-stuff",
			wants: ListOptions{
				Limit:  defaultLimit,
				Branch: "new-cool-stuff",
			},
		},
		{
			name: "user",
			cli:  "--user bak1an",
			wants: ListOptions{
				Limit: defaultLimit,
				Actor: "bak1an",
			},
		},
		{
			name: "status",
			cli:  "--status completed",
			wants: ListOptions{
				Limit:  defaultLimit,
				Status: "completed",
			},
		},
		{
			name: "event",
			cli:  "--event push",
			wants: ListOptions{
				Limit: defaultLimit,
				Event: "push",
			},
		},
		{
			name: "created",
			cli:  "--created >=2023-04-24",
			wants: ListOptions{
				Limit:   defaultLimit,
				Created: ">=2023-04-24",
			},
		},
>>>>>>> origin/trunk
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
<<<<<<< HEAD
			io, _, _, _ := iostreams.Test()
			io.SetStdinTTY(tt.tty)
			io.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: io,
=======
			ios, _, _, _ := iostreams.Test()
			ios.SetStdinTTY(tt.tty)
			ios.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: ios,
>>>>>>> origin/trunk
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *ListOptions
			cmd := NewCmdList(f, func(opts *ListOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
<<<<<<< HEAD
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)
=======
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
>>>>>>> origin/trunk

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wants.Limit, gotOpts.Limit)
<<<<<<< HEAD
=======
			assert.Equal(t, tt.wants.WorkflowSelector, gotOpts.WorkflowSelector)
			assert.Equal(t, tt.wants.Branch, gotOpts.Branch)
			assert.Equal(t, tt.wants.Actor, gotOpts.Actor)
			assert.Equal(t, tt.wants.Status, gotOpts.Status)
			assert.Equal(t, tt.wants.Event, gotOpts.Event)
			assert.Equal(t, tt.wants.Created, gotOpts.Created)
>>>>>>> origin/trunk
		})
	}
}

func TestListRun(t *testing.T) {
	tests := []struct {
		name       string
		opts       *ListOptions
<<<<<<< HEAD
		wantOut    string
		wantErrOut string
		stubs      func(*httpmock.Registry)
		nontty     bool
	}{
		{
			name: "blank tty",
			opts: &ListOptions{
				Limit: defaultLimit,
			},
=======
		wantErr    bool
		wantOut    string
		wantErrOut string
		wantErrMsg string
		stubs      func(*httpmock.Registry)
		isTTY      bool
	}{
		{
			name: "default arguments",
			opts: &ListOptions{
				Limit: defaultLimit,
				now:   shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
			},
			isTTY: true,
>>>>>>> origin/trunk
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.TestRuns,
					}))
<<<<<<< HEAD
			},
			wantOut: "X  cool commit  timed out    trunk  push  1\n-  cool commit  in progress  trunk  push  2\n✓  cool commit  successful   trunk  push  3\n✓  cool commit  cancelled    trunk  push  4\nX  cool commit  failed       trunk  push  1234\n✓  cool commit  neutral      trunk  push  6\n✓  cool commit  skipped      trunk  push  7\n-  cool commit  requested    trunk  push  8\n-  cool commit  queued       trunk  push  9\nX  cool commit  stale        trunk  push  10\n\nFor details on a run, try: gh run view <run-id>\n",
		},
		{
			name: "blank nontty",
			opts: &ListOptions{
				Limit:       defaultLimit,
				PlainOutput: true,
			},
			nontty: true,
=======
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
			},
			wantOut: heredoc.Doc(`
				STATUS  TITLE        WORKFLOW  BRANCH  EVENT  ID    ELAPSED  AGE
				X       cool commit  CI        trunk   push   1     4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   2     4m34s    about 4 minutes ago
				✓       cool commit  CI        trunk   push   3     4m34s    about 4 minutes ago
				X       cool commit  CI        trunk   push   4     4m34s    about 4 minutes ago
				X       cool commit  CI        trunk   push   1234  4m34s    about 4 minutes ago
				-       cool commit  CI        trunk   push   6     4m34s    about 4 minutes ago
				-       cool commit  CI        trunk   push   7     4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   8     4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   9     4m34s    about 4 minutes ago
				X       cool commit  CI        trunk   push   10    4m34s    about 4 minutes ago
			`),
		},
		{
			name: "inactive disabled workflow selected",
			opts: &ListOptions{
				Limit:            defaultLimit,
				now:              shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
				WorkflowSelector: "d. inact",
				All:              false,
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				// Uses abbreviated names and commit messages because of output column limit
				workflow := workflowShared.Workflow{
					Name:  "d. inact",
					ID:    1206,
					Path:  ".github/workflows/disabledInactivity.yml",
					State: workflowShared.DisabledInactivity,
				}

				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							workflow,
						},
					}))
			},
			wantErr:    true,
			wantErrMsg: "could not find any workflows named d. inact",
		},
		{
			name: "inactive disabled workflow selected and all states applied",
			opts: &ListOptions{
				Limit:            defaultLimit,
				now:              shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
				WorkflowSelector: "d. inact",
				All:              true,
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				// Uses abbreviated names and commit messages because of output column limit
				workflow := workflowShared.Workflow{
					Name:  "d. inact",
					ID:    1206,
					Path:  ".github/workflows/disabledInactivity.yml",
					State: workflowShared.DisabledInactivity,
				}

				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							workflow,
						},
					}))
				reg.Register(
					httpmock.REST("GET", fmt.Sprintf("repos/OWNER/REPO/actions/workflows/%d/runs", workflow.ID)),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: []shared.Run{
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 101, shared.Completed, shared.TimedOut, "dicto"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 102, shared.InProgress, shared.TimedOut, "diito"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 103, shared.Completed, shared.Success, "dics"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 104, shared.Completed, shared.Cancelled, "dicc"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 105, shared.Completed, shared.Failure, "dicf"),
						},
					}))
			},
			wantOut: heredoc.Doc(`
				STATUS  TITLE  WORKFLOW  BRANCH  EVENT  ID   ELAPSED  AGE
				X       dicto  d. inact  trunk   push   101  4m34s    about 4 minutes ago
				*       diito  d. inact  trunk   push   102  4m34s    about 4 minutes ago
				✓       dics   d. inact  trunk   push   103  4m34s    about 4 minutes ago
				X       dicc   d. inact  trunk   push   104  4m34s    about 4 minutes ago
				X       dicf   d. inact  trunk   push   105  4m34s    about 4 minutes ago
			`),
		},
		{
			name: "manually disabled workflow selected",
			opts: &ListOptions{
				Limit:            defaultLimit,
				now:              shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
				WorkflowSelector: "d. man",
				All:              false,
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				// Uses abbreviated names and commit messages because of output column limit
				workflow := workflowShared.Workflow{
					Name:  "d. man",
					ID:    456,
					Path:  ".github/workflows/disabled.yml",
					State: workflowShared.DisabledManually,
				}

				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							workflow,
						},
					}))
			},
			wantErr:    true,
			wantErrMsg: "could not find any workflows named d. man",
		},
		{
			name: "manually disabled workflow selected and all states applied",
			opts: &ListOptions{
				Limit:            defaultLimit,
				now:              shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
				WorkflowSelector: "d. man",
				All:              true,
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				// Uses abbreviated names and commit messages because of output column limit
				workflow := workflowShared.Workflow{
					Name:  "d. man",
					ID:    456,
					Path:  ".github/workflows/disabled.yml",
					State: workflowShared.DisabledManually,
				}

				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							workflow,
						},
					}))
				reg.Register(
					httpmock.REST("GET", fmt.Sprintf("repos/OWNER/REPO/actions/workflows/%d/runs", workflow.ID)),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: []shared.Run{
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 201, shared.Completed, shared.TimedOut, "dmcto"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 202, shared.InProgress, shared.TimedOut, "dmito"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 203, shared.Completed, shared.Success, "dmcs"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 204, shared.Completed, shared.Cancelled, "dmcc"),
							shared.TestRunWithWorkflowAndCommit(workflow.ID, 205, shared.Completed, shared.Failure, "dmcf"),
						},
					}))
			},
			wantOut: heredoc.Doc(`
				STATUS  TITLE  WORKFLOW  BRANCH  EVENT  ID   ELAPSED  AGE
				X       dmcto  d. man    trunk   push   201  4m34s    about 4 minutes ago
				*       dmito  d. man    trunk   push   202  4m34s    about 4 minutes ago
				✓       dmcs   d. man    trunk   push   203  4m34s    about 4 minutes ago
				X       dmcc   d. man    trunk   push   204  4m34s    about 4 minutes ago
				X       dmcf   d. man    trunk   push   205  4m34s    about 4 minutes ago
			`),
		},
		{
			name: "default arguments nontty",
			opts: &ListOptions{
				Limit: defaultLimit,
				now:   shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
			},
			isTTY: false,
>>>>>>> origin/trunk
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.TestRuns,
					}))
<<<<<<< HEAD
			},
			wantOut: "completed\ttimed_out\tcool commit\ttimed out\ttrunk\tpush\t4m34s\t1\nin_progress\t\tcool commit\tin progress\ttrunk\tpush\t4m34s\t2\ncompleted\tsuccess\tcool commit\tsuccessful\ttrunk\tpush\t4m34s\t3\ncompleted\tcancelled\tcool commit\tcancelled\ttrunk\tpush\t4m34s\t4\ncompleted\tfailure\tcool commit\tfailed\ttrunk\tpush\t4m34s\t1234\ncompleted\tneutral\tcool commit\tneutral\ttrunk\tpush\t4m34s\t6\ncompleted\tskipped\tcool commit\tskipped\ttrunk\tpush\t4m34s\t7\nrequested\t\tcool commit\trequested\ttrunk\tpush\t4m34s\t8\nqueued\t\tcool commit\tqueued\ttrunk\tpush\t4m34s\t9\ncompleted\tstale\tcool commit\tstale\ttrunk\tpush\t4m34s\t10\n",
=======
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
			},
			wantOut: heredoc.Doc(`
				completed	timed_out	cool commit	CI	trunk	push	1	4m34s	2021-02-23T04:51:00Z
				in_progress		cool commit	CI	trunk	push	2	4m34s	2021-02-23T04:51:00Z
				completed	success	cool commit	CI	trunk	push	3	4m34s	2021-02-23T04:51:00Z
				completed	cancelled	cool commit	CI	trunk	push	4	4m34s	2021-02-23T04:51:00Z
				completed	failure	cool commit	CI	trunk	push	1234	4m34s	2021-02-23T04:51:00Z
				completed	neutral	cool commit	CI	trunk	push	6	4m34s	2021-02-23T04:51:00Z
				completed	skipped	cool commit	CI	trunk	push	7	4m34s	2021-02-23T04:51:00Z
				requested		cool commit	CI	trunk	push	8	4m34s	2021-02-23T04:51:00Z
				queued		cool commit	CI	trunk	push	9	4m34s	2021-02-23T04:51:00Z
				completed	stale	cool commit	CI	trunk	push	10	4m34s	2021-02-23T04:51:00Z
			`),
		},
		{
			name: "org ruleset workflow in runs list shows with empty workflow name",
			opts: &ListOptions{
				Limit: defaultLimit,
				now:   shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.TestRunsWithOrgRequiredWorkflows,
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows/456"),
					httpmock.StatusStringResponse(404, "not found"),
				)
			},
			wantOut: heredoc.Doc(`
				STATUS  TITLE        WORKFLOW  BRANCH  EVENT  ID  ELAPSED  AGE
				X       cool commit            trunk   push   1   4m34s    about 4 minutes ago
				*       cool commit            trunk   push   2   4m34s    about 4 minutes ago
				✓       cool commit            trunk   push   3   4m34s    about 4 minutes ago
				X       cool commit            trunk   push   4   4m34s    about 4 minutes ago
				X       cool commit  CI        trunk   push   5   4m34s    about 4 minutes ago
				-       cool commit  CI        trunk   push   6   4m34s    about 4 minutes ago
				-       cool commit  CI        trunk   push   7   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   8   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   9   4m34s    about 4 minutes ago
			`),
>>>>>>> origin/trunk
		},
		{
			name: "pagination",
			opts: &ListOptions{
				Limit: 101,
<<<<<<< HEAD
			},
			stubs: func(reg *httpmock.Registry) {
				runID := 0
				runs := []shared.Run{}
				for runID < 103 {
					runs = append(runs, shared.TestRun(fmt.Sprintf("%d", runID), runID, shared.InProgress, ""))
=======
				now:   shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				var runID int64
				runs := []shared.Run{}
				for runID < 103 {
					runs = append(runs, shared.TestRun(runID, shared.InProgress, ""))
>>>>>>> origin/trunk
					runID++
				}
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
<<<<<<< HEAD
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: runs[0:100],
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: runs[100:],
					}))
			},
			wantOut: longRunOutput,
		},
		{
			name: "no results nontty",
			opts: &ListOptions{
				Limit:       defaultLimit,
				PlainOutput: true,
			},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			nontty:  true,
			wantOut: "",
		},
		{
			name: "no results tty",
=======
					httpmock.WithHeader(httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: runs[0:100],
					}), "Link", `<https://api.github.com/repositories/123/actions/runs?per_page=100&page=2>; rel="next"`))
				reg.Register(
					httpmock.REST("GET", "repositories/123/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: runs[100:],
					}))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows"),
					httpmock.JSONResponse(workflowShared.WorkflowsPayload{
						Workflows: []workflowShared.Workflow{
							shared.TestWorkflow,
						},
					}))
			},
			wantOut: heredoc.Doc(`
				STATUS  TITLE        WORKFLOW  BRANCH  EVENT  ID   ELAPSED  AGE
				*       cool commit  CI        trunk   push   0    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   1    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   2    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   3    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   4    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   5    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   6    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   7    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   8    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   9    4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   10   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   11   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   12   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   13   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   14   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   15   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   16   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   17   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   18   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   19   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   20   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   21   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   22   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   23   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   24   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   25   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   26   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   27   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   28   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   29   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   30   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   31   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   32   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   33   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   34   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   35   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   36   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   37   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   38   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   39   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   40   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   41   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   42   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   43   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   44   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   45   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   46   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   47   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   48   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   49   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   50   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   51   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   52   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   53   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   54   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   55   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   56   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   57   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   58   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   59   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   60   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   61   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   62   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   63   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   64   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   65   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   66   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   67   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   68   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   69   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   70   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   71   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   72   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   73   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   74   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   75   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   76   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   77   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   78   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   79   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   80   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   81   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   82   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   83   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   84   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   85   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   86   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   87   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   88   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   89   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   90   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   91   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   92   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   93   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   94   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   95   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   96   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   97   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   98   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   99   4m34s    about 4 minutes ago
				*       cool commit  CI        trunk   push   100  4m34s    about 4 minutes ago
			`),
		},
		{
			name: "no results",
>>>>>>> origin/trunk
			opts: &ListOptions{
				Limit: defaultLimit,
			},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/runs"),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
<<<<<<< HEAD
			wantOut:    "",
			wantErrOut: "No runs found\n",
=======
			wantErr:    true,
			wantErrMsg: "no runs found",
		},
		{
			name: "workflow selector",
			opts: &ListOptions{
				Limit:            defaultLimit,
				WorkflowSelector: "flow.yml",
				now:              shared.TestRunStartTime.Add(time.Minute*4 + time.Second*34),
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows/flow.yml"),
					httpmock.JSONResponse(workflowShared.AWorkflow))
				reg.Register(
					httpmock.REST("GET", "repos/OWNER/REPO/actions/workflows/123/runs"),
					httpmock.JSONResponse(shared.RunsPayload{
						WorkflowRuns: shared.WorkflowRuns,
					}))
			},
			wantOut: heredoc.Doc(`
				STATUS  TITLE        WORKFLOW    BRANCH  EVENT  ID    ELAPSED  AGE
				*       cool commit  a workflow  trunk   push   2     4m34s    about 4 minute...
				✓       cool commit  a workflow  trunk   push   3     4m34s    about 4 minute...
				X       cool commit  a workflow  trunk   push   1234  4m34s    about 4 minute...
			`),
		},
		{
			name: "branch filter applied",
			opts: &ListOptions{
				Limit:  defaultLimit,
				Branch: "the-branch",
			},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.QueryMatcher("GET", "repos/OWNER/REPO/actions/runs", url.Values{
						"branch": []string{"the-branch"},
					}),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			wantErr:    true,
			wantErrMsg: "no runs found",
		},
		{
			name: "actor filter applied",
			opts: &ListOptions{
				Limit: defaultLimit,
				Actor: "bak1an",
			},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.QueryMatcher("GET", "repos/OWNER/REPO/actions/runs", url.Values{
						"actor": []string{"bak1an"},
					}),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			wantErr:    true,
			wantErrMsg: "no runs found",
		},
		{
			name: "status filter applied",
			opts: &ListOptions{
				Limit:  defaultLimit,
				Status: "queued",
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.QueryMatcher("GET", "repos/OWNER/REPO/actions/runs", url.Values{
						"status": []string{"queued"},
					}),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			wantErr:    true,
			wantErrMsg: "no runs found",
		},
		{
			name: "event filter applied",
			opts: &ListOptions{
				Limit: defaultLimit,
				Event: "push",
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.QueryMatcher("GET", "repos/OWNER/REPO/actions/runs", url.Values{
						"event": []string{"push"},
					}),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			wantErr:    true,
			wantErrMsg: "no runs found",
		},
		{
			name: "created filter applied",
			opts: &ListOptions{
				Limit:   defaultLimit,
				Created: ">=2023-04-24",
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.QueryMatcher("GET", "repos/OWNER/REPO/actions/runs", url.Values{
						"created": []string{">=2023-04-24"},
					}),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			wantErr:    true,
			wantErrMsg: "no runs found",
		},
		{
			name: "commit filter applied",
			opts: &ListOptions{
				Limit:  defaultLimit,
				Commit: "1234567890123456789012345678901234567890",
			},
			isTTY: true,
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.QueryMatcher("GET", "repos/OWNER/REPO/actions/runs", url.Values{
						"head_sha": []string{"1234567890123456789012345678901234567890"},
					}),
					httpmock.JSONResponse(shared.RunsPayload{}),
				)
			},
			wantErr:    true,
			wantErrMsg: "no runs found",
>>>>>>> origin/trunk
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &httpmock.Registry{}
<<<<<<< HEAD
=======
			defer reg.Verify(t)
>>>>>>> origin/trunk
			tt.stubs(reg)

			tt.opts.HttpClient = func() (*http.Client, error) {
				return &http.Client{Transport: reg}, nil
			}

<<<<<<< HEAD
			io, _, stdout, stderr := iostreams.Test()
			io.SetStdoutTTY(!tt.nontty)
			tt.opts.IO = io
=======
			ios, _, stdout, stderr := iostreams.Test()
			ios.SetStdoutTTY(tt.isTTY)
			tt.opts.IO = ios
>>>>>>> origin/trunk
			tt.opts.BaseRepo = func() (ghrepo.Interface, error) {
				return ghrepo.FromFullName("OWNER/REPO")
			}

			err := listRun(tt.opts)
<<<<<<< HEAD
			assert.NoError(t, err)

			assert.Equal(t, tt.wantOut, stdout.String())
			assert.Equal(t, tt.wantErrOut, stderr.String())
			reg.Verify(t)
		})
	}
}

const longRunOutput = "-  cool commit  0    trunk  push  0\n-  cool commit  1    trunk  push  1\n-  cool commit  2    trunk  push  2\n-  cool commit  3    trunk  push  3\n-  cool commit  4    trunk  push  4\n-  cool commit  5    trunk  push  5\n-  cool commit  6    trunk  push  6\n-  cool commit  7    trunk  push  7\n-  cool commit  8    trunk  push  8\n-  cool commit  9    trunk  push  9\n-  cool commit  10   trunk  push  10\n-  cool commit  11   trunk  push  11\n-  cool commit  12   trunk  push  12\n-  cool commit  13   trunk  push  13\n-  cool commit  14   trunk  push  14\n-  cool commit  15   trunk  push  15\n-  cool commit  16   trunk  push  16\n-  cool commit  17   trunk  push  17\n-  cool commit  18   trunk  push  18\n-  cool commit  19   trunk  push  19\n-  cool commit  20   trunk  push  20\n-  cool commit  21   trunk  push  21\n-  cool commit  22   trunk  push  22\n-  cool commit  23   trunk  push  23\n-  cool commit  24   trunk  push  24\n-  cool commit  25   trunk  push  25\n-  cool commit  26   trunk  push  26\n-  cool commit  27   trunk  push  27\n-  cool commit  28   trunk  push  28\n-  cool commit  29   trunk  push  29\n-  cool commit  30   trunk  push  30\n-  cool commit  31   trunk  push  31\n-  cool commit  32   trunk  push  32\n-  cool commit  33   trunk  push  33\n-  cool commit  34   trunk  push  34\n-  cool commit  35   trunk  push  35\n-  cool commit  36   trunk  push  36\n-  cool commit  37   trunk  push  37\n-  cool commit  38   trunk  push  38\n-  cool commit  39   trunk  push  39\n-  cool commit  40   trunk  push  40\n-  cool commit  41   trunk  push  41\n-  cool commit  42   trunk  push  42\n-  cool commit  43   trunk  push  43\n-  cool commit  44   trunk  push  44\n-  cool commit  45   trunk  push  45\n-  cool commit  46   trunk  push  46\n-  cool commit  47   trunk  push  47\n-  cool commit  48   trunk  push  48\n-  cool commit  49   trunk  push  49\n-  cool commit  50   trunk  push  50\n-  cool commit  51   trunk  push  51\n-  cool commit  52   trunk  push  52\n-  cool commit  53   trunk  push  53\n-  cool commit  54   trunk  push  54\n-  cool commit  55   trunk  push  55\n-  cool commit  56   trunk  push  56\n-  cool commit  57   trunk  push  57\n-  cool commit  58   trunk  push  58\n-  cool commit  59   trunk  push  59\n-  cool commit  60   trunk  push  60\n-  cool commit  61   trunk  push  61\n-  cool commit  62   trunk  push  62\n-  cool commit  63   trunk  push  63\n-  cool commit  64   trunk  push  64\n-  cool commit  65   trunk  push  65\n-  cool commit  66   trunk  push  66\n-  cool commit  67   trunk  push  67\n-  cool commit  68   trunk  push  68\n-  cool commit  69   trunk  push  69\n-  cool commit  70   trunk  push  70\n-  cool commit  71   trunk  push  71\n-  cool commit  72   trunk  push  72\n-  cool commit  73   trunk  push  73\n-  cool commit  74   trunk  push  74\n-  cool commit  75   trunk  push  75\n-  cool commit  76   trunk  push  76\n-  cool commit  77   trunk  push  77\n-  cool commit  78   trunk  push  78\n-  cool commit  79   trunk  push  79\n-  cool commit  80   trunk  push  80\n-  cool commit  81   trunk  push  81\n-  cool commit  82   trunk  push  82\n-  cool commit  83   trunk  push  83\n-  cool commit  84   trunk  push  84\n-  cool commit  85   trunk  push  85\n-  cool commit  86   trunk  push  86\n-  cool commit  87   trunk  push  87\n-  cool commit  88   trunk  push  88\n-  cool commit  89   trunk  push  89\n-  cool commit  90   trunk  push  90\n-  cool commit  91   trunk  push  91\n-  cool commit  92   trunk  push  92\n-  cool commit  93   trunk  push  93\n-  cool commit  94   trunk  push  94\n-  cool commit  95   trunk  push  95\n-  cool commit  96   trunk  push  96\n-  cool commit  97   trunk  push  97\n-  cool commit  98   trunk  push  98\n-  cool commit  99   trunk  push  99\n-  cool commit  100  trunk  push  100\n\nFor details on a run, try: gh run view <run-id>\n"
=======
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantOut, stdout.String())
			assert.Equal(t, tt.wantErrOut, stderr.String())
		})
	}
}
>>>>>>> origin/trunk
