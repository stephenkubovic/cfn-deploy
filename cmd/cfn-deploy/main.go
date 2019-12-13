package main

import (
	"os"

	"github.com/stephenkubovic/cfn-deploy/internal/orchestrator"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "cfn-deploy",
		Usage: "get realtime stack events for your cloudformation deploy command",
		Action: func(c *cli.Context) error {
			cmd := orchestrator.Command{
				StackName:    c.String("stack-name"),
				TemplateFile: c.String("template-file"),
				Profile:      c.String("profile"),
				Params:       c.StringSlice("parameter-overrides"),
				Debug:        c.Bool("debug"),
				KmsKeyID:     c.String("kms-key-id"),
				S3Bucket:     c.String("s3-bucket"),
				S3Prefix:     c.String("s3-prefix"),
				ForceUpload:  c.Bool("force-upload"),
			}
			if err := cmd.Execute(); err != nil {
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "stack-name",
				Aliases: []string{"s"},
				Usage:   "The name of the AWS CloudFormation stack you're deploying to.",
			},
			&cli.StringFlag{
				Name:    "template-file",
				Aliases: []string{"t"},
				Usage:   "The path where your AWS CloudFormation template is located.",
			},
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "Use a specific profile from your credential file.",
			},
			&cli.StringSliceFlag{
				Name:    "parameter-overrides",
				Aliases: []string{"params"},
				Usage:   "A list of parameter structures that specify input parameters for your stack template.",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Turn on debug logging. This applies to this program and the underlying aws cli commands.",
			},
			&cli.StringFlag{
				Name:  "kms-key-id",
				Usage: "The ID of an AWS KMS key that the command uses to encrypt artifacts that are at rest in the S3 bucket.",
			},
			&cli.StringFlag{
				Name:  "s3-bucket",
				Usage: "The name of the S3 bucket where this command uploads your CloudFormation template.",
			},
			&cli.StringFlag{
				Name:  "s3-prefix",
				Usage: "A prefix name that the command adds to the artifacts' name when it uploads them to the S3 bucket.",
			},
			&cli.BoolFlag{
				Name:  "force-upload",
				Usage: "Indicates whether to override existing files in the S3 bucket.",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
