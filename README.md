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


## How to use
* download a binary for your architecture from the releases page
* rename the binary to `kubectl-login`, `chmod +x kubectl-login` and place the binary somewhere in on your PATH
* create .kubectl-login.json in your home directory, or execute `kubectl login -u` to download the config file from this repo
* execute `kubectl login -l` to see a list of clusters/namespaces
* execute `kubectl login [alias]` - e.g. `kubectl login staging`
* If your token has expired, you will be prompted for credentials in a browser window. After login, your kube config will be updated with your new token and your context will be changed to the selected cluster/namespace.
