# gogh

GO GitHub project manager

[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gogh)](https://goreportcard.com/report/github.com/kyoh86/gogh)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/gogh.svg)](https://codecov.io/gh/kyoh86/gogh)

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
gogh list
gogh dump
gogh find
gogh where
gogh repos

gogh get
gogh bulk-get
gogh pipe-get
gogh create
gogh fork
gogh remove
```

See `gogh --help-long` for details.

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

## CONFIGURATIONS

It's possible to change targets by a preference **YAML file**.
If you don't set `--config` flag or `GOGH_CONFIG` environment variable,
`gogh` loads configurations from `${XDG_CONFIG_HOME:-$HOME/.config}/gogh/config.yaml`

To set new configure, you should use `gogh config set <name> <value>`.

Each of propoerties are able to be overwritten by environment variables.

### `roots`

The paths to directory under which cloned repositories are placed.
See [DIRECTORY STRUCTURES](#DIRECTORY+STRUCTURES) below. Default: `~/go/src`.

This property can have multiple values.
If so, the first one becomes primary one i.e. new repository clones are always created under it.
You may want to specify `$GOPATH/src` as a secondary root.

### `github.user`

A name of your GitHub user (i.e. `kyoh86`).

If an environment variable `GOGH_GITHUB_USER` is set, its value is used instead.

If an environment variable `GOGH_ROOTS` is set, its value is used instead.

### `github.token`

The token to connect GitHub API.

`gogh` saves this one in `keyring`.
It is saved for a service `<github.host>.gogh.kyoh86.dev` and a user `github.user`.

If an environment variable `GOGH_GITHUB_TOKEN` is set, its value is used instead.

### `github.host`

The host name to connect to GitHub. Default: `github.com`.

If an environment variable `GOGH_GITHUB_HOST` is set, its value is used instead.

### `hooks`

The directory name to store [hooks](#HOOKS-1). Default: 

If an environment variable `GOGH_HOOKS` is et, its value is used instead.

If you don't set this one, `gogh` searches hooks from `${XDG_CONFIG_HOME:-$HOME/.config}/gogh/hooks`

## COMMANDS

See `gogh --long-help` for details.

### `list`

List locally cloned repositories.
If a query argument is given, only repositories whose names contain that query text are listed.
`-f` (`--format`) specifies format of the item. Default: `relative`.

* `-f=full` is given, the full paths to the repository will be printed.
* `-f=short` is given, gogh prints each project as short as possible.
* `-f=relative` is given, the relative paths from gogh root to the repository will be printed.
* `-f=url` is given, the urls of the repository will be printed.

### `dump`

Dump local repository list.
This is shorthand for `gogh list --format=url`.
It can be used to backup and restore projects.

e.g.

```
$ gogh dump > projects.txt

# copy projects.txt to another machine

$ cat projects.txt | gogh bulk-get
```

### `find`

Look into a locally cloned repository with the shell.

### `where`

Show where a local repository is.

### `repos`

Show a list of repositories for a user.

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

### `bulk-get`

Reads repository URLs from stdin line by line and performs 'get' for each of them.

### `pipe-get`

Reads repository URLs from other command output by line and performs 'get' for each of them.

### `create`

Create a new repository on remote GitHub and clone it into local project.

### `fork`

Clone a remote repository under gogh root direcotry if the project is not exist in local.
And fork int on remote GitHub repository with calling `hub fork`

### `remove`

Remove a repository on remote GitHub and local project.

## ENVIRONMENT VARIABLES

Some environment variables are used for flags.

### GOGH_CONFIG

You can set it instead of `--config` flag (configuration file path).
Default:  `${XDG_CONFIG_HOME:-$HOME/.config}/gogh/config.yaml`.

### GOGH_FLAG_ROOT_ALL

You can set it truely value and `gogh root` shows all of the roots like `gogh root --all`.
If we want show only primary root, call `gogh root --no-all`.

e.g.

```
$ echo $GOGH_ROOTS
/Users/kyoh86/Projects:/Users/kyoh86/go/src
$ gogh roots
/Users/kyoh86/Projects
$ gogh roots --all
/Users/kyoh86/Projects
/Users/kyoh86/go/src
$ GOGH_FLAG_ROOT_ALL=1 gogh root
/Users/kyoh86/Projects
/Users/kyoh86/go/src
$ GOGH_FLAG_ROOT_ALL=1 gogh root --no-all
/Users/kyoh86/Projects
```

## DIRECTORY STRUCTURES

Local repositories are placed under `gogh.roots` with named github.com/*user*/*repo*.

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

If you set github.host, they will be placed with *github.host*/*user*/*repo*.

## HOOKS

Like Git Hooks, gogh has a way to fire off custom scripts when some commands be called.

### Installing a hook

The hooks are all stored in the directory that configured in `hooks` and `.gogh/hooks` in the each repository.

Any properly named executable scripts will work fine -
you can write them in Shellscript, Ruby, Python or whatever language you are familiar with.

### Hook names

| Name             | Trigger                                                       |
|------------------|---------------------------------------------------------------|
| post-get-each    | `get`, `bulk-get` or `pipe-get` is processed each repository  |
| post-create      | `create` is processed                                         |
| post-fork        | `fork` is processed                                           |
| pre-remove-each  | `remove` is processing each repository                        |

### Context

All of the hooks will be called in the root of the target project.

### Example

```console
$ mkdir -p ~/.config/gogh/hooks
$ echo '#!/bin/sh' > ~/.config/gogh/hooks/post-create
$ echo 'yo go-project' >> ~/.config/gogh/hooks/post-create
$ chmod +x ~/.config/gogh/hooks/post-create
$ gogh create foobar
Creating new project and a remote repository foobar
Checking existing project
Creating a directory
Initializing a repository
Creating a new repository in GitHub
Executing post-create hook

     _-----_     ╭──────────────────────────╮
    |       |    │      Welcome to the      │
    |--(o)--|    │        funkadelic        │
   `---------´   │ generator-go-project/app │
    ( _´U`_ )    │        generator!        │
    /___A___\   /╰──────────────────────────╯
     |  ~  |     
   __'.___.'__   
 ´   `  |° ´ Y ` 

? Project name 
```

## SHELL EXTENTIONS

To be enabled shell extentions (for zsh / bash), set up gogh in your shell-rc file (`.bashrc` / `.zshrc`).

* `gogo cd` command
  * shorthand for `cd $(gogh find <project name>)`
* `gogo get --cd` option
* auto-completions

```
eval "$(gogh setup)"
```

If you have not set `SHELL` envar right, tell your shell explicitly.

```
eval "$(gogh setup --shell bash)"
```

NOTE: Now gogh supports `bash` or `zsh` only.

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
* `gogh` holds tokens for GitHub in the `keyring` for security.
    * `hub`, `gh` and `ghq` hold them in raw config file.
* I wanted to merge functions of [ghq](https://github.com/motemen/ghq) and [hub](https://github.com/github/hub).
    * `gogh new` creates a new one with make both of a local project and a remote repository.
        * It calls `git init` and `hub create` in the gogh.roots directory.
    * `gogh fork` clones a remote repository into the gogh.roots directory and fork it GitHub (with calling `hub fork`).
    * But there may be some collision in configurations of **ghq** and **hub**. It offers a challenge for me to resolve them by gogh.
* (nits) I don't like `github.com/onsi/gomega` and `github.com/urfave/cli`. But I love `github.com/stretchr/testify` and `github.com/alecthomas/kingpin`.
* (nits) I want gogh to be able to be used as a library (`github.com/kyoh86/gogh/gogh` package).

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see LICENSE.
And this software is based on [ghq](https://github.com/motemen/ghq).
