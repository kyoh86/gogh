## gogh hook add

Add a new hook

```
gogh hook add [flags]
```

### Options

```
  -h, --help                    help for add
      --name string             Name of the hook
      --operation-id string     Operation resource ID
      --operation-type string   Operation type; it can accept "overlay" or "script"
      --repo-pattern string     Repository pattern
      --trigger-event string    event that triggers the hook; it can accept "", "post-clone", "post-fork" or "post-create"
```

### SEE ALSO

* [gogh hook](gogh_hook.md)	 - Manage repository hooks

