## gogh overlay extract

Extract untracked files as overlays

```
gogh overlay extract [flags] [[[<host>/]<owner>/]<name>...]
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
      --for-init gogh create   Register the overlay for gogh create command
      --force                  Do NOT confirm to extract for each file
  -h, --help                   help for extract
      --pattern string         Pattern to match repositories (e.g., 'github.com/owner/repo', '**/gogh'; default: repository reference)
```

### SEE ALSO

* [gogh overlay](gogh_overlay.md)	 - Manage repository overlay files

