package main

import (
	"io/ioutil"
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

	helheimClientOptions := []helheim_go.HttpClientOption{
		helheim_go.WithWokou("chrome"),
		// add here other options like proxy or debug
		// helheim_go.WithDebug(),
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

	// you currently need to use a UA of chrome 100+ according to venom
	for key, value := range map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
	} {
		req.Header.Set(key, value)
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
	log.Println(string(bodyResponse))

	log.Println("response header:")
	log.Println(resp.Header)
}
