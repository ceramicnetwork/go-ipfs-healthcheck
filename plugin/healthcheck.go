// Package plugin implements the interface for a go-ipfs daemon plugin.
package plugin

import (
	"os"

	healthcheck "github.com/ceramicnetwork/go-ipfs-healthcheck"
	"github.com/ipfs/go-ipfs/plugin"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

var Plugins = []plugin.PluginDaemon{
	&HealthcheckPlugin{},
}

type HealthcheckPlugin struct{}

var port = "8011"
var portEnvVar = "HEALTHCHECK_PORT"

// Name returns the plugin's name, satisfying the plugin.Plugin interface.
func (*HealthcheckPlugin) Name() string {
	return "healthcheck"
}

// Version returns the plugin's version, satisfying the plugin.Plugin interface.
func (*HealthcheckPlugin) Version() string {
	return "0.0.1"
}

// Init initializes plugin, satisfying the plugin.Plugin interface.
func (*HealthcheckPlugin) Init(env *plugin.Environment) error {
	_port := os.Getenv(portEnvVar)
	if _port != "" {
		port = _port
		return nil
	}

	cfg, ok := env.Config.(map[string]interface{})
	if ok {
		_port, ok := cfg["port"].(string)
		if ok {
			port = _port
			return nil
		}
	}

	return nil
}

func (*HealthcheckPlugin) Start(ipfs coreiface.CoreAPI) error {
	go func() {
		healthcheck.StartServer(port, ipfs)
	}()
	return nil
}
