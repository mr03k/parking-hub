package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/swaggest/openapi-go/openapi3"
)

func main() {
	directory := flag.String("directory", "./docs/swagger/", "directory of the application")
	flag.Parse()

	ctx := context.Background()

	oapi3Reflector := openapi3.NewReflector()
	spec, err := wireApp(ctx, slog.Default(), oapi3Reflector)
	if err != nil {
		panic(err)
	}

	jsonData, err := spec.MarshalJSON()
	if err != nil {
		panic(err)
	}

	yamldata, err := spec.MarshalYAML()
	if err != nil {
		panic(err)
	}

	yamlFile := *directory + "swagger.yaml"
	jsonFile := *directory + "swagger.json"

	err = os.WriteFile(yamlFile, yamldata, 0o644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(jsonFile, jsonData, 0o644)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(yamldata))
}
