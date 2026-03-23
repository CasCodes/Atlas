package main

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type Tracer struct {
	maxHops int
	graph   *Graph
	lookup  func(string) *GeoInfo // inject lookup function from scanner
}

func NewTracer(maxHops int, graph *Graph, lookup func(string) *GeoInfo) *Tracer {
	return &Tracer{
		maxHops: maxHops,
		graph:   graph,
		lookup:  lookup,
	}
}

func (t *Tracer) Trace(ctx context.Context, ip string) {
	// runs traceroute on the given ip
	// and adds discovered IPs to graph
	cmd := exec.CommandContext(ctx, "traceroute", "-m", strconv.Itoa(t.maxHops), ip)
	// listen on stdout
	out, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	// start command
	if cmd.Start() != nil {
		return
	}

	// read from stdout
	scanner := bufio.NewScanner(out)

	// store previous node
	var prevNode *Node = nil

	for scanner.Scan() {
		line := scanner.Text()
		traceIP, err := extractIP(line)
		if err != nil {
			// couldnt pass this line
			continue
		}

		// call lookup function (cache)
		geoInfo := t.lookup(traceIP)
		// construct node
		curNode := Node{
			IP:      geoInfo.IP,
			Country: geoInfo.Country,
			City:    geoInfo.City,
			Org:     geoInfo.Org,
		}

		// add nodes to graph as soon as we have a pair
		if prevNode == nil {
			// first node
			prevNode = &curNode
			continue
		}

		// else we have a pair
		t.graph.Add(*prevNode, curNode)

		// set curNode as prevNode
		prevNode = &curNode
	}

}

var ipRegex = regexp.MustCompile(`\((\d+\.\d+\.\d+\.\d+)\)`)

func extractIP(line string) (string, error) {
	match := ipRegex.FindStringSubmatch(line)
	if match == nil {
		return "", fmt.Errorf("No IP found") // *** line or header
	}
	return match[1], nil
}
