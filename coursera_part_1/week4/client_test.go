package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCases struct {
	Req SearchRequest
	Err bool
}

func UnknownJsonErrReq(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": "11111111"}`))
}

func InternalErr(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TimeOut(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 10)
}

func BadJson(w http.ResponseWriter, r *http.Request) {
	order_field := r.FormValue("order_field")
	if order_field != "id" && order_field != "age" && order_field != "name" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "ErrorBadOrderField"`))
		return
	}
	w.Write([]byte(`{"error"`))
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	order_field := r.FormValue("order_field")
	if order_field != "id" && order_field != "age" && order_field != "name" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "ErrorBadOrderField"}`))
		return
	}
	token := r.Header.Get("AccessToken")
	if token != "rand_string" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(`[{"name":"bob"}, {"name": "ja"}]`))
}

func TestSearchServer(t *testing.T) {

	cases := []TestCases{
		{
			Req: SearchRequest{Limit: -1}, //bad limit
			Err: true,
		},
		{
			Req: SearchRequest{Limit: 26}, //big limit
			Err: true,
		},
		{
			Req: SearchRequest{Offset: -1}, //bad offset
			Err: true,
		},
		{
			Req: SearchRequest{Offset: 10, Limit: 10}, //bad json
			Err: true,
		},
		{
			Req: SearchRequest{Offset: 10, Limit: 10, OrderField: "age"}, //bad json with more limit, than data
			Err: false,
		},
		{
			Req: SearchRequest{Offset: 10, Limit: 1, OrderField: "age"}, //bad json with less limit, than data
			Err: false,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for _, testCase := range cases {
		c := &SearchClient{
			URL:         ts.URL,
			AccessToken: "rand_string",
		}
		_, err := c.FindUsers(testCase.Req)
		if err != nil && !testCase.Err {
			t.Errorf("[%s]", err)
		}
		if err == nil && testCase.Err {
			t.Errorf("[%s]", err)
		}
	}

	ts.Close()

}

func TestWrongUrl(t *testing.T) {
	c := &SearchClient{
		URL:         "32323",
		AccessToken: "rand_string",
	}
	_, err := c.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Client passed error URL value")
	}
}

func TestWrongAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "rand_stringdddd",
	}
	_, err := c.FindUsers(SearchRequest{Query: "alice", OrderField: "age"})
	if err == nil {
		t.Errorf("Client passed error Token value")
	}

	ts.Close()
}

func TestTimeOut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(TimeOut))

	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "rand_string",
	}
	_, err := c.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Client passed error URL value")
	}

}

func TestBadJson(t *testing.T) {
	cases := []TestCases{
		{
			Req: SearchRequest{Query: "Aloxa", OrderField: "Aaa"}, //bad json when 400
			Err: true,
		},
		{
			Req: SearchRequest{Query: "Aloxa", OrderField: "age"}, //bad json
			Err: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(BadJson))

	for _, testCase := range cases {
		c := &SearchClient{
			URL:         ts.URL,
			AccessToken: "rand_string",
		}
		_, err := c.FindUsers(testCase.Req)
		if err != nil && !testCase.Err {
			t.Errorf("[%s]", err)
		}
		if err == nil && testCase.Err {
			t.Errorf("[%s]", err)
		}
	}
	ts.Close()

	ts1 := httptest.NewServer(http.HandlerFunc(UnknownJsonErrReq))
	c1 := &SearchClient{
		URL:         ts1.URL,
		AccessToken: "rand_string",
	}
	_, err := c1.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Client passed unknown bad request")
	}
	ts1.Close()

	internalServ := httptest.NewServer(http.HandlerFunc(InternalErr))
	c2 := &SearchClient{
		URL:         internalServ.URL,
		AccessToken: "rand_string",
	}
	_, err = c2.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Client passed internal error")
	}
}
