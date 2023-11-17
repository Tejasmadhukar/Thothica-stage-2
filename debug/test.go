package debug

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/headzoo/surf/agent"
	"gopkg.in/headzoo/surf.v1"
)

var (
	bow = surf.NewBrowser()
)

func Handledoi() {
	bow.SetUserAgent(agent.Chrome())
	err := bow.Open("https://doi.org/10.1088/1748-9326/ac0ac8")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bow.Url())
	// Now pass this url to goose
}

func downloadPdf(url, filepath string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("d&s successfully")
}

func Test() {
	url := "https://cp.copernicus.org/articles/17/2031/2021/cp-17-2031-2021.pdf"

	downloadPdf(url, "test.pdf")

	return 

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Header["Content-Type"][0] == "application/pdf")

	reader := bufio.NewReader(resp.Body)

	var pdfBuffer []byte

	for {
		chunk, err := reader.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		pdfBuffer = append(pdfBuffer, chunk)
	}

	fmt.Println(len(pdfBuffer))

}
