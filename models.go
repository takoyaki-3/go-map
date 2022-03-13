package gomap

import ()

type Graph struct {
	Nodes        []Node
	Edges        []Edge
	FromEdges    [][]int
	ToEdges      [][]int
	Stops        []Stop
	stopId2index map[string]int
	stopId2node  map[string]int
}

type Stop struct {
	ID          string  `csv:"stop_id" json:"trip_id"`
	Code        string  `csv:"stop_code" json:"stop_code"`
	Name        string  `csv:"stop_name" json:"stop_name"`
	Description string  `csv:"stop_desc" json:"stop_desc"`
	Latitude    float64 `csv:"stop_lat" json:"stop_lat"`
	Longitude   float64 `csv:"stop_lon" json:"stop_lon"`
	ZoneID      string  `csv:"zone_id" json:"zone_id"`
	Type        string  `csv:"location_type" json:"location_type"`
	Parent      string  `csv:"parent_station" json:"parent_station"`
}

func (g *Graph) GetStop(stopId string) Stop {
	return g.Stops[g.stopId2index[stopId]]
}

func (g *Graph) GetStopIndex(stopId string) int {
	if v, ok := g.stopId2index[stopId]; ok {
		return v
	}
	return -1
}

func (g *Graph) GetStopNode(stopId string) int {
	if v, ok := g.stopId2node[stopId]; ok {
		return v
	}
	return -1
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

func (g *Graph) SetIndex() {
	g.stopId2index = map[string]int{}
	g.stopId2node = map[string]int{}

	for i, s := range g.Stops {
		g.stopId2index[s.ID] = i
	}
	for i, n := range g.Nodes {
		g.stopId2node[n.PlaceID] = i
	}
}
