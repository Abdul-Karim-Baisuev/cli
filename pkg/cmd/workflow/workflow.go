package workflow

import (
<<<<<<< HEAD
	cmdList "github.com/cli/cli/pkg/cmd/workflow/list"
	"github.com/cli/cli/pkg/cmdutil"
=======
	cmdDisable "github.com/cli/cli/v2/pkg/cmd/workflow/disable"
	cmdEnable "github.com/cli/cli/v2/pkg/cmd/workflow/enable"
	cmdList "github.com/cli/cli/v2/pkg/cmd/workflow/list"
	cmdRun "github.com/cli/cli/v2/pkg/cmd/workflow/run"
	cmdView "github.com/cli/cli/v2/pkg/cmd/workflow/view"
	"github.com/cli/cli/v2/pkg/cmdutil"
>>>>>>> origin/trunk
	"github.com/spf13/cobra"
)

func NewCmdWorkflow(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
<<<<<<< HEAD
		Use:   "workflow <command>",
		Short: "View details about GitHub Actions workflows",
		Long:  "List, view, and run workflows in GitHub Actions.",
		// TODO i'd like to have all the actions commands sorted into their own zone which i think will
		// require a new annotation
=======
		Use:     "workflow <command>",
		Short:   "View details about GitHub Actions workflows",
		Long:    "List, view, and run workflows in GitHub Actions.",
		GroupID: "actions",
>>>>>>> origin/trunk
	}
	cmdutil.EnableRepoOverride(cmd, f)

	cmd.AddCommand(cmdList.NewCmdList(f, nil))
<<<<<<< HEAD
=======
	cmd.AddCommand(cmdEnable.NewCmdEnable(f, nil))
	cmd.AddCommand(cmdDisable.NewCmdDisable(f, nil))
	cmd.AddCommand(cmdView.NewCmdView(f, nil))
	cmd.AddCommand(cmdRun.NewCmdRun(f, nil))
>>>>>>> origin/trunk

	return cmd
}
