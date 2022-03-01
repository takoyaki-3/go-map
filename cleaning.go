package gomap

import ()

func (g *Graph) GetLargestGraph() {

	flag := make([]int, len(g.Nodes))
	nums := []int{}

	for pos, _ := range g.Nodes {
		if flag[int(pos)] != 0 {
			continue
		}
		nums = append(nums, 0)
		stack := []int{pos}
		for {
			if len(stack) == 0 {
				break
			}
			pos := stack[0]
			stack = stack[1:]
			if flag[pos] != 0 {
				continue
			}
			flag[pos] = len(nums)

			if len(g.Edges) <= int(pos) {
				continue
			}
			for _, ei := range g.FromEdges[pos] {
				e := g.Edges[ei]
				if flag[e.ToNode] == 0 {
					stack = append(stack, e.ToNode)
				}
			}
		}
	}
	nums = append(nums, 0)
	for k, _ := range g.Nodes {
		nums[flag[k]]++
	}

	maxFlag := 0
	for k, v := range nums {
		if nums[maxFlag] < v {
			maxFlag = k
		}
	}

	newG := Graph{}
	old2new := map[int]int{}
	for k, v := range g.Nodes {
		if flag[k] == maxFlag {
			old2new[k] = len(newG.Nodes)
			newG.Nodes = append(newG.Nodes, v)
		}
	}
	for _, e := range g.Edges {
		if flag[e.FromNode] == maxFlag {
			e.FromNode = old2new[e.FromNode]
			e.ToNode = old2new[e.ToNode]
			newG.Edges = append(newG.Edges, e)
		}
	}

	newG.FromEdges = make([][]int, len(newG.Nodes))
	newG.ToEdges = make([][]int, len(newG.Nodes))
	for ei, e := range newG.Edges {
		newG.FromEdges[e.FromNode] = append(newG.FromEdges[e.FromNode], ei)
		newG.ToEdges[e.ToNode] = append(newG.ToEdges[e.ToNode], ei)
	}

	*g = newG
}
