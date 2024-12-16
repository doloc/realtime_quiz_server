package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
)

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

func GenerateJWT(data interface{}, secretKey string) (string, error) {
	dataMap, err := ConvertStructToMap(data)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"data": dataMap,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
