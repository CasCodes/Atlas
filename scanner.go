package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Scanner struct {
	device    string
	cache     *GeoCache
	client    *http.Client
	graph     *Graph
	tracer    *Tracer
	tracedSet map[string]bool // to remember already traced IPs
}

func NewScanner(device string, maxHops int) *Scanner {
	// intantiate cache and http client
	cache := NewGeoCache()
	client := &http.Client{Timeout: 3 * time.Second}

	// package graph
	graph := NewGraph()

	s := &Scanner{
		device:    device,
		cache:     cache,
		client:    client,
		graph:     graph,
		tracedSet: make(map[string]bool),
	}
	// add tracer with lookup function injected
	s.tracer = NewTracer(maxHops, graph, s.lookup)

	return s
}

func (s *Scanner) Scan(durationMS int, print bool) {
	// create context object to pass on timeout to other threads
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(durationMS)*time.Millisecond)
	defer cancel() // cancel on end of function

	// loop to capture packages on en0
	handle, err := pcap.OpenLive(s.device, 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// each iteration, check for context timeout, else process packet
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-packetSource.Packets():
			ip := packet.Layer(layers.LayerTypeIPv4)
			if ip == nil {
				continue
			}

			ip4, _ := ip.(*layers.IPv4)
			srcIP := ip4.SrcIP.String()
			destIP := ip4.DstIP.String()

			// cache lookup
			srcInfo := s.lookup(srcIP)
			dstInfo := s.lookup(destIP)

			// start trace in new thread
			if !isPrivate(destIP) && !s.tracedSet[destIP] {
				s.tracedSet[destIP] = true
				go s.tracer.Trace(ctx, destIP)
			}

			if print {
				printPacket(ip4, srcInfo, dstInfo)
			}
		}
	}
}

func (s *Scanner) lookup(ip string) *GeoInfo {
	// private is not cached
	if isPrivate(ip) {
		return &GeoInfo{IP: ip, Country: "local", City: "local network", Org: "LAN"}
	}

	// check if ip exists in cache, else fetch from api
	if info, ok := s.cache.Get(ip); ok {
		return info // cache hit
	}
	// cache miss: fetch ip info
	info := fetchFromAPI(ip, s.client)
	// add to cache
	s.cache.Put(ip, info)
	return info
}

func isPrivate(ipStr string) bool {
	// check if ip adress is private
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true
	}
	private := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
	}
	for _, cidr := range private {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func fetchFromAPI(ip string, client *http.Client) *GeoInfo {
	url := "http://ip-api.com/json/" + ip + "?fields=query,country,city,org,lat,lon"
	resp, err := client.Get(url)
	if err != nil {
		return &GeoInfo{IP: ip, Country: "unknown", City: "lookup failed"}
	}
	defer resp.Body.Close()

	var info GeoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return &GeoInfo{IP: ip, Country: "unknown", City: "parse failed"}
	}
	return &info
}

func printPacket(ip4 *layers.IPv4, src *GeoInfo, dst *GeoInfo) {
	fmt.Printf("%-16s (%s, %s) -> %-16s (%s, %s)\n",
		ip4.SrcIP, src.City, src.Country,
		ip4.DstIP, dst.City, dst.Country,
	)
}
