package command_test

func ExampleConfigGetAll() {
	// UNDONE: mock and example
	// if err := command.ConfigGetAll(&config.Config{
	// 	GitHub: config.GitHubConfig{
	// 		Token: "tokenx1",
	// 		Host:  "hostx1",
	// 		User:  "kyoh86",
	// 	},
	// 	VRoot: []string{"/foo", "/bar"},
	// }); err != nil {
	// 	panic(err)
	// }

	// Unordered output:
	// root: /foo:/bar
	// github.host: hostx1
	// github.user: kyoh86
	// github.token: *****
}
