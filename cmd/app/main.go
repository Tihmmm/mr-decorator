package main

import (
	"log"
	"os"

	"github.com/Tihmmm/mr-decorator-core/client"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/validator"
	"github.com/Tihmmm/mr-decorator/cmd/app/cli"
	"github.com/Tihmmm/mr-decorator/cmd/app/list"
	"github.com/Tihmmm/mr-decorator/cmd/app/server"
	"github.com/Tihmmm/mr-decorator/cmd/opts"
	"github.com/spf13/cobra"
)

var (
	configPath string

	rootCmd = &cobra.Command{
		Use:     "",
		Version: "Alpha",
		Long: `A merge request decorator for Gitlab. Can be used in either 'cli' or 'server' mode.
In either mode, don't forget to fill the configuration file.
`,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yml", "path to configuration file. configuration will be overwritten by cli arguments if provided")

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	c := client.NewGitlabClient(cfg.GitlabClient)
	v := validator.NewValidator()
	cmdOpts := &opts.CmdOpts{
		ConfigPath:      configPath,
		ParserConfig:    &cfg.Parser,
		DecoratorConfig: &cfg.Decorator,
		C:               c,
		V:               v,
	}

	rootCmd.AddCommand(list.NewCmd())
	rootCmd.AddCommand(cli.NewCmd(cmdOpts))
	rootCmd.AddCommand(server.NewCmd(cmdOpts))

	completion := completionCommand()
	completion.Hidden = true
	rootCmd.AddCommand(completion)
}

func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "no",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal("this command will not be implemented.")
		},
	}
}
