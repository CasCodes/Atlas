package capture

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"atlas/geo"
	"atlas/graph"
)

type Tracer struct {
	MaxHops   int
	graph     *graph.Graph
	geoLookup func(string) *geo.GeoInfo // inject lookup function from scanner
}

func NewTracer(maxHops int, graph *graph.Graph, lookup func(string) *geo.GeoInfo) *Tracer {
	return &Tracer{
		MaxHops:   maxHops,
		graph:     graph,
		geoLookup: lookup,
	}
}

func (t *Tracer) Trace(ctx context.Context, ip string) {
	// runs traceroute on the given ip
	// and adds discovered IPs to graph
	fmt.Printf("	START trace on %s\n", ip)
	cmd := exec.CommandContext(ctx, "traceroute", "-m", strconv.Itoa(t.MaxHops), ip)
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
	var prevNode *graph.Node = nil

	for scanner.Scan() {
		line := scanner.Text()
		traceIP, err := extractIP(line)
		if err != nil {
			// couldnt pass this line
			continue
		}

		// call lookup function (cache)
		geoInfo := t.geoLookup(traceIP)
		// construct node
		curNode := graph.Node{
			IP:      geoInfo.IP,
			Country: geoInfo.Country,
			City:    geoInfo.City,
			Org:     geoInfo.Org,
			Lat:     geoInfo.Lat,
			Lon:     geoInfo.Lon,
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
	fmt.Printf("	END trace on %s\n", ip)
}

var ipRegex = regexp.MustCompile(`\((\d+\.\d+\.\d+\.\d+)\)`)

func extractIP(line string) (string, error) {
	match := ipRegex.FindStringSubmatch(line)
	if match == nil {
		return "", fmt.Errorf("No IP found") // *** line or header
	}
	return match[1], nil
}
