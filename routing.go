package gomap

import (
	"errors"
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

type cfb struct {
	cost   float64
	flag   bool
	before int
}

//
func Search(g *Graph, query Query) (Output, error) {

	q := prioque.NewMinSet()
	cfbs := make([]cfb, len(g.Nodes))

	for k, _ := range cfbs {
		cfbs[k].cost = math.MaxFloat64
		cfbs[k].before = -1
	}

	if len(g.Nodes) <= int(query.From) || len(g.Nodes) <= int(query.To) {
		return Output{}, errors.New("query node index not found.")
	}

	cfbs[query.From].cost = 0.0
	cfbs[query.From].before = -2

	q.AddVal(query.From, 0.0)

	var pos int
	for q.Len() > 0 {
		pos = q.GetMin()
		if cfbs[pos].flag {
			continue
		}
		cfbs[pos].flag = true

		if pos == query.To {
			break
		}

		for _, eid := range g.FromEdges[pos] {
			e := g.Edges[eid]
			eto := e.ToNode
			if cfbs[eto].flag {
				continue
			}
			etoCost := cfbs[pos].cost + e.Weight
			if cfbs[eto].cost <= etoCost {
				continue
			}
			cfbs[eto].cost = etoCost
			cfbs[eto].before = pos
			if len(e.ViaNodes) == 0 {
				cfbs[eto].before = pos
			} else {
				if eto != e.ViaNodes[len(e.ViaNodes)-1] {
					cfbs[eto].before = e.ViaNodes[len(e.ViaNodes)-1]
				}
				if e.ViaNodes[0] != pos {
					cfbs[e.ViaNodes[0]].before = pos
				}
				for k, v := range e.ViaNodes {
					if k == 0 {
						continue
					}
					if v != e.ViaNodes[k-1] {
						cfbs[v].before = e.ViaNodes[k-1]
					}
				}
			}
			q.AddVal(eto, etoCost)
		}
	}

	// 出力
	pos = query.To
	out := Output{}
	if pos == query.To {
		out.Cost = cfbs[pos].cost
		out.Nodes = append([]int{pos}, out.Nodes...)

		bef := cfbs[pos].before
		if bef == -1 {
			return Output{}, errors.New("path not found")
		}
		for bef != -2 {
			out.Nodes = append([]int{bef}, out.Nodes...)
			bef = cfbs[bef].before
		}
	}

	return out, nil
}

func Voronoi(g *Graph, bases []int) map[int]int {
	// initialization
	q := prioque.NewMinSet()
	cfbs := make([]cfb, len(g.Nodes))
	start_group := map[int]int{}

	counter := 0
	for k, _ := range cfbs {
		cfbs[k].before = -1
		cfbs[k].cost = math.MaxFloat64
	}

	for _, v := range bases {
		cfbs[v].cost = 0.0
		q.AddVal(v, 0.0)
		start_group[int(v)] = counter % 20
		counter++
	}

	for q.Len() > 0 {
		pos := q.GetMin()
		if cfbs[pos].flag {
			continue
		}
		cfbs[pos].flag = true

		// グラフ拡張処理
		for _, eid := range g.FromEdges[pos] {
			e := g.Edges[eid]
			eto := e.ToNode
			if cfbs[eto].flag {
				continue
			}
			if cfbs[eto].cost <= cfbs[pos].cost+e.Weight {
				continue
			}
			cfbs[eto].cost = cfbs[pos].cost + e.Weight
			start_group[eto] = start_group[pos]
			q.AddVal(eto, cfbs[pos].cost+e.Weight)
		}
	}

	return start_group
}

func AllDistance(g *Graph, base []int) []float64 {
	q := prioque.NewMinSet()
	cfbs := make([]cfb, len(g.Nodes))

	for k, _ := range cfbs {
		cfbs[k].cost = math.MaxFloat64
	}

	for _, v := range base {
		cfbs[v].cost = 0.0
		q.AddVal(v, 0.0)
	}

	for q.Len() > 0 {
		pos := q.GetMin()
		if cfbs[pos].flag {
			continue
		}
		cfbs[pos].flag = true
		for _, eid := range g.FromEdges[pos] {
			e := g.Edges[eid]
			eto := e.ToNode
			if cfbs[eto].flag {
				continue
			}
			if cfbs[eto].cost <= cfbs[pos].cost+e.Weight {
				continue
			}
			cfbs[eto].cost = cfbs[pos].cost + e.Weight
			q.AddVal(eto, cfbs[eto].cost)
		}
	}
	cost := make([]float64, len(g.Nodes))
	for k, v := range cfbs {
		cost[k] = v.cost
	}
	return cost
}
