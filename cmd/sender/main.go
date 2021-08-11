package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var (
	tmplFilePath, tmplString, tmplBase64String *string
	envVars                                    map[string]string
)

func main() {
	tgBotToken := flag.String("tg-bot-token", "", "Telegram bot token value")
	tgChatID := flag.String("tg-chat-id", "", "Telegram chat id value")
	envVarPrefix := flag.String("env-var-prefix", "TG_MSG_", "Environment variables name prefix")
	tmplFilePath = flag.String("template-file-path", "", "Go template file path")
	tmplString = flag.String("template-string", "", "Go template string")
	tmplBase64String = flag.String("template-base64-string", "", "Go template base64 encoded string")
	flag.Parse()

	if *tgBotToken == "" {
		log.Fatalf("bot token must be given")
	}

	if *tgChatID == "" {
		log.Fatalf("chat id must be given")
	}

	t, err := getTemplateValue()
	if err != nil {
		log.Fatal(err)
	}

	evs := getEnvVars(*envVarPrefix)

	msg, err := getParsedMessage(t, evs)
	if err != nil {
		log.Fatal(err)
	}

	err = sendMessage(*tgBotToken, *tgChatID, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTemplateValue() (string, error) {
	var t string

	switch {
	case *tmplFilePath != "":
		f, err := os.Open(*tmplFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to open template file %s: %s", *tmplFilePath, err)
		}
		tb, err := ioutil.ReadAll(f)
		if err != nil {
			return "", fmt.Errorf("failed to read template file %s: %s", *tmplFilePath, err)
		}
		t = string(tb)
	case *tmplString != "":
		t = *tmplString
	case *tmplBase64String != "":
		tb, err := base64.StdEncoding.DecodeString(*tmplBase64String)
		if err != nil {
			return "", fmt.Errorf("failed decode given base64 ecoded string template: %s", err)
		}
		t = string(tb)
	default:
		return "", errors.New("Message template should be given")
	}

	return t, nil
}

func getEnvVars(prefix string) map[string]string {
	evs := make(map[string]string)

	for _, vl := range os.Environ() {
		vs := strings.Split(vl, "=")
		k, v := vs[0], vs[1]
		if strings.HasPrefix(k, prefix) {
			evs[k] = v
		}
	}

	return evs
}

func getParsedMessage(t string, evs map[string]string) (string, error) {
	b := new(bytes.Buffer)

	tmpl, err := template.New("message").Funcs(tmplFuncs).Parse(t)
	if err != nil {
		return "", fmt.Errorf("failed to parse given template: %s", err)
	}

	err = tmpl.Execute(b, evs)
	if err != nil {
		return "", fmt.Errorf("failed to execute given template: %s", err)
	}

	return b.String(), nil
}

func sendMessage(token, chat, msg string) error {
	c := &http.Client{}
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create http request: %s", err)
	}

	q := req.URL.Query()
	q.Add("chat_id", chat)
	q.Add("text", msg)
	q.Add("parse_mode", "HTML")
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make http request: %s", err)
	}

	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read response body: %s", err)
		}

		return fmt.Errorf("request failed with status code %d and body: \"%s\"", resp.StatusCode, string(b))
	}

	return nil
}
