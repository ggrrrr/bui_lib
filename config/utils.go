package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type ParamValue struct {
	Name     string
	DefValue interface{}
	Info     string
	IsSecret bool
}

var (
	params []ParamValue
)

func init() {
	viper.AutomaticEnv()
	// viper.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

}

func Configure(p []ParamValue) {
	params = append(params, p...)
	for _, v := range params {
		viper.SetDefault(v.Name, v.DefValue)
	}
}

func ConfigureGroup(p []ParamValue, group string) {
	// params = append(params, p...)
	for _, v := range p {
		param := ParamValue{
			Name:     fmt.Sprintf(v.Name, group),
			Info:     fmt.Sprintf(v.Info, group),
			DefValue: v.DefValue,
		}
		params = append(params, param)
	}
}

func ErrorParamNotSet(paramName string) error {
	return fmt.Errorf("param: %s not set", paramName)
}

func ErrorParamInvalid(paramName string, val interface{}) error {
	return fmt.Errorf("invalid %s=%v", paramName, val)
}

func PrintConfig() string {

	var sb strings.Builder
	sb.WriteString("\nConfiguration parameters: \n")
	for _, v := range params {
		// sb.WriteString(fmt.Sprintf("%+20s:%5s%s\n", v.Name, "", v.Info))
		// sb.WriteString(fmt.Sprintf("%20s:%10s\"%+v\"\n", "default", "", v.DefValue))
		val := "***"
		if !v.IsSecret {
			val = viper.GetString(v.Name)
		}
		sb.WriteString(fmt.Sprintf("%20s:%10s\"%+v\"\n", v.Name, "", val))

	}
	return sb.String()

}

func Help() string {

	var sb strings.Builder
	sb.WriteString("\nConfiguration parameters: \n")
	for _, v := range params {
		n := strings.Replace(v.Name, ".", "_", -1)
		// sb.WriteString(fmt.Sprintf("%+20s:%5s%s\n", strings.ToUpper(n), "", ""))
		sb.WriteString(fmt.Sprintf("%+20s:%5s%s\n", v.Name, "", v.Info))
		sb.WriteString(fmt.Sprintf("%20s:%10s\"%+v\"\n", "default", "", v.DefValue))
		val := "***"
		if !v.IsSecret {
			val = viper.GetString(v.Name)
		}
		sb.WriteString(fmt.Sprintf("%20s=%s\n\n", strings.ToUpper(n), val))

	}
	return sb.String()

}
