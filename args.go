package main

import (
	"github.com/alexflint/go-arg"
	ic "github.com/egirna/icap-client"
	"strings"
)

type arguments struct {
	Url     string `arg:"positional" default:"https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"`
	Mode    string `arg:"-m,--mode" help:"ICAP mode - REQMODE, RESPMOD or OPTIONS"`
	Timeout int    `arg:"-t,--timeout" help:"Timeout in seconds" default:"5"`
	IcapUrl string `arg:"--icap" help:"ICAP server URL" default:"icap://localhost:1344/kasperskiy"`
	File    string `arg:"-f,--file" help:"File to upload" default:"./files-for-testing/sample.pdf"`
}

var args arguments

func parseArgs() *arguments {
	arg.MustParse(&args)

	return &args
}

func (a arguments) getMode() string {
	switch strings.ToUpper(a.Mode) {
	case "REQMOD":
	case "REQUEST":
		return ic.MethodREQMOD
	case "RESPMOD":
	case "RESPONSE":
		return ic.MethodRESPMOD
	case "OPTIONS":
		return ic.MethodOPTIONS
	}

	return ic.MethodRESPMOD
}
