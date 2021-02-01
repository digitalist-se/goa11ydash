package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	http.HandleFunc("/", A11yServer)
	http.ListenAndServe(":8080", nil)
}

// A11yServer runs the server
func A11yServer(w http.ResponseWriter, r *http.Request) {

	//lines := []string{}
	files, err := ReadDir("reports/")

	if err != nil {
		panic(err)
	}

	output, err := json.Marshal(files)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println(string(jsonData))
	fmt.Fprintf(w, string(output))
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			base := filepath.Base(path)
			l := "2000-01-02"
			_, err := time.Parse(l, base)
			if err != nil {
				return filepath.SkipDir
			}
			files = append(files, path)
		}
		log.Output(-1, files[0])
		return nil
	})
	return files, err
}

func ReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		var first string
		first = root + file.Name()
		f, err := os.Open(first)
		if err != nil {
			return files, err
		}
		fileInfo, err := f.Readdir(-1)
		for _, file := range fileInfo {
			var second string
			second = first + "/" + file.Name() + "/pages"
			f, err := os.Open(second)
			if err != nil {
				return files, err
			}
			fileInfo, err := f.Readdir(-1)
			for _, file := range fileInfo {
				var third string
				third = second + "/" + file.Name()
				f, err := os.Open(third)
				if err != nil {
					return files, err
				}
				fileInfo, err := f.Readdir(-1)
				for _, file := range fileInfo {
					var fourth string
					fourth = third + "/" + file.Name()
					files = append(files, fourth)
					// What I get:
					//["reports/bar/2021-01-25/pages/kulturhusetstadsteatern.dgstage.se/node-1","reports/bar/2021-01-26/pages/kulturhusetstadsteatern.dgstage.se/node-2"]
					// What I want explained in
					// reports:
					//	 bar:
					//	   2021-01-25:
					//	 	 pages:
					//			kulturhusetstadsteatern.dgstage.se:
					//			- node-1
					//			- node-2
					//			- node-3
				}
			}
		}

	}

	return files, nil
}

// https://stackoverflow.com/questions/32962128/create-iterative-json-directory-tree-golang
// https://dev.to/manigandand/list-files-in-a-directory-using-golang-3k78
