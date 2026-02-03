## gogh hook update

Update an existing hook

```
gogh hook update [flags] <hook-id>
```

### Options

```
  -h, --help                    help for update
      --name string             Name of the hook
      --operation-id string     Operation resource ID (overlay ID or script ID). It can be a partial ID as it is matched by prefix.
      --operation-type string   Operation type; it can accept "overlay" or "script"
      --repo-pattern string     Repository pattern
      --trigger-event string    event to hook automatically; it can accept "post-clone", "post-fork" or "post-create"
```

### SEE ALSO

* [gogh hook](gogh_hook.md)	 - Manage repository hooks

