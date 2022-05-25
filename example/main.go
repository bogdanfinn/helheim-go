package main

import (
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

	//wokouResp, err = session.Wokou("chrome")

	reqOpts := helheim_go.RequestOptions{
		Method:  http.MethodGet,
		Url:     "https://www.genx.co.nz/iuam/",
		Options: make(map[string]string),
	}

	//
	//exampleForPost := helheim_go.RequestOptions{
	//	Method: http.MethodPost,
	//	Url:    "https://www.genx.co.nz/iuam/",
	//	Options: map[string]string{
	//		"data": `{"key":"myBodyPayload"}`, // set body here depends on your content type.
	//		"json": `{"key":"myBodyPayload"}`, // or set body here depends on your content type.
	//	},
	//}

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
