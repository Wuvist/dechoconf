package dechoconf

import (
	"io"

	"gopkg.in/yaml.v2"
)

func decodeYaml(data string, obj interface{}) (err error) {
	return yaml.Unmarshal([]byte(data), obj)
}

func encodeYaml(w io.Writer) encodeWriter {
	return yaml.NewEncoder(w)
}

var yamlCoder = &ConfCoder{
	defaultPrefixTagName,
	decodeYaml,
	encodeYaml,
}

// DecodeYaml accept yaml data bytes, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeYaml(data string, objs ...interface{}) (err error) {
	return yamlCoder.Decode(data, objs...)
}

// DecodeYamlFile accept path of a yaml file, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeYamlFile(path string, objs ...interface{}) (err error) {
	return yamlCoder.DecodeFile(path, objs...)
}
