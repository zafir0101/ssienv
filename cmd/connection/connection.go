package connection

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/connection/accept"
	"github.com/zafir0101/ssienv/cmd/connection/create"
	"github.com/zafir0101/ssienv/cmd/connection/list"
)

var (
	ConnectionCmd = &cobra.Command{
		Use:   "connection",
		Short: "Peer-to-Peer DIDComm connection",
	}
)

func init() {
	ConnectionCmd.AddCommand(create.CreateCmd)
	ConnectionCmd.AddCommand(accept.AcceptCmd)
	ConnectionCmd.AddCommand(list.ListCmd)
}
