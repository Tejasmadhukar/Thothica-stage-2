package utils

import (
	"bytes"
	"net/http"

	goose "github.com/Tejasmadhukar/GoOse"
)

var (
	g = goose.New()
)

func HandleHtml(Response *http.Response, obj map[string]interface{}) error {
	read_link := obj["readLink"].(string)
	title := obj["title"].(string)

	var buf bytes.Buffer

	buf.ReadFrom(Response.Body)

	article, err := g.ExtractFromRawHTML(buf.String(), "")
	if err != nil {
		return err
	}

	finalArticle := &Article{
		Title:     title,
		Content:   article.CleanedText,
		Author:    obj["authors"].([]interface{}),
		Publisher: obj["publisher"].(string),
		Title_URL: read_link,
	}

	newFilePath := output_dir + title + ".json"

	err = SaveArticle(finalArticle, newFilePath)
	if err != nil {
		return err
	}

	return nil
}
