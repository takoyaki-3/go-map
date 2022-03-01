package gomap

import (
	gtfs "github.com/takoyaki-3/go-gtfs"
)

type Graph struct {
	Nodes        []Node
	Edges        []Edge
	FromEdges    [][]int
	ToEdges      [][]int
	Stops        []gtfs.Stop
	stopId2index map[string]int
}

func (g Graph) GetStop(stopId string) gtfs.Stop {
	return g.Stops[g.stopId2index[stopId]]
}

type Node struct {
	Lat     float64 `csv:"lat"`
	Lon     float64 `csv:"lon"`
	Type    string  `csv:"type"`
	PlaceID string  `csv:"place_id"`
}

type Edge struct {
	FromNode int     `csv:"from_node_index"`
	ToNode   int     `csv:"to_node_index"`
	Weight   float64 `csv:"weight"`
	Type     string  `csv:"type"`
	ViaNodes []int
}
