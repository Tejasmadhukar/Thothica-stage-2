package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/ledongthuc/pdf"
)

const (
	temp_dir   = "./temp/"
	output_dir = "./output/"
)

type Article struct {
	Title     string
	Content   string
	Author    []interface{}
	Publisher string
	Title_URL string
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	content := ""

	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				content += word.S + " "
			}
			content += "\n"
		}
	}

	return content, nil
}

func HandlePdf(obj map[string]interface{}, Response *http.Response) error {
	title := obj["title"].(string)

	filepath := temp_dir + title + ".pdf"

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, Response.Body)
	if err != nil {
		return err
	}

	content, err := readPdf(filepath)
	if err != nil {
		return err
	}

	finalArticle := &Article{
		Title:     title,
		Content:   content,
		Author:    obj["authors"].([]interface{}),
		Publisher: obj["publisher"].(string),
		Title_URL: obj["readLink"].(string),
	}

	newFilePath := output_dir + title + ".json"

	err = SaveArticle(finalArticle, newFilePath)
	if err != nil {
		return err
	}

	return nil
}
