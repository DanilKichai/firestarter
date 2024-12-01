package main

import (
	"firestarter/internal/app/dhcp6ir/config"
	"firestarter/internal/pkg/generator"
	"flag"
	"log"
)

func main() {
	srcport := flag.Int("source-port", 546, "source port")
	dstport := flag.Int("destination-port", 547, "destination port")
	dstaddr := flag.String("destination-address", "ff02::01", "destination address")
	clientid := flag.String("client-id", "dhcp6ir", "client ID")
	retryint := flag.Int("retry-interval", 1, "retry interval seconds")
	timeout := flag.Int("timeout", 30, "timeout seconds")
	genfile := flag.String("generator", "/etc/dhcp6ir.gen", "generator file")

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
