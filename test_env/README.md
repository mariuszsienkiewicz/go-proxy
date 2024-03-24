# Test environment

## Local + Docker

The `go-proxy` service should be run locally, the `mysql` containing `primary` and `replica` should be run via `Docker`.

### How to run `Docker`

You can set up the test environment via `docker-compose`.

1. Run from command line: `docker compose up -d`
2. Run `docker exec main_db bash /setup/setup.sh` and copy File and Position from Master Status
3. Run `docker exec replica_db bash /setup/setup.sh <File> <Position>`, replace `<File>` with File and `<Position>` with Position from the command above

Test environment should be now up and running.