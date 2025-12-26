# Gogh

Gogh is a tool to manage GitHub repositories efficiently, inspired by [`ghq`](https://github.com/motemen/ghq).

[![`PkgGoDev`](https://pkg.go.dev/github.com/kyoh86/gogh/v4)](https://pkg.go.dev/github.com/kyoh86/gogh/v4)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/gogh/v4)](https://goreportcard.com/report/github.com/kyoh86/gogh/v4)
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

`gogh` provides a way to organize remote repository clones, like `go clone` does. When you clone a
remote repository by `gogh clone`, `gogh` makes a directory under a specific root directory (by default
`~/Projects`) using the remote repository URL's host and path. And creating new one by `gogh create`,
`gogh` make both of a local repositories and a remote repository.

```console
$ gogh clone https://github.com/kyoh86/gogh
# Runs `git clone https://github.com/kyoh86/gogh ~/Projects/github.com/kyoh86/gogh`
```

You can also do:

- List repositories (local repositories) (`gogh list`).
- Create a new repository (`gogh create`).
- Fork a repository (`gogh fork`).
- Clone a repository (`gogh clone`).
- Delete a repository (`gogh delete`).
- List remote repositories (`gogh repos`).
- Show the current working directory's repository (`gogh cwd`).
- Manage [overlay files](#overlay-feature) (`gogh overlay`), [scripts](#script-feature) (`gogh script`), [hooks](#hook-feature) (`gogh hook`), and [extras](#extra-feature) (`gogh extra`).

See [#Available commands](#available-commands), [#Overlay Feature](#overlay-feature), [#Script Feature](#script-feature), [#Hook Feature](#hook-feature), and [#Extra Feature](#extra-feature) for more information.

## Install

### For Go developers

Ensure you have Go installed before running the following commands.

```console
$ go install github.com/kyoh86/gogh/v4/cmd/gogh@latest
```

### `Homebrew`/`Linuxbrew`

```console
$ brew tap kyoh86/tap
$ brew update
$ brew install kyoh86/tap/gogh
```

### Shell completions

You can generate the autocompletion script for your shell with the following command:

```console
$ gogh completion <shell>
```

If you want to use the generated script, you can save it to a file and source it in your shell configuration file.
For example, to generate the autocompletion script for `bash` and save it to `~/.gogh-completion.bash`, you can run:

```console
$ gogh completion bash > ~/.gogh-completion.bash
```

Then, add the following line to your `~/.bashrc` or `~/.bash_profile`:

```bash
source ~/.gogh-completion.bash
```

Or, if you want to use `zsh`, you can run:

```console
$ gogh completion zsh > ~/.config/zsh/completions/_gogh.zsh
```

Then, add the following line to your `~/.zshrc`:
(zsh completions are loaded from `fpath`)

```zsh
fpath=("~/.config/zsh/completions" $fpath)

autoload -Uz compinit && compinit
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

### Automation

| Command   | Description                                   |
| --        | --                                            |
| `extra`   | Manage overlay-hook packages (extras)         |
| `hook`    | Manage repository automation hooks            |
| `overlay` | Manage repository overlay files               |
| `script`  | Manage Lua scripts for repository actions     |

### Configurations

| Command   | Description                               |
| --        | --                                        |
| `auth`    | Manage authentication tokens              |
| `config`  | Show / Change configurations              |
| `roots`   | Manage root directories                   |

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
    - Enable debug mode
    - Default: `` (empty)
    - Set to any non-empty value to enable debug mode
- `GOGH_DEFAULT_NAMES_PATH`
    - The path for the default names
    - Default: `${XDG_CONFIG_HOME}/gogh/default_names.v4.toml`
- `GOGH_FLAG_PATH`
    - The path for values for each `gogh` flags
    - Default: `${XDG_CONFIG_HOME}/gogh/flags.v4.toml`
- `GOGH_HOOK_CONTENT_PATH`
    - The path to store hook content
    - Default: `${XDG_CONFIG_HOME}/gogh/hook.v4/`
- `GOGH_HOOK_PATH`
    - The path to store hook configuration
    - Default: `${XDG_CONFIG_HOME}/gogh/hook.v4.toml`
- `GOGH_OVERLAY_CONTENT_PATH`
    - The path to store overlay file contents
    - Default: `${XDG_CONFIG_HOME}/gogh/overlay.v4/`
- `GOGH_OVERLAY_PATH`
    - The path to store overlay configuration
    - Default: `${XDG_CONFIG_HOME}/gogh/overlay.v4.toml`
- `GOGH_SCRIPT_PATH`
    - The path to store script configuration
    - Default: `${XDG_CONFIG_HOME}/gogh/script.v4.toml`
- `GOGH_TOKENS_PATH`
    - The path for the authentication tokens
    - Default: `${XDG_CACHE_HOME}/gogh/tokens.v4.toml`
- `GOGH_WORKSPACE_PATH`
    - The path for the workspaces
    - Default: `${XDG_CONFIG_HOME}/gogh/workspace.v4.toml`

## Configurations

### Roots

`gogh` manages repositories under the `roots` directories.

See also: [Directory structures](#directory-structures)

You can change the roots with `roots add <path>` or `roots remove <path>` and see all of them by
`roots list`. `gogh` uses the primary one to `create`, `fork` or `clone` to put a local repository
under it. If you want to change the primary, use `roots set-primary <path>`.

Default: `${HOME}/Projects`.

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

You can set flags for each command in the configuration file. The flags are used to set the default
values for each command. You can set the flags in the configuration file like this:

```toml
[repos]
    limit = 7
    archive = "not-archived"
[bundle-restore]
    request-timeout = 5
```

The configuration file is located at `${XDG_CONFIG_HOME}/gogh/flags.v4.toml` by default, and you can
change the path with the `GOGH_FLAG_PATH` environment variable.

NOTE: If you set the boolean flags to `true` in the configuration file, you can disable them in the command line by
using `--<flag>=false`. For example, if you set `--private=true` in the configuration file, you can disable it by
using `--private=false` in the command line.

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
/path/to/another/root  -- another root
    +-- github.com/
        |-- kyoh86/
        |   +-- gogh/
        |   +-- git-branches/
        |   +-- vim-wipeout/
        |   +-- bare/
        |   +-- tryouts/
        |-- anybody/
            +-- yyy/
/...
```

## Overlay Feature

### What are Overlays?

Overlays are template files that can be applied to repositories. They are particularly useful for:

- Adding untracked files (like editor configurations or scripts)
- Applying consistent settings across multiple repositories
- Creating templates for new projects

### How Overlays Work

Overlays are never applied automatically. You must either:
1. Apply them manually using `gogh overlay apply`
2. Configure hooks to apply them automatically during repository operations

### Use Cases

1. **Editor Configuration**: Add your favorite editor settings to repositories
2. **Project Templates**: Apply language-specific configurations to new projects
3. **License Files**: Ensure all your repositories have the correct license file
4. **CI/CD Templates**: Add standard workflow files to repositories

### Basic Overlay Commands

Example commands to manage overlays:
See each `--help` for more details.

```console
# gogh overlay add <name> <source-path> <target-path>
$ gogh overlay add vscode-settings /path/to/source/vscode/settings.json .vscode/settings.json
$ gogh overlay list
# gogh overlay remove <overlay-id>
$ gogh overlay remove f8be36a27fa682b7b8d3c4117086851c74e47142705eba633cd91715c315d96b
# gogh overlay apply <overlay-id> [[host/]owner/]repo...
$ gogh overlay apply f8be36a27fa682b7b8d3c4117086851c74e47142705eba633cd91715c315d96b github.com/owner/repo
# Apply to current directory repository
$ gogh overlay apply f8be36a27fa682b7b8d3c4117086851c74e47142705eba633cd91715c315d96b .
```

### Practical Usages

#### Extracting Untracked Files as Overlays

Extract files from a repository that aren't tracked by git:

```console
$ gogh overlay extract [repo-refs...]
```

#### Showing Overlay Content

View the content of registered overlays:

```console
$ gogh overlay show <overlay-id>
```

## Script Feature

### What are Scripts?

Scripts in Gogh are Lua scripts that can be executed within repository contexts. They provide a powerful way to automate repository-specific tasks.

### How Scripts Work

Scripts are written in Lua and have access to repository information through the `gogh` global table. They can be invoked manually or automatically through hooks.

### Available Context Variables

When scripts are executed, they have access to:

```lua
-- Repository information
gogh.repo.host      -- e.g., "github.com"
gogh.repo.owner     -- e.g., "kyoh86"
gogh.repo.name      -- e.g., "gogh"
gogh.repo.path      -- Repository path relative to workspace
gogh.repo.full_path -- Full absolute path to the repository

-- Hook information (when invoked via hooks)
gogh.hook.id            -- Hook UUID
gogh.hook.name          -- Hook name
gogh.hook.repoPattern   -- Pattern that matched
gogh.hook.triggerEvent  -- Event that triggered the hook
gogh.hook.operationType -- Type of operation
gogh.hook.operationId   -- ID of the operation
```

### Basic Script Commands

```console
# Add a script
$ gogh script add setup-deps /path/to/setup-deps.lua

# List all scripts
$ gogh script list

# Invoke a script in repositories
$ gogh script invoke <script-id> [[host/]owner/]repo...
# Invoke in current directory repository
$ gogh script invoke <script-id> .

# Edit a script
$ gogh script edit <script-id>

# Remove a script
$ gogh script remove <script-id>
```

### Example Scripts

#### Setting Up Dependencies

```lua
-- setup-deps.lua
if os.execute("test -f ./package.json") == 0 then
  print("Installing Node.js dependencies...")
  os.execute("npm install")
elseif os.execute("test -f ./go.mod") == 0 then
  print("Downloading Go dependencies...")
  os.execute("go mod download")
end
```

#### Custom Git Configuration

```lua
-- project-git-config.lua
print("Setting repository-specific Git configuration...")
os.execute("git config user.email 'work@example.com'")
```

## Hook Feature

### What are Hooks?

Hooks in Gogh are automation triggers that execute operations (overlays or scripts) at specific points in the repository lifecycle.

### How Hooks Work

Hooks combine:
- **Trigger Events**: When to run (post-clone, post-fork, post-create)
- **Repository Patterns**: Which repositories to target
- **Operations**: What to execute (overlay application or script execution)

### Hook Configuration

- **Repository Pattern**: Controls which repositories the hook applies to
  - Works with glob patterns: `github.com/owner/*`, exact matches, etc.
- **Trigger Event**: When the hook should run
  - `post-clone`: After cloning a repository
  - `post-fork`: After forking a repository
  - `post-create`: After creating a new repository
- **Operation Type**: What action to perform
  - `overlay`: Apply overlay files
  - `script`: Execute a Lua script

### Basic Hook Commands

```console
# Add a hook to apply an overlay after cloning
$ gogh hook add --name "apply-vscode" \
  --repo-pattern "github.com/myorg/*" \
  --trigger-event "post-clone" \
  --operation-type "overlay" \
  --operation-id "<overlay-id>"

# Add a hook to run a script after creating
$ gogh hook add --name "setup-new-repo" \
  --repo-pattern "github.com/myorg/*" \
  --trigger-event "post-create" \
  --operation-type "script" \
  --operation-id "<script-id>"

# List all hooks
$ gogh hook list

# Manually invoke a hook
$ gogh hook invoke <hook-id> [[host/]owner/]repo
# Invoke for current directory repository
$ gogh hook invoke <hook-id> .

# Remove a hook
$ gogh hook remove <hook-id>
```

### Practical Examples

#### Automatic Project Setup

1. Create a setup script:
```console
$ gogh script add project-setup /path/to/setup.lua
```

2. Create a hook to run it after cloning:
```console
$ gogh hook add --name "auto-setup" \
  --repo-pattern "github.com/mycompany/*" \
  --trigger-event "post-clone" \
  --operation-type "script" \
  --operation-id "<script-id>"
```

#### Apply Templates to New Projects

1. Create overlays for project templates:
```console
$ gogh overlay add gitignore /path/to/template/.gitignore .gitignore
```

2. Create a hook to apply them:
```console
$ gogh hook add --name "apply-templates" \
  --repo-pattern "github.com/myorg/*" \
  --trigger-event "post-create" \
  --operation-type "overlay" \
  --operation-id "<overlay-id>"
```

## Extra Feature

### What are Extras?

Extras are higher-level configurations that combine overlays and hooks into reusable packages. They simplify the process of managing repository templates and automation.

### Types of Extras

1. **Auto-apply Extras**: Automatically applied to specific repositories when cloned
   - Created from existing repositories with their ignored files
   - Applied automatically via hooks when the repository is cloned

2. **Named Extras**: Reusable templates that can be applied to any repository
   - Created with custom names for easy reference
   - Applied manually using the `extra apply` command

### How Extras Work

Each extra contains:
- One or more overlay-hook pairs
- Metadata about the source and creation time
- Type information (auto-apply or named)

When you save an extra from a repository, it:
1. Extracts all untracked files as overlays
2. Creates hooks to apply these overlays
3. Bundles them together as a single extra

### Basic Extra Commands

```console
# Save a repository's untracked files as an auto-apply extra
$ gogh extra save github.com/owner/repo

# Create a named extra from a specific repository
$ gogh extra create my-template --source github.com/owner/repo --overlay <overlay-name>

# List all extras
$ gogh extra list

# Show details of an extra
$ gogh extra show <extra-id-or-name>

# Apply a named extra to repositories
$ gogh extra apply <extra-name> [[host/]owner/]repo...

# Remove an extra
$ gogh extra remove <extra-id-or-name>
```

### Practical Examples

#### Auto-apply Repository Configuration

Save your local development environment setup to be automatically restored when cloning:

```console
# Save untracked files from a repository as auto-apply extra
$ gogh extra save github.com/owner/repo
# Save from current directory repository
$ gogh extra save .
# Now when you clone this repository again, all untracked files will be restored
```

#### Create Reusable Project Templates

Create a template from a well-configured repository:

```console
# Create a template with common development files
$ gogh extra create python-template github.com/myorg/python-starter

# Apply it to new projects
$ gogh extra apply python-template github.com/myorg/new-python-project
```

#### Managing Extras

View and manage your extras:

```console
# List all extras with their types
$ gogh extra list --type all

# Show detailed information about an extra
$ gogh extra show my-template

# Remove an extra that's no longer needed
$ gogh extra remove old-template
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This software is released under the [MIT License](http://www.opensource.org/licenses/MIT), see
LICENSE. And this software is based on [`ghq`](https://github.com/motemen/ghq).
