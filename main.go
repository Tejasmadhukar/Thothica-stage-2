package main

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/GoOse"
	"log"
	"net/http"
	"os"
)

const (
	input_dir  = "./data"
	output_dir = "./output"
)

var (
	total_articles = 0
	bad_articles   = 0
	data           []map[string]interface{}
	client         = &http.Client{}
	g              = goose.New()
)

func process_article(obj map[string]interface{}) {
	read_link := obj["readLink"].(string)

	fmt.Println(read_link)

	article, err := g.ExtractFromURL(read_link)
	if err != nil {
		bad_articles += 1
		fmt.Println("Could not get article titled", obj["title"].(string))
		return
	}

	println("title", article.Title)
	println("description", article.MetaDescription)
	println("keywords", article.MetaKeywords)
	println("content", article.CleanedText)
	println("url", article.FinalURL)
	println("top image", article.TopImage)
}

func main() {
	fmt.Println("Reading files from", input_dir)

	files, err := os.ReadDir(input_dir)
	if err != nil {
		log.Fatal(err)
	}

	keysToCheck := []string{"title", "readLink", "authors", "publisher"}

	for _, file := range files {
		file_path := input_dir + "/" + file.Name()

		f, err := os.Open(file_path)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		err = json.NewDecoder(f).Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

	Loop:
		for _, obj := range data {
			if obj == nil {
				continue
			}

			total_articles += 1

			for _, key := range keysToCheck {
				if _, exists := obj[key]; !exists {
					bad_articles += 1
					continue Loop
				}
			}

			process_article(obj)

		}
	}

	fmt.Println(total_articles, bad_articles)
}
