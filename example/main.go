package main

import (
	"github.com/tlog/tlog"
)

func main() {
	log := tlog.NewLogger("YOUR_KEY", "YOUR_CHANEL_ID")

	log.Info("Simple info message, don't matter")
	log.Infof("Simple formatted info message, don't matter: %s", "format")
	log.Error("Error message")
	log.Errorf("Error formatted message: %s", "format")
}
