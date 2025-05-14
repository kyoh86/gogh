## gogh repos

List remote repositories

```
gogh repos [flags]
```

### Options

```
      --archive string     Show only archived/not-archived repositories; it can accept "archived" or "not-archived"
      --color string       Colorize the output; it can accept "auto", "always" or "never" (default "auto")
      --fork string        Show only forked/not-forked repositories; it can accept "forked" or "not-forked"
  -f, --format string      
                           Print each repository in a given format, where [format] can be one of "table", "ref",
                           "url" or "json".
                           
  -h, --help               help for repos
      --limit int          Max number of repositories to list. -1 means unlimited (default 30)
      --order sort         Directions in which to order a list of items when provided a sort flag; it can accept "asc", "ascending", "desc" or "descending"
      --privacy string     Show only public/private repositories; it can accept "private" or "public"
      --relation strings   The relation of user to each repository; it can accept "owner", "organization-member" or "collaborator" (default [owner,organizationMember])
      --sort string        Property by which repository be ordered; it can accept "created-at", "name", "pushed-at", "stargazers" or "updated-at"
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

