# `stackit`

## This README is full of lies. Right now it serves as more of a TODO for the developers.

`stackit` is a CLI tool to synchronously and idempotently operate on AWS
CloudFormation stacks - a perfect complement for continuous integration systems
and developers who prefer the comfort of the command line.

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

## Usage

* `stackit <stack-name> up`
* `stackit <stack-name> down`
* `stackit <stack-name> outputs`
* `stackit <stack-name> tail`
* `stackit <stack-name> cancel`
* `stackit <stack-name> change create`
* `stackit <stack-name> change execute`
* ~~`stackit <stack-name> resources`~~
* ~~`stackit <stack-name> signal <logical-name>`~~

### Flags

global:

* `--profile VAL`
* `--region VAL`
* `--stack-name VAL`

* `--service-role VAL`
* `--param-value NAME=VAL` (multiple)
* `--previous-param-value NAME`
* `--tag NAME=VAL` (multiple)
* `--notification-arn` (multiple)
* `--stack-policy VAL`
* `--template VAL`
* `--previous-template`
* `--no-destroy`
* `--no-cancel-on-exit`

for changes:
* `--name VAL`
* `--execute-if-no-destroy`

### Notes

* Change-sets return special exit code to indicate destructive (replacement,
  deletion) actions
* Human output to stderr, machine output to stdout
* Pipeable JSON everywhere
* MFA support

## Installation

