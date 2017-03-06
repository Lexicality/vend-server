package main

import (
	"github.com/lexicality/vending/client/vendio"
	"github.com/lexicality/vending/shared/vending"
)

func vendItem(hw vendio.Hardware, location uint8) vending.Result {
	// TODO: Actually know if it vends??
	err := hw.Vend(location)
	if err != nil {
		log.Criticalf("Unable to vend: %s", err)
		return vending.ResultHardwareFailure
	}

	return vending.ResultSuccess
}
