package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"syscall"
)

func main() {
	file, err := ioutil.ReadFile("torrc")
	if err != nil {
		panic(err)
	}

	ipAddrs, err := net.LookupHost("nginx")
	if err != nil {
		panic(err)
	}
	if len(ipAddrs) == 0 {
		panic("could not resolve nginx")
	}

	file = bytes.Replace(file, []byte("$NGINX_IP_ADDRESS$"), []byte(ipAddrs[0]), 1)
	err = ioutil.WriteFile("torrc", file, 0700)
	if err != nil {
		panic(err)
	}

	syscall.Exec("tor", []string{"tor", "-f", "torrc"}, os.Environ())
}
