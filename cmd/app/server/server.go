package server

import (
	"github.com/Tihmmm/mr-decorator-core/decorator"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator/cmd/opts"
	"github.com/Tihmmm/mr-decorator/internal/server"
	"github.com/spf13/cobra"
	"log"
)

var (
	port         string
	apiKey       string
	promptApiKey bool
)

func NewCmd(opts *opts.CmdOpts) *cobra.Command {
	d := decorator.NewDecorator(decorator.ModeServer, opts.Cfg.Decorator, opts.C)

	run := func(cmd *cobra.Command, args []string) {
		prsrs := parser.List()
		for _, k := range prsrs {
			prsr, _ := parser.Get(k)
			prsr.SetConfig(&opts.Cfg.Parser)
		}

		if opts.Cfg.Server.ApiKey == "" {
			opts.Cfg.Server.ApiKey = apiKey
		}

		s := server.NewEchoServer(opts.Cfg.Server, opts.V, d)
		if err := s.Start(port); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launches decorator in server mode.",
		Run:   run,
	}

	initArgs(cmd)

	return cmd
}

func initArgs(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&port, "port", "p", "3000", "Server port")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "Server api key")
	cmd.Flags().BoolVarP(&promptApiKey, "prompt-api-key", "a", false, "Prompt for server api key")
	cmd.MarkFlagsMutuallyExclusive("api-key", "prompt-api-key")
}
