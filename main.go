package main

import (
	"fmt"
	"log"
	"os"
)

const (
  input_dir = "./data"
  output_dir = "./output"
)

func main() {
  fmt.Println("Reading files from ",input_dir)

  files, err := os.ReadDir(input_dir)
  if err != nil {
    log.Fatal(err)
  }
  
  for _,el := range files {
    fmt.Println(el)
  }
}
