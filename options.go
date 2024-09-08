package endpoint

import "net/http"

// -- endpoint --
type endpointOptions struct {
	doer   Doer
	header http.Header
}

type endpointOptionSetter func(*endpointOptions)

func WithClient(client Doer) endpointOptionSetter {
	return func(o *endpointOptions) {
		o.doer = client
	}
}

func WithHeader(name string, values ...string) endpointOptionSetter {
	return func(o *endpointOptions) {
		for _, v := range values {
			o.header.Add(name, v)
		}
	}
}

// -- request --
type requestOptions struct {
	pathParams          map[string]any
	header              http.Header
	successJSONReceiver any
	errorJSONReceiver   any
	jsonBody            any
	rawBody             []byte
}

type requestOptionSetter func(*requestOptions)

func WithPathParam(key string, value any) requestOptionSetter {
	return func(o *requestOptions) {
		o.pathParams[key] = value
	}
}

func WithJSONReceivers(successReceiver, errorReceiver any) requestOptionSetter {
	return func(o *requestOptions) {
		o.successJSONReceiver = successReceiver
		o.errorJSONReceiver = errorReceiver
	}
}

func WithRequestHeader(name string, values ...string) requestOptionSetter {
	return func(o *requestOptions) {
		for _, v := range values {
			o.header.Add(name, v)
		}
	}
}

func WithJSONBody(v any) requestOptionSetter {
	return func(o *requestOptions) {
		o.jsonBody = v
	}
}

func WithBody(body []byte) requestOptionSetter {
	return func(o *requestOptions) {
		o.rawBody = body
	}
}
