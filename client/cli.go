package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const url = "http://localhost:8080/grayscale"

func main() {
	f, err := os.Open("source.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)

	client := http.Client{}

	response, err := client.Post(
		url, "", bytes.NewReader(data),
	)
	if err != nil {
		panic(err)
	}

	resFile, err := os.Create("res.png")
	if err != nil {
		panic(err)
	}
	defer resFile.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(resFile, response.Body)
	if err != nil {
		panic(err)
	}
}
