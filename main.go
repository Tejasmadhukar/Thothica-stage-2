package main

import (
	"encoding/json"
	"fmt"
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
)

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

			read_link := obj["readLink"].(string)

			fmt.Println(read_link)

			req, err := http.NewRequest("GET", read_link, nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			defer res.Body.Close()

			if res.StatusCode != 200 {
				bad_articles += 1
				fmt.Println("Could not get article", f.Name())
				continue Loop
			}
			fmt.Println(res.StatusCode)

		}

		f.Close()
	}
	fmt.Println(total_articles, bad_articles)
}
