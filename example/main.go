package main

import (
	"log"
	"net/http"

	helheim_go "github.com/bogdanfinn/helheim-go"
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

	// you currently need to use a UA of chrome 100+ according to venom
	headerResp, err := session.SetHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("set headers response:")
	log.Println(headerResp)

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
