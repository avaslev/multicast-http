package resolver

import (
    "encoding/json"
    "fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"os"
	"errors"
)

const K8S_POD_SCHEMA string = "http"
const TOKEN_FILE string = "/var/run/secrets/kubernetes.io/serviceaccount/token"
// const K8S_POD_LABEL string = "app: symfony"

type podListResponse struct {
	Kind string
	Items []item
}

type item struct {
	Metadata metadata
	Status status
}

type metadata struct {
	Name string
	Labels map[string]string
}

type status struct {
	PodIP string
	Phase string
}

func Resolve(label string) []string  {
	k8sMasterHost := "https://" + getCriticalEnv("KUBERNETES_SERVICE_HOST")
	k8sMasterHost += ":" + getCriticalEnv("KUBERNETES_SERVICE_PORT_HTTPS")
	return handlePodList(podList(k8sMasterHost), getCriticalEnv("MULTICAST_HTTP_K8S_POD_LABEL"))
}

func podList(k8sMasterHost string) []byte {
	var request *http.Request
	request, _ = http.NewRequest("GET", k8sMasterHost + "/api/v1/pods", nil )
	request.Header.Add("Authorization", "Bearer " + getCriticalEnv("MULTICAST_HTTP_K8S_TOKEN") )

	httpClient := http.Client{}
    response, err := httpClient.Do(request)
    if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var bodyBytes []byte

	if response.StatusCode != http.StatusOK {
		log.Printf("K8s response with %d code", response.StatusCode)
		return bodyBytes
	}

	bodyBytes, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(err)
	}

	return bodyBytes
}

func handlePodList (body []byte, label string) []string {
	podListResponse := podListResponse{}
    if err := json.Unmarshal(body, &podListResponse); err != nil {
		log.Println(err)
	}

	var hosts []string
	if (podListResponse.Kind != "PodList") {
		log.Printf("Unexpected kind: %s", podListResponse.Kind)
		return hosts
	}

	if len(podListResponse.Items) == 0 {
		return hosts
	}
	
	
	s := strings.Split(label, ":")
	key, value := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
	for i := 0; i < len(podListResponse.Items); i++ {
		item := podListResponse.Items[i]
		fmt.Println(item)
		if item.Status.Phase == "Running" && contains(item.Metadata.Labels, key, value) {
			hosts = append(hosts, fmt.Sprintf("%s://%s:%s", K8S_POD_SCHEMA, item.Status.PodIP, getCriticalEnv("MULTICAST_HTTP_K8S_POD_PORT")))
		}
	}
	return hosts
} 

func contains(s map[string]string, key string, value string) bool {
	if mapValue, ok := s[key]; ok && mapValue == value {    
        return true
    }
    return false
}

func getCriticalEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if debug, _ := os.LookupEnv("MULTICAST_HTTP_DEBUG"); debug != "0" {
		log.Printf(fmt.Sprintf("Variable %s not exist", key))
		return ""
	}

	err := errors.New(fmt.Sprintf("Variable %s not exist", key))
	log.Fatal(err)

	return ""
}

func getToken() string {
	if token, err := ioutil.ReadFile(TOKEN_FILE); err != nil {
		log.Fatal(err)
		return ""
	} else {
		return string(token)
	}
	
}