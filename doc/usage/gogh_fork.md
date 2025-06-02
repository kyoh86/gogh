## gogh fork

Fork a repository

```
gogh fork [flags] <owner>/<name>
```

### Options

```
      --clone-retry-limit int           (default 3)
  -t, --clone-retry-timeout duration   Timeout for each clone attempt (default 5m0s)
      --default-branch-only[=false]    Only fork the default branch
  -h, --help                           help for fork
      --to string                      Fork to the specified repository. It accepts a notation like 'OWNER/NAME' or 'OWNER/NAME=ALIAS'. If not specified, it will be forked to the default owner and same name as the original repository. If the alias is specified, it will be set as the local repository name.
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

