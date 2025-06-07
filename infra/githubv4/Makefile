generated.go: schema.graphql genqlient.graphql genqlient.yaml graphql.config.toml
	go tool genqlient

# Smart download function that only updates if content changed
define smart_download
	@echo "Downloading $(2)..."
	@mkdir -p $$(dirname $(2))
	@curl -sSL $(1) -o $(2).tmp
	@if [ ! -f $(2) ] || ! cmp -s $(2) $(2).tmp; then \
		mv $(2).tmp $(2); \
		echo "Updated $(2)"; \
	else \
		rm $(2).tmp; \
		echo "No changes to $(2)"; \
	fi
endef

schema.graphql:
	$(call smart_download,https://docs.github.com/public/fpt/schema.docs.graphql,./schema.graphql)

.PHONY: clean
clean:
	rm -f ./generated.go
	rm -f ./schema.graphql
