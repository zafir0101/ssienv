package connection

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/connection/accept"
	"github.com/zafir0101/SSI-ENV/cmd/connection/create"
	"github.com/zafir0101/SSI-ENV/cmd/connection/list"
)

var (
	ConnectionCmd = &cobra.Command{
		Use:   "connection",
		Short: "",
	}
)

func init() {
	ConnectionCmd.AddCommand(create.CreateCmd)
	ConnectionCmd.AddCommand(accept.AcceptCmd)
	ConnectionCmd.AddCommand(list.ListCmd)
}
