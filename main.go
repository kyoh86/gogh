package main

import (
	"fmt"
	"log"
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
	app := kingpin.New("gogh", "A Project Manager").Version(fmt.Sprintf("%s-%s (%s)", version, commit, date)).Author("kyoh86")

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

	cmdNew := app.Command("new", "Create a repository in local and remote.")
	_ = cmdNew
	/*TODO: hub create options
	arg: [[ORGANIZATION/]NAME]
	-p     Create a private repository.

	-d DESCRIPTION
	      Use this text as the description of the GitHub repository.

	-h HOMEPAGE
	      Use this text as the URL of the GitHub repository.

	-o, --browse
	      Open the new repository in a web browser.

	-c, --copy
	      Put the URL of the new repository to clipboard instead of printing it.

	[ORGANIZATION/]NAME
	      The name for the repository on GitHub (default: name of the current working directory).

	      Optionally, create the repository within ORGANIZATION.
	*/
	/*TODO: git init options
	arg: [directory]
	-q, --quiet
	    Only print error and warning messages; all other output will be suppressed.

	--bare
	    Create a bare repository. If GITT_DIR environment is not set, it is set to the current working directory.

	--template=<template_directory>
	    Specify the directory from which templates will be used. (See the "TEMPLATE DIRECTORY" section below.)

	--separate-git-dir=<git dir>
	    Instead of initializing the repository as a directory to either $GITT_DIR or ./.git/, create a text file there containing the path to the actual repository. This file acts as filesystem-agnostic Git symbolic link to the repository.

	    If this is reinitialization, the repository will be moved to the specified path.

	--shared[=(false|true|umask|group|all|world|everybody|0xxx)]
	    Specify that the Git repository is to be shared amongst several users. This allows users belonging to the same group to push into that repository. When specified, the config variable "core.sharedRepository" is set so that files and directories under $GITT_DIR are created with
	    the requested permissions. When not specified, Git will use permissions reported by umask(2).

	    The option can have the following values, defaulting to _group if no value is given:

	    _umask (or _false)
	        Use permissions reported by umask(2). The default, when --shared is not specified.

	    _group (or _true)
	        Make the repository group-writable, (and g+sx, since the git group may be not the primary group of all users). This is used to loosen the permissions of an otherwise safe umask(2) value. Note that the umask still applies to the other permission bits (e.g. if umask is
	        _002, using _group will not remove read privileges from other (non-group) users). See _0xx for how to exactly specify the repository permissions.

	    _al (or _world or _everybody)
	        Same as _group, but make the repository readable by all users.

	    _0xx
	        _0xx is an octal number and each file will have mode _0xx.  _0xx will override users' umask(2) value (and not only loosen permissions as _group and _al does).  _0640 will create a repository which is group-readable, but not group-writable or accessible to others.  _0660
	        will create a repo that is readable and writable to the current user and group, but inaccessible to others.

	By default, the configuration flag receive.denyNonFastForwards is enabled in shared repositories, so that you cannot force a non fast-forwarding push into it.

	If you provide a _directory, the command is run inside it. If this directory does not exist, it will be created.
	*/

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
