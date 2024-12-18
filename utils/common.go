package utils

import "github.com/mitchellh/mapstructure"

// ConvertStruct converts a struct to another struct
func ConvertStruct(i interface{}, o interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           o,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(i); err != nil {
		return err
	}

	return nil
}

func ConvertStructToMap(i interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &result,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(i); err != nil {
		return nil, err
	}

	return result, nil
}

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
