package yandex

import (
	"net/http"
	"net/url"
	"fmt"
	"bytes"
	"encoding/json"
	"golang.org/x/text/language"
)

const (
	defaultBaseURL = "https://translate.yandex.net/api/v1.5/tr.json/"
	apiKey         = "trnsl.1.1.20171102T195151Z.0ed6e46b065fb5c5.5fb7f63e7f6d9bb06e348bb27093408ba9b00618"
)

type Client struct {
	client *http.Client

	BaseURL *url.URL
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{httpClient, baseURL}
}

type translateResult struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

func (c *Client) TranslateString(from, to language.Tag, text string) (string, error) {
	urlStr := fmt.Sprintf("translate?lang=%s-%s&key=%s", from.String(), to.String(), apiKey)
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return "", err
	}

	payload := url.Values{}
    payload.Set("text", text)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var trResult translateResult
	json.NewDecoder(resp.Body).Decode(&trResult)
	return trResult.Text[0], nil
}
