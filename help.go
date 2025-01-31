package main

import (
	"fmt"
	"os"
)

func showHelp() {

	fmt.Println(`traffic-generator is a tool to generate traffic for demos and labs.

  Usage: traffic-generator [command] [arguments]

  Commands:
	open-listener   Open a listener on a port. requires on argument of port number to open.
	send-traffic    Sends traffic based on input csv file. requires one argument of the csv file to use.
	continuous		Opens ports and continually sends traffic based on csv file. requires on argument for the csv file to use.

  CSV format notes:
	- The CSV requires three headers: src, dst, and port. Other columns are acceptable but ignored. The order of the columns does not matter.
	- src_ip, dst_ip can be used in place of src and dst. both are acceptable`)
	os.Exit(0)
}
