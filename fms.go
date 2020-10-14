package fms

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Generate is a function that generates the dom using LINE's Flex Message Simulator
func Generate(json io.Reader) (string, error) {

	apiURL := "https://developers.line.biz/api/v1/fx/render"
	extractor := cookieExtractor{
		Email:    os.Getenv("EMAIL"),
		Password: os.Getenv("PASSWORD"),
	}

	// get the cookies
	cookie, err := extractor.getCookie()
	if err != nil {
		return "", err
	}

	// create new request
	req, err := http.NewRequest(
		"POST",
		apiURL,
		json,
	)
	if err != nil {
		return "", err
	}

	// set the headers
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// read the response
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
