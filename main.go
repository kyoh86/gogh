package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/gogh"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	log.SetOutput(os.Stderr)
	app := kingpin.New("gogh", "GO GitHub project manager").Version(fmt.Sprintf("%s-%s (%s)", version, commit, date)).Author("kyoh86")

	cmds := map[string]func() error{}
	for _, f := range []func(*kingpin.Application) (string, func() error){
		get,
		bulk,
		pipe,
		fork,
		create,
		list,
		find,
		root,
		setup,
	} {
		key, run := f(app)
		cmds[key] = run
	}
	if err := cmds[kingpin.MustParse(app.Parse(os.Args[1:]))](); err != nil {
		log.Fatal(err)
	}
}

func wrapContext(f func(gogh.Context) error) func() error {
	return func() error {
		ctx, err := gogh.CurrentContext(context.Background())
		if err != nil {
			return err
		}
		return f(ctx)
	}
}

func get(app *kingpin.Application) (string, func() error) {
	var (
		update    bool
		withSSH   bool
		shallow   bool
		repoSpecs gogh.Specs
	)
	cmd := app.Command("get", "Clone/sync with a remote repository")
	cmd.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Arg("repositories", "Target repositories (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&repoSpecs)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.GetAll(ctx, update, withSSH, shallow, repoSpecs)
	})
}

func bulk(app *kingpin.Application) (string, func() error) {
	var (
		update  bool
		withSSH bool
		shallow bool
	)
	cmd := app.Command("bulk", "Bulk get repositories specified in stdin")
	cmd.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.Bulk(ctx, update, withSSH, shallow)
	})
}

func pipe(app *kingpin.Application) (string, func() error) {
	var (
		update      bool
		withSSH     bool
		shallow     bool
		command     string
		commandArgs []string
	)
	cmd := app.Command("pipe", "Bulk get repositories specified from other command output")
	cmd.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Arg("command", "Subcommand calling to get import paths").StringVar(&command)
	cmd.Arg("command-args", "Arguments that will be passed to subcommand").StringsVar(&commandArgs)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.Pipe(ctx, update, withSSH, shallow, command, commandArgs)
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
		repoSpec     gogh.RepoSpec
	)
	cmd := app.Command("fork", "Clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin")
	cmd.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&update)
	cmd.Flag("ssh", "Clone with SSH").BoolVar(&withSSH)
	cmd.Flag("shallow", "Do a shallow clone").BoolVar(&shallow)
	cmd.Flag("no-remote", "Skip adding a git remote for the fork").BoolVar(&noRemote)
	cmd.Flag("remote-name", "Set the name for the new git remote").PlaceHolder("REMOTE").StringVar(&remoteName)
	cmd.Flag("org", "Fork the repository within this organization").PlaceHolder("ORGANIZATION").StringVar(&organization)
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&repoSpec)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.Fork(ctx, update, withSSH, shallow, noRemote, remoteName, organization, repoSpec)
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
		shared         gogh.Shared
		repoName       gogh.RepoName
	)
	cmd := app.Command("new", "Create a repository in local and remote.").Alias("create")
	cmd.Flag("private", "Create a private repository").BoolVar(&private)
	cmd.Flag("description", "Use this text as the description of the GitHub repository").StringVar(&description)
	cmd.Flag("homepage", "Use this text as the URL of the GitHub repository").URLVar(&homepage)
	cmd.Flag("browse", "Open the new repository in a web browser").Short('o').BoolVar(&browse)
	cmd.Flag("copy", "Put the URL of the new repository to clipboard instead of printing it").Short('c').BoolVar(&clip)
	cmd.Flag("bare", "Create a bare repository. If GIT_DIR environment is not set, it is set to the current working directory").BoolVar(&bare)
	cmd.Flag("template", "Specify the directory from which templates will be used").ExistingDirVar(&template)
	cmd.Flag("separate-git-dir", `Instead of initializing the repository as a directory to either $GIT_DIR or ./.git/`).StringVar(&separateGitDir)
	cmd.Flag("shared", "Specify that the Git repository is to be shared amongst several users.").SetValue(&shared)
	cmd.Arg("repository name", "<user>/<name>").Required().SetValue(&repoName)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.New(ctx, private, description, homepage, browse, clip, bare, template, separateGitDir, shared, repoName)
	})
}

func list(app *kingpin.Application) (string, func() error) {
	var (
		exact    bool
		fullPath bool
		short    bool
		primary  bool
		query    string
	)
	cmd := app.Command("list", "List local repositories")
	cmd.Flag("exact", "Perform an exact match").Short('e').BoolVar(&exact)
	cmd.Flag("full-path", "Print full paths").Short('f').BoolVar(&fullPath)
	cmd.Flag("short", "Print short names").Short('s').BoolVar(&short)
	cmd.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&primary)
	cmd.Arg("query", "Repository name query").StringVar(&query)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.List(ctx, exact, fullPath, short, primary, query)
	})
}

func find(app *kingpin.Application) (string, func() error) {
	var name string
	cmd := app.Command("find", "Find a path of a local repository")
	cmd.Arg("name", "Target repository name").Required().StringVar(&name)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.Find(ctx, name)
	})
}

func root(app *kingpin.Application) (string, func() error) {
	var all bool
	cmd := app.Command("root", "Show repositories' root")
	cmd.Flag("all", "Show all roots").BoolVar(&all)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.Root(ctx, all)
	})
}

func setup(app *kingpin.Application) (string, func() error) {
	var (
		cdFuncName string
		shell      string
	)
	cmd := app.Command("setup", "Generate shell script to setup gogh").Hidden()
	cmd.Flag("cd-function-name", "Name of the function to define").Default("gogogh").Hidden().StringVar(&cdFuncName)
	cmd.Flag("shell", "Target shell path").Envar("SHELL").Hidden().StringVar(&shell)

	return cmd.FullCommand(), wrapContext(func(ctx gogh.Context) error {
		return gogh.Setup(ctx, cdFuncName, shell)
	})
}
