package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jeffbmartinez/log"
)

func Forward(defaultEndoint string) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		forwardUrl := defaultEndoint + request.URL.String()
		forwardRequest(response, request, forwardUrl)
	}
}

func ForwardTo(theUrl string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		forwardRequest(response, request, theUrl)
	}
}

func forwardRequest(response http.ResponseWriter, request *http.Request, forwardedUrl string) {
	newUrl, err := url.Parse(forwardedUrl)
	if err != nil {
		log.Errorf("Couldn't parse url (%v): %v", forwardedUrl, err)
		return
	}

	newRequest := &http.Request{URL: newUrl}

	intermediateResponse, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		log.Errorf("Had a problem with a response: %v", err)
		return
	}

	responseBody, err := ioutil.ReadAll(intermediateResponse.Body)
	if err != nil {
		log.Errorf("Had a problem reading response body: %v", err)
		return
	}

	response.WriteHeader(intermediateResponse.StatusCode)
	response.Write(responseBody)
}
