package daemon

import (
	"path/filepath"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
)

// Creates a new go-daemon context
//  based on the program name
//  (which is taken from viper)
func GetContext() *daemon.Context {
	SetHandlers(true)

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
		WorkDir: ".",

		Umask: 022,
		// Args:        nil,
		// Env:         nil,
	}

	return context
}

// sets up signal handlers as in the example
// given in https://github.com/sevlyar/go-daemon/blob/3fdf7dcbb9d92331eaa91e649c4755da76c64382/examples/cmd/gd-signal-handling/signal-handling.go#L21
func SetHandlers(quit bool) {
	daemon.AddCommand(daemon.BoolFlag(&quit), syscall.SIGINT, stopHandler)
}
