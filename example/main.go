package main

import (
	"fmt"
	helheim_go "github.com/bogdanfinn/helheim-go"
	"log"
	"net/http"
)

const YourApiKey = "INSERT_HERE"

func main() {
	helheimClient, err := helheim_go.ProvideClient(YourApiKey, false, true, nil)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("helheim client initiated")

	versionResponse, err := helheimClient.Version()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(fmt.Sprintf("you are using helheim version %s", versionResponse.Version))

	options := helheim_go.CreateSessionOptions{
		Browser: helheim_go.BrowserOptions{
			Browser:  "chrome",
			Mobile:   false,
			Platform: "windows",
		},
		Captcha: helheim_go.CaptchaOptions{
			Provider: "vanaheim",
		},
	}
	session, err := helheimClient.NewSession(options)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("session response:")
	log.Println(session)

	balance, err := helheimClient.GetBalance()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("balance response:")
	log.Println(balance)

	reqOpts := helheim_go.RequestOptions{
		Method:  http.MethodGet,
		Url:     "https://www.genx.co.nz/iuam/",
		Options: make(map[string]string),
	}

	/*
	exampleForPost := helheim_go.RequestOptions{
		Method: http.MethodPost,
		Body:   `{"key":"myBodyPayload"}`, //either set body here
		Url:    "https://www.genx.co.nz/iuam/",
		Options: map[string]string{
			"body": `{"key":"myBodyPayload"}`, // or set body here.
		},
	}
	*/

	resp, err := session.Request(reqOpts)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("request response status code:")
	log.Println(resp.Response.StatusCode)

	log.Println("request response:")
	log.Println(resp)

	err = session.Delete()

	if err != nil {
		log.Println(err)
		return
	}
}
