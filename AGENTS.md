# Agent Instructions

## Project Goals
This repository should provide a Go-based CLI (`node-creator`) that scaffolds a monorepo-style Node.js backend project with configurable modules and infrastructure.

## Required Features
- Initialize a Go module and implement the CLI entry point (e.g., using Cobra) under `cmd/node-creator`.
- CLI command to generate Node.js project templates with selectable modules: routes, controllers, services, repositories, middleware, helpers.
- Scaffold a shared `packages` (or similar) directory that contains reusable monorepo resources: database clients, libraries, and Redis caching utilities.
- Support environment-based configuration files for `local`, `dev`, `stg`, and `prod` (`.env` variants keyed by `NODE_ENV`).
- Generate database configuration using an ORM with choice between MySQL, PostgreSQL, or MongoDB.
- Provide Redis caching support with customizable namespaced keys (e.g., `namespace:group:key`) and configurable TTL, without per-route cache strategy toggles (to be added later).
- Offer options to scaffold both `apps/api` and `apps/worker` folders that consume the shared monorepo resources.
- Include logging setup that writes HTTP access logs to files and separates `log`, `error`, and `warn` outputs with daily file rotation.
- Provide authentication scaffolding that can be configured to use either JWT or PASETO.
- Include Prettier and ESLint configuration in generated projects by default.
- Provide template files for all layers (routes, controllers, services, repositories, middleware, helpers, DB config, caching helpers, logging, auth, env files).
- Implement option parsing for module selection, database choice, shared resource setup, authentication method, and other configuration toggles.
- Render templates into a target directory, wiring monorepo structure, DB, caching, logging, and authentication selections accordingly.
- Include unit tests covering template rendering for key option combinations.
- Document installation, usage, and configuration in `README.md` (including CLI commands and troubleshooting).

