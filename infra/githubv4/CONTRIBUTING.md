# CONTRIBUTING (infra/githubv4)

If you want to change internal github v4 client library in this package,
you *SHOULD* take the following steps.

1. Write query (`*.graphql`)
2. Re-generate client code with `go generate ./...`

## To edit GrqphQL

- You can use [GraphiQL](https://github.com/skevy/graphiql-app) or the [API Explorer](https://docs.github.com/en/graphql/overview/explorer).

NOTE: GraphiQL does not work on linux now
https://github.com/skevy/graphiql-app/issues/175

## To generate client code

- You **MUST** set a PAT for GITHUB in the `GITHUB_TOKEN` envar.
- The PAT to be generated does not need any scope

If you manage it in password-manager (e.g. `1password-cli`)

```console
$ export GITHUB_TOKEN="$(op get item gqlgen --fields password)"
```
