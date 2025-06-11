## gogh hook add

Add an existing Lua script as hook

```
gogh hook add [flags] <lua-script-path>
```

### Options

```
      --event string          event to hook automatically; it can accept "", "clone", "fork", "create" or "never" (default "never")
  -h, --help                  help for add
      --name string           Name of the hook
      --repo-pattern string   Repository pattern
      --use-case string       Use case to hook automatically; it can accept "", "clone", "fork", "create" or "never" (default "never")
```

### SEE ALSO

* [gogh hook](gogh_hook.md)	 - Manage repository hooks (Lua scripts)

