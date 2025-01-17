package main

import (
	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"
)

// TestSuite represents a group of tests that should be run together
type TestSuite interface {
	// RunTests runs all tests of this testsuite
	RunTests()
}

type Config struct {
	waitUser bool

	nodeAhost      string
	nodeAhttpPort  string
	nodeApeeringID string

	nodeBhost      string
	nodeBhttpPort  string
	nodeBpeeringID string

	nodeChost      string
	nodeChttpPort  string
	nodeCpeeringID string

	triggerNodeHost   string
	triggerNodeWSHost string
	triggerNodeID     string

	certFile string
	keyFile  string

	littleCertFile string
	littleKeyFile  string

	nodeA 		 *grpcclient.Node
	nodeB 		 *grpcclient.Node
	nodeC 		 *grpcclient.Node
	littleClient *grpcclient.Node
}
