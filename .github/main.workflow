workflow "Release" {
  on = "push"
  resolves = ["goreleaser"]
}

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

action "test" {
  uses = "docker://golang:1.12"
  args = "go test ./..."
  needs = ["not-tag"]
}

action "goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  secrets = [
    "GORELEASER_GITHUB_TOKEN"
  ]
  args = "release"
  needs = ["is-tag"]
}
