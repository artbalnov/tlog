# tlog
Log your events directly to Telegram

### Usage:
~~~golang
err := tlog.Init("YOUR_KEY", "YOUR_CHANEL_ID")
if err != nil {
    log.Fatal(err)
}

tlog.Info("Simple info message, don't matter")
tlog.Infof("Simple formatted info message, don't matter: %s", "format")
tlog.Error("Error message")
tlog.Errorf("Error formatted message: %s", "format")
~~~
