package dechoconf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// PrefixTagName defines the tagname used in "-" field for prefix
var PrefixTagName string = "prefix"

type decodeObjFunc func(val interface{}, obj interface{}) error

func decode(configs map[string]interface{}, decodeObj decodeObjFunc, objs ...interface{}) (err error) {
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

	for configPrefix, configVal := range configs {

	FINDOBJ:
		for obj, objPrefix := range objToPrefix {
			if configPrefix == objPrefix {
				if err != decodeObj(configVal, obj) {
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
							if err != decodeObj(v, obj) {
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
