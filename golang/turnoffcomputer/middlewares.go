package main

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

type Authentication struct {
	Password string `json:"password"`
}

func CheckAuthentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading body", 400)
			return
		}

		auth := Authentication{}
		unmarshalErr := json.Unmarshal(body, &auth)
		if unmarshalErr != nil {
			http.Error(w, "error decoding json", 400)
			return
		}

		if auth.Password != viper.Get("app.password") {
			http.Error(w, http.StatusText(403), 403)
			return
		}
		handler.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
