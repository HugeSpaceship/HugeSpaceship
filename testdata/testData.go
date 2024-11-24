// Package testdata contains all the data used in unit tests
package testdata

import _ "embed"

//go:embed image/test.dds.compressed
var TestDDSCompressed []byte

//go:embed npticket/example-psn-ticket
var ExamplePSNTicket []byte

//go:embed npticket/example-rpcn-ticket
var ExampleRPCNTicket []byte

//go:embed npticket/ezonpticket.bin
var Wack []byte
