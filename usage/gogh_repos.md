## gogh repos

List remote repositories

```
gogh repos [flags]
```

### Options

```
      --color string       Colorize the output; It can accept 'auto', 'always' or 'never' (default "auto")
      --fork               Show only forks
      --format string      The formatting style for each repository; it can accept "spec", "url", "json" or "table" (default "table")
  -h, --help               help for repos
      --limit int          Max number of repositories to list. -1 means unlimited (default -1)
      --no-fork            Omit forks (default true)
      --order sort         Directions in which to order a list of items when provided an sort flag; it can accept "ASC", "DESC"
      --private            Show only private repositories
      --public             Show only public repositories
      --relation strings   The relation of user to each repository; it can accept "owner", "organizationMember", "collaborator" (default [owner])
      --sort string        Property by which repository be ordered; it can accept "CREATED_AT", "UPDATED_AT", "PUSHED_AT", "NAME", "STARGAZERS"
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub project manager

