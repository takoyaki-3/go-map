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

	if err := gm.DumpCSV(g, "./sample/nodes.csv", "./sample/edges.csv"); err != nil {
		log.Fatalln(err)
	}

	indexes := gm.MakeH3Index(g, 9)
	n := gm.Node{}
	gm.FindNode(g, indexes, n, 9)
}
