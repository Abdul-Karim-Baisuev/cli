package list

import (
	"fmt"
	"net/http"

<<<<<<< HEAD
	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
	"github.com/spf13/cobra"
)

const (
	defaultLimit = 10

	Active           WorkflowState = "active"
	DisabledManually WorkflowState = "disabled_manually"
)
=======
	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmd/workflow/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

const defaultLimit = 50
>>>>>>> origin/trunk

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (ghrepo.Interface, error)
<<<<<<< HEAD

	PlainOutput bool
=======
	Exporter   cmdutil.Exporter
>>>>>>> origin/trunk

	All   bool
	Limit int
}

<<<<<<< HEAD
=======
var workflowFields = []string{
	"id",
	"name",
	"path",
	"state",
}

>>>>>>> origin/trunk
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
<<<<<<< HEAD
		Use:    "list",
		Short:  "List GitHub Actions workflows",
		Args:   cobra.NoArgs,
		Hidden: true,
=======
		Use:     "list",
		Short:   "List workflows",
		Long:    "List workflow files, hiding disabled workflows by default.",
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

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", defaultLimit, "Maximum number of workflows to fetch")
<<<<<<< HEAD
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Show all workflows, including disabled workflows")

=======
	cmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Include disabled workflows")
	cmdutil.AddJSONFlags(cmd, &opts.Exporter, workflowFields)
>>>>>>> origin/trunk
	return cmd
}

func listRun(opts *ListOptions) error {
	repo, err := opts.BaseRepo()
	if err != nil {
<<<<<<< HEAD
		return fmt.Errorf("could not determine base repo: %w", err)
=======
		return err
>>>>>>> origin/trunk
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("could not create http client: %w", err)
	}
	client := api.NewClientFromHTTP(httpClient)

	opts.IO.StartProgressIndicator()
<<<<<<< HEAD
	workflows, err := getWorkflows(client, repo, opts.Limit)
=======
	workflows, err := shared.GetWorkflows(client, repo, opts.Limit)
>>>>>>> origin/trunk
	opts.IO.StopProgressIndicator()
	if err != nil {
		return fmt.Errorf("could not get workflows: %w", err)
	}

<<<<<<< HEAD
	if len(workflows) == 0 {
		if !opts.PlainOutput {
			fmt.Fprintln(opts.IO.ErrOut, "No workflows found")
		}
		return nil
	}

	tp := utils.NewTablePrinter(opts.IO)
	cs := opts.IO.ColorScheme()

	for _, workflow := range workflows {
		if workflow.Disabled() && !opts.All {
			continue
		}
		tp.AddField(workflow.Name, nil, cs.Bold)
		tp.AddField(string(workflow.State), nil, nil)
		tp.AddField(fmt.Sprintf("%d", workflow.ID), nil, cs.Cyan)
=======
	var filteredWorkflows []shared.Workflow
	if opts.All {
		filteredWorkflows = workflows
	} else {
		for _, workflow := range workflows {
			if !workflow.Disabled() {
				filteredWorkflows = append(filteredWorkflows, workflow)
			}
		}
	}

	if len(filteredWorkflows) == 0 {
		return cmdutil.NewNoResultsError("no workflows found")
	}

	if err := opts.IO.StartPager(); err == nil {
		defer opts.IO.StopPager()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "failed to start pager: %v\n", err)
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, filteredWorkflows)
	}

	cs := opts.IO.ColorScheme()
	tp := tableprinter.New(opts.IO, tableprinter.WithHeader("Name", "State", "ID"))

	for _, workflow := range filteredWorkflows {
		tp.AddField(workflow.Name)
		tp.AddField(string(workflow.State))
		tp.AddField(fmt.Sprintf("%d", workflow.ID), tableprinter.WithColor(cs.Cyan))
>>>>>>> origin/trunk
		tp.EndRow()
	}

	return tp.Render()
}
<<<<<<< HEAD

type WorkflowState string

type Workflow struct {
	Name  string
	ID    int
	State WorkflowState
}

func (w *Workflow) Disabled() bool {
	return w.State != Active
}

type WorkflowsPayload struct {
	Workflows []Workflow
}

func getWorkflows(client *api.Client, repo ghrepo.Interface, limit int) ([]Workflow, error) {
	perPage := limit
	page := 1
	if limit > 100 {
		perPage = 100
	}

	workflows := []Workflow{}

	for len(workflows) < limit {
		var result WorkflowsPayload

		path := fmt.Sprintf("repos/%s/actions/workflows?per_page=%d&page=%d", ghrepo.FullName(repo), perPage, page)

		err := client.REST(repo.RepoHost(), "GET", path, nil, &result)
		if err != nil {
			return nil, err
		}

		if len(result.Workflows) == 0 {
			break
		}

		for _, workflow := range result.Workflows {
			workflows = append(workflows, workflow)
			if len(workflows) == limit {
				break
			}
		}

		if len(result.Workflows) < perPage {
			break
		}

		page++
	}

	return workflows, nil
}
=======
>>>>>>> origin/trunk
