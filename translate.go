package gtranslate

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Define constants for API paths and user agent
const (
	googleTranslateAPI      = "https://translate.google.com/translate_a/single"
	googleTranslateBatchAPI = "https://translate.googleapis.com/translate_a/t"
	userAgent               = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15"
)

// Define an HTTP client with a timeout
var client = &http.Client{
	Timeout: time.Second * 10, // Timeout after 10 seconds
}

// Function to prepare a URL for API request
func prepareURL(apiPath string, data map[string]string) (*url.URL, error) {
	u, err := url.Parse(apiPath)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	// Create URL parameters
	parameters := url.Values{}
	for k, v := range data {
		parameters.Add(k, v)
	}
	for _, v := range []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"} {
		parameters.Add("dt", v)
	}

	u.RawQuery = parameters.Encode()
	return u, nil
}

// Function to execute a HTTP request
func doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error in calling Google translate API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

// Function to translate a single piece of content
func Translate(ctx context.Context, content string, sourceLanguage, targetLanguage language.Tag) (*TranslationResult, error) {
	sourceLanguageStr := "auto"
	if !sourceLanguage.IsRoot() {
		sourceLanguageStr = sourceLanguage.String()
	}
	targetLanguageStr := "fa"
	if !targetLanguage.IsRoot() {
		targetLanguageStr = targetLanguage.String()
	}

	token := getToken(content)

	// Prepare the data for the request
	data := map[string]string{
		"client": "gtx",
		"sl":     sourceLanguageStr,
		"tl":     targetLanguageStr,
		"hl":     targetLanguageStr,
		"ie":     "UTF-8",
		"oe":     "UTF-8",
		"otf":    "1",
		"ssel":   "0",
		"tsel":   "0",
		"kc":     "7",
		"q":      content,
		"tk":     token,
	}

	u, err := prepareURL(googleTranslateAPI, data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	body, err := doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	t, err := parseTranslationJSON(body)
	if err != nil {
		return nil, fmt.Errorf("parse translation JSON: %w", err)
	}

	return t, nil
}

// Function to translate a batch of content
func TranslateBatch(ctx context.Context, contents []string, from string, to string) ([]string, error) {
	preparedText := encodeForBatch(contents)
	token := getToken(strings.Join(preparedText, ""))

	data := map[string]string{
		"anno":   "3",
		"client": "te",
		"v":      "1.0",
		"format": "html",
		"sl":     from,
		"tl":     to,
		"tk":     token,
	}

	u, err := prepareURL(googleTranslateBatchAPI, data)
	if err != nil {
		return nil, err
	}

	body := strings.Join(preparedText, "&q=")
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	respBody, err := doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return []string{string(respBody)}, nil
}

// Function to prepare the content for batch translation
func encodeForBatch(textList []string) []string {
	encodedText := make([]string, len(textList))
	for i, text := range textList {
		encodedText[i] = fmt.Sprintf("<pre><a i=\"%d\">%s</a></pre>", i, text)
	}
	return encodedText
}
