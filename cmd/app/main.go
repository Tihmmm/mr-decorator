package main

import (
	"github.com/Tihmmm/mr-decorator/internal/client"
	"github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/internal/decorator"
	"github.com/Tihmmm/mr-decorator/internal/models"
	"github.com/Tihmmm/mr-decorator/internal/parser"
	"github.com/Tihmmm/mr-decorator/internal/server"
	"github.com/Tihmmm/mr-decorator/internal/validator"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var (
	mode                string // either `cli` or `server`
	authToken           string
	promptToken         bool
	vulnerabilityMgmtId int
	path                string
	projectId           int
	jobId               int
	artifactFormat      string
	artifactFileName    string
	mergeRequestIid     int

	port string

	rootCmd = &cobra.Command{
		Use:     "",
		Version: "2.0",
		Long: `A merge request decorator for Gitlab. Can be used in either 'cli' or 'server' mode.
In either mode don't forget to set the following environment variables:
	SCA_VULN_MGMT_PROJECT_BASE_URL
	SCA_VULN_MGMT_INSTANCE_SUBPATH_TEMPLATE
	SCA_VULN_MGMT_REPORT_SUBPATH_TEMPLATE
	SAST_VULN_MGMT_PROJECT_BASE_URL,unset"          // e.g. https://fortify-ssc.company.com/html/ssc/version/%d
	SAST_VULN_MGMT_INSTANCE_SUBPATH_TEMPLATE		// e.g. audit?q=instance_id%3A
	SAST_VULN_MGMT_REPORT_SUBPATH_TEMPLATE			// e.g. audit?q=analysis_type%3Asca
	GITLAB_IP
	GITLAB_DOMAIN
        `,
		Run: run,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "server", "Accepts either `cli` or `server`")
	switch mode {
	case "cli":
		rootCmd.Flags().StringVarP(&authToken, "token", "t", "", "Gitlab auth token with `api` scope")
		rootCmd.Flags().BoolVarP(&promptToken, "prompt-token", "p", false, "Prompt for Gitlab token")
		rootCmd.MarkFlagsOneRequired("token", "prompt-token")
		rootCmd.MarkFlagsMutuallyExclusive("token", "prompt-token")
		rootCmd.Flags().IntVarP(&vulnerabilityMgmtId, "vid", "v", -1, "Some identifier in your vulnerability management system")
		rootCmd.Flags().StringVarP(&path, "file", "f", "", "Path to locally stored report file")
		rootCmd.Flags().IntVar(&projectId, "project-id", -1, "Gitlab project ID")
		rootCmd.MarkFlagsOneRequired("file", "project-id")
		rootCmd.MarkFlagsMutuallyExclusive("file", "project-id")
		rootCmd.Flags().IntVar(&jobId, "job-id", -1, "Gitlab job id")
		rootCmd.Flags().StringVar(&artifactFormat, "artifact-format", "", "Format of report file")
		rootCmd.Flags().StringVar(&artifactFileName, "artifact-file", "", "Filename of artifact")
		rootCmd.Flags().IntVar(&mergeRequestIid, "mr-iid", -1, "Merge request internal ID")
		rootCmd.MarkFlagsRequiredTogether("project-id", "job-id", "artifact-format", "artifact-file", "mr-iid")
	case "server":
		rootCmd.Flags().StringVarP(&port, "port", "p", "-1", "Server port. If not specified, it will use the SERVER_PORT environment variable")
	default:
		log.Fatalf("Invalid mode: %s. Only `cli` and `server` are allowed", mode)
	}
}

func run(cmd *cobra.Command, args []string) {
	cfg := config.NewConfig()
	v := validator.NewValidator()
	c := client.NewHttpClient(cfg.HttpClient)
	p := parser.NewParser(cfg.Parser)
	d := decorator.NewDecorator(mode, c, p)
	switch mode {
	case "cli":
		mr := &models.MRRequest{
			FilePath:            path,
			ProjectId:           projectId,
			JobId:               jobId,
			ArtifactFormat:      artifactFormat,
			ArtifactFileName:    artifactFileName,
			MergeRequestIid:     mergeRequestIid,
			VulnerabilityMgmtId: vulnerabilityMgmtId,
		}
		if !v.IsValidAll(mr) {
			os.Exit(127)
		}
		if promptToken {
			tokenBytes, err := terminal.ReadPassword(0)
			if err != nil {
				log.Fatal(err)
			}
			authToken = string(tokenBytes)
		}
		mr.AuthToken = authToken
		if err := d.DecorateCli(mr); err != nil {
			log.Fatal(err)
		}
	case "server":
		s := server.NewEchoServer(cfg.Server, v, d)
		if err := s.Start(port); err != nil {
			log.Fatal(err)
		}
	}
}
