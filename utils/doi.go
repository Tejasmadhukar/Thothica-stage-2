package utils

import (
	"net/url"

	"github.com/headzoo/surf/agent"
	"gopkg.in/headzoo/surf.v1"
)

var (
	bow = surf.NewBrowser()
)

func init() {
	bow.SetUserAgent(agent.Chrome())
}

func Handledoi(readLink string) (*url.URL, error) {
	err := bow.Open(readLink)
	if err != nil {
		return nil, err
	}

	return bow.Url(), nil
}
