# tlog
Log your events directly to Telegram

### Usage:
~~~golang
t, err := tlog.NewLogger("YOUR_KEY", "YOUR_CHANEL_ID")
if err != nil {
    log.Fatal(err)
}

t.Info("Simple info message, don't matter")
t.Infof("Simple formatted info message, don't matter: %s", "format")
t.Error("Error message")
t.Errorf("Error formatted message: %s", "format")
~~~
