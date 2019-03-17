# gogh

[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gogh)](https://goreportcard.com/report/github.com/kyoh86/gogh)
[![CircleCI](https://img.shields.io/circleci/project/github/kyoh86/gogh.svg)](https://circleci.com/gh/kyoh86/gogh)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/gogh.svg)](https://codecov.io/gh/kyoh86/gogh)

GO GitHub project manager

![](./image/gogh.jpg)

## DESCRIPTION

**`gogh` is forked from [ghq](https://github.com/motemen/ghq).**

`gogh` provides a way to organize remote repository clones, like `go get` does.
When you clone a remote repository by `gogh get`, gogh makes a directory under a specific root directory (by default `~/go/src`) using the remote repository URL's host and path.
And creating new one by `gogh new`, gogh make both of a local project and a remote repository.

```
$ gogh get https://github.com/kyoh86/gogh
# Runs `git clone https://github.com/kyoh86/gogh ~/go/src/github.com/kyoh86/gogh`
```

You can also list projects (local repositories) (`gogh list`), find a project (`gogh find`).

## SYNOPSIS

```
gogh get [--update,u] [--ssh] [--shallow] [(<repository URL> | <user>/<project> | <project>)...]
gogh bulk [--update,u] [--ssh] [--shallow]
gogh pipe [--update,u] [--ssh] [--shallow] <command> <command-args>...
gogh fork [--update,u] [--ssh] [--shallow] [--no-remote] [--remote-name=<REMOTE>] [--org=<ORGANIZATION] (<repository URL> | <user>/<project> | <project>)
gogh new [--update,u] [--ssh] [--shallow] [--no-remote] [--remote-name=<REMOTE>] [--org=<ORGANIZATION] (<repository URL> | <user>/<project> | <project>)
gogh list [--format,f=short|full|relative|url] [--primary,p] [<query>]
gogh dump [--primary,p] [<query>]
gogh find (<project>)
gogh where [--primary,p] [<query>]
gogh repo [--user=<USER>] [--own] [--collaborate] [--member] [--visibility=<VISIBILITY>] [--sort=<SORT>] [--direction=<DIRECTION>]
gogh root [--all]
```

## INSTALLATION

For Golang developers:

```
go get github.com/kyoh86/gogh
```

For [Homebrew](https://brew.sh/) users:

```
brew tap kyoh86/tap
brew update
brew install gogh
```

## COMMANDS

```
gogh --long-help
```

### `get`

Clone a remote repository under gogh root directory (see [DIRECTORY STRUCTURES](#DIRECTORY+STRUCTURES) below).
If the repository is already cloned to local project, nothing will happen unless `-u` (`--update`) flag is supplied,
in which case the project (local repository) is updated (`git pull --ff-only` eg.).
When you use `-p` option, the repository is cloned via SSH protocol.

If there are multiple `gogh.root`s, existing local clones are searched first.
Then a new repository clone is created under the primary root if none is found.

With `--shallow` option, a "shallow clone" will be performed (for Git repositories only, `git clone --depth 1 ...` eg.).
Be careful that a shallow-cloned repository cannot be pushed to remote.

Currently Git and Mercurial repositories are supported.

### `bulk`

Reads repository URLs from stdin line by line and performs 'get' for each of them.

### `pipe`

Reads repository URLs from other command output by line and performs 'get' for each of them.

### `fork`

Clone a remote repository under gogh root direcotry if the project is not exist in local.
And fork int on remote GitHub repository with calling `hub fork`

### `list`

List locally cloned repositories.
If a query argument is given, only repositories whose names contain that query text are listed.
`-f` (`--format`) specifies format of the item. Default: `relative`.

* `-f=full` is given, the full paths to the repository will be printed.
* `-f=short` is given, gogh prints each project as short as possible.
* `-f=relative` is given, the relative paths from gogh root to the repository will be printed.
* `-f=url` is given, the urls of the repository will be printed.

### `list`

Dump local repository list.
This is shorthand for `gogh list --format=url`.
It can be used to backup and restore projects.

e.g.

```
$ gogh dump > projects.txt

# copy projects.txt to another machine

$ cat projects.txt | gogh bulk
```

### `find`

Look into a locally cloned repository with the shell.

### `where`

Show where a local repository is.

### `repo`

Show a list of repositories for a user.

### `root`

Prints repositories' root (i.e. `gogh.root`). Without `--all` option, the primary one is shown.

## ENVIRONMENT VARIABLES

### GOGH_ROOT

The paths to directory under which cloned repositories are placed.
See [DIRECTORY STRUCTURES](#DIRECTORY+STRUCTURES) below. Defaults to `~/go/src`.

This variable can have multiple values.
If so, the first one becomes primary one i.e. new repository clones are always created under it.
You may want to specify `$GOPATH/src` as a secondary root (environment variables should be expanded.)

### GOGH_GITHUB_USER

A name of your GitHub user (i.e. `kyoh86`).
If it is not set, gogh uses `GITHUB_USER` envar or OS user name from envar (`USERNAME` in windows, `USER` in others) instead.

### GOGH_LOG_LEVEL

The level to output logs (debug, info, warn, error or panic). Default: warn

### GOGH_GHE_HOST

Hostnames of your GitHub Enterprise installation.
This variable can have multiple values that separated with spaces.

### GOGH_GITHUB_TOKEN

The token to connect GitHub API.
If it is not set, gogh uses `GITHUB_TOKEN` envar instead.

### GOGH_GITHUB_HOST

The host to connect GitHub on default.
If it is not set, gogh uses `GITHUB_HOST` envar or `github.com` instead.

### GOGH_FLAG_ROOT_ALL

We can set it truely value and `gogh root` shows all of the roots like `gogh root --all`.
If we want show only primary root, call `gogh root --no-all`.

e.g.

```
$ echo $GOGH_ROOT
/Users/kyoh86/Projects:/Users/kyoh86/go/src
$ gogh root
/Users/kyoh86/Projects
$ gogh root --all
/Users/kyoh86/Projects
/Users/kyoh86/go/src
$ GOGH_FLAG_ROOT_ALL=1 gogh root
/Users/kyoh86/Projects
/Users/kyoh86/go/src
$ GOGH_FLAG_ROOT_ALL=1 gogh root --no-all
/Users/kyoh86/Projects
```

## DIRECTORY STRUCTURES

Local repositories are placed under `gogh.root` with named github.com/*user*/*repo*.

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

* `gogogh` command
  * shorthand for `cd $(gogh find <project name>)`
* auto-completions

```
eval "$(gogh setup)"
```

If you want to rename `gogogh` command, specify `--cd-function-name=<NAME>` like below.

```
eval "$(gogh setup --cd-function-name=foobar)"
```

# DEFERENCES TO `ghq`

* `ghq` is too complex for me. That's why I forked this project from it.
    * `ghq look` runs new shell only to change working directory to a project.
        * So I cannot back to previous by `cd -`. I need to `exit` to do it.
        * But `ghq look` says `cd <project to path>` when I run it.
    * `ghq list --unique` returns a bizarre and complex list when I set multiple root.
    * `ghq import xxx` needs to setup `ghq.import.xxx` option to specify what command to run.
    * `git config ghq.xxx`...
        * `gogh` can be configured with envars instead of git-config(1).
* `gogh` doesn't support VCSs other than GitHub
    * If I want to manage projects in VCSs other than GitHub, I should use other tool, I think so.
* I wanted to merge functions of [ghq](https://github.com/motemen/ghq) and [hub](https://github.com/github/hub).
    * `gogh new` creates a new one with make both of a local project and a remote repository.
        * It calls `git init` and `hub create` in the gogh.root directory.
    * `gogh fork` clones a remote repository into the gogh.root directory and fork it GitHub (with calling `hub fork`).
    * But there may be some collision in configurations of **ghq** and **hub**. It offers a challenge for me to resolve them by gogh.
* (nits) I don't like `github.com/onsi/gomega` and `github.com/urfave/cli`. But I love `github.com/stretchr/testify` and `github.com/alecthomas/kingpin`.
* (nits) I want gogh to be able to be used as a library (`github.com/kyoh86/gogh/gogh` package).

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see LICENSE.
And this software is based on [ghq](https://github.com/motemen/ghq).
