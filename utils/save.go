package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

var (
	mu sync.Mutex
)

func SaveArticle(finalArticle *Article, path string) error {
	jsonArticle, err := json.Marshal(finalArticle)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(jsonArticle)
	if err != nil {
		return err
	}

	return nil
}
