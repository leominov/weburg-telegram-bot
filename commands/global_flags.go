package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/oleiade/reflections"
)

type GlobalFlagsConfig struct {
	Debug   bool
	NoColor bool
}

var DebugFlag = cli.BoolFlag{
	Name:   "debug, d",
	Usage:  "Enable debug mode",
	EnvVar: "WEBURG_BOT_DEBUG",
}

var NoColorFlag = cli.BoolFlag{
	Name:   "no-color, nc",
	Usage:  "Don't show colors in logging",
	EnvVar: "WEBURG_BOT_NO_COLOR",
}

func HandleGlobalFlags(cfg interface{}) {
	debug, err := reflections.GetField(cfg, "Debug")
	if debug == true && err == nil {
		logrus.SetLevel(logrus.DebugLevel)
	}

	noColor, err := reflections.GetField(cfg, "NoColor")
	if noColor == true && err == nil {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
