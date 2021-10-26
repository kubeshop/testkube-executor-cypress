package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubeshop/testkube-executor-cypress/pkg/runner"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/storage/minio"
)

var (
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	location        string
	token           string
	ssl             bool
	scrapperEnabled bool
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("missing input argument")
		os.Exit(1)
	}

	script := args[1]

	e := testkube.Execution{}
	json.Unmarshal([]byte(script), &e)
	runner := runner.NewCypressRunner()
	result := runner.Run(e)
	fmt.Println(result)
	if scrapperEnabled {
		err := scrapeArtefacts(e.Id)
		if err != nil {
			fmt.Println("error occured while scrapping artefacts") // maybe we should consider the run failed since it is not able to save artefacts
		}
	}
	fmt.Printf("$$$%s$$$", e.Id)
}

func scrapeArtefacts(id string) error {
	client, err := minio.NewClient(endpoint, accessKeyID, secretAccessKey, location, token, ssl) // create storage client
	if err != nil {
		return err
	}

	err = client.CreateBucket(id) // create bucket name it by execution ID
	if err != nil {
		return fmt.Errorf("failed to create a bucket %s: %w", id, err)
	}
	err = filepath.Walk("cypress/", // cypress stores test artefacts in cypress directory by default
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				fmt.Println(path, info.Size())
				err = client.SaveFile(id, path) //The function will detect if there is a subdirectory and store accordingly
				if err != nil {
					return err
				}
			}

			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
