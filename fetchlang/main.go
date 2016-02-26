package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
)

const MinFileSize = 100

type Language struct {
	Name       string
	Extensions []string
}

var Languages = []Language{
	{"ActionScript", []string{"as"}},
	{"C", []string{"c"}},
	{"C#", []string{"cs"}},
	{"C++", []string{"cpp", "c++", "C"}},
	{"Clojure", []string{"clj"}},
	{"CoffeeScript", []string{"coffee"}},
	{"CSS", []string{"css"}},
	{"Go", []string{"go"}},
	{"Haskell", []string{"hs"}},
	{"HTML", []string{"html", "htm"}},
	{"Java", []string{"java"}},
	{"JavaScript", []string{"js"}},
	{"Lua", []string{"lua"}},
	{"Matlab", []string{"m"}},
	{"Objective-C", []string{"m"}},
	{"Perl", []string{"pl"}},
	{"PHP", []string{"php"}},
	{"Python", []string{"py"}},
	{"R", []string{"r"}},
	{"Ruby", []string{"rb"}},
	{"Scala", []string{"scala"}},
	{"Shell", []string{"sh", "bash"}},
	{"Swift", []string{"swift"}},
	{"TeX", []string{"tex", "TeX"}},
	{"VimL", []string{"vim"}},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: fetchlang <sample dir>")
		os.Exit(1)
	}
	sampleDir := os.Args[1]

	fmt.Print("Username: ")
	username := readInput()

	fmt.Print("Password: ")
	password, _ := gopass.GetPasswd()

	client := &GithubClient{User: username, Pass: string(password)}

	for _, lang := range Languages {
		fmt.Println("Fetching samples for", lang.Name, "...")
		repos, err := client.LanguageRepositories(lang.Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if err := os.Mkdir(filepath.Join(sampleDir, lang.Name), 0755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		for i, repo := range repos {
			file, err := client.FirstFile(repo, MinFileSize, lang.Extensions)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			if file == nil {
				continue
			}
			targetFile := filepath.Join(sampleDir, lang.Name, strconv.Itoa(i)+"."+
				lang.Extensions[0])
			ioutil.WriteFile(targetFile, file, 0755)
		}
	}
}

func readInput() string {
	res := ""
	for {
		var ch [1]byte
		if _, err := os.Stdin.Read(ch[:]); err != nil || ch[0] == '\n' {
			break
		}
		res += string(ch[0])
	}
	return strings.TrimSpace(res)
}
