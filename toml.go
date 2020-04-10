package dechoconf

import (
	"bytes"
	"io/ioutil"
	"reflect"

	"github.com/BurntSushi/toml"
)

// PrefixTag defines the tagname used in "-" field for prefix
var PrefixTag string = "prefix"

func encodeToml(obj interface{}) (result string, err error) {
	var buf bytes.Buffer

	e := toml.NewEncoder(&buf)
	err = e.Encode(obj)
	if err != nil {
		return
	}

	return buf.String(), nil
}

// DecodeToml accept toml data string, and unmarshal it to multiple structs
// according to their prefix tag defined in "-" field
func DecodeToml(data string, objs ...interface{}) (err error) {
	configs := make(map[string]interface{})
	if _, err = toml.Decode(data, &configs); err != nil {
		return err
	}

	for k, v := range configs {
		for _, obj := range objs {
			tt := reflect.ValueOf(obj).Elem().Type()
			prefix, _ := tt.FieldByName("_")

			if prefix.Tag.Get(PrefixTag) == k {
				// todo: re-encode & decode is kind of stupid, but works for now
				tomlString, err := encodeToml(v)
				if err != nil {
					return err
				}
				_, err = toml.Decode(tomlString, obj)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
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
