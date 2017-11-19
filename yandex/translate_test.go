package yandex

import (
	"testing"
	"golang.org/x/text/language"
	"net/http"
	"net/http/httptest"
	"net/url"
	"fmt"
	"reflect"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Yandex client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	server = httptest.NewServer(apiHandler)

	// yandex client configured to use test server
	client = NewClient(nil, "secret")
	u, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = u
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func TestClient_TranslateString(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		testFormValues(t, r, values{
			"key":  "secret",
			"lang": "en-ru",
			"text": "Hello",
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "lang": "en-ru", "text": ["Привет"]}`)
	})

	got, err := client.TranslateString(language.English, language.Russian, "Hello")
	if err != nil {
		t.Fatalf("TranslateString returned unexpected error: %v", err)
	}
	want := "Привет"
	if got != want {
		t.Errorf("TranslateString result = %v, want %v", got, want)
	}
}

func TestClient_TranslateString_returnFirst(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "lang": "en-ru", "text": ["Привет", "мир"]}`)
	})

	got, err := client.TranslateString(language.English, language.Russian, "Hello")
	if err != nil {
		t.Fatalf("TranslateString returned unexpected error: %v", err)
	}
	want := "Привет"
	if got != want {
		t.Errorf("TranslateString result = %v, want %v", got, want)
	}
}

func TestClient_TranslateString_determineLang(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		testFormValues(t, r, values{
			"key":  "secret",
			"lang": "ru",
			"text": "Hello",
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "lang": "ru", "text": ["Привет"]}`)
	})

	got, err := client.TranslateString(language.Und, language.Russian, "Hello")
	if err != nil {
		t.Fatalf("TranslateString returned unexpected error: %v", err)
	}
	want := "Привет"
	if got != want {
		t.Errorf("TranslateString result = %v, want %v", got, want)
	}
}

func TestClient_TranslateString_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	got, err := client.TranslateString(language.English, language.Russian, "Hello")

	if err == nil {
		t.Fatal("Expected HTTP 500 error, got no error.")
	}
	if want := ""; got != want {
		t.Errorf("Expected empty string, got %v", got)
	}
}

func TestClient_TranslateString_jsonError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"code":401,"message":"API key is invalid"}`)
	})

	got, err := client.TranslateString(language.English, language.Russian, "Hello")

	if err == nil {
		t.Fatal("Expected HTTP 403 error, got no error.")
	}

	if err.(*ErrorResponse).Message != "API key is invalid" {
		t.Errorf("Expected correct error message, got %v", err.(*ErrorResponse).Message)
	}

	if want := ""; got != want {
		t.Errorf("Expected empty string, got %v", got)
	}
}
