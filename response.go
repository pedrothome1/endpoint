package endpoint

import (
	"io"
	"net/http"
)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func readResponse(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	finalResp := &Response{StatusCode: resp.StatusCode, Header: resp.Header.Clone()}

	if resp.StatusCode == http.StatusNoContent || resp.ContentLength == 0 {
		return finalResp, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	finalResp.Body = b
	return finalResp, nil
}
