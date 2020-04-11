package dechoconf

import (
	"io/ioutil"
	"log"
	"os"
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

func TestException(t *testing.T) {
	tomlData := `
	[db]
	host = "localhost"
	port = 3306
	username = "root"
	password = ""

	[api]
	url = "https://localhost:8080"
	`

	type config struct {
		URL string
	}

	var apiConfig config
	err := TOMLConf.Decode(tomlData, &apiConfig)
	if err == nil || err.Error() != "No `-` field is found on struct: config" {
		t.Error(err)
	}
}

func TestWrap(t *testing.T) {
	file, err := ioutil.TempFile("", "example.*.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

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

	file.WriteString(tomlData)
	file.Close()

	var dbConfig dbWrap
	if err := DecodeFile(file.Name(), &dbConfig); err != nil {
		t.Error(err)
	}

	if dbConfig.Host != "localhost" {
		t.Errorf("Invalid host: %s", dbConfig.Host)
	}

	if dbConfig.Port != 3306 {
		t.Errorf("Invalid port: %d", dbConfig.Port)
	}
}

func TestDuplicate(t *testing.T) {
	type dbWrap struct {
		_ string `prefix:"db"`
	}
	var dbConfig DBConfig
	var apiConfig dbWrap
	err := TOMLConf.Decode("", &dbConfig, &apiConfig)

	if err.Error() != "Duplicated prefix db on struct DBConfig and dbWrap" {
		t.Error("Failed to detect duplicate prefix")
	}
}

func TestMissing(t *testing.T) {
	var dbConfig DBConfig
	var apiConfig APIConfig

	err := TOMLConf.DecodeFile("", &dbConfig, &apiConfig)
	if err == nil {
		t.Error("Failed to detect missing file")
	}

	tomlData := `
[db]
host = "localhost"
port = 3306
username = "root"
password = ""

`

	err = TOMLConf.Decode(tomlData, &dbConfig, &apiConfig)

	if dbConfig.Host != "localhost" {
		t.Errorf("Invalid host: %s", dbConfig.Host)
	}

	if dbConfig.Port != 3306 {
		t.Errorf("Invalid port: %d", dbConfig.Port)
	}

	if err.Error() != "No config found for: type APIConfig with prefix [api]\n" {
		t.Error("Fail to detect missing config")
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
	if err := TOMLConf.Decode(tomlData, &dbConfig, &apiConfig); err != nil {
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
	if err := TOMLConf.Decode(tomlData, &dbConfig, &apiConfig, &authConfig); err != nil {
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
