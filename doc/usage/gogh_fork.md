## gogh fork

Fork a repository

```
gogh fork [flags] [<host>/]<owner>/<name>
```

### Examples

```
  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host set by "config set-default-host"
```

### Options

```
      --clone-retry-limit int          The number of retries to clone a repository (default 3)
  -t, --clone-retry-timeout duration   Timeout for each clone attempt (default 5m0s)
      --default-branch-only[=false]    Only fork the default branch
  -h, --help                           help for fork
      --to string                      Fork to the specified repository. It accepts a notation like '<owner>/<name>' or '<owner>/<name>=<alias>'. If not specified, it will be forked to the default owner and same name as the original repository. If the alias is specified, it will be set as the local repository name
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

