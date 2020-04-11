package dechoconf

import (
	"errors"
	"io"
	"testing"
)

func TestDecodeFile(t *testing.T) {
	err := DecodeFile("application.properties")
	if err.Error() != "application.properties doesn't have toml / yaml file extensions" {
		t.Error("Failed to detect unsupported file: application.properties")
	}
}

type mockEncoder struct{}

func (m *mockEncoder) Encode(v interface{}) error {
	return errors.New("MockEncoder can't encode")
}

func TestErrorHandling(t *testing.T) {
	coder := &ConfCoder{
		defaultPrefixTagName,
		func(data string, obj interface{}) (err error) {
			return errors.New("decode error")
		},
		getTomlEncoder,
	}

	var conf APIConfig
	err := coder.Decode("", &conf)
	if err.Error() != "decode error" {
		t.Error("Failed to handle decode error")
	}

	coder = &ConfCoder{
		defaultPrefixTagName,
		func(data string, obj interface{}) (err error) {
			configs := *obj.(*map[string]interface{})
			configs["api"] = ""

			return nil
		},
		func(w io.Writer) encoder {
			return &mockEncoder{}
		},
	}

	err = coder.Decode("", &conf)
	if err == nil {
		t.Error("Failed to detect mockEncoder")
	} else if err.Error() != "MockEncoder can't encode" {
		t.Error("Failed to handle decode error: " + err.Error())
	}

	var backendAPI BackendAPIConfig
	flag := false
	coder = &ConfCoder{
		defaultPrefixTagName,
		func(data string, obj interface{}) (err error) {
			if flag {
				return nil
			}
			configs := *obj.(*map[string]interface{})
			configs["service"] = map[string]interface{}{
				"backend": map[string]interface{}{
					"api": "",
				},
			}
			flag = true

			return nil
		},
		func(w io.Writer) encoder {
			return &mockEncoder{}
		},
	}

	err = coder.Decode("", &backendAPI)
	if err == nil {
		t.Error("Failed to detect mockEncoder")
	} else if err.Error() != "MockEncoder can't encode" {
		t.Error("Failed to handle decode error: " + err.Error())
	}
}
