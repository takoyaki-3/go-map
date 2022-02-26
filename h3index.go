package gomap

import (
	"math"

	uh3 "github.com/uber/h3-go"
)

type UH3index uh3.H3Index
type UH3indexes map[UH3index][]int

func MakeH3Index(g *Graph, resolution int) UH3indexes {
	h3index := UH3indexes{}

	for k, v := range g.Nodes {
		hex := uh3.FromGeo(uh3.GeoCoord{
			Latitude:  v.Lat,
			Longitude: v.Lon}, resolution)

		if _, ok := h3index[UH3index(hex)]; !ok {
			h3index[UH3index(hex)] = []int{}
		}
		h3index[UH3index(hex)] = append(h3index[UH3index(hex)], k)
	}
	return h3index
}

func FindNode(g *Graph, h3indexes UH3indexes, latlon Node, resolution int) int {
	h3index := uh3.FromGeo(uh3.GeoCoord{
		Latitude:  latlon.Lat,
		Longitude: latlon.Lon}, resolution)

	hexes, _ := uh3.HexRing(h3index, 1)
	hexes = append(hexes, h3index)

	min_node := -1
	min_d := math.MaxFloat64

	for _, v := range hexes {
		if _, ok := h3indexes[UH3index(v)]; !ok {
			continue
		}
		for _, v := range h3indexes[UH3index(v)] {
			d := HubenyDistance(g.Nodes[v], latlon)
			if min_d > d {
				min_node = v
				min_d = d
			}
		}
	}

	return min_node
}

func FindNodes(g *Graph, h3indexes UH3indexes, latlon Node, resolution int, r float64) []int {
	h3index := uh3.FromGeo(uh3.GeoCoord{
		Latitude:  latlon.Lat,
		Longitude: latlon.Lon}, resolution)

	hexes, _ := uh3.HexRing(h3index, 1)
	hexes = append(hexes, h3index)

	nodes := []int{}

	for _, v := range hexes {
		if _, ok := h3indexes[UH3index(v)]; !ok {
			continue
		}
		for _, v := range h3indexes[UH3index(v)] {
			d := HubenyDistance(g.Nodes[v], latlon)
			if r > d {
				nodes = append(nodes, v)
			}
		}
	}

	return nodes
}
