package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	// Be careful to use constant time comparison to prevent timing attacks
	// 	if subtle.ConstantTimeCompare([]byte(username), []byte("foppa")) == 1 &&
	// 		subtle.ConstantTimeCompare([]byte(password), []byte("T0ffl0r")) == 1 {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))
	// Routes
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.GET("/", A11y)
	e.HEAD("/", A11y)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))

}

// A11yServer runs the server

func A11y(c echo.Context) error {

	files, err := ReadDir("reports/")

	if err != nil {
		panic(err)
	}
	output, err := json.Marshal(files)
	if err != nil {
		log.Fatal(err)
	}
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, string(output))
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
				}
			}
		}

	}

	return files, nil
}

// https://stackoverflow.com/questions/32962128/create-iterative-json-directory-tree-golang
// https://dev.to/manigandand/list-files-in-a-directory-using-golang-3k78
