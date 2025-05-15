# Gogh

Gogh is a tool to manage GitHub repositories efficiently, inspired by [`ghq`](https://github.com/motemen/ghq).

[![`PkgGoDev`](https://pkg.go.dev/badge/kyoh86/gogh)](https://pkg.go.dev/kyoh86/gogh)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gogh)](https://goreportcard.com/report/github.com/kyoh86/gogh)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/gogh.svg)](https://codecov.io/gh/kyoh86/gogh)
[![GitHub release](https://github.com/kyoh86/gogh/actions/workflows/release.yml/badge.svg)](https://github.com/kyoh86/gogh/releases)

![](./doc/image/gogh.jpg)

## Description

**`gogh` is forked from [`ghq`](https://github.com/motemen/ghq).**

```console
$ gogh list
github.com/kyoh86/git-branches
github.com/kyoh86/gogh
github.com/kyoh86/vim-wipeout
github.com/kyoh86-tryouts/bare
github.com/nvim-telescope/telescope.nvim
...
```

`gogh` provides a way to organize remote repository clones, like `go clone` does.  When you clone a
remote repository by `gogh clone`, `gogh` makes a directory under a specific root directory (by default
`~/go/src`) using the remote repository URL's host and path.  And creating new one by `gogh create`,
`gogh` make both of a local repositories and a remote repository.

```console
$ gogh clone https://github.com/kyoh86/gogh
# Runs `git clone https://github.com/kyoh86/gogh ~/go/src/github.com/kyoh86/gogh`
```

You can also do:

- List repositories (local repositories) (`gogh list`).
- Create a new repository (`gogh create`).
- Fork a repository (`gogh fork`).
- Clone a repository (`gogh clone`).
- Delete a repository (`gogh delete`).
- List remote repositories (`gogh repos`).

See [#Available commands](#available-commands) for more information.

## Install

### For Golang developers

Ensure you have Go installed before running the following commands.

```console
$ go install github.com/kyoh86/gogh/v3/cmd/gogh
```

If you want zsh-completions, you can create completions file like this:

```console
$ echo "autoload -Uz compinit && compinit" >> ~/.zshrc
$ gogh completion zsh > $fpath[1]/_gogh
```

### `Homebrew`/`Linuxbrew`

```console
$ brew tap kyoh86/tap
$ brew update
$ brew install kyoh86/tap/gogh
```

## Setup

`gogh` manages repositories in multiple servers that is pairs of an owner and a host name.
To login in new server or logout, you should use `auth login`.

## Available commands

See [doc/usage/gogh.md](./doc/usage/gogh.md) for detailed command usage.

### Show repositories

| Command | Description                                                               |
| --      | --                                                                        |
| `cwd`   | Print the local repository which the current working directory belongs to |
| `list`  | List local repositories                                                   |
| `repos` | List remote repositories                                                  |

### Manipulate repositories

| Command  | Description                              |
| --       | --                                       |
| `clone`  | Clone remote repositories to local       |
| `create` | Create a new local and remote repository |
| `delete` | Delete local and remote repository       |
| `fork`   | Fork a repository                        |

### Configurations

| Command   | Description            |
| --        | --                     |
| `auth`    | Manage tokens          |
| `config`  | Show configurations    |
| `migrate` | Migrate configurations |
| `roots`   | Manage roots           |

### Others

| Command      | Description                                                |
| --           | --                                                         |
| `bundle`     | Manage bundle                                              |
| `completion` | Generate the autocompletion script for the specified shell |
| `help`       | Help about any command                                     |

Use `gogh [command] --help` for more information about a command.
Or see the manual in [doc/usage/gogh.md](./doc/usage/gogh.md).

## Environment variables

- `GOGH_CONFIG_PATH`
    - **(DEPRECATED)** The path to the configuration file.
    - Default: `${XDG_CONFIG_HOME}/gogh/config.yaml`.
- `GOGH_DEBUG`
    - Enable debug mode.
    - Default: ``.
    - Set to `1` to enable debug mode.
- `GOGH_DEFAULT_NAMES_PATH`
    - The path for the default names.
    - Default: `${XDG_CONFIG_HOME}/gogh/default_names.v4.toml`
- `GOGH_FLAG_PATH`
    - The path for values for each `gogh` flags.
    - Default: `${XDG_CONFIG_HOME}/gogh/flags.v4.toml`
- `GOGH_TOKENS_PATH`
    - The path for the tokens.
    - Default: `${XDG_CACHE_HOME}/gogh/tokens.v4.toml`
- `GOGH_WORKSPACE_PATH`
    - The path for the workspaces.
    - Default: `${XDG_CONFIG_HOME}/gogh/workspace.v4.toml`

## Configurations

### Roots

`gogh` manages repositories under the `roots` directories.

See also: [Directory structures](#Directory+structures)

You can change the roots with `roots add <path>` or `roots remove <path>` and see all of them by
`roots list`.  `gogh` uses the primary one to `create`, `fork` or `clone` to put a local repository
under it. If you want to change the primary, use `roots set-primary <path>`.

Default: `~/Projects`.

### Default Host and Owner

When you specify a repository with ambiguous user or host, it will be interpolated with a default
value. You may set them with `config set-default-host <host>` and `config set-default-owner <host> <owner>`.

If you set them like below:

| key     | value         |
| -       | -             |
| `host`  | `example.com` |
| `owner` | `kyoh86`      |

ambiguous repository names will be interpolated:

| Ambiguous name | Interpolated name       |
| --             | --                      |
| `gogh`         | example.com/kyoh86/gogh |
| `foobar/gogh`  | example.com/foobar/gogh |

NOTE: default host will be "github.com" if you don't set it.

### Flags

You can set flags for each command in the configuration file.  The flags are used to set the default
values for each command.  You can set the flags in the configuration file like this:

```toml
[repos]
    limit = 7
    archive = "not-archived"
[create]
    license-template = "mit"
```

## Directory structures

Local repositories are placed under `gogh.roots` with named `*host*/*user*/*repo*.

```
~/Projects             -- primary root
+-- github.com/
    |-- google/
    |   +-- go-github/
    |-- kyoh86/
    |   +-- gogh/
    +-- alecthomas/
        +-- kingpin/
/path/to/another/root
+-- github.com/
    |-- kyoh86/
    |   +-- xxx/
    +-- anybody/
        +-- yyy/
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see
LICENSE.  And this software is based on [`ghq`](https://github.com/motemen/ghq).
