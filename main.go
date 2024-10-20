package main

import (
	"github.com/kohmebot/manager/manager"
	"github.com/kohmebot/plugin"
)

func NewPlugin() plugin.Plugin {
	return manager.NewPlugin()
}
