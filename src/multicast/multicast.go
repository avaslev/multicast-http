package multicast

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
    "fmt"
    "github.com/avaslev/multicast-http/src/resolver"
)

func HandleRequest(res http.ResponseWriter, req *http.Request) {
	hosts := hostList()

	log.Printf("Multicast urls: %s\n", strings.Trim(strings.Join(strings.Split(fmt.Sprint(hosts), " "), ", "), "[]"))

    if len(hosts) > 0 {
        for i := 0; i < len(hosts); i++ {
            go transmit(hosts[i], res, req)
        }
    } else {
        log.Printf("List of hosts to multicast is empty.")
    }


    if isDebugMode() {
        requestDump, _ := httputil.DumpRequest(req, true)
        log.Printf(string(requestDump))
    }
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

    log.Printf("Transmited %s, response status code: %d\n", proxyUrl, resp.StatusCode)
    if isDebugMode() {
        requestDump, _ := httputil.DumpRequest(proxyReq, true)
        log.Printf(string(requestDump))
    }
}

func hostList() []string {
    var hosts []string

    if definedHosts := getEnv("MULTICAST_HTTP_HOSTS", ""); definedHosts != "" {
		for _, v := range strings.Split(definedHosts, ",") {
			hosts = append(hosts, strings.TrimSpace(v))
		}
    }

    if podLabel := getEnv("MULTICAST_HTTP_K8S_POD_LABEL", ""); podLabel != "" {
        hosts = append(hosts, resolver.Resolve(podLabel)...)
    }

    return hosts
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// "0" is false
func isDebugMode() bool {
	return "0" != getEnv("MULTICAST_HTTP_DEBUG", "0")
}