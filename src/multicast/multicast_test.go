package multicast

import (
	"testing"
	"fmt"
	"os"
)

func TestHostList(t *testing.T) {
    port := "https://ya.ru, goole.com, 127.0.0.1:5000"
    os.Setenv("MULTICAST_HTTP_HOSTS", port)
    hosts := hostList()
    if len(hosts) != 3 {
        t.Error(fmt.Sprintf("Expected 3 hosts, actual %d", len(hosts)))
    }
    os.Unsetenv("MULTICAST_HTTP_HOSTS")
}

func TestGetEnv(t *testing.T) {
    key := "TEST"
    value := "value"
    os.Setenv(key, value)

    if getEnv(key, "exist") == "exist" {
        t.Error(fmt.Sprintf("Expected %s", value))
    }
    if getEnv("NO_EXIST", value) != value {
        t.Error(fmt.Sprintf("Expected %s", value))
    }
    os.Unsetenv(key)
}