package gomap

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/takoyaki-3/go-map/pb"
)

// Load Protocol Buffer
func LoadFromPath(filename string) (*Graph, error) {

	g := new(Graph)

	// Read the existing graph.
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return &Graph{}, err
	}
	graph := &pb.Graph{}
	if err := proto.Unmarshal(in, graph); err != nil {
		return &Graph{}, err
	}

	for _, v := range graph.Node {
		g.Nodes = append(g.Nodes, Node{
			Type:    v.Type,
			Lat:     v.Lat,
			Lon:     v.Lon,
			PlaceID: v.PlaceId,
		})
	}

	g.FromEdges = make([][]int, len(g.Nodes))
	g.ToEdges = make([][]int, len(g.Nodes))
	for ei, v := range graph.Edge {
		nodes := []int{}
		for _, v := range v.ViaNodes {
			nodes = append(nodes, int(v))
		}
		g.Edges = append(g.Edges, Edge{
			Type:     v.Type,
			FromNode: int(v.From),
			ToNode:   int(v.To),
			Weight:   v.Weight,
		})
		g.FromEdges[v.From] = append(g.FromEdges[v.From], ei)
		g.ToEdges[v.To] = append(g.ToEdges[v.To], ei)
	}

	return g, nil
}

// Write to Protocol Buffer
func DumpToFile(g *Graph, filename string) error {
	edges := []*pb.Edge{}
	for _, v := range g.Edges {
		edges = append(edges, &pb.Edge{
			Type:   v.Type,
			From:   int64(v.FromNode),
			To:     int64(v.ToNode),
			Weight: v.Weight,
		})
	}

	nodes := []*pb.Node{}
	for _, v := range g.Nodes {
		nodes = append(nodes, &pb.Node{
			Type:    v.Type,
			Lat:     v.Lat,
			Lon:     v.Lon,
			PlaceId: v.PlaceID,
		})
	}

	graph := &pb.Graph{
		Node: nodes,
		Edge: edges,
	}

	// Write the new address book back to disk.
	out, err := proto.Marshal(graph)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, 0644)
}
