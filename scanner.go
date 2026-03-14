package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Scanner struct {
	cache  *GeoCache
	client *http.Client
	device string
}

func NewScanner(device string) *Scanner {
	// intantiate cache and http client
	cache := NewGeoCache()
	client := &http.Client{Timeout: 3 * time.Second}

	return &Scanner{
		cache:  cache,
		client: client,
		device: device,
	}
}

func (s *Scanner) Scan(print bool) {
	// loop to capture packages on en0
	handle, err := pcap.OpenLive(s.device, 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		ip := packet.Layer(layers.LayerTypeIPv4)
		if ip == nil {
			continue
		}

		ip4, _ := ip.(*layers.IPv4)

		// cache lookup
		srcInfo := s.lookup(ip4.SrcIP.String())
		dstInfo := s.lookup(ip4.DstIP.String())

		if print {
			printPacket(ip4, srcInfo, dstInfo)
		}
	}
}

func (s *Scanner) lookup(ip string) *GeoInfo {
	// check if ip exists in cache, else fetch from api
	if info, ok := s.cache.Get(ip); ok {
		return info // cache hit
	}

	info := fetchFromAPI(ip, s.client) // cache miss
	s.cache.Put(ip, info)
	return info
}

func isPrivate(ip string) {
	// check if ip adress is private
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
