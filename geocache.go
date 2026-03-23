package main

import "sync"

type GeoInfo struct {
	IP      string  `json:"query"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Org     string  `json:"org"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Status  string  `json:"status"` // "success" or "fail"
}

type GeoCache struct {
	mu   sync.RWMutex
	data map[string]*GeoInfo
}

func NewGeoCache() *GeoCache {
	return &GeoCache{
		data: make(map[string]*GeoInfo),
	}
}

func (g *GeoCache) Get(ip string) (*GeoInfo, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	info, ok := g.data[ip]
	return info, ok
}

func (g *GeoCache) Put(ip string, info *GeoInfo) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.data[ip] = info
}
