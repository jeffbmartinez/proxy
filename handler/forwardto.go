package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jeffbmartinez/log"
)

func ForwardTo(url string) {
	return func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		request.URL = url
		r, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Errorf("Had a problem with a response: %v", err)
			return
		}

		responseBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Had a problem reading response body: %v", err)
			return
		}

		fmt.Sprintf(response, string(responseBody))
	}
}
