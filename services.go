package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/findWords", findWordsHandler)
	router.HandleFunc("/grayscale", grayscaleHandler)
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		return
	}
}

func findWordsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Char  string   `json:"char"`
		Words []string `json:"words"`
	}
	req := request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	res := make([]string, 0, len(req.Words))

	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	for _, word := range req.Words {
		wg.Add(1)
		go func(word string, letter string) {
			defer wg.Done()
			if strings.Contains(word, letter) {
				m.Lock()
				defer m.Unlock()
				res = append(res, word)
			}
		}(word, req.Char)
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func grayscaleHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	img, _, _ := image.Decode(bytes.NewReader(data))

	newImg := image.NewRGBA64(img.Bounds())

	wg := sync.WaitGroup{}
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		wg.Add(1)
		func(x int) {
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				curr := img.At(x, y)
				r, g, b, a := curr.RGBA()
				//gray := uint8(float32(r)*0.2126 + float32(g)*0.7152 + float32(b)*0.0722)
				gray := uint16(float32(r+g+b) / 3.)
				newColor := color.RGBA64{R: gray, G: gray, B: gray, A: uint16(a)}
				newImg.Set(x, y, newColor)
			}
			wg.Done()
		}(x)
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)
	err := png.Encode(w, newImg)
	if err != nil {
		return
	}
}
