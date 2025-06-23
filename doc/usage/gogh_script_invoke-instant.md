## gogh script invoke-instant

Run a temporary script in a repository without storing it

```
gogh script invoke-instant [flags] [[[<host>/]<owner>/]<name>...]
```

### Examples

```
  invoke-instant --file script.lua repo1 repo2
  invoke-instant --file - repo1 < script.lua
  echo 'print(gogh.repo.name)' | gogh script invoke-instant --file - repo1
  invoke-instant --file script.lua .  # Use current directory repository
  invoke-instant --file script.lua --all
  invoke-instant --file script.lua --pattern <pattern>

  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".
```

### Options

```
      --all               Apply to all repositories in the workspace
  -f, --file string       Path to script file to invoke (use '-' for stdin)
  -h, --help              help for invoke-instant
  -p, --pattern strings   Patterns for selecting repositories
```

### SEE ALSO

* [gogh script](gogh_script.md)	 - Manage repository script files

