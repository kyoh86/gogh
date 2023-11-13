# Gogh

GO GitHub project manager

[![`PkgGoDev`](https://pkg.go.dev/badge/kyoh86/gogh)](https://pkg.go.dev/kyoh86/gogh)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gogh)](https://goreportcard.com/report/github.com/kyoh86/gogh)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/gogh.svg)](https://codecov.io/gh/kyoh86/gogh)
[![Release](https://github.com/kyoh86/gogh/workflows/Release/badge.svg)](https://github.com/kyoh86/gogh/releases)

![](./image/gogh.jpg)

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

`gogh` provides a way to organize remote repository clones, like `go get` does.  When you clone a
remote repository by `gogh get`, `gogh` makes a directory under a specific root directory (by default
`~/go/src`) using the remote repository URL's host and path.  And creating new one by `gogh new`,
`gogh` make both of a local project and a remote repository.

```console
$ gogh get https://github.com/kyoh86/gogh
# Runs `git clone https://github.com/kyoh86/gogh ~/go/src/github.com/kyoh86/gogh`
```

You can also list projects (local repositories) (`gogh list`).

## Install

### For Golang developers

```console
$ go install github.com/kyoh86/gogh/cmd/gogh@latest
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

### `Makepkg`

```console
$ mkdir -p gogh_build && \
  cd gogh_build && \
  curl -L --silent https://github.com/kyoh86/gogh/releases/latest/download/gogh_PKGBUILD.tar.gz | tar -xvz
$ makepkg -i
```

## Available commands

### Show projects

| Command        | Description              |
| --             | --                       |
| `gogh list`    | List local projects      |
| `gogh repos`   | List remote repositories |

### Manipulate projects

| Command        | Description                                   |
| --             | --                                            |
| `gogh create`  | Create a new project with a remote repository |
| `gogh delete`  | Delete a repository with a remote repository  |
| `gogh fork`    | Fork a repository                             |
| `gogh clone`   | Clone a repository to local                   |

### Others

| Command       | Description            |
| --            | --                     |
| `gogh roots`  | Manage roots           |
| `gogh auth`   | Manage Authentications |
| `gogh bundle` | Manage bundle          |
| `gogh help`   | Help about any command |

Use `gogh [command] --help` for more information about a command.
Or see the manual in [usage/gogh.md](./usage/gogh.md).

## Roots

`gogh` manages projects under the `roots` directories.

See also: [Directory structures](#Directory+structures)

You can change the roots with `roots add <path>` or `roots remove <path>` and see all of them by
`roots list`.  `gogh` uses the first one as the default one, `create`, `fork` or `clone` will put a
local project under it. If you want to change the default, use `roots set-default <path>`.

Default: `~/Projects`.

## Servers

`gogh` manages repositories in some servers that pairs of an owner and a host name.  To login in new
server or logout, you should use `auth login`.

## Default Host and Owner

When you specify a repository with ambiguous user or host, it will be interpolated with a default
value. You may set them with `set-default`.

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

## Commands

Manual: [usage/gogh.md](./usage/gogh.md).

## Directory structures

Local projects are placed under `gogh.roots` with named `*host*/*user*/*repo*.

```
~/Projects
+-- github.com/
|-- google/
|   +-- go-github/
|-- kyoh86/
|   +-- gogh/
+-- alecthomas/
  +-- kingpin/
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see
LICENSE.  And this software is based on [`ghq`](https://github.com/motemen/ghq).
