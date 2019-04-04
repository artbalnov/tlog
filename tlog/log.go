package tlog

import (
	"fmt"
	"log"
	"net/http"
)

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

const (
	requestTemplate = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=HTML&text=%s"
)

const (
	errorHeaderTemplate = "<b>ERROR:</b>%0A"
	fyiHeaderTemplate   = "<i>FYI:</i>%0A"

	errorStyleTemplate = "code"
	fyiStyleTemplate   = "pre"
)

func NewLogger(apiKey, chanelID string) Logger {
	return &tlog{
		apiKey:   apiKey,
		chanelID: chanelID,
	}
}

type tlog struct {
	apiKey   string
	chanelID string
}

func (rcv *tlog) Info(args ...interface{}) {
	log.Print(args...)
	rcv.sendMessageRequest(formatMessage(fyiHeaderTemplate, fyiStyleTemplate, args...))
}

func (rcv *tlog) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
	rcv.sendMessageRequest(formatfMessage(fyiHeaderTemplate, fyiStyleTemplate, format, args...))
}

func (rcv *tlog) Error(args ...interface{}) {
	log.Print(args...)
	rcv.sendMessageRequest(formatMessage(errorHeaderTemplate, errorStyleTemplate, args...))
}

func (rcv *tlog) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	rcv.sendMessageRequest(formatfMessage(errorHeaderTemplate, errorStyleTemplate, format, args...))
}

func formatMessage(header, style string, args ...interface{}) string {
	return fmt.Sprintf("%s<%s>%s</%s>", header, style, fmt.Sprint(args...), style)
}

func formatfMessage(header, style, format string, args ...interface{}) string {
	return fmt.Sprintf("%s<%s>%s</%s>", header, style, fmt.Sprintf(format, args...), style)
}

func (rcv *tlog) sendMessageRequest(msg string) {
	log.Print(msg)

	resp, err := http.Get(fmt.Sprintf(requestTemplate, rcv.apiKey, rcv.chanelID, msg))
	if err != nil {
		log.Printf("[tlog] can't send message to chanel: %+v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("[tlog] can't send message to chanel: %s", resp.Status)
	}

	defer func() {
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("[tlog] can't close response body: %+v", err)
			}
		}
	}()
}
