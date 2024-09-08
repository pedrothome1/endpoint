package endpoint

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasttemplate"
	"io"
	"net/http"
	"net/url"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Endpoint struct {
	doer        Doer
	urlTemplate string
	options     *endpointOptions
}

func (e *Endpoint) Head(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("HEAD", opts...)
}

func (e *Endpoint) Get(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("GET", opts...)
}

func (e *Endpoint) Post(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("POST", opts...)
}

func (e *Endpoint) Put(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("PUT", opts...)
}

func (e *Endpoint) Patch(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("PATCH", opts...)
}

func (e *Endpoint) Delete(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("DELETE", opts...)
}

func (e *Endpoint) Options(opts ...requestOptionSetter) (*Response, error) {
	return e.doRequest("OPTIONS", opts...)
}

func (e *Endpoint) doRequest(method string, opts ...requestOptionSetter) (*Response, error) {
	reqOptions := defaultRequestOptions()

	for _, opt := range opts {
		opt(&reqOptions)
	}

	reqURL, err := url.ParseRequestURI(e.urlTemplate)
	if err != nil {
		return nil, err
	}

	template := fasttemplate.New(reqURL.Path, "{", "}")
	reqURL.Path = template.ExecuteString(reqOptions.pathParams)

	body, err := func() (io.Reader, error) {
		if reqOptions.rawBody != nil {
			return bytes.NewReader(reqOptions.rawBody), nil
		}

		if reqOptions.jsonBody != nil {
			b, err := json.Marshal(reqOptions.jsonBody)
			if err != nil {
				return nil, err
			}

			return bytes.NewReader(b), nil
		}

		return nil, nil
	}()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	for key, values := range e.options.header {
		for _, v := range values {
			req.Header.Set(key, v)
		}
	}

	for key, values := range reqOptions.header {
		for _, v := range values {
			req.Header.Set(key, v)
		}
	}

	resp, err := e.doer.Do(req)
	if err != nil {
		return nil, err
	}

	finalResp, err := readResponse(resp)
	if err != nil {
		return nil, err
	}

	if len(finalResp.Body) > 0 {
		if reqOptions.successJSONReceiver != nil && finalResp.StatusCode < 400 {
			err = json.Unmarshal(finalResp.Body, reqOptions.successJSONReceiver)
		} else if reqOptions.errorJSONReceiver != nil && finalResp.StatusCode >= 400 {
			err = json.Unmarshal(finalResp.Body, reqOptions.errorJSONReceiver)
		}
		if err != nil {
			return nil, err
		}
	}

	return finalResp, err
}

func New(urlTemplate string, opts ...endpointOptionSetter) *Endpoint {
	options := defaultEndpointOptions()

	for _, opt := range opts {
		opt(&options)
	}

	return &Endpoint{
		urlTemplate: urlTemplate,
		doer:        options.doer,
		options:     &options,
	}
}

func defaultEndpointOptions() endpointOptions {
	return endpointOptions{
		doer:   &http.Client{},
		header: make(http.Header),
	}
}

func defaultRequestOptions() requestOptions {
	return requestOptions{
		pathParams:          make(map[string]any),
		header:              make(http.Header),
		successJSONReceiver: nil,
		errorJSONReceiver:   nil,
		jsonBody:            nil,
		rawBody:             nil,
	}
}
