package notification

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Send(to string, msg string, token string) {
	if to == "" {
		panic(errors.New("No chat_id provided!"))
	} else if msg == "" {
		panic(errors.New("No message provided!"))
	}

	msg = strings.ReplaceAll(msg, "{", "\\{")
	msg = strings.ReplaceAll(msg, "}", "\\}")
	msg = strings.ReplaceAll(msg, ".", "\\.")

	data := url.Values{
		"parse_mode": {"MarkdownV2"},
		"chat_id":    {to},
		"text":       {msg},
	}
	if resp, err := http.PostForm(fmt.Sprintf(`https://api.telegram.org/bot%s/sendMessage`, token), data); err != nil {
		panic(err)
	} else if byte, err := ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	} else {
		println(string(byte))
	}
}
