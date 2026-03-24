package main

import (
	"fmt"
	"sync"
)

type Node struct {
	IP      string  `json:"ip"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Org     string  `json:"org"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Edge struct {
	FromIP string `json:"from"`
	ToIP   string `json:"to"`
	Count  int    `json:"count"`
}

type EdgeKey struct {
	FromIP string
	ToIP   string
}

type Graph struct {
	mu    sync.RWMutex
	Nodes map[string]*Node  // keyed by IP
	Edges map[EdgeKey]*Edge // keyed by EdgeKey
}

type GraphJSON struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func NewGraph() *Graph {
	// allocate maps and return empty graph
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[EdgeKey]*Edge),
	}
}

func (g *Graph) Add(from Node, to Node) {
	g.mu.Lock()         // write lock
	defer g.mu.Unlock() // unlock at end of function

	// assignmend test on edge
	edgeKey := EdgeKey{from.IP, to.IP}
	edge, exists := g.Edges[edgeKey]

	if exists {
		// edge exists => nodes exist => we just need to increment in-place
		edge.Count++
		g.Edges[edgeKey] = edge
	} else {
		// add nodes & add new edge with count 1
		g.Nodes[from.IP] = &from
		g.Nodes[to.IP] = &to

		g.Edges[edgeKey] = &Edge{from.IP, to.IP, 1}
	}
}

func (g *Graph) Print() {
	for n := range g.Edges {
		fmt.Printf(" (%s) -- %d --> (%s) \n", n.FromIP, g.Edges[n].Count, n.ToIP)
	}
}

func (g *Graph) ToJSON() GraphJSON {
	g.mu.RLock()         // read lock
	defer g.mu.RUnlock() // unlock at end of function

	nodes := make([]Node, 0, len(g.Nodes))
	edges := make([]Edge, 0, len(g.Edges))

	for _, node := range g.Nodes {
		nodes = append(nodes, *node)
	}

	for _, edge := range g.Edges {
		edges = append(edges, *edge)
	}
	return GraphJSON{
		Nodes: nodes,
		Edges: edges,
	}
}
