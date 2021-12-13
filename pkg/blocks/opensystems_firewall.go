package blocks

import (
	"dfe/pkg/helpers"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type OSFW struct {
	Message  <-chan string
	Hostname <-chan string
	Out      chan<- string
}

type osfwLog struct {
	TimestampUsec int64 `json:"timestampUsec"`
}

func osfw2syslog_str(message, hostname string) (string, error) {
	l := &osfwLog{}
	err := json.Unmarshal([]byte(message), l)
	if err != nil {
		return "", err
	}
	p := 0
	tag := ""

	nl := ""
	if strings.HasSuffix(message, "\n") {
		nl = "\n"
	}
	timestamp := time.UnixMicro(l.TimestampUsec).Format(time.RFC3339)
	return fmt.Sprintf("<%d>%s %s %s[%d]: %s%s",
		p, timestamp, hostname,
		tag, os.Getpid(), message, nl), nil
}

func NewOSFWParser() *OSFW {
	return &OSFW{}
}

func (c *OSFW) Process() {
	guard := helpers.NewInputGuard("message", "hostname")

	message := ""
	hostname := ""

	for {
		select {
		case m, ok := <-c.Message:
			if ok {
				message = m
				if hostname != "" {
					s, err := osfw2syslog_str(message, hostname)
					if err == nil {
						c.Out <- s
					}
					message = ""
					hostname = ""
				}
			} else if guard.Complete("message") {
				return
			}
		case h, ok := <-c.Hostname:
			if ok {
				hostname = h
				if message != "" {
					s, err := osfw2syslog_str(message, hostname)
					if err == nil {
						c.Out <- s
					}
					message = ""
					hostname = ""
				}
			} else if guard.Complete("hostname") {
				return
			}
		}
	}
}
