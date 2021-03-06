package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/bbl-state-resource/concourse"
	"github.com/cloudfoundry/bbl-state-resource/storage"
)

func main() {
	rawBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	checkRequest, err := concourse.NewInRequest(rawBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	storageClient, err := storage.NewStorageClient(
		checkRequest.Source.GCPServiceAccountKey,
		checkRequest.Version.Name,
		checkRequest.Source.Bucket,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create storage client: %s\n", err)
		os.Exit(1)
	}

	versions, err := storageClient.GetAllNewerVersions(checkRequest.Version)
	if err == storage.ObjectNotFoundError {
		fmt.Fprintf(os.Stdout, `[]`)
		os.Exit(0)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch bbl state version: %s\n", err)
		os.Exit(1)
	}

	err = json.NewEncoder(os.Stdout).Encode(versions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal version: %s\n", err)
		os.Exit(1)
	}
}
