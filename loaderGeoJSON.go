package gomap

import (
	geojson "github.com/takoyaki-3/go-geojson"
	json "github.com/takoyaki-3/go-json"
)

func DumpGeoJSON(g *Graph, nodeFileName, edgeFileName, stopFileName string) error {

	// 辺情報
	fc := geojson.FeatureCollection{
		Type: "FeatureCollection",
	}
	for _, e := range g.Edges {
		from := g.Nodes[e.FromNode]
		to := g.Nodes[e.ToNode]
		line := [][]float64{
			[]float64{from.Lon, from.Lat},
			[]float64{to.Lon, to.Lat},
		}
		f := geojson.Feature{
			Type:     "Feature",
			Geometry: *geojson.NewLineString(line, nil),
		}
		fc.Features = append(fc.Features, f)
	}
	if err := json.DumpToFile(fc, edgeFileName); err != nil {
		return err
	}

	// 頂点情報
	fc = geojson.FeatureCollection{
		Type: "FeatureCollection",
	}
	for _, n := range g.Nodes {
		geom := geojson.Geometry{
			Type:        "Point",
			Coordinates: []float64{n.Lon, n.Lat},
		}
		f := geojson.Feature{
			Type:     "Feature",
			Geometry: geom,
		}
		fc.Features = append(fc.Features, f)
	}
	if err := json.DumpToFile(fc, nodeFileName); err != nil {
		return err
	}

	// 停留所情報
	fc = geojson.FeatureCollection{
		Type: "FeatureCollection",
	}
	for _, s := range g.Stops {
		geom := geojson.Geometry{
			Type:        "Point",
			Coordinates: []float64{s.Longitude, s.Latitude},
		}
		f := geojson.Feature{
			Type:     "Feature",
			Geometry: geom,
		}
		fc.Features = append(fc.Features, f)
	}
	if err := json.DumpToFile(fc, stopFileName); err != nil {
		return err
	}
	return nil
}
