package main

import (
	"fmt"
	"os"

	"github.com/sh3rp/swan/net"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Printf("Must specify a device IP and a community.\n")
		return
	}
	// device IP, community
	manager := net.NewSwitchManager(os.Args[1], os.Args[2])

	ver, err := manager.GetVersion()

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("%s\n\n", ver.Hostname)
	fmt.Printf("%+v\n\n", ver.OSVersion)

	intfs, err := manager.GetIfs()

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Interface       In     Out\n")
	fmt.Printf("=============== ====== =======\n")

	for _, intf := range intfs {
		stats, _ := manager.GetIfStats(intf)

		fmt.Printf("%-15s %-6d %-6d\n", intf.Name, stats.IfBitsInPerSecond, stats.IfBitsOutPerSecond)
	}
}
