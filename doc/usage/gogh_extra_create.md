## gogh extra create

Create a named extra template

### Synopsis

Create a named extra template from overlays.

This creates a reusable template that can be applied to any repository later.
By default, it uses the current directory's repository as the source.

```
gogh extra create <name> [flags]
```

### Options

```
  -h, --help              help for create
  -o, --overlay strings   Overlay names to include in the extra
  -s, --source string     Source repository (default: current directory)
```

### SEE ALSO

* [gogh extra](gogh_extra.md)	 - Manage repository extra files

