package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jeffbmartinez/log"
)

// Forward returns a handler function which forwards the unmodified request
// to the domain in the param.
func Forward(domain string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		forwardURL := domain + request.URL.String()
		forwardRequest(response, request, forwardURL)
	}
}

// ForwardTo eturns a handler function which intercepts a request and
// replaces what would have been the original request's response with the
// response from an entirely different request at given url. This can be
// a different domain and endpoint entirely.
func ForwardTo(theURL string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		forwardRequest(response, request, theURL)
	}
}

func forwardRequest(response http.ResponseWriter, request *http.Request, forwardedURL string) {
	newURL, err := url.Parse(forwardedURL)
	if err != nil {
		log.Errorf("Couldn't parse url (%v): %v", forwardedURL, err)
		return
	}

	newRequest := &http.Request{
		URL:    newURL,
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
