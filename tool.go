package gomap

func CutGraph(g *Graph, leftUp Node, rightDown Node) error {
	old2new := map[int]int{}
	newG := Graph{}
	for old, v := range g.Nodes {
		if leftUp.Lon <= v.Lon && v.Lon <= rightDown.Lon {
			if rightDown.Lat <= v.Lat && v.Lat <= leftUp.Lat {
				old2new[old] = len(newG.Nodes)
				newG.Nodes = append(newG.Nodes, v)
			}
		}
	}
	for _, e := range g.Edges {
		if newFrom, ok := old2new[e.FromNode]; ok {
			if newTo, ok := old2new[e.ToNode]; ok {
				e.FromNode = newFrom
				e.ToNode = newTo
				ei := len(newG.Edges)
				newG.FromEdges[newFrom] = append(newG.FromEdges[newFrom], ei)
				newG.ToEdges[newTo] = append(newG.ToEdges[newTo], ei)
				newG.Edges = append(newG.Edges, e)
			}
		}
	}

	*g = newG

	return nil
}
