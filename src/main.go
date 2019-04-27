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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func logSetup() {
	log.Printf("Server will run \n")
}

func logRequestPayload(URLs []string) {
    proxyUrls := strings.Trim(strings.Join(strings.Split(fmt.Sprint(URLs), " "), ", "), "[]")
	log.Printf("Multicast urls: %s\n", proxyUrls)
}

func getMulticastHosts() []string {
    var hosts []string
    if "" != getEnv("MULTICAST_HTTP_HOSTS", "") {
        definedHosts := getEnv("MULTICAST_HTTP_HOSTS", "")
        hosts = append(hosts,  strings.Split(definedHosts, ",")...)
    }
    return hosts
}

func transmit(target string, res http.ResponseWriter, req *http.Request) {

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

    log.Printf("Transmited %s\n", proxyUrl)
    if "0" != getEnv("MULTICAST_HTTP_DEBUG", "0") {
        requestDump, _ := httputil.DumpRequest(proxyReq, true)
        log.Printf(string(requestDump))
    }
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
    hosts := getMulticastHosts()
    logRequestPayload(hosts)

    if len(hosts) > 0 {
        for i := 0; i < len(hosts); i++ {
            go transmit(hosts[i], res, req)
        }
    } else {
        log.Printf("List of hosts to multicast is empty.")
    }


    if "0" != getEnv("MULTICAST_HTTP_DEBUG", "0") {
        requestDump, _ := httputil.DumpRequest(req, true)
        log.Printf(string(requestDump))
    }

    fmt.Fprintf(res, "Ok")
}

func main() {
	logSetup()

	// start server on port 80
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}