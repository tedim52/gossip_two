package main

import (	
	"github.com/tedim52/gossip_two/node_impls"
	"github.com/tedim52/gossip_two/node_interface"
	"github.com/tedim52/gossip_two/node_interface/objects"

	"fmt"
	"bufio"
	"os"
	"errors"
	"regexp"
	"strings"
	"strconv"
)

const (
	promptStr = ">> "
	portRegexStr = "^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$"
	ipAddressRegexStr = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
)

var (
	ipAddressRegexPat = regexp.MustCompile(ipAddressRegexStr)
	portRegexPat = regexp.MustCompile(portRegexStr)

	InvalidInput = errors.New("Invalid input format. Please provide './gossip-two <ip-address> <port>'.")
)

func main() {
	// process input arguments
	ip, port, err := processInput(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// fmt.Println(ip)
	// fmt.Println(port)

	// initialize and start healthy gossip node
	node := node_impls.NewHealthyGossipNode(ip, port)
	node.BoostrapNode()

	// start read-eval print loop
	gossipRepl(node)
}

func gossipRepl(node node_interface.GossipNode){
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(promptStr)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		input = strings.TrimSpace(input)
		// fmt.Println(input)

		if (input == "?"){
			fmt.Println("printing gossip database...")
			node.GetDatabase().PrintDatabase()
		} else if (input[0] == '+' && len(input) > 1) {
			fmt.Println("adding peer to gossip node...")
			input = input[1:]
			peerNodeID, err := objects.DeserializeNodeID(input)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = node.AddPeer(peerNodeID)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if intVal, err := strconv.ParseInt(input, 10, 32); err == nil {
			fmt.Printf("updating gossip node to have this value: %d\n", intVal)
			node.UpdateValue(intVal)
		} else {
			fmt.Println("Unrecognized input. Try again.")
		}
	}
}

// processes command line input by asserting the following format and corresponding regexes of args:
// input format: ./... <ip-address> <port>
func processInput(args []string) (objects.IPAddress, objects.Port, error) {
	if args == nil || len(args) < 2 {
		return "", "", InvalidInput
	}
	ipAddressStr := args[1]
	portStr := args[2]
	if !ipAddressRegexPat.Match([]byte(ipAddressStr)) {
		return "", "", objects.InvalidIPAddress
	}
	if !portRegexPat.Match([]byte(portStr)) {
		return "", "", objects.InvalidPortNumber
	}
	return objects.IPAddress(ipAddressStr), objects.Port(portStr), nil
}
