package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var (
	version string
)

const (
	indexYaml = "index.yaml"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Currently plugin doesn't support any external commands and is meant to be used by Helm native commands.")
		fmt.Println("Try: `helm repo add ecr://` and `helm template ecr://my/repository/chart/version`")
		fmt.Printf("Version: %s\n", version)
		return
	}

	repository, err := NewRepository(os.Args[4])
	if err != nil {
		log.Fatalf("failed to parse repository: %v", err)
	}

	session, err := NewSession()
	if err != nil {
		log.Fatalf("failed to create AWS session: %v", err)
	}

	ecrService := ecr.New(session, aws.NewConfig().WithRegion(*repository.Region))

	// Either Helm is asking to add ECR as the Helm repository or Helm index (index.yaml)
	// either way return newly generated index by describing images matching Helm artifacts
	if repository.Filename == nil || (repository.Filename != nil && *repository.Filename == indexYaml) {
		index, err := NewIndex()
		if err != nil {
			log.Fatalf("failed to create Helm index: %v", err)
		}

		err = PopulateIndex(ecrService, repository, index)
		if err != nil {
			log.Fatalf("failed to populate Helm index: %v", err)
		}

		b, err := index.MarshalBinary()
		if err != nil {
			log.Fatalf("failed to marshal Helm index: %v", err)
		}

		fmt.Println(string(b))
		return
	}

	// Otherwise Helm is asking to download the Chart and print it out in binary form
	err = FetchChartAndWrite(ecrService, repository, os.Stdout)
	if err != nil {
		log.Fatalf("failed to fetch chart: %v", err)
	}
}
