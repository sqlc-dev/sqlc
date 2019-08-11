workflow "dinosql test suite" {
  on = "push"
  resolves = ["Setup Go for use with actions"]
}

action "Setup Go for use with actions" {
  uses = "actions/setup-go@419ae75c254126fa6ae3e3ef573ce224a919b8fe"
}
