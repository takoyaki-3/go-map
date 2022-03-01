package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	gm "github.com/takoyaki-3/go-map"
)

func main() {

	// Graph load
	fmt.Println("Graph load")
	g, err := gm.LoadOSM("./sample/Tokyo.osm.pbf")
	// g, err := gm.LoadFromPath("./sample/Tokyo.gomap.pbf")
	if err != nil {
		log.Fatalln(err)
	}
	g.GetLargestGraph()

	g.DumpGeoJSON("./sample/nodes.geojson", "./sample/edges.geojson", "./sample/stops.geojson")
	g.DumpCSV("./sample/nodes.csv", "./sample/edges.csv", "./sample/stops.csv")
	g.DumpToFile("./sample/Tokyo.gomap.pbf")
	// g := osm.Load("./kanto-latest.osm.pbf")
	// g := geojson.Load("./kanto-lines.geojson")

	// Make index
	fmt.Println("Make index")
	h3indexes := g.MakeH3Index(6)

	// Start server
	fmt.Println("start Server")
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		if v == nil {
			return
		}
		parm_v := map[string]float64{}
		must_parms := []string{"from_lat", "from_lon", "to_lat", "to_lon"}
		for _, k := range must_parms {
			if _, ok := v[k]; !ok {
				fmt.Fprintf(w, "{\"ErrorMessage\":\"Required parameters do not exist.\"}")
				return
			}
			f, err := strconv.ParseFloat(v[k][0], 64)
			if err != nil {
				log.Fatal(err)
				fmt.Fprintf(w, "{\"ErrorMessage\":\"Required parameters cannot be converted.\"}")
				return
			}
			parm_v[k] = f
		}

		q := gm.Query{}

		// Find node
		q.To = g.FindNode(h3indexes, gm.Node{
			Lat: parm_v["from_lat"],
			Lon: parm_v["from_lon"]}, 6)
		q.From = g.FindNode(h3indexes, gm.Node{
			Lat: parm_v["to_lat"],
			Lon: parm_v["to_lon"]}, 6)

		if q.To < 0 || q.From < 0 {
			fmt.Println("from node or to node is not found.")
			return
		}

		// Search
		o, err := g.Search(q)
		if err != nil {
			fmt.Println(err)
		}

		if len(o.Nodes) == 0 {
			fmt.Println("route not found")
			return
		}

		// Make GeoJSON
		rawJSON, err := g.MakeLineString(o.Nodes)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, rawJSON)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile("./index.html")
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(w, string(bytes))
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
