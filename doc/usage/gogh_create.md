## gogh create

Create a new local and remote repository

```
gogh create [flags] [[OWNER/]NAME[=ALIAS]]
```

### Options

```
      --auto-init                   Create an initial commit with empty README
      --clone-retry-limit int       The number of retries to clone a repository (default 3)
      --delete-branch-on-merge      Allow automatically deleting head branches when pull requests are merged
      --description string          A short description of the repository
      --disable-downloads           Disable "Downloads" page
      --disable-issues              Disable issues for the repository
      --disable-projects            Disable projects for the repository
      --disable-wiki                Disable Wiki for the repository
      --dryrun                      Displays the operations that would be performed using the specified command without actually running them
      --gitignore-template string   Desired language or platform .gitignore template to apply when "auto-init" flag is set. Use the name of the template without the extension. For example, "Haskell"
  -h, --help                        help for create
      --homepage string             A URL with more information about the repository
      --include-all-branches        Create all branches in the template
      --is-template                 Whether the repository is available as a template
      --license-template string     Choose an open source license template that best suits your needs, and then use the license keyword as the license_template string when "auto-init" flag is set. For example, "mit" or "mpl-2.0"
      --prevent-merge-commit        Prevent merging pull requests with a merge commit
      --prevent-rebase-merge        Prevent rebase-merging pull requests
      --prevent-squash-merge        Prevent squash-merging pull requests
      --private                     Whether the repository is private
      --template string             Create new repository from the template
  -t, --timeout duration            Timeout for the request (default 5s)
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

