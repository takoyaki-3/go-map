package gomap

import csvtag "github.com/artonge/go-csv-tag/v2"

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

func DumpCSV(g *Graph, nodeFileName, edgeFileName, stopFileName string) error {
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

func (g Graph) AddStopsFromCSV(stopFileName string) (err error) {
	// 停留所情報読み込み
	if ok := csvtag.LoadFromPath(stopFileName, &g.Stops); ok != nil {
		return err
	}

	// 停留所indexの作成
	for i, s := range g.Stops {
		g.stopId2index[s.ID] = i
	}

	return nil
}
