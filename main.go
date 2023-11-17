package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	utils "github.com/Tejasmadhukar/Thothica-stage-2/utils"
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
	wg             sync.WaitGroup
	routineChannel = make(chan struct{}, maxRoutines)
)

func process_article(obj map[string]interface{}) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			bad_articles += 1
			return
		}
	}()

	routineChannel <- struct{}{}

	defer func() {
		wg.Done()
		<-routineChannel
	}()

	read_link := obj["readLink"].(string)

	urlPattern := `^(https?://)?(dx\.)?doi\.org/`

	re, err := regexp.Compile(urlPattern)
	if err != nil {
		log.Fatal("Wrong regex")
	}

	if re.MatchString(read_link) {
		newLink, err := utils.Handledoi(read_link)
		if err != nil {
			bad_articles += 1
			color.Red(read_link)
			return
		}

		obj["readLink"] = newLink.String()
		read_link = newLink.String()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Get(read_link)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.Header["Content-Type"][0] == "application/pdf" {
		err = utils.HandlePdf(obj, resp)
		if err != nil {
			bad_articles += 1
			color.Red(read_link)
			return
		}
		return
	}

	err = utils.HandleHtml(resp, obj)
	if err != nil {
		bad_articles += 1
		color.Red(read_link)
		return
	}

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
			go process_article(obj)
		}
	}

	wg.Wait()
	close(routineChannel)

	fmt.Println(total_articles, bad_articles)
}
