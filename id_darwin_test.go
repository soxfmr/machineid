// +build darwin

package machineid

import (
	"log"
	"testing"
)

func Test_machineID(t *testing.T) {
	got, err := machineID()
	if err != nil {
		t.Error(err)
	}
	if got == "" {
		t.Error("Got empty machine id")
	}
	log.Printf("Hardware UUID: %s", got)
}
