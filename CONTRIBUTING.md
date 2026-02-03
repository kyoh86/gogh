# Contributing

## Local gogh binary with mise

If you use mise, this repo is configured to prefer a locally built `gogh`
binary over the installed release by adding `.mise/gogh` to `PATH` when you
are inside this project.

Setup:

1. Enter this repo and trust the config if prompted: `mise trust`
2. Build the local binary: `mise run build`

After that, `gogh` should resolve to `.mise/gogh/gogh` while you are in this
repository.
