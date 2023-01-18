# Development

## Database

Migrations are managed using [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate), with the CLI tool.
If this is your first time using golang-migrate, check out the

Migrations are handled using the [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) package, which includes a command-line interface (CLI) tool. If you are new to golang-migrate, it is recommended to refer to the [Getting Started](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md) guide.

To install the golang-migrate CLI, follow the instructions in the [migrate CLI README](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md).

### Local development database

1. In order to set up local development, install MySQL version 5 on your machine. The default port for MySQL should be set to `3306`. An alternative option is to use a MySQL database within a Docker container. To do this, you can use the following command:

```sh
docker run \
   -p 3306:3306 \
   --env MYSQL_DATABASE=test \
   --env MYSQL_ALLOW_EMPTY_PASSWORD=yes \
   --name mysql \
   --detach \
   --rm \
   mysql:5 \
   --character-set-server=utf8mb4 \
   --collation-server=utf8mb4_unicode_ci
```

### Setting up for tests

If you have set up a MySQL database within a Docker container as outlined in step 1, you should now have a local development database available. However, when running `go test ./...`, it is important to note that the database tests will not run unless the MySQL service is currently running. Therefore, ensure that the MySQL service is running before executing the command to run the tests.

2. The tests rely on the DATABASE_DATASOURCE environment variable as the MySQL connection string. This variable should be set before running the tests.

```sh
export DATABASE_DATASOURCE="root@tcp(localhost:3306)/test?parseTime=true"
```

Make sure to replace `root:password` with your MySQL root user and its corresponding password, also replace test with your desired database name.

Please keep in mind that the above command is for Unix-based systems, for windows use `set` instead of `export`

3. In order for the database to be up-to-date with the current schema, it is necessary to run the database migrations. If the database has not been migrated, use the following command to migrate it:

```sh
migrate -source file:migrations -database "mysql://$DATABASE_DATASOURCE" up
```

4. Run the test

```
go test ./...
```

5. Start the server

```
go run cmd/codescan/main.go
```

### Creating a migration

To create migration:

```
./create_migration.sh <title>
```

This creates two empty files in `/migrations`:

```
{version}_{title}.up.sql
{version}_{title}.down.sql
```

The two migration files are used to migrate "up" to the specified version from the previous version, and to migrate "down" to the previous version. See [golang-migrate/migrate/MIGRATIONS.md](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) for details.

If you are migrating for the first time. You just runs the command, make sure the connection is matched with your local development setup.

```
migrate -source file:migrations -database "mysql://root@tcp(localhost:3306)/test?parseTime=true" up
```

## Environment variables

This program uses the following environment variables:

```
SERVER_ADDRESS=":8080"
MANAGER_NUM_WORKERS=1
DATABASE_DATASOURCE="root:1@tcp(localhost:3306)/test?parseTime=true"
DATABASE_MAX_CONNECTIONS=0
FETCH_DIR="tmp"
FETCH_DIR_PREFIX="codescan"
```
