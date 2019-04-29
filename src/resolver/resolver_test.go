package resolver

import (
	"testing"
	"fmt"
	"os"
	"path"
)

func TestHandlePodList(t *testing.T) {
	port := "1234"
	os.Setenv("MULTICAST_HTTP_K8S_POD_PORT", port)
	responseBody := []byte(`{
		"kind": "PodList",
		"items": [
			{
				"metadata": {
					"labels": {
						"app": "symfony",
						"pod-template-hash": "dc998467d"
					}
				},
				"status": {
					"phase": "Running",
					"podIP": "10.244.0.72"
				}
			}
		]
	}`)

	hosts := handlePodList(responseBody, "app: symfony")
	if len(hosts) != 1 {
		t.Error(fmt.Sprintf("Expected 1 host, actual %d", len(hosts)))
	}

	expected := fmt.Sprintf("%s://%s:%s", K8S_POD_SCHEMA, "10.244.0.72", port)
	if hosts[0] !=  expected {
        t.Error(fmt.Sprintf("Expected %s, actual %s", expected, hosts[0]))
	}
	os.Unsetenv("MULTICAST_HTTP_K8S_POD_PORT")
}

func TestContains(t *testing.T) {
	elements := map[string]string{
		"H": "Hydrogen",
		"He": "Helium",
		"Li": "Lithium",
		"Be": "Beryllium",
	}

	if !contains(elements, "H", "Hydrogen") {
        t.Error("Expected  true")
	}
	
	if contains(elements, "Li", "Hydrogen") {
        t.Error("Expected  false")
	}
	
	if contains(elements, "H", "Lithium") {
        t.Error("Expected  false")
    }
}

func TestGetToken(t *testing.T) {
	token := "secret"
	os.MkdirAll(path.Dir(TOKEN_FILE), os.ModePerm)
	file, _ := os.Create(TOKEN_FILE)
	defer file.Close()
	file.WriteString(token)

	resultToken := getToken()
	if resultToken != token {
		t.Error(fmt.Sprintf("Expected %s, actual %s", token, resultToken))
	}

	os.Remove(TOKEN_FILE)
}