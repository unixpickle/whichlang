package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/unixpickle/whichlang"
	"github.com/unixpickle/whichlang/tokens"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Usage: server <algorithm> <classifier.json> <assets_dir> <port>")
		os.Exit(1)
	}

	classifier := readClassifier()
	assets := os.Args[3]

	http.HandleFunc("/classify", func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		counts := tokens.CountTokens(string(contents))
		freqs := counts.Freqs()
		lang := classifier.Classify(freqs)
		jsonObj := map[string]interface{}{"lang": lang}
		jsonData, _ := json.Marshal(jsonObj)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
	http.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(":"+os.Args[4], nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readClassifier() whichlang.Classifier {
	decoder := whichlang.Decoders[os.Args[1]]
	if decoder == nil {
		fmt.Fprintln(os.Stderr, "Unknown algorithm:", os.Args[1])
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	c, err := decoder(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return c
}
