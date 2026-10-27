package main

import (
	"flag"
	"fmt"
	"os"
)

func runGetAsns(args []string) {
	cmd := flag.NewFlagSet("asn", flag.ContinueOnError)
	cmd.Parse(args)

	if os.Args[1] == "asn" && cmd.NFlag() == 0 {
		getAsns()
	}
	if (cmd.NFlag() > 1) || len(args) > 1 {
		fmt.Fprintf(os.Stderr, "expected asn or aic only")
	}

}

func runGetAics(args []string) {
	cmd := flag.NewFlagSet("aic", flag.ContinueOnError)
	cmd.Parse(args)

	if os.Args[1] == "aic" && cmd.NFlag() == 0 {
		getAics()
	}
	if (cmd.NFlag() > 1) || len(args) > 1 {
		fmt.Fprintf(os.Stderr, "expected asn or aic only")
	}

}

func main() {
	// Initialize
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "expected asn or aic only")
	}

	switch os.Args[1] {
	case "asn":
		runGetAsns(os.Args[2:])
	case "aic":
		runGetAics(os.Args[2:])
	default:
		fmt.Println("expected asn or aic only")
	}
}
