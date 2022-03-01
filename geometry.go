package gomap

import (
	geojson "github.com/takoyaki-3/go-geojson"
	json "github.com/takoyaki-3/go-json"
)

func (g *Graph) MakeLineString(latlons []int) (string, error) {

	var coordinates [][]float64

	for _, v := range latlons {
		coordinates = append(coordinates, []float64{g.Nodes[v].Lon, g.Nodes[v].Lat})
	}

	geom := geojson.NewLineString(coordinates, nil)
	return json.DumpToString(geom)
}
