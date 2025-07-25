## gogh clone

Clone remote repositories to local

```
gogh clone [flags] [[[<host>/]<owner>/]<name>[=<alias>]...]
```

### Examples

```
  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is used for a local repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"
```

### Options

```
  -t, --clone-retry-timeout duration   Timeout for each clone attempt (default 5m0s)
      --dry-run                        Displays the operations that would be performed using the specified command without actually running them
  -h, --help                           help for clone
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

