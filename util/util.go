package util

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

// tokenAPI 获取带 token 的 API 地址
func TokenAPI(api, token string) (string, error) {
	queries := RequestQueries{
		"access_token": token,
	}

	return EncodeURL(api, queries)
}

// encodeURL add and encode parameters.
func EncodeURL(api string, params RequestQueries) (string, error) {
	url, err := url.Parse(api)
	if err != nil {
		return "", err
	}

	query := url.Query()

	for k, v := range params {
		query.Set(k, v)
	}

	url.RawQuery = query.Encode()

	return url.String(), nil
}

// getQuery returns url query value
func GetQuery(req *http.Request, key string) string {
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0]
	}

	return ""
}

// randomString random string generator
//
// ln length of return string
func RandomString(ln int) string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, ln)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}

	return string(b)
}

// postJSON perform a HTTP/POST request with json body
func GetJSON(url string, response interface{}) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}

// postJSONWithBody return with http body.
func PostJSONWithBody(url string, params interface{}) (*http.Response, error) {
	reader := new(bytes.Reader)
	if params != nil {
		raw, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		reader = bytes.NewReader(raw)
	}

	return http.Post(url, "application/json; charset=utf-8", reader)
}

func PostFormByFile(url, field, filename string, response interface{}) error {
	// Add your media file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return postForm(url, field, filename, file, response)
}

func PostForm(url, field, filename string, reader io.Reader, response interface{}) error {
	// Prepare a form that you will submit to that URL.
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile(field, filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fw, reader); err != nil {
		return err
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}
