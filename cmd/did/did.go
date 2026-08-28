package did

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/did/list"
	"github.com/zafir0101/ssienv/cmd/did/publish"
	"github.com/zafir0101/ssienv/cmd/did/resolve"
	"github.com/zafir0101/ssienv/cmd/did/update"
	"github.com/zafir0101/ssienv/cmd/did/whoami"
)

var (
	DIDCmd = &cobra.Command{
		Use: "did",
	}
)

func init() {
	DIDCmd.AddCommand(resolve.ResolveCmd)
	DIDCmd.AddCommand(update.UpdateCmd)
	DIDCmd.AddCommand(publish.PublishCmd)
	DIDCmd.AddCommand(list.ListCmd)
	DIDCmd.AddCommand(whoami.WhoamiCmd)
}
