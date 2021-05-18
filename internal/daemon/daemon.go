package daemon

import (
	"path/filepath"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
)

// Creates a new go-daemon context
//  based on the program name
//  (which is taken from viper)
func GetContext() *daemon.Context {
	// by default this is ~/.heiko/<name>
	workDir := filepath.Join(
		viper.GetString("dataLocation"),
		viper.GetString("name"),
	)

	context := &daemon.Context{
		PidFileName: filepath.Join(
			workDir,
			"daemon.pid",
		),
		PidFilePerm: 0644,

		// TODO: log file apparently captures only stderr
		LogFileName: filepath.Join(
			workDir,
			"daemon.log",
		),
		LogFilePerm: 0644,

		// WorkDir is set to the directory running heiko
		//   this is needed to ensure that the config
		//   file is read correctly
		WorkDir:     ".",

		Umask:       022,
		// Args:        nil,
		// Env:         nil,
	}

	return context
}
