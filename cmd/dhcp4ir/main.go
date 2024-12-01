package main

import (
	"firestarter/internal/app/dhcp4ir/config"
	"firestarter/internal/pkg/generator"
	"flag"
	"log"
)

func main() {
	srcport := flag.Int("source-port", 68, "source port")
	dstport := flag.Int("destination-port", 67, "destination port")
	dstaddr := flag.String("destination-address", "255.255.255.255", "destination address")
	clientid := flag.String("client-id", "dhcp4ir", "client ID")
	retryint := flag.Int("retry-interval", 1, "retry interval seconds")
	timeout := flag.Int("timeout", 30, "timeout seconds")
	genfile := flag.String("generator", "/etc/dhcp4ir.gen", "generator file")

	flag.Parse()

	cfg, err := config.Request(*srcport, *dstport, *dstaddr, *clientid, *retryint, *timeout)
	if err != nil {
		log.Fatalf("load config from efivars: %v", err)
	}

	batch, err := generator.Load(*genfile, cfg)
	if err != nil {
		log.Fatalf("load generator file: %v", err)
	}

	err = batch.Write()
	if err != nil {
		log.Fatalf("extract files from generated batch: %v", err)
	}
}
