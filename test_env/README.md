# Test environment

## Local + Docker

The `go-proxy` service should be run locally, the `mysql` containing `primary` and `replica` should be run via `Docker`.

### How to run `Docker`

Using **Powershell**:

```shell
docker run --name main-mysql -v $pwd/main_cnf/main.cnf
```
