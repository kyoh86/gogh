package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/mainutil"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	log.SetOutput(os.Stderr)

	app := kingpin.New("gogh", "GO GitHub project manager").Version(fmt.Sprintf("%s-%s (%s)", version, commit, date)).Author("kyoh86")
	app.Command("config", "Get and set options")

	cmds := map[string]func() error{}
	for _, f := range []func(*kingpin.Application) (string, func() error){
		configGetAll,
		configGet,
		configPut,
		configUnset,

		get,
		bulk,
		pipe,
		fork,
		create,
		where,
		list,
		dump,
		find,
		root,
		setup,

		repos,
	} {
		key, run := f(app)
		cmds[key] = run
	}
	if err := cmds[kingpin.MustParse(app.Parse(os.Args[1:]))](); err != nil {
		log.Fatalf("error: %s", err)
	}
}

func configGetAll(app *kingpin.Application) (string, func() error) {
	cmd := app.GetCommand("config").Command("get-all", "get all options").Alias("list").Alias("ls")

	return mainutil.WrapConfigurableCommand(cmd, command.ConfigGetAll)
}

func configGet(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("get", "get an option")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return mainutil.WrapConfigurableCommand(cmd, func(cfg *config.Config) error {
		return command.ConfigGet(cfg, name)
	})
}

func configPut(app *kingpin.Application) (string, func() error) {
	var (
		name  string
		value string
	)
	cmd := app.GetCommand("config").Command("put", "put an option").Alias("set")
	cmd.Arg("name", "option name").Required().StringVar(&name)
	cmd.Arg("value", "option value").Required().StringVar(&value)

	return mainutil.WrapConfigurableCommand(cmd, func(cfg *config.Config) error {
		return command.ConfigPut(cfg, name, value)
	})
}

func configUnset(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("unset", "unset an option").Alias("rm")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return mainutil.WrapConfigurableCommand(cmd, func(cfg *config.Config) error {
		return command.ConfigUnset(cfg, name)
	})
}

func get(app *kingpin.Application) (string, func() error) {
	var (
		update    bool
		withSSH   bool
		shallow   bool
		repoNames gogh.Repos
		cd        bool
		// unused: it is dummy to accept option.
		//         its function defined in init.*sh in the /sh/src
		//         if we want to use gogh get --cd,
	)
	cmd := app.Command("get", "Clone/sync with a remote repository")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Flag("cd", "Jump to the local project").BoolVar(&cd)
	cmd.Arg("repositories", "Target repositories (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&repoNames)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.GetAll(ctx, update, withSSH, shallow, repoNames)
	})
}

func bulk(app *kingpin.Application) (string, func() error) {
	var (
		update  bool
		withSSH bool
		shallow bool
	)
	cmd := app.Command("bulk-get", "Bulk get repositories specified in stdin")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Bulk(ctx, update, withSSH, shallow)
	})
}

func pipe(app *kingpin.Application) (string, func() error) {
	var (
		update     bool
		withSSH    bool
		shallow    bool
		srcCmd     string
		srcCmdArgs []string
	)
	cmd := app.Command("pipe-get", "Bulk get repositories specified from other command output")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Arg("command", "Subcommand calling to get import paths").StringVar(&srcCmd)
	cmd.Arg("command-args", "Arguments that will be passed to subcommand").StringsVar(&srcCmdArgs)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Pipe(ctx, update, withSSH, shallow, srcCmd, srcCmdArgs)
	})
}

func fork(app *kingpin.Application) (string, func() error) {
	var (
		update       bool
		withSSH      bool
		shallow      bool
		noRemote     bool
		remoteName   string
		organization string
		repo         gogh.Repo
	)
	cmd := app.Command("fork", "Clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Flag("no-remote", "Skip adding a git remote for the fork").BoolVar(&noRemote)
	cmd.Flag("remote-name", "Set the name for the new git remote").PlaceHolder("REMOTE").StringVar(&remoteName)
	cmd.Flag("org", "Fork the repository within this organization").PlaceHolder("ORGANIZATION").StringVar(&organization)
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&repo)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Fork(ctx, update, withSSH, shallow, noRemote, remoteName, organization, &repo)
	})
}

func create(app *kingpin.Application) (string, func() error) {
	var (
		private        bool
		description    string
		homepage       *url.URL
		browse         bool
		clip           bool
		bare           bool
		template       string
		separateGitDir string
		shared         gogh.ProjectShared
		repo           gogh.Repo
	)
	cmd := app.Command("new", "Create a local project and a remote repository.").Alias("create")
	cmd.Flag("private", "Create a private repository").BoolVar(&private)
	cmd.Flag("description", "Use this text as the description of the GitHub repository").StringVar(&description)
	cmd.Flag("homepage", "Use this text as the URL of the GitHub repository").URLVar(&homepage)
	cmd.Flag("browse", "Open the new repository in a web browser").Short('o').BoolVar(&browse)
	cmd.Flag("copy", "Put the URL of the new repository to clipboard instead of printing it").Short('c').BoolVar(&clip)
	cmd.Flag("bare", "Create a bare repository. If GIT_DIR environment is not set, it is set to the current working directory").BoolVar(&bare)
	cmd.Flag("template", "Specify the directory from which templates will be used").ExistingDirVar(&template)
	cmd.Flag("separate-git-dir", `Instead of initializing the repository as a directory to either $GIT_DIR or ./.git/`).StringVar(&separateGitDir)
	cmd.Flag("shared", "Specify that the Git repository is to be shared amongst several users.").SetValue(&shared)
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&repo)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.New(ctx, private, description, homepage, browse, clip, bare, template, separateGitDir, shared, &repo)
	})
}

func where(app *kingpin.Application) (string, func() error) {
	var (
		primary bool
		exact   bool
		query   string
	)
	cmd := app.Command("where", "Where is a local project")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Flag("exact", "Specifies name of the project in query").Short('e').BoolVar(&exact)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Where(ctx, primary, exact, query)
	})
}

func list(app *kingpin.Application) (string, func() error) {
	var (
		format   command.ProjectListFormat
		primary  bool
		isPublic bool
		query    string
	)
	cmd := app.Command("list", "List projects (local repositories)").Alias("ls")
	cmd.Flag("format", "Format of each repository").Short('f').Default(command.ProjectListFormatLabelRelPath).SetValue(&format)
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Flag("public", "Only projects which are referred to public repositories").BoolVar(&isPublic)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.List(ctx, format.Formatter(), primary, isPublic, query)
	})
}

func dump(app *kingpin.Application) (string, func() error) {
	var (
		primary  bool
		isPublic bool
		query    string
	)
	cmd := app.Command("dump", "Dump list of projects (local repositories)")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Flag("public", "Only projects which are referred to public repositories").BoolVar(&isPublic)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.List(ctx, gogh.URLFormatter(), primary, isPublic, query)
	})
}

func find(app *kingpin.Application) (string, func() error) {
	var (
		primary bool
		query   string
	)
	cmd := app.Command("find", "Find a path of a project. This is shorthand of `gogh where --exact`")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Where(ctx, primary, true, query)
	})
}

func root(app *kingpin.Application) (string, func() error) {
	var all bool
	cmd := app.Command("root", "Show repositories' root")
	cmd.Flag("all", "Show all roots").Envar("GOGH_FLAG_ROOT_ALL").BoolVar(&all)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Root(ctx, all)
	})
}

func setup(app *kingpin.Application) (string, func() error) {
	var (
		shell string
	)
	cmd := app.Command("setup", `Generate shell script to setup gogh

If you want to use "gogh cd", "gogh get --cd" and autocompletions for gogh,

set up gogh in your shell-rc file (".bashrc" / ".zshrc") like below.

eval "$(gogh setup)"`).Hidden()
	cmd.Flag("shell", "Target shell path").Envar("SHELL").Hidden().StringVar(&shell)

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Setup(ctx, "", shell)
	})
}

func repos(app *kingpin.Application) (string, func() error) {
	var (
		user        string
		own         bool
		collaborate bool
		member      bool
		visibility  string
		sort        string
		direction   string
	)
	cmd := app.Command("repo", "List remote repositories").Alias("repos").Alias("search-repo").Alias("search-repos")
	cmd.Flag("user", "Who has the repositories. Empty means the authenticated user").StringVar(&user)
	cmd.Flag("own", "Include repositories that are owned by the user").Default("true").BoolVar(&own)
	cmd.Flag("collaborate", "Include repositories that the user has been added to as a collaborator").Default("true").BoolVar(&collaborate)
	cmd.Flag("member", "Include repositories that the user has access to through being a member of an organization. This includes every repository on every team that the user is on").Default("true").BoolVar(&member)
	cmd.Flag("visibility", "Include repositories that can be access public/private").Default("all").EnumVar(&visibility, "all", "public", "private")
	cmd.Flag("sort", "Sort repositories by").Default("full_name").EnumVar(&sort, "created", "updated", "pushed", "full_name")
	cmd.Flag("direction", "Sort direction").Default("default").EnumVar(&direction, "asc", "desc", "default")

	return mainutil.WrapCommand(cmd, func(ctx gogh.Context) error {
		return command.Repos(ctx, user, own, collaborate, member, visibility, sort, direction)
	})
}
