package core

//go:generate go tool mockgen -source ./auth/authenticate_service.go       -destination ./auth_mock/gen_authenticate_service_mock.go       -package auth_mock
//go:generate go tool mockgen -source ./auth/token_service.go              -destination ./auth_mock/gen_token_service_mock.go              -package auth_mock
//go:generate go tool mockgen -source ./git/git_service.go                 -destination ./git_mock/gen_git_service_mock.go                 -package git_mock
//go:generate go tool mockgen -source ./hosting/hosting_service.go         -destination ./hosting_mock/gen_hosting_service_mock.go         -package hosting_mock
//go:generate go tool mockgen -source ./hosting/repository_format.go       -destination ./hosting_mock/gen_repository_format_mock.go       -package hosting_mock
//go:generate go tool mockgen -source ./repository/default_name_service.go -destination ./repository_mock/gen_default_name_service_mock.go -package repository_mock
//go:generate go tool mockgen -source ./repository/location_format.go      -destination ./repository_mock/gen_location_format_mock.go      -package repository_mock
//go:generate go tool mockgen -source ./repository/parser.go               -destination ./repository_mock/gen_parser_mock.go               -package repository_mock
//go:generate go tool mockgen -source ./store/store.go                     -destination ./store_mock/gen_store_mock.go                     -package store_mock
//go:generate go tool mockgen -source ./workspace/finder_service.go        -destination ./workspace_mock/gen_finder_service_mock.go        -package workspace_mock
//go:generate go tool mockgen -source ./workspace/layout_service.go        -destination ./workspace_mock/gen_layout_service_mock.go        -package workspace_mock
//go:generate go tool mockgen -source ./workspace/workspace_service.go     -destination ./workspace_mock/gen_workspace_service_mock.go     -package workspace_mock
