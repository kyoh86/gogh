# gogh

GO GitHub project manager

![](./image/gogh.jpg)

## DESCRIPTION

**`gogh` is forked from [ghq](https://github.com/motemen/ghq).**

`gogh` provides a way to organize remote repository clones, like `go get` does.
When you clone a remote repository by `gogh get`, gogh makes a directory under a specific root directory (by default `~/go/src`) using the remote repository URL's host and path.
And creating new repository by `gogh new`, gogh make both of local and remote ones.

```
$ gogh get https://github.com/kyoh86/gogh
# Runs `git clone https://github.com/kyoh86/gogh ~/go/src/github.com/kyoh86/gogh`
```

You can also list local repositories (`gogh list`), find a local repositories (`gogh find`).

## SYNOPSIS

```
gogh get [--update,u] [--ssh] [--shallow] [(<repository URL> | <user>/<project> | <project>)...]
gogh bulk [--update,u] [--ssh] [--shallow]
gogh pipe [--update,u] [--ssh] [--shallow] <command> <command-args>...
gogh fork [--update,u] [--ssh] [--shallow] [--no-remote] [--remote-name=<REMOTE>] [--org=<ORGANIZATION] (<repository URL> | <user>/<project> | <project>)
gogh new [--update,u] [--ssh] [--shallow] [--no-remote] [--remote-name=<REMOTE>] [--org=<ORGANIZATION] (<repository URL> | <user>/<project> | <project>)
gogh list [--exect,e] [--full-path,f] [--short,s] [--primary,p] [<query>]
gogh find (<project>)
gogh root [--all]
```

## COMMANDS

### `get`

Clone a remote repository under gogh root directory (see [DIRECTORY STRUCTURES](#DIRECTORY+STRUCTURES) below).
If the repository is already cloned to local, nothing will happen unless `-u` (`--update`) flag is supplied,
in which case the local repository is updated (`git pull --ff-only` eg.).
When you use `-p` option, the repository is cloned via SSH protocol.

If there are multiple `gogh.root` s, existing local clones are searched first.
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
`-e` (`--exact`) forces the match to be an exact one (i.e. the query equals to _project_ or _user_/_project_)
If `-p` (`--full-path`) is given, the full paths to the repository root are printed instead of relative ones.
IF `-s` (`--short`) is given, the names of the repository which as short as possible are printed instead of relative paths.

### `find`

Look into a locally cloned repository with the shell.

### `root`

Prints repositories' root (i.e. `gogh.root`). Without `--all` option, the primary one is shown.

## CONFIGURATION

Configuration uses `git-config` variables.

### `gogh.root`

The path to directory under which cloned repositories are placed.
See [DIRECTORY STRUCTURES](#DIRECTORY+STRUCTURES) below. Defaults to `~/go/src`.

This variable can have multiple values.
If so, the first one becomes primary one i.e. new repository clones are always created under it.
You may want to specify `$GOPATH/src` as a secondary root (environment variables should be expanded.)

## ENVIRONMENT VARIABLES

### `PM_ROOT`

If set to a path, this value is used as the only root directory regardless of other existing `gogh.root` settings.

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

## INSTALLATION

```
go get github.com/kyoh86/gogh
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see LICENSE.
And this software is based on [ghq](https://github.com/motemen/ghq).
