# buf-playground

A playground for experimenting with [Buf](https://buf.build) tooling, including custom lint plugins and protobuf schema definitions.

## Repository structure

```
.
├── plugin/         # Custom buf lint plugins (Go)
│   ├── cmd/
│   │   ├── api-lint-plugin/        # Plugin for API-level lint rules
│   │   └── request-lint-plugin/    # Plugin for request message lint rules
│   ├── internal/
│   │   └── request/                # Lint rule implementations
│   └── dist/                       # Compiled plugin binaries
└── schema/         # Protobuf schema workspace
    └── protos/
        ├── user/v1/                # User API
        │   ├── enums.proto
        │   ├── models.proto
        │   ├── refs.proto
        │   └── service_user.proto
        └── team/v1/                # Team API
            ├── enums.proto
            ├── models.proto
            ├── refs.proto
            └── service_team.proto
```

## Requirements

This project uses [mise](https://mise.jdx.dev) to manage tools. Install it then run:

```bash
mise install
```

This installs Go and Buf, and adds the compiled plugin binaries in `plugin/dist/` to your `PATH`.

## Plugins

Custom buf lint plugins are written in Go using the [bufplugin SDK](https://github.com/bufbuild/bufplugin-go).

### Build

```bash
cd plugin
mise run build
```

Binaries are output to `plugin/dist/` and automatically available on `PATH` via mise.

### Available plugins

| Plugin | Description |
|---|---|
| `request-lint-plugin` | Enforces lint rules on request messages |
| `api-lint-plugin` | Enforces API-level lint rules |

### Lint rules

#### `request-lint-plugin`

| Rule ID | Default | Description |
|---|---|---|
| `REPEATED_FIELD_VALIDATION` | yes | Repeated fields in request messages must have a `max_items` constraint to prevent unbounded input attacks |

## Schema

The `schema` workspace defines protobuf APIs for two domains:

- **user.v1** — `UserService` with Create, Get, Update, Delete, List operations
- **team.v1** — `TeamService` with Create, Get, Update, Delete, List operations, plus `AddTeamMember`, `AddTeamMembers`, and `RemoveTeamMember`

Cross-domain references use ref types (e.g. `user.v1.UserRef`) rather than plain string IDs.

### Lint

```bash
cd schema
mise run lint
```
