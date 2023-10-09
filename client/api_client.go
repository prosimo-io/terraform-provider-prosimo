package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	defaultBaseURL = "https://api.github.com/"
)

type api_client struct {
	client      *http.Client
	baseURL     *url.URL
	apiEndPoint string
	token       string
	insecure    bool
}

type Response struct {
	*http.Response

	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int
}

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

func newApiClient() *api_client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{Transport: tr}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &api_client{client: httpClient, baseURL: baseURL}

	return c

}

func (c *api_client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *api_client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	// log.Println("NewRequest", u.String(), buf)
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	// log.Println("response", req)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Prosimo-ApiToken", c.token)
	return req, nil
}

func (c *api_client) ReqFileUpload(method, urlStr string, inputData map[string]string, Filepath string) (*http.Request, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range inputData {
		_ = writer.WriteField(k, v)
	}
	file, errFile := os.Open(Filepath)
	defer file.Close()
	filepart, errFile := writer.CreateFormFile("details", filepath.Base(Filepath))
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}
	_, errFile = io.Copy(filepart, file)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("Prosimo-ApiToken", c.token)
	return req, nil
}
func (c *api_client) ReqFileUploadIDP(method, urlStr string, inputData map[string]string, Filepath string) (*http.Request, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range inputData {
		_ = writer.WriteField(k, v)
	}
	file, errFile := os.Open(Filepath)
	defer file.Close()
	filepart, errFile := writer.CreateFormFile("apiFile", filepath.Base(Filepath))
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}
	_, errFile = io.Copy(filepart, file)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("Prosimo-ApiToken", c.token)
	return req, nil
}
func (c *api_client) ReqCertUpload(method, urlStr string, certpath string, keypath string) (*http.Request, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	Certfile, errFile := os.Open(certpath)
	defer Certfile.Close()

	Certfilepart, errFile := writer.CreateFormFile("certificate", filepath.Base(certpath))
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}

	_, errFile = io.Copy(Certfilepart, Certfile)

	Keyfile, errFile := os.Open(keypath)
	defer Keyfile.Close()

	Keyfilepart, errFile := writer.CreateFormFile("privateKey", filepath.Base(keypath))
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}
	_, errFile = io.Copy(Keyfilepart, Keyfile)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(payload)
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Prosimo-ApiToken", c.token)
	return req, nil
}

func (c *api_client) ReqCACertUpload(method, urlStr string, certpath string) (*http.Request, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	Certfile, errFile := os.Open(certpath)
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}
	defer Certfile.Close()

	Certfilepart, errFile := writer.CreateFormFile("certificate", filepath.Base(certpath))
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}

	_, errFile = io.Copy(Certfilepart, Certfile)
	if errFile != nil {
		fmt.Println(errFile)
		return nil, errFile
	}

	// Keyfile, errFile := os.Open(keypath)
	// defer Keyfile.Close()

	// Keyfilepart, errFile := writer.CreateFormFile("privateKey", filepath.Base(keypath))
	// if errFile != nil {
	// 	fmt.Println(errFile)
	// 	return nil, errFile
	// }
	// _, errFile = io.Copy(Keyfilepart, Keyfile)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(payload)
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Prosimo-ApiToken", c.token)
	return req, nil
}

func (c *api_client) Do(ctx context.Context, req *http.Request, returnObj interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	body, _ := io.ReadAll(resp.Body)
	// log.Printf("respBody   %s", string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if returnObj != nil {
		if w, ok := returnObj.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(returnObj)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				log.Printf("Error in marshall response %s", string(decErr.Error()))
				err = decErr
			}
		}
	}

	return response, err
}

func (c *api_client) PostRequest(ctx context.Context, postEndpoint string, postData interface{}) (*ResourcePostResponseData, error) {
	req, err := c.NewRequest("POST", postEndpoint, postData)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = c.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (c *api_client) PutRequest(ctx context.Context, postEndpoint string, postData interface{}) (*ResourcePostResponseData, error) {
	req, err := c.NewRequest("PUT", postEndpoint, postData)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = c.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (c *api_client) AppOffboarding_DeleteRequest(ctx context.Context, deleteEndpoint string) (*ResourcePostResponseData, error) {

	req, err := c.NewRequest("DELETE", deleteEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = c.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (c *api_client) DeleteRequest(ctx context.Context, deleteEndpoint string) error {

	req, err := c.NewRequest("DELETE", deleteEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (c *api_client) GetRequest(ctx context.Context, getEndpoint string) (*Task_Status_res, error) {

	req, err := c.NewRequest("GET", getEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Task_Status_res{}
	_, err = c.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func CheckResponse(r *http.Response) error {

	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse

}
