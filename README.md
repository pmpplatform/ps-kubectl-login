# kubectl-login

A kubectl plugin for dex

## Config file

`$HOME/.kubectl-login.json`

```
{
  "clstr01": {
    "client-id": "clstr01-dex",
    "client-secret": "LKJfsdfhlKLhfsd0jhsdf",
    "issuer": "https://clstr01-dex.yourdomain/dex",
    "redirect-uri": "http://127.0.0.1:5555/callback",
    "aliases": ["development", "staging", "qa"]
  },
  "clstr02": {
    "client-id": "clstr02-dex",
    "client-secret": "LKhkljhdf0JHh3hfolUlkj",
    "issuer": "https://clstr01-dex.yourdomain.com/dex",
    "redirect-uri": "http://127.0.0.1:5555/callback",
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
* `GOOS=linux GOARCH=amd64 go build -o kubectl-login-linux .`
* `GOOS=darwin GOARCH=amd64 go build -o kubectl-login-darwin .`
* `GOOS=windows GOARCH=amd64 go build -o kubectl-login-windows.exe .`

### Create Github Release
* Upload the binaries on the release

## How to use
* rename the binary to kubectl-login and place somewhere in on your PATH
* create .kubectl-login.json in your home directory
* execute `kubectl login [alias]` - e.g. `kubectl login development`
* If your token has expired, the dex login page will be opened in your browser. After login, your kube config will be updated with your new token and your context will be changed to the selected cluster/namespace.

