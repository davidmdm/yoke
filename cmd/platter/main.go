package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	io.Copy(os.Stdout, os.Stdin)
	// fmt.Println(`{
	// 	"apiVersion": "apps/v1",
	// 	"kind": "Deployment",
	// 	"metadata": {
	// 	  "name": "sample-app-prod"
	// 	},
	// 	"spec": {
	// 	  "replicas": 6,
	// 	  "selector": {
	// 	    "matchLabels": {
	// 	      "app": "sample-app"
	// 	    }
	// 	  },
	// 	  "template": {
	// 	    "metadata": {
	// 	      "labels": {
	// 	        "app": "sample-app"
	// 	      }
	// 	    },
	// 	    "spec": {
	// 	      "containers": [
	// 	        {
	// 	          "name": "web-app",
	// 	          "image": "alpine:latest",
	// 	          "command": ["watch", "echo", "hello", "world", "no"]
	// 	        }
	// 	      ]
	// 	    }
	// 	  }
	// 	}
	// }`)

	return nil
}
