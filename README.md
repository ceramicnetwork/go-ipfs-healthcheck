# go-ipfs healthcheck plugin

**WIP**

This tiny go server plugs in to the go-ipfs daemon to allow healthchecks that return whether the IPFS instance is online or not.

# Installation and Usage

```
go build server.go
./server
```

```
git clone https://github.com/ipfs/go-ipfs

cd go-ipfs

# Pull in the plugin (you can specify a version other than "latest" if you'd like)
go get github.com/ceramicnetwork/go-ipfs-healthcheck/plugin@latest

# Add the plugin to the [preload list](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md#preloaded-plugins)
echo "\nhealthcheck github.com/ceramicnetwork/go-ipfs-healthcheck/plugin 0" >> plugin/loader/preload_list

go mod download

# Enable the plugin by including it in the IPFS_PLUGINS variable when building IPFS
make build IPFS_PLUGINS="healthcheck"
```

Visit `http://localhost:8080`

# Resources

[go-ipfs Plugins](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md)

# Maintainers

[@v-stickykeys](https://github.com/v-stickykeys)

# License

Fully open source and dual-licensed under MIT and Apache 2.

