package dechoconf

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/BurntSushi/toml"
)

// PrefixTagName defines the tagname used in "-" field for prefix
var PrefixTagName string = "prefix"

func encodeToml(obj interface{}) (result string, err error) {
	var buf bytes.Buffer

	e := toml.NewEncoder(&buf)
	err = e.Encode(obj)
	if err != nil {
		return
	}

	return buf.String(), nil
}

func decodeObj(tomlVal interface{}, obj interface{}) error {
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

	prefixToObj := make(map[string]interface{}, len(objs))
	objToPrefix := make(map[interface{}]string, len(objs))
	for _, obj := range objs {
		tt := reflect.ValueOf(obj).Elem().Type()
		prefixTag, found := tt.FieldByName("_")
		if !found {
			return errors.New("No `-` field is found on struct: " + tt.Name())
		}
		prefix := prefixTag.Tag.Get(PrefixTagName)

		if o, found := prefixToObj[prefix]; found {
			return fmt.Errorf("Duplicated prefix %s on struct %s and %s", prefix,
				reflect.ValueOf(o).Elem().Type().Name(), tt.Name())
		}
		prefixToObj[prefix] = obj
		objToPrefix[obj] = prefix
	}

	for currentPrefix, v := range configs {
		for obj, objPrefix := range objToPrefix {
			if currentPrefix == objPrefix {
				if err != decodeObj(v, obj) {
					return err
				}

				delete(prefixToObj, objPrefix)
			}
		}
	}

	if len(prefixToObj) > 0 {
		msg := ""
		for k, o := range prefixToObj {
			msg += fmt.Sprintf("No config found for: type %s with prefix [%s]\n",
				reflect.ValueOf(o).Elem().Type().Name(), k)
		}
		return errors.New(msg)
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
