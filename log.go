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
	errorHeaderTemplate = "<b>ERROR:</b>"
	fyiHeaderTemplate   = "<i>FYI:</i>"

	errorStyleTemplate = "code"
	fyiStyleTemplate   = "pre"
)

type tlog struct {
	botKey  string
	chatID  string
	baseURL string
}

func NewLogger(botKey, chatID string) (Logger, error) {
	u, err := url.Parse(fmt.Sprintf(baseURLTemplate, botKey))
	if err != nil {
		return nil, fmt.Errorf("can't init logger: %v", err)
	}

	return &tlog{
		botKey:  botKey,
		chatID:  chatID,
		baseURL: u.String(),
	}, nil
}

func (rcv *tlog) Info(args ...interface{}) {
	log.Print(args...)
	rcv.sendMessageRequestAsync(formatMessage(fyiHeaderTemplate, fyiStyleTemplate, args...))
}

func (rcv *tlog) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
	rcv.sendMessageRequestAsync(formatfMessage(fyiHeaderTemplate, fyiStyleTemplate, format, args...))
}

func (rcv *tlog) Error(args ...interface{}) {
	log.Print(args...)
	rcv.sendMessageRequestAsync(formatMessage(errorHeaderTemplate, errorStyleTemplate, args...))
}

func (rcv *tlog) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	rcv.sendMessageRequestAsync(formatfMessage(errorHeaderTemplate, errorStyleTemplate, format, args...))
}

func formatMessage(header, style string, args ...interface{}) string {
	return fmt.Sprintf("%s\n<%s>%s</%s>", header, style, fmt.Sprint(args...), style)
}

func formatfMessage(header, style, format string, args ...interface{}) string {
	return fmt.Sprintf("%s\n<%s>%s</%s>", header, style, fmt.Sprintf(format, args...), style)
}

func (rcv *tlog) sendMessageRequestAsync(msg string) {
	go rcv.sendMessageRequest(msg)
}

func (rcv *tlog) sendMessageRequest(msg string) {
	u, err := url.Parse(rcv.baseURL)
	if err != nil {
		log.Printf("[tlog] can't parse base URL: %v", err)
		return
	}

	params := url.Values{}
	params.Add(chatIDKey, rcv.chatID)
	params.Add(parseModeKey, "HTML")
	params.Add(textKey, msg)

	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("[tlog] can't send message to chanel: %+v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("[tlog] can't send message to chanel: %s", resp.Status)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("[tlog] can't close response body: %+v", err)
		}
	}()
}
