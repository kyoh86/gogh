# Define variables for common commands
MOCKGEN = go tool mockgen

# Main target to generate all mocks
.PHONY: mocks
mocks: \
	auth_mock/gen_authenticate_service_mock.go \
	auth_mock/gen_token_service_mock.go \
	git_mock/gen_git_service_mock.go \
	hosting_mock/gen_hosting_service_mock.go \
	hosting_mock/gen_repository_format_mock.go \
	repository_mock/gen_default_name_service_mock.go \
	repository_mock/gen_location_format_mock.go \
	repository_mock/gen_parser_mock.go \
	store_mock/gen_store_mock.go \
	workspace_mock/gen_finder_service_mock.go \
	workspace_mock/gen_layout_service_mock.go \
	workspace_mock/gen_workspace_service_mock.go \
	overlay_mock/gen_service_mock.go \
	overlay_mock/gen_overlay_mock.go \
	script_mock/gen_service_mock.go \
	script_mock/gen_script_mock.go \
	hook_mock/gen_service_mock.go \
	hook_mock/gen_hook_mock.go \
	extra_mock/gen_service_mock.go

# Rules for each mock file
auth_mock/gen_authenticate_service_mock.go: auth/authenticate_service.go
	$(MOCKGEN) -source $< -destination $@ -package auth_mock

auth_mock/gen_token_service_mock.go: auth/token_service.go
	$(MOCKGEN) -source $< -destination $@ -package auth_mock

git_mock/gen_git_service_mock.go: git/git_service.go
	$(MOCKGEN) -source $< -destination $@ -package git_mock

hosting_mock/gen_hosting_service_mock.go: hosting/hosting_service.go
	$(MOCKGEN) -source $< -destination $@ -package hosting_mock

hosting_mock/gen_repository_format_mock.go: hosting/repository_format.go
	$(MOCKGEN) -source $< -destination $@ -package hosting_mock

repository_mock/gen_default_name_service_mock.go: repository/default_name_service.go
	$(MOCKGEN) -source $< -destination $@ -package repository_mock

repository_mock/gen_location_format_mock.go: repository/location_format.go
	$(MOCKGEN) -source $< -destination $@ -package repository_mock

repository_mock/gen_parser_mock.go: repository/parser.go
	$(MOCKGEN) -source $< -destination $@ -package repository_mock

store_mock/gen_store_mock.go: store/store.go
	$(MOCKGEN) -source $< -destination $@ -package store_mock

workspace_mock/gen_finder_service_mock.go: workspace/finder_service.go
	$(MOCKGEN) -source $< -destination $@ -package workspace_mock

workspace_mock/gen_layout_service_mock.go: workspace/layout_service.go
	$(MOCKGEN) -source $< -destination $@ -package workspace_mock

workspace_mock/gen_workspace_service_mock.go: workspace/workspace_service.go
	$(MOCKGEN) -source $< -destination $@ -package workspace_mock

overlay_mock/gen_service_mock.go: overlay/service.go
	$(MOCKGEN) -source $< -destination $@ -package overlay_mock

overlay_mock/gen_overlay_mock.go: overlay/overlay.go
	$(MOCKGEN) -source $< -destination $@ -package overlay_mock

script_mock/gen_service_mock.go: script/service.go
	$(MOCKGEN) -source $< -destination $@ -package script_mock

script_mock/gen_script_mock.go: script/script.go
	$(MOCKGEN) -source $< -destination $@ -package script_mock

hook_mock/gen_service_mock.go: hook/service.go
	$(MOCKGEN) -source $< -destination $@ -package hook_mock

hook_mock/gen_hook_mock.go: hook/hook.go
	$(MOCKGEN) -source $< -destination $@ -package hook_mock

extra_mock/gen_service_mock.go: extra/extra_service.go
	$(MOCKGEN) -source $< -destination $@ -package extra_mock

# Clean target to remove all generated mocks
.PHONY: clean
clean:
	rm -f ./*_mock/gen_*.go
