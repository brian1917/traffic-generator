package main

import (
	"fmt"
	"os"
)

func showHelp() {

	fmt.Println(`traffic-generator is a tool to generate traffic for demos and labs.

  Usage: traffic-generator [command] [arguments]

  Commands:
	open-http-listener    open an http istener on a port. requires on argument of port number to open.
	send-traffic          sends traffic based on input csv file. requires one argument of the csv file to use.
	continuous            coninously send (when host matches source) and listen (when host matches destination) for traffic based on csv. requires one argument of the csv file to use.
	version               prints the version and last commit of the tool.

  csv format notes:
	- open-http-listener csv requires three headers: src, dst, and port.
	- continuous csv requiest four columns: src, dst, port, and proto.
	- other columns are acceptable but igored.
	- order of columns do not matter
	- src_ip, dst_ip can be used in place of src and dst. both are acceptable.
	- proto and protocol can be used interchangeably.`)
	os.Exit(0)
}
