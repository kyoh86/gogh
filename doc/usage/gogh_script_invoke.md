## gogh script invoke

Invoke an script in a repository

```
gogh script invoke [flags] <script-id> [[[<host>/]<owner>/]<name>...]
```

### Examples

```
  invoke [flags] <script-id> [[[<host>/]<owner>/]<name>...]
  invoke [flags] <script-id> --all
  invoke [flags] <script-id> --pattern <pattern> [--pattern <pattern>]...

  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is a local name for the remote repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"
```

### Options

```
      --all               Apply to all repositories in the workspace
  -h, --help              help for invoke
  -p, --pattern strings   Patterns for selecting repositories
```

### SEE ALSO

* [gogh script](gogh_script.md)	 - Manage repository script files

