## gogh clone

Clone repositories to local

```
gogh clone [flags] [[OWNER/]NAME[=ALIAS]]...
```

### Examples

```
  It accepts a shortly notation for a repository
  (for example, "github.com/kyoh86/example") like below.
    - "NAME": e.g. "example"; 
    - "OWNER/NAME": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default".

  It accepts an alias for each repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example"
  into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"
```

### Options

```
      --dryrun   Displays the operations that would be performed using the specified command without actually running them
  -h, --help     help for clone
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub project manager

