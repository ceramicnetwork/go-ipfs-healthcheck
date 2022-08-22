// Package plugin implements the interface for a go-ipfs daemon plugin.
package plugin

import (
	"os"

	healthcheck "github.com/ukstv/go-ipfs-healthcheck"
	"github.com/ipfs/kubo/plugin"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

var Plugins = []plugin.PluginDaemon{
	&HealthcheckPlugin{},
}

type HealthcheckPlugin struct{}

const portEnvVar = "HEALTHCHECK_PORT"

var port = "8011"

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
	envPort := os.Getenv(portEnvVar)
	if envPort != "" {
		port = envPort
		return nil
	}

	if cfg, ok := env.Config.(map[string]interface{}); ok {
		if cfgPort, ok := cfg["port"].(string); ok {
			port = cfgPort
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
