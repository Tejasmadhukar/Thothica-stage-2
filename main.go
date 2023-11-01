package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	input_dir  = "./data"
	output_dir = "./output"
)

func main() {
	var data []map[string]interface{}

	fmt.Println("Reading files from ", input_dir)

	files, err := os.ReadDir(input_dir)
	if err != nil {
		log.Fatal(err)
	}

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

		for _, obj := range data {
			if obj == nil {
				continue
			}

			for key, value := range obj {
				fmt.Println(key, value)
			}
		}

		f.Close()
	}
}
