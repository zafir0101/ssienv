package schema

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/schema/create"
	"github.com/zafir0101/SSI-ENV/cmd/schema/list"
)

var (
	SchemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "",
	}
)

func init() {
	SchemaCmd.AddCommand(create.CreateCmd)
	SchemaCmd.AddCommand(list.ListCmd)
}
