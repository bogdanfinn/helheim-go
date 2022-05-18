package main

import (
	helheim_go "github.com/bogdanfinn/helheim-go"
	"io/ioutil"
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

	helheimClientOptions := []helheim_go.HttpClientOption{
		helheim_go.WithWokou("chrome"),
		// add here other options like proxy, debug, or bifrost
		// helheim_go.WithDebug(),
		// helheim_go.WithBifrost("path/to/lib"),
		// helheim_go.WithProxyUrl("http://username:password@host:port"),
	}

	httpCLient, err := helheimClient.NewHttpClient(options, helheimClientOptions...)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("created helheim http client")

	req, err := http.NewRequest(http.MethodGet, "https://www.genx.co.nz/iuam/", nil)

	if err != nil {
		log.Println(err)
		return
	}

	resp, err := httpCLient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("request response status code:")
	log.Println(resp.StatusCode)

	defer resp.Body.Close()

	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("request response:")
	log.Println(bodyResponse)

	log.Println("response header:")
	log.Println(resp.Header)
}
