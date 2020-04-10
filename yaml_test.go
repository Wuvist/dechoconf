package dechoconf

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestYaml(t *testing.T) {
	data := `
db:
  host: localhost
  port: 3306
  username: root
  password: ""

api:
  url: https://localhost:8080
`

	file, err := ioutil.TempFile("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(data)
	file.Close()

	var dbConfig DBConfig
	var apiConfig APIConfig
	if err := DecodeYamlFile(file.Name(), &dbConfig, &apiConfig); err != nil {
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
