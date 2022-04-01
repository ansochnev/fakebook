package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"fakebook/internal/backend"
)

type ShowProfile struct {
	Backend *backend.Backend
}

func (h ShowProfile) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	username := path.Base(req.URL.Path)

	userProfile, err := h.Backend.GetProfileByUsername(username)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	if userProfile == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	body, _ := json.MarshalIndent(userProfile, "", "  ")
	body = append(body, byte('\n'))

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(body)
}
