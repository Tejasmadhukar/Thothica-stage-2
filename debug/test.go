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
	"github.com/ledongthuc/pdf"
	"gopkg.in/headzoo/surf.v1"
)

var (
	bow = surf.NewBrowser()
)

func readPdf(path string) string {
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Fatal(err)
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

	return content
}

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

	fmt.Println("Download good")
}

func Test() {
	url := "https://cp.copernicus.org/articles/17/2031/2021/cp-17-2031-2021.pdf"

	downloadPdf(url, "test.pdf")
	content := readPdf("test.pdf")
	fmt.Println(content)
	err := os.Remove("test.pdf")
	if err != nil {
		log.Fatal("Error deleting pdf", err)
	}

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
