package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/git"
	"github.com/kyoh86/gogh/internal/hub"
	"github.com/kyoh86/gogh/internal/mainutil"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	// init logs
	command.InitLog()

	app := kingpin.New("gogh", "GO GitHub project manager").Version(fmt.Sprintf("%s-%s (%s)", version, commit, date)).Author("kyoh86")
	app.Command("config", "Get and set options")

	cmds := map[string]func() error{}
	for _, f := range []func(*kingpin.Application) (string, func() error){
		configGetAll,
		configGet,
		configSet,
		configUnset,
		setup,

		roots,
		initialize,

		list,
		statuses,
		dump,
		find,
		where,
		repos,

		get,
		bulkGet,
		pipeGet,

		create,
		fork,
		remove,
	} {
		key, run := f(app)
		cmds[key] = run
	}
	if err := cmds[kingpin.MustParse(app.Parse(os.Args[1:]))](); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

func configGetAll(app *kingpin.Application) (string, func() error) {
	cmd := app.GetCommand("config").Command("get-all", "get all options").Alias("list").Alias("ls")

	return mainutil.WrapConfigurableCommand(cmd, func(_ gogh.Env, cfg *env.Config) error {
		return command.ConfigGetAll(cfg)
	})
}

func configGet(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("get", "get an option")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return mainutil.WrapConfigurableCommand(cmd, func(_ gogh.Env, cfg *env.Config) error {
		return command.ConfigGet(cfg, name)
	})
}

func configSet(app *kingpin.Application) (string, func() error) {
	var (
		name  string
		value string
	)
	cmd := app.GetCommand("config").Command("set", "set an option")
	cmd.Arg("name", "option name").Required().StringVar(&name)
	cmd.Arg("value", "option value").Required().StringVar(&value)

	return mainutil.WrapConfigurableCommand(cmd, func(ev gogh.Env, cfg *env.Config) error {
		return command.ConfigSet(ev, cfg, name, value)
	})
}

func configUnset(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("unset", "unset an option").Alias("rm")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return mainutil.WrapConfigurableCommand(cmd, func(ev gogh.Env, cfg *env.Config) error {
		return command.ConfigUnset(ev, cfg, name)
	})
}

func setup(app *kingpin.Application) (string, func() error) {
	var (
		force bool
	)
	cmd := app.Command("setup", "Setup gordon with wizards")
	cmd.Flag("force", "Ask even though that the option has already set").BoolVar(&force)

	return mainutil.WrapConfigurableCommand(cmd, func(ev gogh.Env, cfg *env.Config) error {
		return command.Setup(context.Background(), ev, cfg, force)
	})
}

func roots(app *kingpin.Application) (string, func() error) {
	var all bool
	cmd := app.Command("roots", "Show repositories' root").Alias("root")
	cmd.Flag("all", "Show all roots").Envar("GOGH_FLAG_ROOT_ALL").BoolVar(&all)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Roots(ev, all)
	})
}

func initialize(app *kingpin.Application) (string, func() error) {
	var (
		shell string
	)
	cmd := app.Command("init", `Generate shell script to initialize gogh

If you want to use "gogh cd", "gogh get --cd" and autocompletions for gogh,

set up gogh in your shell-rc file (".bashrc" / ".zshrc") like below.

eval "$(gogh init)"`).Hidden()
	cmd.Flag("shell", "Target shell path").Envar("SHELL").Hidden().StringVar(&shell)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Initialize(ev, "", shell)
	})
}

func list(app *kingpin.Application) (string, func() error) {
	var (
		format  command.ProjectListFormat
		primary bool
		query   string
	)
	cmd := app.Command("list", "List projects (local repositories)").Alias("ls")
	cmd.Flag("format", "Format of each repository").Short('f').Default(command.ProjectListFormatLabelRelPath).SetValue(&format)
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.List(ev, format.Formatter(), primary, query)
	})
}

func dump(app *kingpin.Application) (string, func() error) {
	var (
		primary bool
		query   string
	)
	cmd := app.Command("dump", "Dump list of projects (local repositories)")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.List(ev, gogh.URLFormatter(), primary, query)
	})
}

func find(app *kingpin.Application) (string, func() error) {
	var (
		primary bool
		spec    gogh.RepoSpec
	)
	cmd := app.Command("find", "Find a path of a project")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Find(ev, primary, &spec)
	})
}

func where(app *kingpin.Application) (string, func() error) {
	var (
		primary bool
		query   string
	)
	cmd := app.Command("where", "Where is a local project")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Where(ev, primary, query)
	})
}

func statuses(app *kingpin.Application) (string, func() error) {
	var (
		format  command.ProjectListFormat
		primary bool
		query   string
		detail  bool
	)
	cmd := app.Command("statuses", "List project statuses").Alias("status")
	cmd.Flag("format", "Format of each repository").Short('f').Default(command.ProjectListFormatLabelRelPath).SetValue(&format)
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Flag("detail", "Show status detail").Short('d').BoolVar(&detail)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Statuses(ev, new(git.Client), format.Formatter(), primary, query, detail)
	})
}

func repos(app *kingpin.Application) (string, func() error) {
	var (
		user        string
		own         bool
		collaborate bool
		member      bool
		archived    bool
		visibility  string
		sort        string
		direction   string
	)
	cmd := app.Command("repos", "List remote repositories").Alias("repo").Alias("search-repo").Alias("search-repos").Alias("list-repos")
	cmd.Flag("user", "Who has the repositories. Empty means the authenticated user").StringVar(&user)
	cmd.Flag("own", "Include repositories that are owned by the user").Default("true").BoolVar(&own)
	cmd.Flag("collaborate", "Include repositories that the user has been added to as a collaborator").Default("true").BoolVar(&collaborate)
	cmd.Flag("member", "Include repositories that the user has access to through being a member of an organization. This includes every repository on every team that the user is on").Default("true").BoolVar(&member)
	cmd.Flag("archived", "Include archived repository").BoolVar(&archived)
	cmd.Flag("visibility", "Include repositories that can be access public/private").Default("all").EnumVar(&visibility, "all", "public", "private")
	cmd.Flag("sort", "Sort repositories by").Default("full_name").EnumVar(&sort, "created", "updated", "pushed", "full_name")
	cmd.Flag("direction", "Sort direction").Default("default").EnumVar(&direction, "asc", "desc", "default")

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		ctx := context.Background()
		hubClient, err := hub.New(ctx, ev)
		if err != nil {
			return err
		}
		return command.Repos(ctx, ev, hubClient, user, own, collaborate, member, archived, visibility, sort, direction)
	})
}

func get(app *kingpin.Application) (string, func() error) {
	var (
		update  bool
		withSSH bool
		shallow bool
		specs   gogh.RepoSpecs
		cd      bool
		// unused: it is dummy to accept option.
		//         its function defined in init.*sh in the /sh/src
		//         if we want to use gogh get --cd,
	)
	cmd := app.Command("get", "Clone/sync with a remote repository")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Flag("cd", "Jump to the local project").BoolVar(&cd)
	cmd.Arg("repositories", "Target repositories (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&specs)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.GetAll(ev, new(git.Client), update, withSSH, shallow, specs)
	})
}

func bulkGet(app *kingpin.Application) (string, func() error) {
	var (
		update  bool
		withSSH bool
		shallow bool
	)
	cmd := app.Command("bulk-get", "Bulk get repositories specified in stdin")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Bulk(ev, new(git.Client), update, withSSH, shallow)
	})
}

func pipeGet(app *kingpin.Application) (string, func() error) {
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

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Pipe(ev, new(git.Client), update, withSSH, shallow, srcCmd, srcCmdArgs)
	})
}

func create(app *kingpin.Application) (string, func() error) {
	var (
		private        bool
		description    string
		homepage       *url.URL
		bare           bool
		template       string
		separateGitDir string
		shared         command.RepoShared
		spec           gogh.RepoSpec
	)
	cmd := app.Command("create", "Create a local project and a remote repository.").Alias("new")
	cmd.Flag("private", "Create a private repository").BoolVar(&private)
	cmd.Flag("description", "Use this text as the description of the GitHub repository").StringVar(&description)
	cmd.Flag("homepage", "Use this text as the URL of the GitHub repository").URLVar(&homepage)
	cmd.Flag("bare", "Create a bare repository. If GIT_DIR environment is not set, it is set to the current working directory").BoolVar(&bare)
	cmd.Flag("template", "Specify the directory from which templates will be used").ExistingDirVar(&template)
	cmd.Flag("separate-git-dir", `Instead of initializing the repository as a directory to either $GIT_DIR or ./.git/`).StringVar(&separateGitDir)
	cmd.Flag("shared", "Specify that the Git repository is to be shared amongst several users.").SetValue(&shared)
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		ctx := context.Background()
		hubClient, err := hub.New(ctx, ev)
		if err != nil {
			return err
		}
		return command.New(ctx, ev, new(git.Client), hubClient, private, description, homepage, bare, template, separateGitDir, shared, &spec)
	})
}

func fork(app *kingpin.Application) (string, func() error) {
	var (
		update       bool
		withSSH      bool
		shallow      bool
		organization string
		spec         gogh.RepoSpec
	)
	cmd := app.Command("fork", "Clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin")
	cmd.Flag("update", "Update the local project if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Flag("org", "Fork the repository within this organization").PlaceHolder("ORGANIZATION").StringVar(&organization)
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		ctx := context.Background()
		hubClient, err := hub.New(ctx, ev)
		if err != nil {
			return err
		}
		return command.Fork(ctx, ev, new(git.Client), hubClient, update, withSSH, shallow, organization, &spec)
	})
}

func remove(app *kingpin.Application) (string, func() error) {
	var (
		primary bool
		query   string
	)
	cmd := app.Command("remove", "Delete projects").Alias("rm").Alias("delete").Alias("del")
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("query", "Project name query").StringVar(&query)

	return mainutil.WrapCommand(cmd, func(ev gogh.Env) error {
		return command.Delete(ev, primary, query)
	})
}
