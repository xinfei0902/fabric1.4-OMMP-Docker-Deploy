package dcmd

import (
	"github.com/spf13/cobra"
)

var globalCommand = &cobra.Command{
	Use:   "root",
	Short: "root",
	Long:  "root",
}

func init() {
	// append service
	globalCommand.AddCommand(initService())
	//globalCommand.AddCommand(initVersion())
}

//Execute 子命令
func Execute() error {
	return globalCommand.Execute()
}
