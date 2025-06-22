## gogh overlay add

Add an overlay file

```
gogh overlay add [flags] <name> <target-path> <source-path>
```

### Examples

```
   Add an overlay file to a repository.
   The <name> is the name of the overlay, which is used to identify it.
   The <target-path> is the path where the overlay file will be copied to in the repository.
   The <source-path> is the path to the file you want to add as an overlay.

   For example, to add a custom VSCode settings file to a repository, you can run:

     gogh overlay add vsc-setting /path/to/source/vscode/settings.json .vscode/settings.json

   The overlay file will be copied to the repository when you run `gogh overlay apply`.
```

### Options

```
      --for-init   Register the overlay for 'gogh create' command
  -h, --help       help for add
```

### SEE ALSO

* [gogh overlay](gogh_overlay.md)	 - Manage repository overlay files

