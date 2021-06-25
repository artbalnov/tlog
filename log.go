package tlog

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

const baseURLTemplate = "https://api.telegram.org/bot%s/sendMessage"

const (
	chatIDKey    = "chat_id"
	parseModeKey = "parse_mode"
	textKey      = "text"
)

const (
	errorHeaderTemplate = "<b>%s</b>"
	infoHeaderTemplate  = "<i>%s</i>"

	codeStyleTemplate = "<code>%s</code>"
)

var logger *tlog

type tlog struct {
	botKey            string
	chatID            string
	baseURL           string
	errHeaderContent  string
	infoHeaderContent string
	debug             bool
}

func Init(botKey, chatID string) error {
	u, err := url.Parse(fmt.Sprintf(baseURLTemplate, botKey))
	if err != nil {
		return fmt.Errorf("can't init logger: %v", err)
	}

	logger = &tlog{
		botKey:            botKey,
		chatID:            chatID,
		baseURL:           u.String(),
		errHeaderContent:  "ERROR:",
		infoHeaderContent: "FYI:",
	}

	return nil
}

func SetErrorHeader(header string) {
	logger.errHeaderContent = header
}

func SetInfoHeader(header string) {
	logger.infoHeaderContent = header
}

func SetDebug(enabled bool) {
	logger.debug = enabled
}

func Info(args ...interface{}) {
	logger.sendMessageRequestAsync(formatMessage(logger.infoHeader(), "", args...))
}

func Infof(format string, args ...interface{}) {
	logger.sendMessageRequestAsync(formatfMessage(logger.infoHeader(), "", format, args...))
}

func Error(args ...interface{}) {
	logger.sendMessageRequestAsync(formatMessage(logger.errorHeader(), codeStyleTemplate, args...))
}

func Errorf(format string, args ...interface{}) {
	logger.sendMessageRequestAsync(formatfMessage(logger.errorHeader(), codeStyleTemplate, format, args...))
}

func formatMessage(header, style string, args ...interface{}) string {
	var msg = fmt.Sprint(args...)
	if style != "" {
		msg = fmt.Sprintf(style, msg)
	}
	return fmt.Sprintf("%s\n%s", header, msg)
}

func formatfMessage(header, style, format string, args ...interface{}) string {
	var msg = fmt.Sprintf(format, args...)
	if style != "" {
		msg = fmt.Sprintf(style, msg)
	}
	return fmt.Sprintf("%s\n%s", header, msg)
}

func (t *tlog) errorHeader() string {
	return fmt.Sprintf(errorHeaderTemplate, t.errHeaderContent)
}

func (t *tlog) infoHeader() string {
	return fmt.Sprintf(infoHeaderTemplate, t.infoHeaderContent)
}

func (t *tlog) sendMessageRequestAsync(msg string) {
	t.printDebugMessage("[tlog] message: %s", msg)

	go t.sendMessageRequest(msg)
}

func (t *tlog) sendMessageRequest(msg string) {
	u, err := url.Parse(t.baseURL)
	if err != nil {
		t.printDebugMessage("[tlog] can't parse base URL: %v", err)
		return
	}

	params := url.Values{}
	params.Add(chatIDKey, t.chatID)
	params.Add(parseModeKey, "HTML")
	params.Add(textKey, msg)

	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		if t.debug {
			t.printDebugMessage("[tlog] can't send message to chanel: %+v", err)
		}
		return
	}
	if resp.StatusCode != http.StatusOK {
		if t.debug {
			t.printDebugMessage("[tlog] can't send message to chanel: %s", resp.Status)
		}
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.printDebugMessage("[tlog] can't close response body: %+v", err)
		}
	}()
}

func (t *tlog) printDebugMessage(format string, v ...interface{}) {
	if t.debug {
		log.Printf(format, v)
	}
}
