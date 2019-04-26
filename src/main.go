package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
    "fmt"
)

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Log the env variables required for a reverse proxy
func logSetup() {
	log.Printf("Server will run \n")
}

// Log the typeform payload and redirect url
func logRequestPayload(URLs []string) {
    proxyUrls := strings.Trim(strings.Join(strings.Split(fmt.Sprint(URLs), " "), ", "), "[]")
	log.Printf("proxy_urls: %s\n", proxyUrls)
}

// Get the url for a given proxy condition
func getProxyUrls() []string {
    var URLs []string
    if "" != getEnv("MULTICAST_HTTP_HOSTS", "") {
        hosts := getEnv("MULTICAST_HTTP_HOSTS", "")
        URLs = append(URLs,  strings.Split(hosts, ",")...)
    }
    return URLs
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {

    // create a new url from the raw RequestURI sent by the client
    url, _ := url.Parse(target)
    proxyUrl := fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, req.RequestURI)
    if url.Scheme == "" {
        proxyUrl = fmt.Sprintf("http://%s%s", url, req.RequestURI)
    }

    proxyReq, _ := http.NewRequest(req.Method, proxyUrl, nil)
    if value, ok := os.LookupEnv("MULTICAST_HTTP_HEADER"); ok {
        proxyReq.Header.Set(value, getEnv("MULTICAST_HTTP_HEADER_VALUE", ""))
	}

    httpClient := http.Client{}
    resp, err := httpClient.Do(proxyReq)
    if err != nil {
        http.Error(res, err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    log.Printf("serve: %s\n", proxyUrl)
    if "0" != getEnv("MULTICAST_HTTP_HEADER", "0") {
        // Save a copy of this request for debugging.
        requestDump, _ := httputil.DumpRequest(proxyReq, true)
        log.Printf(string(requestDump))
    }
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
    URLs := getProxyUrls()
    logRequestPayload(URLs)

    if len(URLs) > 0 {
        for i := 0; i < len(URLs); i++ {
            go serveReverseProxy(URLs[i], res, req)
        }
    } else {
        log.Printf("List of hosts to multicast is empty.")
    }


    if "0" != getEnv("MULTICAST_HTTP_HEADER", "0") {
        // Save a copy of this request for debugging.
        requestDump, _ := httputil.DumpRequest(req, true)
        log.Printf(string(requestDump))
    }

    // serveReverseProxy(url, res, req)
    fmt.Fprintf(res, "Ok")
}

func main() {
	// Log setup values
	logSetup()

	// start server on port 80
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}