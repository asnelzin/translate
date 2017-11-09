package yandex

import (
	"testing"
	"golang.org/x/text/language"
	"net/http"
	"fmt"
	"os"
	"net/http/httptest"
	"net/url"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/api"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server = httptest.NewServer(apiHandler)

	// github client configured to use test server
	client = NewClient(nil)
	u, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = u
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestClient_TranslateString(t *testing.T) {
	c := NewClient(nil)
	got, err := c.TranslateString(language.Russian, language.English, "Гибон")
	if err != nil {
		t.Fatalf("TranslateString returned unexpected error: %v", err)
	}
	want := "Hello"
	if got != want {
		t.Errorf("TranslateString result = %v, want %v", got, want)
	}
}
