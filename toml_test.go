package dechoconf

import "testing"

type app struct {
	_       string `prefix:"app"`
	Address string
	Mysql   string
}

type backendAPI struct {
	_   string `prefix:"backend_api"`
	URL string
}

func TestMultiDecode(t *testing.T) {
	tomlData := `
[app]
address = ":1323"
mysql = "root:root@tcp(127.0.0.1:3306)/decho"

[backend_api]
url = "https://localhost:8080"
`

	var appConfig app
	var apiConfig backendAPI
	if err := DecodeToml(tomlData, &appConfig, &apiConfig); err != nil {
		t.Error(err)
	}

	if appConfig.Address != ":1323" {
		t.Errorf("Invalid address: %s", appConfig.Address)
	}

	if apiConfig.URL != "https://localhost:8080" {
		t.Errorf("Invalid url: %s", appConfig.Address)
	}
}
