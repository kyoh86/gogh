## gogh bundle restore

Get dumped local repositoiries

```
gogh bundle restore [flags]
```

### Options

```
      --clone-retry-limit int          The number of retries to clone a repository (default 3)
      --clone-retry-timeout duration   Timeout for each clone attempt (default 5m0s)
      --dry-run[=false]                Displays the operations that would be performed using the specified command without actually running them
  -f, --file string                    Read the file as input; if it's empty("") or hyphen("-"), read from stdin (default "/home/kyoh86/.config/gogh/bundle.txt")
  -h, --help                           help for restore
```

### SEE ALSO

* [gogh bundle](gogh_bundle.md)	 - Manage bundle

