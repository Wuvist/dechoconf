package dechoconf

import (
	"io"

	"github.com/BurntSushi/toml"
)

func decodeToml(data string, obj interface{}) (err error) {
	_, err = toml.Decode(data, obj)
	return
}

func encodeToml(w io.Writer) encodeWriter {
	return toml.NewEncoder(w)
}

var tomlCoder = &ConfCoder{
	defaultPrefixTagName,
	decodeToml,
	encodeToml,
}

// DecodeToml accept toml data string, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeToml(data string, objs ...interface{}) (err error) {
	return tomlCoder.Decode(data, objs...)
}

// DecodeTomlFile accept path of a toml file, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeTomlFile(path string, objs ...interface{}) (err error) {
	return tomlCoder.DecodeFile(path, objs...)
}
