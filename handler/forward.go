package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jeffbmartinez/log"
)

func Forward(defaultEndoint string) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		newUrl, err := url.Parse(defaultEndoint + vars["pathname"])
		if err != nil {
			log.Errorf("Couldn't parse url (%v): %v", defaultEndoint, err)
			return
		}

		newRequest := &http.Request{URL: newUrl}

		r, err := http.DefaultClient.Do(newRequest)
		if err != nil {
			log.Errorf("Had a problem with a response: %v", err)
			return
		}

		responseBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Had a problem reading response body: %v", err)
			return
		}

		response.WriteHeader(r.StatusCode)
		response.Write(responseBody)
	}
}

func ForwardTo(theUrl string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		newUrl, err := url.Parse(theUrl)
		if err != nil {
			log.Errorf("Couldn't parse url (%v): %v", theUrl, err)
			return
		}

		newRequest := &http.Request{URL: newUrl}

		r, err := http.DefaultClient.Do(newRequest)
		if err != nil {
			log.Errorf("Had a problem with a response: %v", err)
			return
		}

		responseBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Had a problem reading response body: %v", err)
			return
		}

		response.WriteHeader(r.StatusCode)
		response.Write(responseBody)
	}
}
