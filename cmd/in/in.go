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
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"not enough args - usage: %s <target directory>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	rawBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	inRequest, err := concourse.NewInRequest(rawBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	// versions encode env names, but they can be overridden.
	// if you provide an env config, you will not necessarily fetch
	// the version that concourse has provided.
	storageClient, err := storage.NewStorageClient(
		inRequest.Source.GCPServiceAccountKey,
		inRequest.Version.Name,
		inRequest.Source.Bucket,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create storage client: %s\n", err)
		os.Exit(1)
	}

	version, err := storageClient.Download(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to download bbl state: %s\n", err)
		os.Exit(1)
	}

	outMap := map[string]storage.Version{"version": version}
	err = json.NewEncoder(os.Stdout).Encode(outMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal version: %s\n", err)
		os.Exit(1)
	}
}
