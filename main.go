package main

import (
	"atlas/api"
	"atlas/capture"
)

func main() {
	// create scanner
	scanner := capture.NewScanner("en0", 10)
	// create router
	r := api.SetupRouter(scanner)
	// run the server
	r.Run(":8080")
}
