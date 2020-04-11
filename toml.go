package dechoconf

import (
	"io"

	"github.com/BurntSushi/toml"
)

func decodeToml(data string, obj interface{}) (err error) {
	_, err = toml.Decode(data, obj)
	return
}

func getTomlEncoder(w io.Writer) encoder {
	return toml.NewEncoder(w)
}

// TOMLConf is the ConfCoder instance handling toml configs
var TOMLConf = NewTOMLConf(defaultPrefixTagName)

// NewTOMLConf return is the ConfCoder instance handling toml configs with given prefixTagName
func NewTOMLConf(prefixTagName string) *ConfCoder {
	return &ConfCoder{
		prefixTagName,
		decodeToml,
		getTomlEncoder,
	}
}
