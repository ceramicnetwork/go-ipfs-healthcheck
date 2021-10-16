// Package plugin implements the interface for a go-ipfs daemon plugin.
package plugin

import (
	"fmt"
	"os"
	"reflect"

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
	maybePort := os.Getenv(portEnvVar)
	if maybePort != "" {
		port = maybePort
		return nil
	}

	val := reflect.ValueOf(env.Config).Elem()
	_port, ok := val.FieldByName("port").Interface().(string)
	if !ok {
		fmt.Println("Healthcheck plugin is defaulting to port " + port)
	} else {
		port = _port
	}

	return nil
}

func (*HealthcheckPlugin) Start(ipfs coreiface.CoreAPI) error {
	healthcheck.StartServer(port, ipfs)
	return nil
}
