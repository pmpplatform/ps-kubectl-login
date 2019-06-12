# kubectl-login

A kubectl plugin for dex

## Config file

`$HOME/.kubectl-login.json`

```
{
  "clstr01": {
    "dex-url": "https://clstr01-login.privatedns.zone/",
    "aliases": ["development", "staging", "qa"]
  },
  "clstr02": {
    "dex-url": "https://clstr02-login.privatedns.zone/",
    "aliases": ["testing"]
  }
}
```
'aliases' should be a list of unique namespaces on the cluster

# Releases

## Install dep

### Mac
`brew install dep` / `brew upgrade dep`

### Other Platforms
`curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`

## Run dep ensure
`dep ensure`

## Release
### Build binaries
* `GOOS=linux GOARCH=amd64 go build -o bin/kubectl-login-linux .`
* `GOOS=darwin GOARCH=amd64 go build -o bin/kubectl-login-darwin .`

### Create Github Release
* Upload the binaries on the release

## How to use
* rename the binary to kubectl-login and place somewhere in on your PATH
* create .kubectl-login.json in your home directory
* execute `kubectl login [alias]` - e.g. `kubectl login development`
* If your token has expired, you will be prompted for credentials. After login, your kube config will be updated with your new token and your context will be changed to the selected cluster/namespace.

