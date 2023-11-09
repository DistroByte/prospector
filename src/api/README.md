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
go run .
```

### Generate docs

Follow [this section](https://github.com/swaggo/swag#getting-started) of swaggo's repo.

```bash
swag fmt && swag init
```
