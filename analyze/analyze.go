package analyze

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const stringToFind = "Go"
const method = "GET"

func CountingURL(url string) (c uint64, err error) {

	client := &http.Client{}
	req, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		return 0, reqErr
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		return 0, resErr
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if err != nil {
			err = closeErr
		}
	}(res.Body)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return 0, resErr
	}
	count := strings.Count(string(body), stringToFind)
	return uint64(count), nil
}

func CountingFile(filename string) (c uint64, err error) {
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		return 0, fmt.Errorf("open file error: %w", fileErr)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(file)

	fileContent, readErr := ioutil.ReadAll(file)

	if readErr != nil {
		return 0, fmt.Errorf("reading file(%s) error: %w", filename, readErr)
	}

	count := strings.Count(string(fileContent), stringToFind)
	return uint64(count), nil
}
