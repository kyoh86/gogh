# CONTRIBUTING (internal/githubv4)

If you want to change internal github v4 client library in this package,
you *SHOULD* take the following steps.

1. Write query (`*.graphql`)
  - You can use [GraphiQL](https://github.com/skevy/graphiql-app) or the [API Explorer](https://docs.github.com/en/graphql/overview/explorer).
2. Re-generate client code with `go generate ./...`
  - You **MUST** set a PAT for GITHUB in the `GITHUB_TOKEN` envar.
  - The PAT to be generated does not need any scope
