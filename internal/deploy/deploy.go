package deploy

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	AWS_CLI            = "aws"
	AWS_CLI_CFN        = "cloudformation"
	AWS_CLI_CFN_DEPLOY = "deploy"
)

var execCommand = exec.Command

type DeployCommand struct {
	Args map[string]string
}

func New() *DeployCommand {
	return &DeployCommand{
		Args: make(map[string]string, 0),
	}
}

func (c *DeployCommand) AwsCliArgs() []string {
	return append([]string{AWS_CLI_CFN, AWS_CLI_CFN_DEPLOY}, c.args()...)
}

// Execute the `aws cloudformation deploy` command with the given args,
// returning the command along with it's stdout and stderr pipes.
func (c *DeployCommand) Execute() (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	cmd := execCommand(AWS_CLI, c.AwsCliArgs()...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd, stdout, nil, err
	}

	return cmd, stdout, stderr, err
}

func (c *DeployCommand) SetStackName(stackName string) *DeployCommand {
	c.Args["stack-name"] = stackName
	return c
}

func (c *DeployCommand) SetTemplateFile(template string) *DeployCommand {
	c.Args["template-file"] = template
	return c
}

func (c *DeployCommand) SetProfile(profile string) *DeployCommand {
	c.Args["profile"] = profile
	return c
}

func (c *DeployCommand) SetParameterOverrides(params []string) *DeployCommand {
	c.Args["parameter-overrides"] = strings.Join(params, " ")
	return c
}

func (c *DeployCommand) SetDebug() *DeployCommand {
	c.Args["debug"] = ""
	return c
}

func (c *DeployCommand) SetKmsKeyID(id string) *DeployCommand {
	c.Args["kms-key-id"] = id
	return c
}

func (c *DeployCommand) SetS3Bucket(bucket string) *DeployCommand {
	c.Args["s3-bucket"] = bucket
	return c
}

func (c *DeployCommand) SetS3Prefix(prefix string) *DeployCommand {
	c.Args["s3-prefix"] = prefix
	return c
}

func (c *DeployCommand) SetForceUpload() *DeployCommand {
	c.Args["force-upload"] = ""
	return c
}

func (c *DeployCommand) args() (args []string) {
	for k, v := range c.Args {
		args = append(args, fmt.Sprintf("--%s", k))
		if v != "" {
			args = append(args, v)
		}
	}
	return
}
