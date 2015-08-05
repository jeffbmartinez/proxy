package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jeffbmartinez/log"
)

func Forward(defaultEndoint string) func(http.ResponseWriter, *http.Request) {
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

	newRequest := &http.Request{
		URL:    newUrl,
		Header: request.Header,
	}

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

	// There can be multiple values for a given header key. Here I am
	// clearing any values that may pre-exist in the header and replacing
	// them with the values from the response.
	for key, values := range intermediateResponse.Header {
		response.Header().Del(key)

		for _, value := range values {
			response.Header().Add(key, value)
		}
	}

	response.WriteHeader(intermediateResponse.StatusCode)
	response.Write(responseBody)
}
