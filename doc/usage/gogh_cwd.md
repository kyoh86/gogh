## gogh cwd

Print the local repository which the current working directory belongs to

```
gogh cwd [flags]
```

### Options

```
  -f, --format string   
                        Print local repository in a given format, where [format] can be one of "path",
                        "full-path", "json", "fields" and "fields:[separator]".
                        
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
                        
  -h, --help            help for cwd
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager

