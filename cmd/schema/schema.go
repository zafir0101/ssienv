package schema

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/schema/create"
	"github.com/zafir0101/ssienv/cmd/schema/list"
)

var (
	SchemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Digital blueprint that defines the structure, attribute names, and data types for a verifiable credential",
	}
)

func init() {
	SchemaCmd.AddCommand(create.CreateCmd)
	SchemaCmd.AddCommand(list.ListCmd)
}
