package dechoconf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

var defaultPrefixTagName string = "prefix"

type decodeObjFunc func(val interface{}, obj interface{}) error

type encoder interface {
	Encode(v interface{}) error
}

// ConfCoder accepts prefixTagName / decode / getEncoder func to make them support mulitple decode
// according to their prefix tag defined in "-" field
// Default prefixTagName is "prefix"
type ConfCoder struct {
	prefixTagName string
	decode        func(string, interface{}) error
	getEncoder    func(io.Writer) encoder
}

func (c *ConfCoder) encode(obj interface{}) (result string, err error) {
	var buf bytes.Buffer

	err = c.getEncoder(&buf).Encode(obj)
	if err != nil {
		return
	}

	return buf.String(), nil
}

func (c *ConfCoder) redecode(val interface{}, obj interface{}) error {
	// todo: re-encode & decode is kind of stupid, but works for now
	data, err := c.encode(val)
	if err != nil {
		return err
	}
	return c.decode(data, obj)
}

// Decode data string, and unmarshal it to multiple structs
func (c *ConfCoder) Decode(data string, objs ...interface{}) (err error) {
	configs := make(map[string]interface{})
	if err = c.decode(data, &configs); err != nil {
		return err
	}

	return c.multiDecode(configs, objs...)
}

// DecodeFile accept path of a config file, and unmarshal it to multiple structs
func (c *ConfCoder) DecodeFile(path string, objs ...interface{}) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return c.Decode(string(data), objs...)
}

func (c *ConfCoder) checkPrefix(objs ...interface{}) (prefixToObj map[string]interface{}, objToPrefix map[interface{}]string, err error) {
	prefixToObj = make(map[string]interface{}, len(objs))
	objToPrefix = make(map[interface{}]string, len(objs))
	for _, obj := range objs {
		tt := reflect.ValueOf(obj).Elem().Type()
		prefixTag, found := tt.FieldByName("_")
		if !found {
			return nil, nil, errors.New("No `-` field is found on struct: " + tt.Name())
		}
		prefix := prefixTag.Tag.Get(c.prefixTagName)

		if o, found := prefixToObj[prefix]; found {
			return nil, nil, fmt.Errorf("Duplicated prefix %s on struct %s and %s", prefix,
				reflect.ValueOf(o).Elem().Type().Name(), tt.Name())
		}
		prefixToObj[prefix] = obj
		objToPrefix[obj] = prefix
	}
	return
}

func (c *ConfCoder) multiDecode(configs map[string]interface{}, objs ...interface{}) (err error) {
	prefixToObj, objToPrefix, err := c.checkPrefix(objs...)
	if err != nil {
		return err
	}

	for configPrefix, configVal := range configs {

	FINDOBJ:
		for obj, objPrefix := range objToPrefix {
			if configPrefix == objPrefix {
				if err = c.redecode(configVal, obj); err != nil {
					return err
				}

				delete(prefixToObj, objPrefix)
				continue FINDOBJ
			} else if strings.HasPrefix(objPrefix, configPrefix+".") {
				currentPrefix := configPrefix + "."
				val, ok := configVal.(map[string]interface{})

			MATCHFIELD:
				for ok {
					ok = false

					for k, v := range val {
						if objPrefix == currentPrefix+k {
							if err = c.redecode(v, obj); err != nil {
								return err
							}

							delete(prefixToObj, objPrefix)
							continue FINDOBJ
						} else {
							if strings.HasPrefix(objPrefix, currentPrefix+k+".") {
								val, ok = v.(map[string]interface{})
								currentPrefix = currentPrefix + k + "."
								continue MATCHFIELD
							}
						}
					}
				}
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

// DecodeFile accept path of toml / yaml config file, and unmarshal to multiple structs
// according to their prefix tag defined in "-" field
func DecodeFile(path string, objs ...interface{}) (err error) {
	if strings.HasSuffix(path, ".toml") {
		return DecodeTomlFile(path, objs...)
	}
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return DecodeYamlFile(path, objs...)
	}

	return fmt.Errorf("%s doesn't have toml / yaml file extensions", path)
}
