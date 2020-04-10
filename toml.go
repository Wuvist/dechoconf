package dechoconf

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

func encodeToml(obj interface{}) (result string, err error) {
	var buf bytes.Buffer

	e := toml.NewEncoder(&buf)
	err = e.Encode(obj)
	if err != nil {
		return
	}

	return buf.String(), nil
}

func decodeObjToml(tomlVal interface{}, obj interface{}) error {
	// todo: re-encode & decode is kind of stupid, but works for now
	tomlString, err := encodeToml(tomlVal)
	if err != nil {
		return err
	}
	_, err = toml.Decode(tomlString, obj)
	return err
}

// DecodeToml accept toml data string, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeToml(data string, objs ...interface{}) (err error) {
	configs := make(map[string]interface{})
	if _, err = toml.Decode(data, &configs); err != nil {
		return err
	}

	return decode(configs, decodeObjToml, objs...)
}

// DecodeTomlFile accept path of a toml file, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeTomlFile(path string, objs ...interface{}) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return DecodeToml(string(data), objs...)
}
