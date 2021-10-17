# go-ipfs-healthcheck

**NOTE: This currently works with go-ipfs version 0.8. Other versions have not been tested.**

**TODO**

- [ ] Compatiblity with go-ipfs 9 and 10
- [ ] Use multiaddress for healthcheck endpoint

A plugin for [go-ipfs](https://github.com/ipfs/go-ipfs) that serves a healthcheck endpoint which returns the status of the IPFS node.

# Installation

This is a preloaded plugin built in-tree into go-ipfs when it's compiled.

```sh
git clone https://github.com/ipfs/go-ipfs

cd go-ipfs

# Pull in the plugin (you can specify a version other than "latest" if you'd like)
go get github.com/ceramicnetwork/go-ipfs-healthcheck/plugin@latest

# Add the plugin to the [preload list](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md#preloaded-plugins)
echo "\nhealthcheck github.com/ceramicnetwork/go-ipfs-healthcheck/plugin 0" >> plugin/loader/preload_list

# Download dependencies
go mod download

# Build go-ipfs with the plugin
make build
```

# Usage

Run IPFS and check its status.
```sh
./cmd/ipfs/ipfs daemon
curl -X GET http://localhost:8011
```

# Resources

[go-ipfs Plugins](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md)

# Maintainers

[@v-stickykeys](https://github.com/v-stickykeys)

# License

Fully open source and dual-licensed under MIT and Apache 2.
