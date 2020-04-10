package dechoconf

import (
	"bytes"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func encodeYaml(obj interface{}) (result []byte, err error) {
	var buf bytes.Buffer

	e := yaml.NewEncoder(&buf)
	err = e.Encode(obj)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

func decodeObjYaml(yamlVal interface{}, obj interface{}) error {
	// todo: re-encode & decode is kind of stupid, but works for now
	yamlBytes, err := encodeYaml(yamlVal)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlBytes, obj)
}

// DecodeYaml accept yaml data bytes, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeYaml(data []byte, objs ...interface{}) (err error) {
	configs := make(map[string]interface{})
	if err = yaml.Unmarshal(data, &configs); err != nil {
		return err
	}

	return decode(configs, decodeObjYaml, objs...)
}

// DecodeYamlFile accept path of a yaml file, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeYamlFile(path string, objs ...interface{}) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return DecodeYaml(data, objs...)
}
