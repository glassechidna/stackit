# `stackit`

`stackit` is a CLI tool to synchronously and idempotently operate on AWS
CloudFormation stacks - a perfect complement for continuous integration systems
and developers who prefer the comfort of the command line.

![Animated gif demonstrating functionality](https://user-images.githubusercontent.com/369053/51422660-92d00a80-1c06-11e9-81b2-da296b36cf6d.gif)

## Where

The latest compiled `stackit` binaries for Linux, macOS and Windows can be
downloaded from the project's [**GitHub Releases**](https://github.com/glassechidna/stackit/releases)
page.

## Why

CloudFormation is inherently asychronous and this is reflected in the usage
of the AWS CLI tools - a `create-stack` or `update-stack` operation exits long
before the stack has reached its final state. `stackit` treats a stack update
synchronously, streaming stack events to the CLI until the stack reaches a
steady state.

AWS CLI commands for CloudFormation aren't idempotent. If you call `create-stack`
when a stack already exists, the behaviour is different to if it doesn't.
Likewise with `update-stack`. This means you either have to manually create a
stack before putting it under CI, or script up a "does it exist yet?" check
before deciding which command to invoke. `stackit` abstracts over these with
an `up` facade.

## How

### `up`

```
stackit up --stack-name some-other-name # use this stack name, fallback to yml for rest
stackit up \
  --stack-name some-other-name \
  --template sample.yml \
  --param-value DockerImage=redis \
  --param-value Cluster=some-ecs-cluster # no yml necessary
```

Note that there is JSON printed at the end of the `up` command. This is all the
_Outputs_ defined in your CloudFormation template file. These are printed to
stdout. The event lines above them are printed to stderr.

This separation makes it easy to pipe output from `stackit up` to another
command without having to skip the log lines. Likewise, a non-zero exit code
indicates stack update/creation failure.

### `outputs`

`stackit outputs --stack-name <name>` prints the stack's Outputs in JSON form,
without making any modifications to the stack.

### `tail`

If an existing stack creation or update is in progress, `stackit tail --stack-name <name>`
will poll for events, similar to the `up` command.

### `down`

`stackit down --stack-name <name>` will delete the named stack if it exists,
otherwise it will do nothing. Non-zero exit code indicates failure to delete
an existing stack.

### More

All commands can be passed a `--profile <name>` parameter. This will use alternative
AWS credentials defined in a profile named in `~/.aws/config` if it exists. If your
profile requires MFA credentials in order to assume a role, `stackit` will prompt
for those to be entered on `stdin`.

All commands can be passed a `--region <region>` parameter if you want to deploy
your stack in a different region.

## TODO

* `stackit <stack-name> cancel`
* `stackit <stack-name> signal <logical-name>`

## Additional Flags

* `--stackit-s3-suffix "random"`
If the s3 bucket stackit creates requires a unique S3 bucket name.

* `--stackit-tags "key=val,key2=val2"`
If the s3 bucket stackit creates requires tags, you can add them with this flag.

TODO: Document these properly

* `--service-role VAL`
* `--previous-param-value NAME`
* `--tag NAME=VAL` (multiple)
* `--notification-arn` (multiple)
* `--stack-policy VAL`
* `--previous-template`
* `--no-cancel-on-exit`
* `--no-destroy` (not yet implemented)

for changes: (not yet implemented)
* `--name VAL`
* `--execute-if-no-destroy`

### Notes

* Change-sets return special exit code to indicate destructive (replacement,
  deletion) actions
* MFA support
