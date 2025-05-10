## gogh list

List local repositories

```
gogh list [flags]
```

### Options

```
  -f, --format string   
                        Print each local repository in a given format, where [format] can be one of "path",
                        "full-path", "fields" and "fields:[separator]".
                        
                        - path:
                        
                        	A part of the URL to specify a repository.  For example: "github.com/kyoh86/gogh"
                        
                        - full-path
                        
                        	A full path of the local repository.  For example:
                        	"/root/Projects/github.com/kyoh86/gogh".
                        
                        - fields
                        
                        	Tab separated all formats and properties of the local repository.
                        	i.e. [full-path]\t[path]\t[host]\t[owner]\t[name]
                        
                        - fields:[separator]
                        
                        	Like "fields" but with the explicit separator.
                        
  -h, --help            help for list
      --limit int       Max number of repositories to list. -1 means unlimited (default 100)
      --primary         List up repositories in just a primary root
  -q, --query string    Query for selecting repositories
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

