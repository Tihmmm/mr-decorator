package server

import (
	"log"

	"github.com/Tihmmm/mr-decorator-core/decorator"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator/cmd/opts"
	"github.com/Tihmmm/mr-decorator/config"
	"github.com/Tihmmm/mr-decorator/internal/server"
	"github.com/Tihmmm/mr-decorator/pkg"
	"github.com/spf13/cobra"
)

var (
	serverConfigPath string
	port             string
	apiKey           string
	promptApiKey     bool
)

func NewCmd(opts *opts.CmdOpts) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		prsrs := parser.List()
		for _, k := range prsrs {
			prsr, _ := parser.Get(k)
			prsr.SetConfig(opts.ParserConfig)
		}

		serverCfg, err := config.NewConfig(serverConfigPath)
		if err != nil {
			log.Fatalf("Error parsing server config: %s\n", err)
		}

		if serverCfg.ApiKey == "" {
			if promptApiKey {
				if err := pkg.ReadSecretStdinToString("Enter API key for this server: ", &serverCfg.ApiKey); err != nil {
					log.Fatalln(err)
				}
			}
			serverCfg.ApiKey = apiKey
		}

		d := decorator.NewDecorator(decorator.ModeServer, *opts.DecoratorConfig, opts.C)
		s := server.NewEchoServer(serverCfg, opts.V, d)

		if err := s.Start(port); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launches decorator in server mode",
		Run:   run,
	}

	initArgs(cmd, opts)

	return cmd
}

func initArgs(cmd *cobra.Command, opts *opts.CmdOpts) {
	cmd.Flags().StringVar(&serverConfigPath, "server-config", opts.ConfigPath, "path to server configuration file")
	cmd.Flags().StringVarP(&port, "port", "p", "3000", "server port")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "this server's api key. this cli option is only used if the `api_key` config field is not filled")
	cmd.Flags().BoolVarP(&promptApiKey, "prompt-api-key", "a", false, "prompt for server api key")
	cmd.MarkFlagsMutuallyExclusive("api-key", "prompt-api-key")
}
