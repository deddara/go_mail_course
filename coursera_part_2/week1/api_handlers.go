package main
import "net/http"

// *MyApi
func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.ProfileParamsHandler(w, r)
	case "/user/create":
		h.CreateParamsHandler(w, r)
	default:
		return
	}
}

// *OtherApi
func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/create":
		h.OtherCreateParamsHandler(w, r)
	default:
		return
	}
}

func (h *MyApi) ProfileParamsHandler(w http.ResponseWriter, r *http.Request) {

func (h *MyApi) CreateParamsHandler(w http.ResponseWriter, r *http.Request) {

func (h *OtherApi) OtherCreateParamsHandler(w http.ResponseWriter, r *http.Request) {
