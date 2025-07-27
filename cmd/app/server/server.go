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
	port string
)

func NewCmd(opts *opts.CmdOpts) *cobra.Command {
	d := decorator.NewDecorator(decorator.ModeServer, opts.Cfg.Decorator, opts.C)
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launches decorator in cli mode",
		Run: func(cmd *cobra.Command, args []string) {
			prsrs := parser.List()
			for _, k := range prsrs {
				prsr, _ := parser.Get(k)
				prsr.SetConfig(&opts.Cfg.Parser)
			}
			s := server.NewEchoServer(opts.Cfg.Server, opts.V, d)
			if err := s.Start(port); err != nil {
				log.Fatalf("Error starting server: %s", err)
			}
		},
	}

	initArgs(cmd)

	return cmd
}

func initArgs(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&port, "port", "p", "-1", "Server port. If not specified, it will use the SERVER_PORT environment variable")
	cmd.MarkFlagRequired("port")
}
