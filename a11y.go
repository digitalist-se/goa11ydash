package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", A11yServer)
	http.ListenAndServe(":8080", nil)
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

// A11yServer runs the server
func A11yServer(w http.ResponseWriter, r *http.Request) {
	lines := []string{}
	files, err := ioutil.ReadDir("reports/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		lines = append(lines, file.Name())
	}

	output, err := json.Marshal(lines)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println(string(jsonData))
	fmt.Fprintf(w, string(output))
}
