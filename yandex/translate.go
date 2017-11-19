package yandex

import (
	"net/http"
	"net/url"
	"fmt"
	"bytes"
	"encoding/json"
	"golang.org/x/text/language"
	"io/ioutil"
)

const defaultBaseURL = "https://translate.yandex.net/api/v1.5/tr.json/"

type Client struct {
	client *http.Client

	BaseURL *url.URL
	APIKey  string
}

func NewClient(httpClient *http.Client, APIKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{httpClient, baseURL, APIKey}
}

type translateResult struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

func (c *Client) TranslateString(from, to language.Tag, text string) (string, error) {
	req, err := c.newTranslateRequest(from, to, text)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return "", err
	}

	var trResult translateResult
	err = json.NewDecoder(resp.Body).Decode(&trResult)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %v", err)
	}

	if len(trResult.Text) == 0 {
		return "", fmt.Errorf("got empty text array in JSON")
	}
	return trResult.Text[0], nil
}

func (c *Client) newTranslateRequest(from, to language.Tag, text string) (*http.Request, error) {
	var urlStr string
	if from != language.Und {
		urlStr = fmt.Sprintf("translate?lang=%s-%s&key=%s", from.String(), to.String(), c.APIKey)
	} else {
		urlStr = fmt.Sprintf("translate?lang=%s&key=%s", to.String(), c.APIKey)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	payload := url.Values{}
	payload.Set("text", text)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

type ErrorResponse struct {
	Response *http.Response
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("got an error from Yandex API: HTTP %v %d, %v %v",
		r.Response.Request.Method, r.Response.StatusCode, r.Code, r.Message)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c == 200 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
