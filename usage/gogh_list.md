## gogh list

List local projects

```
gogh list [flags]
```

### Options

```
  -f, --format string   
                        Print each project in a given format, where [format] can be one of "rel-path", "rel-file-path",
                        "full-file-path", "url", "fields" and "fields:[separator]".
                        
                        - rel-path:
                        
                        	A part of the URL to specify a repository.  For example: "github.com/kyoh86/gogh"
                        
                        - rel-file-path:
                        
                        	A relative file path of the project from gogh roots.  For example in windows:
                        	"github.com\kyoh86\gogh"; in other case: "github.com/kyoh86/gogh".
                        
                        - full-file-path
                        
                        	A full file path of the project.  For example in Windows:
                        	"C:\Users\kyoh86\Projects\github.com\kyoh86\gogh"; in other case:
                        	"/root/Projects/github.com/kyoh86/gogh".
                        
                        - url
                        
                        	A URL of the repository.
                        
                        - fields
                        
                        	Tab separated all formats and properties of the project.
                        	i.e. [full-file-path]\t[rel-file-path]\t[url]\t[rel-path]\t[host]\t[owner]\t[name]
                        
                        - fields:[separator]
                        
                        	Like "fields" but with the explicit separator.
                        
  -h, --help            help for list
      --primary         List up projects in just a primary root
  -q, --query string    Query for selecting projects
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub project manager

