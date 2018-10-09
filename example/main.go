package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	stackdriver "github.com/yfuruyama/stackdriver-request-context-log"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		appLogger := stackdriver.LoggerFromRequest(r)
		appLogger.Debugf("This is a debug log")
		appLogger.Infof("This is an info log")
		appLogger.Warnf("This is a warning log")
		appLogger.Errorf("This is an error log")
		fmt.Fprintf(w, "OK\n")
	})

	projectId, _ := getDefaultProjectId()
	logger := stackdriver.NewLogger(os.Stderr, os.Stdout, projectId, stackdriver.SeverityInfo, stackdriver.AdditionalFields{
		"service": "foo",
		"version": 1.0,
	})
	handler := stackdriver.Handler(logger, mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}

func getDefaultProjectId() (string, error) {
	out, err := exec.Command("gcloud", "config", "list", "--format", "value(core.project)").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}
