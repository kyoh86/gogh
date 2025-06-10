## gogh overlay apply

Apply overlays to specified repositories

```
gogh overlay apply [flags] [[[<host>/]<owner>/]<name>...]
```

### Examples

```
  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is a local name for the remote repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"
```

### Options

```
      --for-init gogh create   Apply overlays only for gogh create command (useful for templates)
  -h, --help                   help for apply
      --repo-pattern string    Force apply overlays having this pattern, ignoring automatic repository name matching (useful for applying specific overlays or templates that would not normally match)
```

### SEE ALSO

* [gogh overlay](gogh_overlay.md)	 - Manage repository overlay files

