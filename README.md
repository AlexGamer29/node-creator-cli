# node-creator-cli

`node-creator` is a Go-based CLI that scaffolds a monorepo-style Node.js backend with opinionated defaults for API and worker apps, shared infrastructure packages, and environment-aware configuration.

## Features

- Selectable modules (routes, controllers, services, repositories, middleware, helpers) for the API workspace.
- Shared `packages/shared` library with database, Redis cache, logging, authentication, and helper utilities.
- Database scaffolding using TypeORM for MySQL/PostgreSQL or the MongoDB driver.
- Redis cache utilities with customizable namespace/group/TTL key prefixing.
- Optional API and worker applications that both consume the shared libraries.
- Environment presets for `local`, `dev`, `stg`, and `prod`.
- HTTP access logging to files plus JSON structured logs with daily rotation for `info`, `warn`, and `error`.
- Authentication templates supporting either JWT or PASETO tokens.
- Prettier and ESLint configuration included by default.

## Installation

```bash
go install github.com/example/node-creator-cli/cmd/node-creator@latest
```

The command installs a `node-creator` binary in your `$GOBIN` (default `$GOPATH/bin`).

## Usage

```bash
node-creator generate <project-name> [flags]
```

### Key Flags

| Flag | Description | Default |
| ---- | ----------- | ------- |
| `--modules` | Comma-separated list of API modules (`routes,controllers,services,repositories,middleware,helpers`). | all modules |
| `--database` | ORM target: `mysql`, `postgres`, or `mongodb`. | `mysql` |
| `--auth` | Authentication provider: `jwt` or `paseto`. | `jwt` |
| `--api` | Generate the API workspace. | `true` |
| `--worker` | Generate the worker workspace. | `true` |
| `--cache-namespace` | Redis namespace prefix. | `app` |
| `--cache-group` | Redis group segment. | `default` |
| `--cache-ttl` | Default Redis TTL (seconds). | `300` |
| `--output` | Destination directory (defaults to `<project-name>`). | *(project name)* |

### Example

```bash
node-creator generate awesome-app \
  --modules routes,controllers,services \
  --database postgres \
  --auth paseto \
  --cache-namespace awesome \
  --cache-group api \
  --cache-ttl 900
```

The above command creates the following top-level structure:

```
awesome-app/
  apps/
    api/
    worker/
  packages/
    shared/
  config/env/
```

## Generated Workspaces

- **apps/api** – Express-based HTTP service pre-wired with selected modules, shared logging, and Redis helpers.
- **apps/worker** – BullMQ-based worker ready to share Redis connections and loggers with the API.
- **packages/shared** – Contains:
  - `database`: ORM configuration tailored to the selected database.
  - `cache/redis`: Redis helper with namespace/group/TTL helpers.
  - `libs/logger` & `libs/httpLogger`: Winston and Morgan-based loggers with daily rotating file output.
  - `auth`: Wrapper exposing JWT or PASETO utilities depending on CLI choice.
  - `helpers`: Miscellaneous helper utilities.

## Environment Configuration

Four `.env` files are generated under `config/env` (`.env.local`, `.env.dev`, `.env.stg`, `.env.prod`). Each file includes the Redis namespace/group/TTL selections plus placeholders for database and secret values. Override `ENV_FILE` at runtime to select a different environment preset.

## Testing the CLI

Run the generator tests to verify template rendering for different configuration combinations:

```bash
go test ./...
```

## Contributing

1. Fork the repository and create a feature branch.
2. Add or update templates under `templates/` as needed.
3. Add unit tests in `internal/generator` to cover new scenarios.
4. Run `go fmt ./...` and `go test ./...` before submitting a PR.
