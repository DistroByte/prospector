# Go REST API

## Getting started

### Installing Go

Installing Go binaries

```bash
curl -L -O https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
```

Add Go bin to your PATH

```bash
echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> ~/.profile
source ~/.profile
```

### Building the project

```bash
go build -o bin/
```

### Running the project

```bash
./bin/prospector
```

### Or build and run

```bash
go run main.go
```

### Generate docs

Follow [this section](https://github.com/swaggo/swag#getting-started) of swaggo's repo.

```bash
swag fmt && swag init -g command/server/http.go
```

## CI/CD

Every push to a branch that is _not_ `master` will trigger a build, test and review pipeline.

The testing stage will run `go vet` and `go test` on the codebase. The test stage will fail if either of these steps fail.

The build stage will then build a docker image and push it to the registry at `git.dbyte.xyz`.

It will then trigger a review application to be deployed and viewed at the link from the PR, using gitlab's "environment" feature.
