package list

import (
	"fmt"
	"net/http"
<<<<<<< HEAD

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmd/run/shared"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
=======
	"time"

	"github.com/MakeNowJust/heredoc"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmd/run/shared"
	workflowShared "github.com/cli/cli/v2/pkg/cmd/workflow/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
>>>>>>> origin/trunk
	"github.com/spf13/cobra"
)

const (
<<<<<<< HEAD
	defaultLimit = 10
=======
	defaultLimit = 20
>>>>>>> origin/trunk
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (ghrepo.Interface, error)
<<<<<<< HEAD

	PlainOutput bool

	Limit int
=======
	Prompter   iprompter

	Exporter cmdutil.Exporter

	Limit            int
	WorkflowSelector string
	Branch           string
	Actor            string
	Status           string
	Event            string
	Created          string
	Commit           string
	All              bool

	now time.Time
}

type iprompter interface {
	Select(string, string, []string) (int, error)
>>>>>>> origin/trunk
}

func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
<<<<<<< HEAD
	}

	cmd := &cobra.Command{
		Use:    "list",
		Short:  "List recent workflow runs",
		Args:   cobra.NoArgs,
		Hidden: true,
=======
		Prompter:   f.Prompter,
		now:        time.Now(),
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recent workflow runs",
		Long: heredoc.Docf(`
			List recent workflow runs.

			Note that providing the %[1]sworkflow_name%[1]s to the %[1]s-w%[1]s flag will not fetch disabled workflows.
			Also pass the %[1]s-a%[1]s flag to fetch disabled workflow runs using the %[1]sworkflow_name%[1]s and the %[1]s-w%[1]s flag.

			Runs created by organization and enterprise ruleset workflows will not display a workflow name due to GitHub API limitations.

			To see runs associated with a pull request, users should run %[1]sgh pr checks%[1]s.
		`, "`"),
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
>>>>>>> origin/trunk
		RunE: func(cmd *cobra.Command, args []string) error {
			// support `-R, --repo` override
			opts.BaseRepo = f.BaseRepo

<<<<<<< HEAD
			terminal := opts.IO.IsStdoutTTY() && opts.IO.IsStdinTTY()
			opts.PlainOutput = !terminal

			if opts.Limit < 1 {
				return &cmdutil.FlagError{Err: fmt.Errorf("invalid limit: %v", opts.Limit)}
=======
			if opts.Limit < 1 {
				return cmdutil.FlagErrorf("invalid limit: %v", opts.Limit)
>>>>>>> origin/trunk
			}

			if runF != nil {
				return runF(opts)
			}

			return listRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", defaultLimit, "Maximum number of runs to fetch")
<<<<<<< HEAD
=======
	cmd.Flags().StringVarP(&opts.WorkflowSelector, "workflow", "w", "", "Filter runs by workflow")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Filter runs by branch")
	cmd.Flags().StringVarP(&opts.Actor, "user", "u", "", "Filter runs by user who triggered the run")
	cmd.Flags().StringVarP(&opts.Event, "event", "e", "", "Filter runs by which `event` triggered the run")
	cmd.Flags().StringVarP(&opts.Created, "created", "", "", "Filter runs by the `date` it was created")
	cmd.Flags().StringVarP(&opts.Commit, "commit", "c", "", "Filter runs by the `SHA` of the commit")
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Include disabled workflows")
	cmdutil.StringEnumFlag(cmd, &opts.Status, "status", "s", "", shared.AllStatuses, "Filter runs by status")
	cmdutil.AddJSONFlags(cmd, &opts.Exporter, shared.RunFields)

	_ = cmdutil.RegisterBranchCompletionFlags(f.GitClient, cmd, "branch")
>>>>>>> origin/trunk

	return cmd
}

func listRun(opts *ListOptions) error {
	baseRepo, err := opts.BaseRepo()
	if err != nil {
		return fmt.Errorf("failed to determine base repo: %w", err)
	}

	c, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create http client: %w", err)
	}
	client := api.NewClientFromHTTP(c)

<<<<<<< HEAD
	opts.IO.StartProgressIndicator()
	runs, err := shared.GetRuns(client, baseRepo, opts.Limit)
=======
	filters := &shared.FilterOptions{
		Branch:  opts.Branch,
		Actor:   opts.Actor,
		Status:  opts.Status,
		Event:   opts.Event,
		Created: opts.Created,
		Commit:  opts.Commit,
	}

	opts.IO.StartProgressIndicator()
	defer opts.IO.StopProgressIndicator()

	if opts.WorkflowSelector != "" {
		// initially the workflow state is limited to 'active'
		states := []workflowShared.WorkflowState{workflowShared.Active}
		if opts.All {
			// the all flag tells us to add the remaining workflow states
			// note: this will be incomplete if more workflow states are added to `workflowShared`
			states = append(states, workflowShared.DisabledManually, workflowShared.DisabledInactivity)
		}
		if workflow, err := workflowShared.ResolveWorkflow(opts.Prompter, opts.IO, client, baseRepo, false, opts.WorkflowSelector, states); err == nil {
			filters.WorkflowID = workflow.ID
			filters.WorkflowName = workflow.Name
		} else {
			return err
		}
	}
	runsResult, err := shared.GetRuns(client, baseRepo, filters, opts.Limit)
>>>>>>> origin/trunk
	opts.IO.StopProgressIndicator()
	if err != nil {
		return fmt.Errorf("failed to get runs: %w", err)
	}
<<<<<<< HEAD

	tp := utils.NewTablePrinter(opts.IO)

	cs := opts.IO.ColorScheme()

	if len(runs) == 0 {
		if !opts.PlainOutput {
			fmt.Fprintln(opts.IO.ErrOut, "No runs found")
		}
		return nil
	}

	out := opts.IO.Out

	for _, run := range runs {
		if opts.PlainOutput {
			tp.AddField(string(run.Status), nil, nil)
			tp.AddField(string(run.Conclusion), nil, nil)
		} else {
			symbol, symbolColor := shared.Symbol(cs, run.Status, run.Conclusion)
			tp.AddField(symbol, nil, symbolColor)
		}

		tp.AddField(run.CommitMsg(), nil, cs.Bold)

		tp.AddField(run.Name, nil, nil)
		tp.AddField(run.HeadBranch, nil, cs.Bold)
		tp.AddField(string(run.Event), nil, nil)

		if opts.PlainOutput {
			elapsed := run.UpdatedAt.Sub(run.CreatedAt)
			if elapsed < 0 {
				elapsed = 0
			}
			tp.AddField(elapsed.String(), nil, nil)
		}

		tp.AddField(fmt.Sprintf("%d", run.ID), nil, cs.Cyan)

		tp.EndRow()
	}

	err = tp.Render()
	if err != nil {
		return err
	}

	if !opts.PlainOutput {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "For details on a run, try: gh run view <run-id>")
	}

=======
	runs := runsResult.WorkflowRuns
	if len(runs) == 0 && opts.Exporter == nil {
		return cmdutil.NewNoResultsError("no runs found")
	}

	if err := opts.IO.StartPager(); err == nil {
		defer opts.IO.StopPager()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "failed to start pager: %v\n", err)
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, runs)
	}

	tp := tableprinter.New(opts.IO, tableprinter.WithHeader("STATUS", "TITLE", "WORKFLOW", "BRANCH", "EVENT", "ID", "ELAPSED", "AGE"))

	cs := opts.IO.ColorScheme()

	for _, run := range runs {
		if tp.IsTTY() {
			symbol, symbolColor := shared.Symbol(cs, run.Status, run.Conclusion)
			tp.AddField(symbol, tableprinter.WithColor(symbolColor))
		} else {
			tp.AddField(string(run.Status))
			tp.AddField(string(run.Conclusion))
		}
		tp.AddField(run.Title(), tableprinter.WithColor(cs.Bold))
		tp.AddField(run.WorkflowName())
		tp.AddField(run.HeadBranch, tableprinter.WithColor(cs.Bold))
		tp.AddField(string(run.Event))
		tp.AddField(fmt.Sprintf("%d", run.ID), tableprinter.WithColor(cs.Cyan))
		tp.AddField(run.Duration(opts.now).String())
		tp.AddTimeField(opts.now, run.StartedTime(), cs.Muted)
		tp.EndRow()
	}

	if err := tp.Render(); err != nil {
		return err
	}

>>>>>>> origin/trunk
	return nil
}
