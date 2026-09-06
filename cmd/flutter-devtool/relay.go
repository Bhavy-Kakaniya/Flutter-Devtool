package main

import (
	"devtool/internal/relay"
	"fmt"
)

// start local rela server
func runRelay() {

	// start tcp relay server
	if err:= relay.StartServer(); err != nil {
		fmt.Println("Error:", err)
	}
	
}