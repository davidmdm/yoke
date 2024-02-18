package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	fmt.Println(`[{
  		"apiVersion": "apps/v1",
  		"kind": "Deployment",
  		"metadata": {
  		  "name": "sample-app-prod"
  		},
  		"spec": {
  		  "replicas": 3,
  		  "selector": {
  		    "matchLabels": {
  		      "app": "sample-app"
  		    }
  		  },
  		  "template": {
  		    "metadata": {
  		      "labels": {
  		        "app": "sample-app"
  		      }
  		    },
  		    "spec": {
  		      "containers": [
  		        {
  		          "name": "web-app",
  		          "image": "alpine:latest",
  		          "command": ["watch", "echo", "hello", "world", "yes?"]
  		        }
  		      ]
  		    }
  		  }
  		}
	}]`)

	return nil
}
