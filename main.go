package main

import (
	"errors"
	"io"
	"log"
	"os"
	"regexp"

	"git.sr.ht/~adnano/go-gemini"
	"golang.org/x/net/context"
)

const URL string = "gemini://konpeito.media/"

func getData(URI string) *gemini.Response {
	client := &gemini.Client{}
	ctx := context.Background()
	resp, err := client.Get(ctx, URI)
	handle(err)
	return resp
}

func readResponseBody(body io.ReadCloser) []byte {
	bytes, err := io.ReadAll(body)
	handle(err)
	return bytes
}

func writeToFile(bytes []byte, path string) {
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	handle(err)
	defer file.Close()

	_, err = file.Write(bytes)
	handle(err)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}

func replace(reg string, orig []byte, neu []byte) []byte {
	query := regexp.MustCompile(reg)
	result := query.ReplaceAll(orig, neu)
	return result
}

func handle(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func checkPath(path string) string {
	if len(path) > 1 {
		matched, err := regexp.Match(`.*\/$`, []byte(path))
		handle(err)
		if !matched {
			path = path + "/"
		}
		return path
	}
	return ""
}

func main() {
	path := ""
	if len(os.Args) > 1 {
		path = checkPath(os.Args[1])
	}
	resp := getData(URL)
	defer resp.Body.Close()
	if resp.Status.Class() == gemini.StatusSuccess {
		bytes := readResponseBody(resp.Body)
		response := regexp.MustCompile(`\=\>.*\.mp3`).FindAll(bytes, -1)
		for i := 0; i < len(response); i++ {
			link := string(replace(`\=\>\s*`, response[i], []byte("")))
			filename := path + string(replace(URL, []byte(link), []byte("")))
			if fileExists(filename) {
				resp = getData(link)
				defer resp.Body.Close()
				if resp.Status.Class() == gemini.StatusSuccess {
					bytes = readResponseBody(resp.Body)
					log.Println("Saving to " + filename)
					writeToFile(bytes, filename)
				}
			} else {
			}
		}
	} else {
		handle(errors.New("Bad response"))
	}
}
