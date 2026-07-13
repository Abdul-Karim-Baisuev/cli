package run

import (
<<<<<<< HEAD
	cmdList "github.com/cli/cli/pkg/cmd/run/list"
	cmdView "github.com/cli/cli/pkg/cmd/run/view"
	"github.com/cli/cli/pkg/cmdutil"
=======
	cmdCancel "github.com/cli/cli/v2/pkg/cmd/run/cancel"
	cmdDelete "github.com/cli/cli/v2/pkg/cmd/run/delete"
	cmdDownload "github.com/cli/cli/v2/pkg/cmd/run/download"
	cmdList "github.com/cli/cli/v2/pkg/cmd/run/list"
	cmdRerun "github.com/cli/cli/v2/pkg/cmd/run/rerun"
	cmdView "github.com/cli/cli/v2/pkg/cmd/run/view"
	cmdWatch "github.com/cli/cli/v2/pkg/cmd/run/watch"
	"github.com/cli/cli/v2/pkg/cmdutil"
>>>>>>> origin/trunk
	"github.com/spf13/cobra"
)

func NewCmdRun(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
<<<<<<< HEAD
		Use:   "run <command>",
		Short: "View details about workflow runs",
		Long:  "List, view, and watch recent workflow runs from GitHub Actions.",
		// TODO i'd like to have all the actions commands sorted into their own zone which i think will
		// require a new annotation
=======
		Use:     "run <command>",
		Short:   "View details about workflow runs",
		Long:    "List, view, and watch recent workflow runs from GitHub Actions.",
		GroupID: "actions",
>>>>>>> origin/trunk
	}
	cmdutil.EnableRepoOverride(cmd, f)

	cmd.AddCommand(cmdList.NewCmdList(f, nil))
	cmd.AddCommand(cmdView.NewCmdView(f, nil))
<<<<<<< HEAD
=======
	cmd.AddCommand(cmdRerun.NewCmdRerun(f, nil))
	cmd.AddCommand(cmdDownload.NewCmdDownload(f, nil))
	cmd.AddCommand(cmdWatch.NewCmdWatch(f, nil))
	cmd.AddCommand(cmdCancel.NewCmdCancel(f, nil))
	cmd.AddCommand(cmdDelete.NewCmdDelete(f, nil))
>>>>>>> origin/trunk

	return cmd
}
