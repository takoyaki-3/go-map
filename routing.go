package gomap

import (
	"math"

	"github.com/takoyaki-3/go-map/prioque"
)

type Query struct {
	From int
	To   int
}

type Output struct {
	Nodes []int
	Cost  float64
}

//
func Search(g *Graph, query Query) Output {

	l := len(g.Edges)
	if l < len(g.Nodes) {
		l = len(g.Nodes)
	}

	q := prioque.NewMinSet()
	cost := make([]float64, l)
	flag := make([]bool, l)
	before := make([]int, l)

	for k, _ := range cost {
		cost[k] = math.MaxFloat64
	}

	if len(g.Edges) <= int(query.From) || len(g.Edges) <= int(query.To) {
		return Output{}
	}

	cost[query.From] = 0.0
	before[query.From] = -2

	q.AddVal(query.From, 0.0)

	var pos int
	for q.Len() > 0 {
		pos = q.GetMin()
		if flag[pos] {
			continue
		}
		flag[pos] = true

		if pos == query.To {
			break
		}

		for _, eid := range g.FromEdges[pos] {
			e := g.Edges[eid]
			eto := e.ToNode
			if flag[eto] {
				continue
			}
			if cost[eto] <= cost[pos]+e.Weight {
				continue
			}
			cost[eto] = cost[pos] + e.Weight
			if len(e.ViaNodes) == 0 {
				before[eto] = pos
			} else {
				if eto != e.ViaNodes[len(e.ViaNodes)-1] {
					before[eto] = e.ViaNodes[len(e.ViaNodes)-1]
				}
				if e.ViaNodes[0] != pos {
					before[e.ViaNodes[0]] = pos
				}
				for k, v := range e.ViaNodes {
					if k == 0 {
						continue
					}
					if v != e.ViaNodes[k-1] {
						before[v] = e.ViaNodes[k-1]
					}
				}
			}
			q.AddVal(eto, cost[eto])
		}
	}

	// 出力
	out := Output{}
	if pos == query.To {
		out.Cost = cost[pos]
		out.Nodes = append(out.Nodes, pos)

		bef := before[pos]
		if bef == -1 {
			return Output{}
		}
		for bef != -2 {
			out.Nodes = append([]int{bef}, out.Nodes...)
			bef = before[bef]
		}
	}

	return out
}

func Voronoi(g *Graph, bases []int) map[int]int {
	// initialization
	q := prioque.NewMinSet()
	cost := make([]float64, len(g.Edges))
	flag := make([]bool, len(g.Edges))
	start_group := map[int]int{}

	counter := 0
	for k, _ := range cost {
		cost[k] = math.MaxFloat64
	}

	for _, v := range bases {
		cost[v] = 0.0
		q.AddVal(v, 0.0)
		start_group[int(v)] = counter % 20
		counter++
	}

	for q.Len() > 0 {
		pos := q.GetMin()
		if flag[pos] {
			continue
		}
		flag[pos] = true

		// グラフ拡張処理
		for _, eid := range g.FromEdges[pos] {
			e := g.Edges[eid]
			eto := e.ToNode
			if flag[eto] {
				continue
			}
			if cost[eto] <= cost[pos]+e.Weight {
				continue
			}
			cost[eto] = cost[pos] + e.Weight
			start_group[eto] = start_group[pos]
			q.AddVal(eto, cost[pos]+e.Weight)
		}
	}

	return start_group
}

func AllDistance(g *Graph, base []int) []float64 {
	q := prioque.NewMinSet()
	cost := make([]float64, len(g.Edges))
	flag := make([]bool, len(g.Edges))

	for k, _ := range cost {
		cost[k] = math.MaxFloat64
	}

	for _, v := range base {
		cost[v] = 0.0
		q.AddVal(v, 0.0)
	}

	for q.Len() > 0 {
		pos := q.GetMin()
		if flag[pos] {
			continue
		}
		flag[pos] = true
		for _, eid := range g.FromEdges[pos] {
			e := g.Edges[eid]
			eto := e.ToNode
			if flag[eto] {
				continue
			}
			if cost[eto] <= cost[pos]+e.Weight {
				continue
			}
			cost[eto] = cost[pos] + e.Weight
			q.AddVal(eto, cost[eto])
		}
	}
	return cost
}
