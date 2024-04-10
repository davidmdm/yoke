package main

import (
	"encoding"
	"encoding/json"

	"github.com/davidmdm/conf"
	"github.com/davidmdm/yoke/internal"
)

type Parameters struct {
	Wasm  string
	Input string
	Args  []string
}

var _ encoding.TextUnmarshaler = new(Parameters)

func (parameters *Parameters) UnmarshalText(data []byte) error {
	type Param struct {
		Name   string   `json:"name"`
		String string   `json:"string"`
		Array  []string `json:"array"`
	}

	var elems []Param
	if err := json.Unmarshal(data, &elems); err != nil {
		return err
	}

	wasm, _ := internal.Find(elems, func(param Param) bool { return param.Name == "wasm" })
	parameters.Wasm = wasm.String

	input, _ := internal.Find(elems, func(param Param) bool { return param.Name == "input" })
	parameters.Input = input.String

	args, _ := internal.Find(elems, func(param Param) bool { return param.Name == "args" })
	parameters.Args = args.Array

	return nil
}

type Config struct {
	Application struct {
		Name      string
		Namespace string
	}
	Flight Parameters
}

func getConfig() (cfg Config, err error) {
	conf.Var(conf.Environ, &cfg.Application.Name, "ARGOCD_APP_NAME", conf.Required[string](true))
	conf.Var(conf.Environ, &cfg.Application.Namespace, "ARGOCD_APP_NAMESPACE", conf.Required[string](true))
	conf.Var(conf.Environ, &cfg.Flight, "ARGOCD_APP_PARAMETERS", conf.Required[Parameters](true))
	err = conf.Environ.Parse()
	return
}
