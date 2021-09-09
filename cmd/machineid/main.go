// Package main provides the command line app for reading the unique machine id of most OSs.
//
// Usage: machineid [options]
//
// Options:
//   --appid    <AppID>    Protect machine id by hashing it together with an app id.
//
// Try:
//   machineid
//   machineid --appid MyAppID
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log"

	"github.com/soxfmr/machineid"
)

const usageStr = `
Usage: machineid [options]

Options:
  --hash    Encrypt machine id by hashing it in specified hash algorithm.

Try:
  machineid
  machineid --hash
`

func usage() {
	log.Fatalln(usageStr)
}

func main() {
	var hashed bool
	flag.BoolVar(&hashed, "hash", false, "Encrypt machine id by hashing it in specified hash algorithm.")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	var id string
	var err error
	if hashed {
		algo := md5.New()
		id, err = machineid.HashID(algo)
	} else {
		id, err = machineid.ID()
	}
	if err != nil {
		log.Fatalf("Failed to read machine id with error: %s\n", err)
	}
	fmt.Println(id)
}
