package gomap

import (
	"io"
	"os"
	"runtime"

	"github.com/cheggaaa/pb"
	"github.com/qedus/osmpbf"
)

type Opts struct {
	ExclusionList map[string]bool
}

type option func(*Opts)

func LoadOSM(filename string, opts ...option) (*Graph, error) {

	g := new(Graph)

	// OSMファイル読み込み
	f, err := os.Open(filename)
	if err != nil {
		return g, err
	}
	defer f.Close()
	stat, _ := f.Stat()
	filesiz := int(stat.Size() / 1024)

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return g, err
	}

	exclusionList := map[string]bool{}
	// exclusionList["motorway"] = true
	// exclusionList["bus_guideway"] = true
	// exclusionList["raceway"] = true
	// exclusionList["busway"] = true
	// exclusionList["cycleway"] = true
	// exclusionList["proposed"] = true
	// exclusionList["construction"] = true
	// exclusionList["motorway_junction"] = true
	// exclusionList["platform"] = true

	parms := &Opts{
		ExclusionList: exclusionList,
	}
	for _, opt := range opts {
		opt(parms)
	}

	// 一時記憶用変数
	latlons := map[int64]Node{}
	ways := map[int64][]int64{}
	edgeTypes := map[int64]string{}
	usednode := map[int64]bool{}

	nc, wc, rc := int64(0), int64(0), int64(0)
	pb.New(filesiz).SetUnits(pb.U_NO)
	for i := 0; ; i++ {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return g, err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				latlons[v.ID] = Node{
					Lat: v.Lat,
					Lon: v.Lon}
				nc++
			case *osmpbf.Way:
				if _, ok := v.Tags["highway"]; !ok {
					continue
				}
				if _, ok := parms.ExclusionList[v.Tags["highway"]]; ok {
					continue
				}

				nodes := []int64{}
				for _, v := range v.NodeIDs {
					nodes = append(nodes, v)
					usednode[v] = true
				}
				ways[v.ID] = nodes
				edgeTypes[v.ID] = v.Tags["highway"]
				wc++
			case *osmpbf.Relation:
				rc++
			default:
				return g, err
			}
		}
	}

	nodeid := NewReplace()

	for wid, v := range ways {
		for i := 1; i < len(v); i++ {
			e := Edge{}
			e.Weight = HubenyDistance(latlons[v[i-1]], latlons[v[i]])
			node1 := nodeid.AddReplace(v[i-1])
			node2 := nodeid.AddReplace(v[i])

			e.FromNode = node1
			e.ToNode = node2
			e.Type = edgeTypes[wid]
			g.Edges = append(g.Edges, e)
		}
	}
	for k, v := range latlons {
		// 未使用頂点は追加する必要なし
		if _, ok := usednode[k]; !ok {
			continue
		}
		id := nodeid.AddReplace(k)

		for len(g.Nodes) <= id {
			g.Nodes = append(g.Nodes, Node{})
		}
		g.Nodes[id] = v
	}

	// 始終点辺を追加
	g.FromEdges = make([][]int, len(g.Nodes))
	g.ToEdges = make([][]int, len(g.Nodes))
	for ei, e := range g.Edges {
		g.FromEdges[e.FromNode] = append(g.FromEdges[e.FromNode], ei)
		g.ToEdges[e.ToNode] = append(g.ToEdges[e.ToNode], ei)
	}

	return g, nil
}

// 置き換え関数
type Replace struct {
	Id2Str map[int]int64
	Str2Id map[int64]int
}

func (s *Replace) AddReplace(str int64) int {
	if val, ok := s.Str2Id[str]; ok {
		return val
	}
	id := len(s.Str2Id)
	s.Str2Id[str] = id
	s.Id2Str[id] = str
	return id
}

func (s *Replace) AddReplaceIndex(str int64, index int) int {
	if val, ok := s.Str2Id[str]; ok {
		return val
	}
	s.Str2Id[str] = index
	s.Id2Str[index] = str
	return index
}

func NewReplace() *Replace {
	s := new(Replace)
	s.Str2Id = map[int64]int{}
	s.Id2Str = map[int]int64{}
	return s
}
