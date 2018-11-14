package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/repo"
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

	type getParams struct {
		Update  bool
		WithSSH bool
		Shallow bool
	}

	var optGet struct {
		getParams

		RepoSpecs repo.Specs
	}
	cmdGet := app.Command("get", "Clone/sync with a remote repository")
	cmdGet.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&optGet.Update)
	cmdGet.Flag("ssh", "Clone with SSH").BoolVar(&optGet.WithSSH)
	cmdGet.Flag("shallow", "Do a shallow clone").BoolVar(&optGet.Shallow)
	cmdGet.Arg("repositories", "Target repositories (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&optGet.RepoSpecs)

	var optBulk struct {
		getParams
	}
	cmdBulk := app.Command("bulk", "Bulk get repositories specified in stdin")
	cmdBulk.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&optBulk.Update)
	cmdBulk.Flag("ssh", "Clone with SSH").BoolVar(&optBulk.WithSSH)
	cmdBulk.Flag("shallow", "Do a shallow clone").BoolVar(&optBulk.Shallow)

	var optPipe struct {
		getParams

		Command     string
		CommandArgs []string
	}
	cmdPipe := app.Command("pipe", "Bulk get repositories specified from other command output")
	cmdPipe.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&optPipe.Update)
	cmdPipe.Flag("ssh", "Clone with SSH").BoolVar(&optPipe.WithSSH)
	cmdPipe.Flag("shallow", "Do a shallow clone").BoolVar(&optPipe.Shallow)
	cmdPipe.Arg("command", "Subcommand calling to get import paths").StringVar(&optPipe.Command)
	cmdPipe.Arg("command-args", "Arguments that will be passed to subcommand").StringsVar(&optPipe.CommandArgs)

	var optFork struct {
		getParams

		NoRemote     bool
		RemoteName   string
		Organization string

		RepoSpec repo.Spec
	}
	cmdFork := app.Command("fork", "Clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin")
	cmdFork.Flag("update", "Update local repository if cloned already").Short('u').BoolVar(&optFork.Update)
	cmdFork.Flag("ssh", "Clone with SSH").BoolVar(&optFork.WithSSH)
	cmdFork.Flag("shallow", "Do a shallow clone").BoolVar(&optFork.Shallow)
	cmdFork.Flag("no-remote", "Skip adding a git remote for the fork").BoolVar(&optFork.NoRemote)
	cmdFork.Flag("remote-name", "Set the name for the new git remote").PlaceHolder("REMOTE").StringVar(&optFork.RemoteName)
	cmdFork.Flag("org", "Fork the repository within this organization").PlaceHolder("ORGANIZATION").StringVar(&optFork.Organization)
	cmdFork.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&optFork.RepoSpec)

	var optNew struct {
		// hub create options
		Private        bool
		Description    string
		Homepage       *url.URL
		Browse         bool
		Copy           bool
		Bare           bool
		Template       string
		SeparateGitDir string
		Shared         repo.Shared
		RepoName       repo.Name
	}
	cmdNew := app.Command("new", "Create a repository in local and remote.").Alias("create")
	cmdNew.Flag("private", "Create a private repository").BoolVar(&optNew.Private)
	cmdNew.Flag("description", "Use this text as the description of the GitHub repository").StringVar(&optNew.Description)
	cmdNew.Flag("homepage", "Use this text as the URL of the GitHub repository").URLVar(&optNew.Homepage)
	cmdNew.Flag("browse", "Open the new repository in a web browser").Short('o').BoolVar(&optNew.Browse)
	cmdNew.Flag("copy", "Put the URL of the new repository to clipboard instead of printing it").Short('c').BoolVar(&optNew.Copy)
	cmdNew.Flag("bare", "Create a bare repository. If GIT_DIR environment is not set, it is set to the current working directory").BoolVar(&optNew.Bare)
	cmdNew.Flag("template", "Specify the directory from which templates will be used").ExistingDirVar(&optNew.Template)
	cmdNew.Flag("separate-git-dir", `Instead of initializing the repository as a directory to either $GIT_DIR or ./.git/`).StringVar(&optNew.SeparateGitDir)
	cmdNew.Flag("shared", "Specify that the Git repository is to be shared amongst several users.").SetValue(&optNew.Shared)
	cmdNew.Arg("repository name", "<user>/<name>").Required().SetValue(&optNew.RepoName)

	var optList struct {
		Exact    bool
		FullPath bool
		Short    bool
		Primary  bool
		Query    string
	}
	cmdList := app.Command("list", "List local repositories")
	cmdList.Flag("exact", "Perform an exact match").Short('e').BoolVar(&optList.Exact)
	cmdList.Flag("full-path", "Print full paths").Short('f').BoolVar(&optList.FullPath)
	cmdList.Flag("short", "Print short names").Short('s').BoolVar(&optList.Short)
	cmdList.Flag("primary", "Only in primary root directory").Short('p').BoolVar(&optList.Primary)
	cmdList.Arg("query", "Repository name query").StringVar(&optList.Query)

	var optFind struct {
		Name string
	}
	cmdFind := app.Command("find", "Find a path of a local repository")
	cmdFind.Arg("name", "Target repository name").Required().StringVar(&optFind.Name)

	var optRoot struct {
		All bool
	}
	cmdRoot := app.Command("root", "Show repositories' root")
	cmdRoot.Flag("all", "Show all roots").BoolVar(&optRoot.All)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdGet.FullCommand():
		if err := gogh.GetAll(optGet.Update, optGet.WithSSH, optGet.Shallow, optGet.RepoSpecs); err != nil {
			log.Fatal(err)
		}
	case cmdBulk.FullCommand():
		if err := gogh.Bulk(optBulk.Update, optBulk.WithSSH, optBulk.Shallow); err != nil {
			log.Fatal(err)
		}
	case cmdPipe.FullCommand():
		if err := gogh.Pipe(optPipe.Update, optPipe.WithSSH, optPipe.Shallow, optPipe.Command, optPipe.CommandArgs); err != nil {
			log.Fatal(err)
		}
	case cmdFork.FullCommand():
		if err := gogh.Fork(optFork.Update, optFork.WithSSH, optFork.Shallow, optFork.NoRemote, optFork.RemoteName, optFork.Organization, optFork.RepoSpec); err != nil {
			log.Fatal(err)
		}
	case cmdNew.FullCommand():
		if err := gogh.New(optNew.Private, optNew.Description, optNew.Homepage, optNew.Browse, optNew.Copy, optNew.Bare, optNew.Template, optNew.SeparateGitDir, optNew.Shared, optNew.RepoName); err != nil {
			log.Fatal(err)
		}
	case cmdList.FullCommand():
		if err := gogh.List(optList.Exact, optList.FullPath, optList.Short, optList.Primary, optList.Query); err != nil {
			log.Fatal(err)
		}
	case cmdFind.FullCommand():
		if err := gogh.Find(optFind.Name); err != nil {
			log.Fatal(err)
		}
	case cmdRoot.FullCommand():
		if err := gogh.Root(optRoot.All); err != nil {
			log.Fatal(err)
		}
	}
}
