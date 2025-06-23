## gogh hook invoke

Manually invoke a hook for a repository

```
gogh hook invoke [flags] <hook-id> [[<host>/]<owner>/]<name>
```

### Examples

```
  invoke <hook-id> github.com/owner/repo
  invoke <hook-id> owner/repo
  invoke <hook-id> repo
  invoke <hook-id> .  # Use current directory repository
  
  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".
```

### Options

```
  -h, --help   help for invoke
```

### SEE ALSO

* [gogh hook](gogh_hook.md)	 - Manage repository hooks

