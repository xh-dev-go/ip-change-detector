package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	notification "gitlab.xh-network.xyz/xeth/ip-changed-detector/notificatioin"
	"io"
	"net/http"
	"os"
)

var FailRequest = errors.New("fail to make request")
var FailReadResponse = errors.New("fail to read response")
var FailCreateFile = errors.New("fail to create file")

const FILE_NAME = ".ip-cur"

func fileExist() {
	_, err := os.Stat(FILE_NAME)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	} else if err != nil {
		writeFile([]byte(""))
	}
}

func readFile() string {
	bytes, err := os.ReadFile(FILE_NAME)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
func writeFile(bytes []byte) {
	f, err := os.Create(".ip-cur")
	if err != nil {
		panic(FailCreateFile)
	}
	f.Write(bytes)
}

func main() {
	tokenFlag := flagString.New("token", "token for tg bot").BindCmd()
	targetId := flagString.New("to", "target id").BindCmd()
	flag.Parse()

	fileExist()

	if respons, err := http.Get("https://api.myip.com"); err != nil {
		panic(FailRequest)
	} else if bytes, err := io.ReadAll(respons.Body); err != nil {
		panic(FailReadResponse)
	} else {
		md5Previouse := fmt.Sprintf("%x", md5.Sum([]byte(readFile())))
		md5Now := fmt.Sprintf("%x", md5.Sum(bytes))

		if md5Previouse != md5Now {
			writeFile(bytes)
			notification.Send(targetId.Value(), string(bytes), tokenFlag.Value())
		}
	}

}
