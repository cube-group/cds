package plugins

import (
	"github.com/go-martini/martini"
	"github.com/spf13/viper"
	"alex/log"
)

func PluginLog(m *martini.ClassicMartini) *log.Logger {
	l := log.NewLogger(
		viper.GetString("app"),
		viper.GetString("path.log"),
		viper.GetString("locate"),
		viper.GetBool("debug"),
	)
	l.Debug("plugin", "PluginLog Loaded")
	m.Map(l)
	return l
}
