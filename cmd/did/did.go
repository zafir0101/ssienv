package did

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/did/list"
	"github.com/zafir0101/SSI-ENV/cmd/did/publish"
	"github.com/zafir0101/SSI-ENV/cmd/did/resolve"
	"github.com/zafir0101/SSI-ENV/cmd/did/update"
	"github.com/zafir0101/SSI-ENV/cmd/did/whoami"
)

var (
	DIDCmd = &cobra.Command{
		Use:   "did",
		Short: "",
	}
)

func init() {
	DIDCmd.AddCommand(resolve.ResolveCmd)
	DIDCmd.AddCommand(update.UpdateCmd)
	DIDCmd.AddCommand(publish.PublishCmd)
	DIDCmd.AddCommand(list.ListCmd)
	DIDCmd.AddCommand(whoami.WhoamiCmd)
}
