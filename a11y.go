// Env. variables that should be set:
// JSONHOST = the host that deliveres the jsonfiles from Sitespeedio, it not set, defaults to http://localhost:9000
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

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
	e.GET("/", Serve)
	e.HEAD("/", Serve)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))

}

// Serve serves the response.
func Serve(c echo.Context) error {

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

// ReadDir reads the directories.
func ReadDir(root string) ([]string, error) {

	// env. variables
	var host string
	host = os.Getenv("JSONHOST")
	if host == "" {
		host = "http://localhost:9000"
	}

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

					var doexist string
					doexist = fourth + "/data"

					if _, err := os.Stat(doexist); err == nil {

						var axe string
						axe = doexist + "/axe.run-3.json"

						var axeSummary string
						axeSummary = doexist + "/axe.pageSummary.json"

						var browsertime string
						browsertime = doexist + "/browsertime.run-3.json"

						var browsertimeHar string
						browsertimeHar = doexist + "/browsertime.har"

						var browsertimeSummary string
						browsertimeSummary = doexist + "/browsertime.pageSummary.json"

						var coach string
						coach = doexist + "/coach.run-3.json"

						var coachSummary string
						coachSummary = doexist + "/coach.pageSummary.json"

						var pageXray string
						pageXray = doexist + "/pagexray.run-3.json"

						var pageXraySummary string
						pageXraySummary = doexist + "/pagexray.pageSummary.json"

						var thirdPartyRun string
						thirdPartyRun = doexist + "/thirdparty.run.json"

						var thridPartySummary string
						thridPartySummary = doexist + "/thirdparty.pageSummary.json"

						files = append(files, axe, browsertime, axeSummary, browsertimeHar, browsertimeSummary, coach, coachSummary, pageXray, pageXraySummary, thirdPartyRun, thridPartySummary)

					}

				}
			}
		}
	}

	// Remove reports dir from string and replace it with domain containing json reports
	for i := range files {
		files[i] = strings.ReplaceAll(files[i], "reports", host)
	}

	return files, nil
}

// https://github.com/unrolled/secure
// https://github.com/labstack/echo/
