## Template Project

This is a template project which can\should be used as a blueprint for new projects

## Config

The application depends on config values that are located in the environment variables.
It means that before execution of any CLI command, the corresponding env. variables should be defined.

In Linux and macOS you can define and set a value of an env. variable by the `export` command

Example. Let's set the value of the `DB_MIGRATIONS_PATH` env. variable

```
export DB_MIGRATIONS_PATH=/Users/levboiko/go_projects/boostylabs/project_template/database/migrations
```

Another handy way for setting of env. variables, is to "import" a corresponding .env file `export $(grep -v '^#' [.env file location] | xargs)`

Example:

```
export $(grep -v '^#' ./.env | xargs)
```


Sample of configuration is in `.env.dist` file

## Console commands

### Main app | cmd/template_project

The CLI command to launch the application

```bash
export $(grep -v '^#' ./.env | xargs)

go run cmd/template_project/main.go run
```

### Migrations | cmd/database 

#### Create a new migration

```
go run cmd/database/main.go create-migration [migration_name]
```

Example:

```bash
export $(grep -v '^#' ./.env | xargs)

go run cmd/database/main.go create-migration init
```

Sample output:

```
new file: 000001_init.up.sql
new file: 000001_init.down.sql
```

#### Apply migrations

```
go run cmd/database/main.go migrate up
```

Example:

```bash
export $(grep -v '^#' ./.env | xargs)

go run cmd/database/main.go migrate up
```

In case of successful execution, there will be an empty output

#### Rollback migrations

```
go run cmd/database/main.go migrate down
```

Example:

```bash
export $(grep -v '^#' ./.env | xargs)

go run cmd/database/main.go migrate down
```

In case of successful execution, there will be an empty output

## Test/dev environment setup

1. Create a `.env` and set all params
2. Run the `docker-compose` command as it shown below

```bash
docker-compose --env-file=.env up
```

or

```bash
docker-compose --env-file=.env up -d
```

3. Apply migrations

```bash
export $(grep -v '^#' ./.env | xargs)

go run cmd/database/main.go migrate up
```

4. Lunch a required application

```bash
export $(grep -v '^#' ./.env | xargs)

go run cmd/template_project/main.go run
```

5. Visit the `http://localhost:3030/` url to open Grafana UI
   1. Credentials are **admin\admin**
   2. Pick a dashboard **Go Metrics**


## Run tests

In order to run test we need "import" env. variables and run `go test` CLI command

For example, let's run tests for the 'dummy' package

```bash
export $(grep -v '^#' ./.env | xargs)

go test ./dummy
```
