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
	d := decorator.NewDecorator(decorator.ModeCli, opts.Cfg.Decorator, opts.C)

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
			if err := pkg.ReadSecretStdinToString(&mr.AuthToken); err != nil {
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
		prsr.SetConfig(&opts.Cfg.Parser)
		if err := d.Decorate(mr, prsr); err != nil {
			log.Fatalf("Error decorating: %s", err)
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
	cmd.Flags().StringVarP(&authToken, "token", "t", "", "Gitlab auth token with `api` scope")
	cmd.Flags().BoolVarP(&promptToken, "prompt-token", "p", false, "Prompt for Gitlab token")
	cmd.MarkFlagsOneRequired("token", "prompt-token")
	cmd.MarkFlagsMutuallyExclusive("token", "prompt-token")
	cmd.Flags().IntVarP(&vulnerabilityMgmtId, "vid", "v", -1, "Some identifier in your vulnerability management system")
	cmd.Flags().StringVarP(&path, "file", "f", "", "Path to locally stored report file")
	cmd.Flags().IntVar(&projectId, "project-id", -1, "Gitlab project ID")
	cmd.MarkFlagsOneRequired("file", "project-id")
	cmd.MarkFlagsMutuallyExclusive("file", "project-id")
	cmd.Flags().IntVar(&jobId, "job-id", -1, "Gitlab job id")
	cmd.Flags().StringVarP(&artifactFormat, "artifact-format", "a", "", "Format of report file")
	cmd.MarkFlagRequired("artifact-format")
	cmd.Flags().StringVar(&artifactFileName, "artifact-file", "", "Filename of artifact")
	cmd.Flags().IntVar(&mergeRequestIid, "mr-iid", -1, "Merge request internal ID")
	cmd.MarkFlagsRequiredTogether("project-id", "job-id", "mr-iid")
}
