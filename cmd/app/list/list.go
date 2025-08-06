package list

import (
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/spf13/cobra"
	"log"
)

func NewCmd() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		prsrs := parser.List()
		log.Printf("Available parsers:\n%v\n", prsrs)
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "Lists available parsers.",
		Run:     run,
	}

	return cmd
}
