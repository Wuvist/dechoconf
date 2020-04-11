package dechoconf

import (
	"io"

	"gopkg.in/yaml.v2"
)

func decodeYaml(data string, obj interface{}) (err error) {
	return yaml.Unmarshal([]byte(data), obj)
}

func getYamlEncoder(w io.Writer) encoder {
	return yaml.NewEncoder(w)
}

// YAMLConf is the ConfCoder instance handling yaml configs
var YAMLConf = NewYAMLConf(defaultPrefixTagName)

// NewYAMLConf return is the ConfCoder instance handling yaml configs with given prefixTagName
func NewYAMLConf(prefixTagName string) *ConfCoder {
	return &ConfCoder{
		prefixTagName,
		decodeYaml,
		getYamlEncoder,
	}
}
