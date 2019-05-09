workflow "CI" {
  on = "push"
  resolves = ["test"]
}

action "is-tag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "not-tag" {
  uses = "actions/bin/filter@master"
  args = "not tag"
}

action "short test" {
  uses = "./"
  args = "go test -test.short ./..."
  needs = ["not-tag"]
}

action "test" {
  uses = "./"
  args = "go test ./..."
  env = {
    AWS_REGION = "ap-southeast-2"
  }
  secrets = [
    "AWS_ACCESS_KEY_ID",
    "AWS_SECRET_ACCESS_KEY"
  ]
  needs = ["short test"]
}

action "goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  secrets = [
    "GORELEASER_GITHUB_TOKEN"
  ]
  args = "release"
  needs = ["is-tag"]
}
