package main

import (
	"flag"
	"fmt"
	"github.com/stgnet/speedtest/api"
	"log"
	"os"
	"time"
)

func version() {
	fmt.Print(api.Version)
}

func usage() {
	fmt.Fprint(os.Stderr, "Command line interface for testing internet bandwidth using speedtest.net.\n\n")
	flag.PrintDefaults()
}

func main() {
	opts := api.ParseOpts()

	switch {
	case opts.Help:
		usage()
		return
	case opts.Version:
		version()
		return
	}

	client := api.NewClient(opts)

	if opts.List {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatalf("Failed to load server list: %v\n", err)
		}
		fmt.Println(servers)
		return
	}

	config, err := client.Config()
	if err != nil {
		log.Fatal(err)
	}

	client.Log("Testing from %s (%s)...\n", config.Client.ISP, config.Client.IP)

	server := selectServer(opts, client)

	downloadSpeed := server.DownloadSpeed()
	reportSpeed(opts, "Download", downloadSpeed)

	uploadSpeed := server.UploadSpeed()
	reportSpeed(opts, "Upload", uploadSpeed)
}

func reportSpeed(opts *api.Opts, prefix string, speed int) {
	if opts.SpeedInBytes {
		fmt.Printf("%s: %.2f MiB/s\n", prefix, float64(speed)/(1<<20))
	} else {
		fmt.Printf("%s: %.2f Mib/s\n", prefix, float64(speed)/(1<<17))
	}
}

func selectServer(opts *api.Opts, client api.Client) (selected *api.Server) {
	if opts.Server != 0 {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatal("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.Find(opts.Server)
		if selected == nil {
			log.Fatalf("Server not found: %d\n", opts.Server)
			return nil
		}
		selected.MeasureLatency(api.DefaultLatencyMeasureTimes, api.DefaultErrorLatency)
	} else {
		servers, err := client.ClosestServers()
		if err != nil {
			log.Fatal("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.MeasureLatencies(
			api.DefaultLatencyMeasureTimes,
			api.DefaultErrorLatency).First()
	}

	if opts.Quiet {
		log.Printf("Ping: %d ms\n", selected.Latency/time.Millisecond)
	} else {
		client.Log("Hosted by %s (%s) [%.2f km]: %d ms\n",
			selected.Sponsor,
			selected.Name,
			selected.Distance,
			selected.Latency/time.Millisecond)
	}

	return selected
}
