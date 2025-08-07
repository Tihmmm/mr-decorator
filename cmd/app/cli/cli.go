package cli

import (
	"github.com/Tihmmm/mr-decorator-core/decorator"
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator/cmd/opts"
	"github.com/Tihmmm/mr-decorator/pkg"
	"github.com/spf13/cobra"
	"log"
)

var (
	authToken           string
	promptToken         bool
	vulnerabilityMgmtId int
	path                string
	projectId           int
	jobId               int
	artifactFormat      string
	artifactFileName    string
	mergeRequestIid     int
)

func NewCmd(opts *opts.CmdOpts) *cobra.Command {
	d := decorator.NewDecorator(decorator.ModeCli, *opts.DecoratorConfig, opts.C)

	run := func(cmd *cobra.Command, args []string) {
		mr := &models.MRRequest{
			FilePath:            path,
			ProjectId:           projectId,
			JobId:               jobId,
			ArtifactFormat:      artifactFormat,
			ArtifactFileName:    artifactFileName,
			MergeRequestIid:     mergeRequestIid,
			VulnerabilityMgmtId: vulnerabilityMgmtId,
		}
		if promptToken {
			if err := pkg.ReadSecretStdinToString("Enter Gitlab auth token (scope: api)", &mr.AuthToken); err != nil {
				log.Fatalln(err)
			}
		}

		if !opts.V.IsValidAll(mr) {
			log.Fatal("Input parameters invalid")
		}
		prsr, err := parser.Get(mr.ArtifactFormat)
		if err != nil {
			log.Fatalf("Error getting parser for format `%s`: %s", mr.ArtifactFormat, err)
		}
		prsr.SetConfig(opts.ParserConfig)
		if err := d.Decorate(mr, prsr); err != nil {
			log.Fatal(err)
		}
	}

	cmd := &cobra.Command{
		Use:   "cli",
		Short: "Launches decorator in cli mode",
		Run:   run,
	}

	initArgs(cmd)

	return cmd
}

func initArgs(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&authToken, "token", "t", "", "gitlab auth token with `api` scope")
	cmd.Flags().BoolVarP(&promptToken, "prompt-token", "p", false, "prompt for Gitlab token")
	cmd.MarkFlagsOneRequired("token", "prompt-token")
	cmd.MarkFlagsMutuallyExclusive("token", "prompt-token")
	cmd.Flags().IntVarP(&vulnerabilityMgmtId, "vid", "v", -1, "some identifier in your vulnerability management system to create a links to vulnerability views")
	cmd.Flags().StringVarP(&path, "file", "f", "", "path to locally stored report file")
	cmd.Flags().IntVar(&projectId, "project-id", -1, "gitlab project ID")
	cmd.MarkFlagsOneRequired("file", "project-id")
	cmd.MarkFlagsMutuallyExclusive("file", "project-id")
	cmd.Flags().IntVar(&jobId, "job-id", -1, "gitlab job id")
	cmd.Flags().StringVarP(&artifactFormat, "artifact-format", "a", "", "format of report file")
	cmd.MarkFlagRequired("artifact-format")
	cmd.Flags().StringVar(&artifactFileName, "artifact-file", "", "filename of artifact")
	cmd.Flags().IntVar(&mergeRequestIid, "mr-iid", -1, "merge request internal ID")
	cmd.MarkFlagsRequiredTogether("project-id", "job-id", "mr-iid")
}
