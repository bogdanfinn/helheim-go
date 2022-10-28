package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
		helheim_go.WithBase64Response(),
		// add here other options like proxy or debug
		// helheim_go.WithDebug(),
		// helheim_go.WithProxyUrl("http://username:password@host:port"),
	}

	httpClient, err := helheimClient.NewHttpClient(options, helheimClientOptions...)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("created helheim http client")

	req, err := http.NewRequest(http.MethodGet, "http://i.imgur.com/m1UIjW1.jpg", nil)

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

	cookie := helheim_go.SessionCookie{
		Name:  "myTestCookie",
		Value: "myTestValue",
		//	Domain:  "",
		//	Path:    "",
		//	Expires: 0,
	}

	err = httpClient.SetCookie(cookie)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("http client cookies:")
	log.Println(httpClient.GetSessionCookies())

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("request response status code:")
	log.Println(resp.StatusCode)

	defer resp.Body.Close()

	//open a file for writing
	file, err := os.Create("./asdf.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")

	return
}
