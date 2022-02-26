package gomap

import (
	json "github.com/takoyaki-3/go-json"
	geojson "github.com/takoyaki-3/go-geojson"
)

func DumpGeoJSON(g *Graph,nodeFileName, edgeFileName string)error{

	fc := geojson.FeatureCollection{
		Type: "FeatureCollection",
	}
	for _,e:=range g.Edges{
		from := g.Nodes[e.FromNode]
		to := g.Nodes[e.ToNode]
		line := [][]float64{
			[]float64{from.Lon,from.Lat},
			[]float64{to.Lon,to.Lat},
		}
		f := geojson.Feature{
			Type: "Feature",
			Geometry: *geojson.NewLineString(line,nil),
		}
		fc.Features = append(fc.Features, f)
	}

	if err := json.DumpToFile(fc,edgeFileName);err!=nil{
		return err
	}

	fc = geojson.FeatureCollection{
		Type: "FeatureCollection",
	}
	for _,n:=range g.Nodes{
		geom := geojson.Geometry{
			Type: "Point",
			Coordinates: []float64{n.Lon,n.Lat},
		}
		f := geojson.Feature{
			Type: "Feature",
			Geometry: geom,
		}
		fc.Features = append(fc.Features, f)
	}

	return json.DumpToFile(fc,nodeFileName)
}
