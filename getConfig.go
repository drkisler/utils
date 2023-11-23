package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type ConfigServ struct {
	FilePath string
	FileName string
	FileType string
}

func (s *ConfigServ) GetConfig(target interface{}) error {
	viper.AddConfigPath(s.FilePath)
	viper.SetConfigName(s.FileName)
	viper.SetConfigType(s.FileType)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	val := reflect.ValueOf(target).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			val.Field(i).SetString(viper.GetString(tag))
		case reflect.Int32:
			val.Field(i).SetInt(int64(viper.GetInt32(tag)))
		default:
			return fmt.Errorf("unsupported field type: %v", field.Type)
		}
	}
	return nil
}
