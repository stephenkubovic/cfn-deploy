# cfn-deploy (experimental)

A thin wrapper around the [`aws cloudformation deploy`][aws-cloudformation-deploy] command that provides realtime CloudFormation stack event updates while the command is running.

## Installation

```
go get github.com/stephenkubovic/cfn-deploy/cmd/cfn-deploy
```

## Usage

```
cfn-deploy -s my-stack -t stack.yml -p dev --params Key=Value
```

The `cfn-deploy` command is intended to be a near drop-in replacement for `aws cloudformation deploy`, so it accepts the same argument list. Additionally, some aliases have been added for common/required CLI options to save some keystrokes.

The following command line options are equivalent:
```
cfn-deploy -s my-stack -t stack.yml -p dev --params Key=Value

cfn-deploy --stack-name my-stack --template-file stack.yml --profile dev --parameter-overrides Key=Value

aws cloudformation deploy --stack-name my-stack --template-file stack.yml --profile dev --parameter-overrides Key=Value
```

## Roadmap

- [ ] Colourized output
- [ ] Nested stack support
- [ ] Test more CLI argument combinations

[aws-cloudformation-deploy]: https://docs.aws.amazon.com/cli/latest/reference/cloudformation/deploy/index.html
