<p align="center">
    <a href="https://github.com/fragoulis/pites/actions/workflows/lint.yml" target="_blank" rel="noopener"><img src="https://github.com/fragoulis/pites/actions/workflows/lint.yml/badge.svg" alt="build" /></a>
    <a href="https://github.com/fragoulis/pites/actions/workflows/test.yml" target="_blank" rel="noopener"><img src="https://github.com/fragoulis/pites/actions/workflows/test.yml/badge.svg" alt="build" /></a>
</p>

# Create a user

```bash
go run main.go user create :name
```

# Development tools

```sh
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.1
```

```sh
go install golang.org/x/tools/cmd/goimports@latest
```

```sh
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

```sh
$ GO111MODULE=on go get -u -f github.com/DarthSim/hivemind
```
