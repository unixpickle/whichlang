package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const RoutineCount = 8

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: fetchlang-pastie <start idx> <end idx> <output_dir>")
		os.Exit(1)
	}
	startIdx, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	endIdx, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pasteIndices := make(chan int)
	go func() {
		for i := startIdx; i <= endIdx; i++ {
			pasteIndices <- i
		}
		close(pasteIndices)
	}()
	outDir := os.Args[3]
	if err := ensureDirectoryPresent(outDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fetchPastes(pasteIndices, outDir)
}

func fetchPastes(indices <-chan int, outDir string) {
	var wg sync.WaitGroup
	for i := 0; i < RoutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range indices {
				if err := fetchPaste(index, outDir); err != nil {
					fmt.Fprintln(os.Stderr, "error for paste", index, err)
				} else {
					fmt.Println("succeeded for paste", index)
				}
			}
		}()
	}
	wg.Wait()
}

func fetchPaste(index int, outDir string) error {
	code, lang, err := fetchPasteCode(index)
	if err != nil {
		return err
	}
	codeDir := filepath.Join(outDir, lang)
	ensureDirectoryPresent(codeDir)
	fileName := strconv.Itoa(index) + ".txt"
	return ioutil.WriteFile(filepath.Join(codeDir, fileName), []byte(code), 0755)
}

func fetchPasteCode(index int) (contents, language string, err error) {
	response, err := http.Get("http://pastie.org/pastes/" + strconv.Itoa(index))
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return
	}
	pageData := string(body)
	exp := regexp.MustCompile("\\<p\\>\n(.*?)\n&nbsp;")
	match := exp.FindStringSubmatch(pageData)
	if match == nil {
		ioutil.WriteFile("/Users/alex/Desktop/foo.html", []byte(pageData), 0755)
		return "", "", errors.New("cannot locate language")
	}
	language = match[1]
	language = strings.Replace(language, "/", ":", -1)

	response, err = http.Get("http://pastie.org/pastes/" + strconv.Itoa(index) + "/text")
	if err != nil {
		return
	}
	root, err := html.Parse(response.Body)
	response.Body.Close()
	if err != nil {
		return
	}
	codeBlock, ok := scrape.Find(root, scrape.ByTag(atom.Pre))
	if !ok {
		return "", "", errors.New("no <pre> tag")
	}
	contents = codeBlockText(codeBlock)
	return
}

func codeBlockText(n *html.Node) string {
	if n.DataAtom == atom.Br {
		return "\n"
	}
	if n.Type == html.TextNode {
		return n.Data
	}

	var res string
	child := n.FirstChild
	for child != nil {
		res += codeBlockText(child)
		child = child.NextSibling
	}
	return res
}

func ensureDirectoryPresent(dirPath string) error {
	if _, err := os.Stat(dirPath); err != nil {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}
