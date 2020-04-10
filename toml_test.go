package dechoconf

import (
	"testing"
)

type DBConfig struct {
	_        string `prefix:"db"`
	Host     string
	Port     int
	Username string
	Password string
}

type APIConfig struct {
	_   string `prefix:"api"`
	URL string
}

type BackendAuthConfig struct {
	_     string `prefix:"service.backend.auth"`
	Token string
}

type BackendAPIConfig struct {
	_   string `prefix:"service.backend.api"`
	URL string
}

func TestWrap(t *testing.T) {
	type db struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	type dbWrap struct {
		_ string `prefix:"db"`
		db
	}

	tomlData := `
[db]
host = "localhost"
port = 3306
username = "root"
password = ""
`

	var dbConfig dbWrap
	if err := DecodeToml(tomlData, &dbConfig); err != nil {
		t.Error(err)
	}

	if dbConfig.Host != "localhost" {
		t.Errorf("Invalid host: %s", dbConfig.Host)
	}

	if dbConfig.Port != 3306 {
		t.Errorf("Invalid port: %d", dbConfig.Port)
	}
}

func TestMultiDecode(t *testing.T) {
	tomlData := `
[db]
host = "localhost"
port = 3306
username = "root"
password = ""

[api]
url = "https://localhost:8080"
`

	var dbConfig DBConfig
	var apiConfig APIConfig
	if err := DecodeToml(tomlData, &dbConfig, &apiConfig); err != nil {
		t.Error(err)
	}

	if dbConfig.Host != "localhost" {
		t.Errorf("Invalid host: %s", dbConfig.Host)
	}

	if dbConfig.Port != 3306 {
		t.Errorf("Invalid port: %d", dbConfig.Port)
	}

	if apiConfig.URL != "https://localhost:8080" {
		t.Errorf("Invalid url: %s", apiConfig.URL)
	}
}

func TestMultiPrefixDecode(t *testing.T) {
	tomlData := `
[db]
host = "localhost"
port = 3306
username = "root"
password = ""

[service]
	[service.backend]
		[service.backend.api]
		url = "https://localhost:8080"

		[service.backend.auth]
		token = "I love animal crossing!"
`

	var dbConfig DBConfig
	var apiConfig BackendAPIConfig
	var authConfig BackendAuthConfig
	if err := DecodeToml(tomlData, &dbConfig, &apiConfig, &authConfig); err != nil {
		t.Error(err)
	}

	if dbConfig.Host != "localhost" {
		t.Errorf("Invalid host: %s", dbConfig.Host)
	}

	if dbConfig.Port != 3306 {
		t.Errorf("Invalid port: %d", dbConfig.Port)
	}

	if apiConfig.URL != "https://localhost:8080" {
		t.Errorf("Invalid url: %s", apiConfig.URL)
	}

	if authConfig.Token != "I love animal crossing!" {
		t.Errorf("Invalid token: %s", authConfig.Token)
	}
}
