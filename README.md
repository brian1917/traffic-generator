# Traffic Generator

## Description
A utility for generating traffic in a lab. The utility opens listening ports and generates traffic.

## Generate Traffic
To generate traffic run `traffic-generator traffic traffic.csv` where `traffic.csv` is an input file with the following headers: `src`, `dst`, and `port`. The name of the header does not actually matter, but the order does. Source must be first column, destination second, and port third.

See below for an example input file:
```
src,dst,port
fin-dev-web01.domain.com,fin-dev-app01.domain.com,8080
fin-dev-app01.domain.com,fin-dev-db01.domain.com,3306
fin-dev-app01.domain.com,fin-prd-db02.domain.com,3306
fin-prd-web01.domain.com,fin-prd-app01.domain.com,8080
```

If the `src` field in the input file matches the hostname of the workload, it makes a connection to the destination on the provided port.

## Opening ports
To open a listening TCP port on a host, run `traffic-generator open 3306` where 3306 is the port to open. The application returns a simple statement showing the connection information.

