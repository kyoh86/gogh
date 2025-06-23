## gogh extra save

Save excluded files as auto-apply extra

### Synopsis

Save files that are excluded by .gitignore as auto-apply extra.
These extra will be automatically applied when the repository is cloned.

```
gogh extra save <repository> [flags]
```

### Examples

```
  save github.com/kyoh86/example
  save .  # Save from current directory repository

  It accepts a short notation for the repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".
```

### Options

```
  -h, --help   help for save
```

### SEE ALSO

* [gogh extra](gogh_extra.md)	 - Manage repository extra files

