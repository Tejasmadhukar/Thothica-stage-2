package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	goose "github.com/Tejasmadhukar/GoOse"
	"github.com/fatih/color"
)

const (
	input_dir   = "./data"
	output_dir  = "./output"
	maxRoutines = 100
)

var (
	total_articles = 0
	bad_articles   = 0
	data           []map[string]interface{}
	g              = goose.New()
	wg             sync.WaitGroup
	routineChannel = make(chan struct{}, maxRoutines)
	mu             sync.Mutex
)

type Article struct {
	Title     string
	Content   string
	Author    []interface{}
	Publisher string
	Title_URL string
}

func process_article(obj map[string]interface{}, wg *sync.WaitGroup) {

	routineChannel <- struct{}{}

	defer func() {
		wg.Done()
		<-routineChannel
	}()

	read_link := obj["readLink"].(string)
	title := obj["title"].(string)

	article, err := g.ExtractFromURL(read_link)
	if err != nil {
		bad_articles += 1
		color.Red(read_link)
		return
	}

	finalArticle := &Article{
		Title:     title,
		Content:   article.CleanedText,
		Author:    obj["authors"].([]interface{}),
		Publisher: obj["publisher"].(string),
		Title_URL: article.FinalURL,
	}

	jsonArticle, err := json.Marshal(finalArticle)
	if err != nil {
		bad_articles += 1
		color.Red("Could not Marshal article", title, "to json")
		return
	}

	newFilePath := output_dir + "/" + title + ".json"

	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(newFilePath)
	if err != nil {
		bad_articles += 1
		color.Red(err.Error())
		return
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(jsonArticle)
	if err != nil {
		fmt.Println("Error writing json to file", title)
		bad_articles += 1
		fmt.Println(err)
		return
	}

	writer.Flush()
	color.Green(read_link)
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

			wg.Add(1)
			go process_article(obj, &wg)
		}
	}

	wg.Wait()
	close(routineChannel)

	fmt.Println(total_articles, bad_articles)
}
