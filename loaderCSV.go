package gomap

import csvtag "github.com/takoyaki-3/go-csv-tag/v3"

func LoadCSV(nodeFileName, edgeFileName string) (g *Graph, err error) {

	g = new(Graph)

	// 頂点情報読み込み
	if ok := csvtag.LoadFromPath(nodeFileName, &g.Nodes); ok != nil {
		return g, err
	}

	// 辺情報読み込み
	if ok := csvtag.LoadFromPath(edgeFileName, &g.Edges); ok != nil {
		return g, err
	}

	// 始点・終点に辺を追加
	for ei, edge := range g.Edges {
		g.FromEdges[edge.FromNode] = append(g.FromEdges[edge.FromNode], ei)
		g.ToEdges[edge.ToNode] = append(g.ToEdges[edge.ToNode], ei)
	}

	return g, nil
}

func LoadCSVWithStop(nodeFileName, edgeFileName, stopFileName string) (g *Graph, err error) {
	if g, err := LoadCSV(nodeFileName, edgeFileName); err != nil {
		return g, err
	}
	g.AddStopsFromCSV(stopFileName)
	return g, err
}

func (g *Graph) DumpCSV(nodeFileName, edgeFileName, stopFileName string) error {
	if err := csvtag.DumpToFile(g.Nodes, nodeFileName); err != nil {
		return err
	}
	if err := csvtag.DumpToFile(g.Edges, edgeFileName); err != nil {
		return err
	}
	if err := csvtag.DumpToFile(g.Stops, stopFileName); err != nil {
		return err
	}
	return nil
}

func (g *Graph) AddStopsFromCSV(stopFileName string) (err error) {
	// 停留所情報読み込み
	if ok := csvtag.LoadFromPath(stopFileName, &g.Stops); ok != nil {
		return err
	}

	h3indexes := g.MakeH3Index(9)

	// 停留所indexの作成
	g.stopId2index = map[string]int{}
	g.stopId2node = map[string]int{}
	for i, s := range g.Stops {
		n := Node{
			Lat:     s.Latitude,
			Lon:     s.Longitude,
			PlaceID: s.ID,
			Type:    "stop",
		}
		ni := len(g.Nodes)
		g.stopId2index[s.ID] = i
		g.stopId2node[s.ID] = ni
		g.Nodes = append(g.Nodes, n)
		g.FromEdges = append(g.FromEdges, []int{})
		g.ToEdges = append(g.ToEdges, []int{})

		// 近くの道路へ接続
		nearestNode := g.FindNode(h3indexes, n, 9)
		if nearestNode < 0 {
			continue
		}
		cost := HubenyDistance(g.Nodes[nearestNode], n)

		ei := len(g.Edges)
		g.Edges = append(g.Edges, Edge{
			FromNode: nearestNode,
			ToNode:   ni,
			Weight:   cost,
		})
		g.FromEdges[nearestNode] = append(g.FromEdges[ni], nearestNode)
		g.ToEdges[ni] = append(g.ToEdges[ni], ei)

		ei = len(g.Edges)
		g.Edges = append(g.Edges, Edge{
			FromNode: ni,
			ToNode:   nearestNode,
			Weight:   cost,
		})
		g.FromEdges[ni] = append(g.FromEdges[ni], ei)
		g.ToEdges[nearestNode] = append(g.FromEdges[nearestNode], ei)
	}

	return nil
}
