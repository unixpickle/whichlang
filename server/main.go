package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/unixpickle/whichlang"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: server <classifier.json> <assets_dir> <port>")
		os.Exit(1)
	}

	classifier := readClassifier()
	assets := os.Args[2]

	http.HandleFunc("/classify", func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		freqs := whichlang.ComputeFrequencies(string(contents))
		lang := classifier.Classify(freqs)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(lang))
	})
	http.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(":"+os.Args[3], nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readClassifier() *whichlang.Classifier {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var classifier whichlang.Classifier
	if err := json.Unmarshal(data, &classifier); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse classifier:", err)
		os.Exit(1)
	}

	return &classifier
}
