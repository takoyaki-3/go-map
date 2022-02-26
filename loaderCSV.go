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

func DumpCSV(g *Graph, nodeFileName, edgeFileName string) error {
	if err := csvtag.DumpToFile(g.Nodes, nodeFileName); err != nil {
		return err
	}
	if err := csvtag.DumpToFile(g.Edges, edgeFileName); err != nil {
		return err
	}
	return nil
}
