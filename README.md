# tlog
Log your events directly to Telegram

### Usage:
~~~golang
log := tlog.NewLogger("YOUR_KEY", "YOUR_CHANEL_ID")

log.Info("Simple info message, don't matter")
log.Infof("Simple formatted info message, don't matter: %s", "format")
log.Error("Error message")
log.Errorf("Error formatted message: %s", "format")
~~~
