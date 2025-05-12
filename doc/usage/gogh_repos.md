## gogh repos

List remote repositories

```
gogh repos [flags]
```

### Options

```
      --archived           Show only archived repositories
      --color string       Colorize the output; It can accept 'auto', 'always' or 'never' (default "auto")
      --fork               Show only forks
      --format string      
                           Print each repository in a given format, where [format] can be one of "table", "ref",
                           "url" or "json".
                           
  -h, --help               help for repos
      --limit int          Max number of repositories to list. -1 means unlimited (default -1)
      --no-archived        Omit archived repositories
      --no-fork            Omit forks (default true)
      --order sort         Directions in which to order a list of items when provided an sort flag; it can accept "asc", "ascending", "ASC", "ASCENDING", "desc", "descending", "DESC" or "DESCENDING"
      --private            Show only private repositories
      --public             Show only public repositories
      --relation strings   The relation of user to each repository; it can accept "owner", "organization-member", "organization_member", "organizationMember" or "collaborator" (default [owner,organizationMember])
      --sort string        Property by which repository be ordered; it can accept "CREATED_AT", "created_at", "created-at", "createdAt", "NAME", "name", "PUSHED_AT", "pushed_at", "pushed-at", "pushedAt", "STARGAZERS", "stargazers", "UPDATED_AT", "updated_at", "updated-at" or "updatedAt"
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

