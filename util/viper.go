package util

import (
	"fmt"
	"go-usage/model"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func viperLoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/org/")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Fatal error config file: %v\n", err.Error()))
	}

	var config model.EnvConfigModel

	configMap := map[string]any{
		"server.mode": &config.Server.Mode,
		"server.port": &config.Server.Port,

		"workspace.cache": &config.Workspace.Cache,
		"workspace.key":   &config.Workspace.Key,

		"performance.maxcpucore": &config.Performance.MaxCpuCore,
		"performance.maxmemory":  &config.Performance.MaxMemory,
		"performance.tasklimit":  &config.Performance.TaskLimit,
	}

	//load from config file
	configToModel(&configMap)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//if config empty, load from env var
	configToModel(&configMap)

	fmt.Println("==== config ====")
	printStructFields(config)
	fmt.Println("==== config ====")
}

func configToModel(configMap *map[string]any) {
	for k, v := range *configMap {
		pointerAdd := reflect.ValueOf(v)
		config := reflect.Indirect(pointerAdd)
		//check if config exist, not need to load again
		if !config.IsZero() {
			continue
		}
		switch p := v.(type) {
		case *string:
			*p = viper.GetString(k)
		case *[]string:
			*p = append(*p, viper.GetStringSlice(k)...)
		case *uint:
			*p = uint(viper.GetInt(k))
		case *int:
			*p = viper.GetInt(k)
		}
	}
}

func printStructFields(s any) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// Check if the field is a struct and print its fields recursively
		if field.Kind() == reflect.Struct {
			fmt.Println(fieldName + ":")
			printStructFields(field.Interface())
		} else {
			fmt.Printf("%s: %v\n", fieldName, field.Interface())
			checkConfigValid(&fieldName, &field)
		}
	}
}

func checkConfigValid(fieldName *string, value *reflect.Value) {
	if value.IsZero() {
		panic(fmt.Errorf("missing config: %v", *fieldName))
	}
}
