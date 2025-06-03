## gogh overlay add

Add an overlay file

```
gogh overlay add [flags] <source-path> <repo-pattern> <target-path>
```

### Examples

```
   Add an overlay file to a repository.
   The <source-path> is the path to the file you want to add as an overlay.
   The <repo-pattern> is the pattern of the repository you want to add the overlay to.
   The <target-path> is the path where the overlay file will be copied to in the repository.

   For example, to add a custom VSCode settings file to a repository, you can run:

     gogh overlay add /path/to/source/vscode/settings.json "github.com/owner/repo" .vscode/settings.json

   The overlay file will be copied to the repository when you run `gogh create`, `gogh clone` or `gogh fork`.

   You can also apply template files only for the `gogh create` command by using the `--for-init` flag:

     gogh overlay add --for-init /path/to/source/deno.jsonc "github.com/owner/deno-*" deno.jsonc

   This will copy the `deno.jsonc` file to the root of the repository only when you run `gogh create`
   if the repository matches the pattern `github.com/owner/deno-*`.
```

### Options

```
      --for-init   Register the overlay for 'gogh create' command
  -h, --help       help for add
```

### SEE ALSO

* [gogh overlay](gogh_overlay.md)	 - Manage repository overlay files

