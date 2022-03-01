package main

import (
	"log"

	gm "github.com/takoyaki-3/go-map"
)

func main() {

	g, err := gm.LoadOSM("./sample/Tokyo.osm.pbf")
	if err != nil {
		log.Fatalln(err)
	}

	if err := g.DumpCSV("./sample/nodes.csv", "./sample/edges.csv"); err != nil {
		log.Fatalln(err)
	}

	indexes := g.MakeH3Index(9)
	n := gm.Node{}
	g.FindNode(indexes, n, 9)
}
