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

### Build binaries
* `make pack`

### Create Github Release
* Upload the binaries on the release

## How to use
* rename the binary to kubectl-login and place somewhere in on your PATH
* create .kubectl-login.json in your home directory, or execute `kubectl login -u` to download the config file from this repo
* execute `kubectl login [alias]` - e.g. `kubectl login development`
* If your token has expired, you will be prompted for credentials. After login, your kube config will be updated with your new token and your context will be changed to the selected cluster/namespace.

