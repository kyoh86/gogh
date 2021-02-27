# gogh

GO GitHub project manager

[![PkgGoDev](https://pkg.go.dev/badge/kyoh86/gogh)](https://pkg.go.dev/kyoh86/gogh)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gogh)](https://goreportcard.com/report/github.com/kyoh86/gogh)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/gogh.svg)](https://codecov.io/gh/kyoh86/gogh)
[![Release](https://github.com/kyoh86/gogh/workflows/Release/badge.svg)](https://github.com/kyoh86/gogh/releases)

![](./image/gogh.jpg)

## Description

**`gogh` is forked from [ghq](https://github.com/motemen/ghq).**

`gogh` provides a way to organize remote repository clones, like `go get` does.  When you clone a
remote repository by `gogh get`, gogh makes a directory under a specific root directory (by default
`~/go/src`) using the remote repository URL's host and path.  And creating new one by `gogh new`,
gogh make both of a local project and a remote repository.

```
$ gogh get https://github.com/kyoh86/gogh
# Runs `git clone https://github.com/kyoh86/gogh ~/go/src/github.com/kyoh86/gogh`
```

You can also list projects (local repositories) (`gogh list`).

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

| Command        | Description              |
| --             | --                       |
| `gogh roots`   | Manage roots             |
| `gogh servers` | Manage servers           |
| `gogh help`    | Help about any command   |

Use `gogh [command] --help` for more information about a command.

## Install

For Golang developers:

```
go get github.com/kyoh86/gogh/cmd/gogh
```

For [Homebrew](https://brew.sh/) users:

```
brew tap kyoh86/tap
brew update
brew install gogh
```

## Roots

`gogh` manages projects under the `roots` directories.

Seealso: [Directory structures](#Directory+structures)

You can change the roots with `roots add <path>` or `roots remove <path>` and see all of them by
`roots list`.  `gogh` uses the first one as the default one, `create`, `fork` or `clone` will put a
local project under it. If you want to change the default, use `roots set-default <path>`.

Default: `~/Projects`.

## Servers

`gogh` manages respositories in some servers that pairs of a user and a host name.  To login in new
server or logout, you should use `servers login`.  `gogh` uses the first server as the default one.
When you specify a repository with amonguous user or host, it will be interpolated with a default
server.

I.E. when servers are:

- github.com:
  - user: kyoh86
- example.com:
  - user: foobar

Amonguas repository names will be interpolated:

| Amonguals name | Interpolated name      |
| --             | --                     |
| gogh           | github.com/kyoh86/gogh |
| foobar/gogh    | github.com/foobar/gogh |

## Commands

See `gogh <command> --help` for details.

### `list [<query>]`

List locally projects.  If a query argument is given, only repositories whose names contain that
query text are listed.  `-f` (`--format`) specifies format of the item. Default: `relative`.

* `-f=full` is given, the full paths to the repository will be printed.
* `-f=relative` is given, the relative paths from gogh root to the repository will be printed.
* `-f=url` is given, the urls of the repository will be printed.

### `repos`

List remote repositories for a user.

### `get`

Clone a remote repository under gogh root directory (see [DIRECTORY STRUCTURES](#DIRECTORY+STRUCTURES) below).
If the repository is already cloned to local project, nothing will happen unless `-u` (`--update`) flag is supplied,
in which case the project (local repository) is updated (`git pull --ff-only` eg.).
When you use `--ssh` option, the repository is cloned via SSH protocol.

If there are multiple `gogh.roots`, existing local clones are searched first.
Then a new repository clone is created under the primary root if none is found.

With `--shallow` option, a "shallow clone" will be performed (for Git repositories only, `git clone --depth 1 ...` eg.).
Be careful that a shallow-cloned repository cannot be pushed to remote.

Currently Git and Mercurial repositories are supported.

### `create`

Create a new repository on remote GitHub and clone it into local project.

### `fork`

Clone a remote repository under gogh root direcotry if the project is not exist in local.
And fork int on remote GitHub repository with calling `hub fork`

### `remove`

Remove a repository on remote GitHub and local project.

## Directory structures

Local projects are placed under `gogh.roots` with named `*host*/*user*/*repo*.

```
~/go/src
+-- github.com/
|-- google/
|   +-- go-github/
|-- kyoh86/
|   +-- gogh/
+-- alecthomas/
  +-- kingpin/
```

## SHELL EXTENTIONS

To be enabled shell extentions (for zsh / bash), set up gogh in your shell-rc file (`.bashrc` / `.zshrc`).

* `gogo cd` command
  * shorthand for `cd $(gogh find <project name>)`
* `gogo get --cd` option
* auto-completions

```
eval "$(gogh init)"
```

If you have not set `SHELL` envar right, tell your shell explicitly.

```
eval "$(gogh init --shell bash)"
```

NOTE: Now gogh supports `bash` or `zsh` only.

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see LICENSE.
And this software is based on [ghq](https://github.com/motemen/ghq).
