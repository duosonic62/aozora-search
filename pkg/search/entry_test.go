package search

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindEntries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))
	defer ts.Close()

	tmp := pageUrlFormat
	pageUrlFormat = ts.URL + "/cards/%s/card%s.html"
	defer func() {
		pageUrlFormat = tmp
	}()

}
