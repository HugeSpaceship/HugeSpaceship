// Package testdata contains all the data used in unit tests
package testdata

import _ "embed"

//go:embed test.dds.compressed
var TestDDSCompressed []byte

//go:embed example-psn-ticket
var ExamplePSNTicket []byte

//go:embed example-rpcn-ticket
var ExampleRPCNTicket []byte
