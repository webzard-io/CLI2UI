package main

import (
	"CLI2UI/pkg/config"
	"encoding/json"
	"fmt"
)

func main() {
	docker := config.CLI{
		Name: "docker",
		Command: config.Command{
			Name: "docker",
			Flags: []config.FlagOrArg{
				{
					Name: "config",
					Type: config.FlagArgTypeString,
				},
				{
					Name:    "log-level",
					Type:    config.FlagArgTypeEnum,
					Default: "info",
					Options: []string{"debug", "info", "warn", "error", "fatal"},
				},
			},
			Subcommands: []config.Command{
				{
					Name: "volume",
					Subcommands: []config.Command{
						{
							Name: "create",
							Flags: []config.FlagOrArg{
								{
									Name: "driver",
									Type: config.FlagArgTypeString,
								},
							},
							Args: []config.FlagOrArg{
								{
									Name: "name",
									Type: config.FlagArgTypeString,
								},
							},
						},
					},
				},
			},
		},
		FlagDelim: "=",
	}

	form := docker.Form()
	b, err := json.Marshal(form)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	raw := `
	{
		"Flags":{
			"config":"this-config.yaml",
			"log-level":"info"
		},
		"Args":{
			
		},
		"Subcommands":{
			"volume":{
				"Flags":{
					
				},
				"Args":{
					
				},
				"Subcommands":{
					"create":{
						"Flags":{
							"driver":null
						},
						"Args":{
							"name":"new-stuff"
						},
						"Subcommands":{
							
						},
						"Choice":""
					}
				},
				"Choice":"create"
			}
		},
		"Choice":"volume"
	}
	`

	var f config.Form
	json.Unmarshal([]byte(raw), &f)

	script := docker.Script(&f)
	fmt.Println(script)
}
