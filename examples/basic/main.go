package main

import (
	"fmt"
	"log"

	"github.com/prashanth1k/gingest"
)

func main() {
	// Example 1: Process local directory
	fmt.Println("=== Processing Local Directory ===")

	config := gingest.Config{
		Source:      "../../", // Process the gingest repo itself
		OutputFile:  "local_digest.md",
		MaxFileSize: 1024 * 1024, // 1MB limit
	}

	err := gingest.ProcessAndWriteDigest(config)
	if err != nil {
		log.Printf("Error processing local directory: %v", err)
	} else {
		fmt.Printf("Local digest created: %s\n", config.OutputFile)
	}

	// Example 2: Process remote repository
	fmt.Println("\n=== Processing Remote Repository ===")

	remoteConfig := gingest.Config{
		Source:       "https://github.com/golang/example.git",
		OutputFile:   "remote_digest.md",
		TargetBranch: "master",
		MaxFileSize:  2 * 1024 * 1024, // 2MB limit
	}

	err = gingest.ProcessAndWriteDigest(remoteConfig)
	if err != nil {
		log.Printf("Error processing remote repository: %v", err)
	} else {
		fmt.Printf("Remote digest created: %s\n", remoteConfig.OutputFile)
	}
}
