## gogh overlay

Manage repository overlay files

### Examples

```
   Overlay files are used to put custom files into repositories.
   They are useful to add files that are not tracked by the repository, such as editor configurations or scripts.

   For example, to add a custom VSCode settings file to a repository, you can run:

     gogh overlay add /path/to/source/vscode/settings.json "github.com/owner/repo" .vscode/settings.json

   Then when you run `gogh create`, `gogh clone` or `gogh fork`, the files will be copied to the repository.

   You can also apply template files only for the `gogh create` command by using the `--for-init` flag:

     gogh overlay add --for-init /path/to/source/deno.jsonc "github.com/owner/deno-*" deno.jsonc

   This will copy the `deno.jsonc` file to the root of the repository only when you run `gogh create`
   if the repository matches the pattern `github.com/owner/deno-*`.

   And then you can use the `gogh overlay apply` command to apply the overlay files manually.

   You can create overlay files that never be applied to the repository automatically,
   (and only be applied manually by `gogh overlay apply` command),
   you can set the `--repo-pattern` flag to never match any repository.
```

### Options

```
  -h, --help   help for overlay
```

### SEE ALSO

* [gogh](gogh.md)	 - GO GitHub local repository manager
* [gogh overlay add](gogh_overlay_add.md)	 - Add an overlay file
* [gogh overlay apply](gogh_overlay_apply.md)	 - Apply an overlay to a repository
* [gogh overlay edit](gogh_overlay_edit.md)	 - Edit an existing overlay (with $EDITOR)
* [gogh overlay extract](gogh_overlay_extract.md)	 - Extract untracked files as overlays
* [gogh overlay list](gogh_overlay_list.md)	 - List registered overlays
* [gogh overlay remove](gogh_overlay_remove.md)	 - Remove an overlay
* [gogh overlay show](gogh_overlay_show.md)	 - Show an overlay
* [gogh overlay update](gogh_overlay_update.md)	 - Update an existing overlay

