/*
Package remote contains commands those get informations from GitHub.

* Commands will be organized like below that follows GitHub API v3 specification.
	See: https://developer.github.com/v3/

| Function                | Subcommand                    |
|-------------------------|-------------------------------|
| List repos              | repo                          |
| Search repos            | search repo                   |
| Search commits          | search commit                 |
| Search code             | search code                   |
| Search pull requests    | search issue                  |
| Search issues           | search pr, pull, pull-request |
| Search users            | search user                   |
| Search topics           | search topic                  |
| Search labels           | search label                  |

All of plural form of each subcommand can be used. (i.e. "repos")

* This should not manage pull-request and issues.
We should request features and enhancements to `hub` (https://github.com/github/hub).

* This should not manage gists.
	We should request features and enhancements to `gist` (http://defunkt.io/gist/).
*/
package remote
