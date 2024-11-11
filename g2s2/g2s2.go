package g2s2

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io/ioutil"
// 	"s2id-experiment/s2id"

// 	"github.com/golang/geo/s2"
// 	geojson "github.com/paulmach/go.geojson"
// )

// func getCoordinates(fc s2id.GeoJSONFeature) []s2.Point {
// 	points := make([]s2.Point, 0)

// 	if len(fc.Geometry.Coordinates) == 0 {
// 		return points
// 	}

// 	// for _, polygon := range fc.Geometry.Coordinates {
// 	// 	for _, ring := range polygon {
// 	// 		for _, coordinate := range ring {
// 	// 			latLong := s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
// 	// 			points = append(points, latLong)
// 	// 		}
// 	// 	}

// 	// }
// 	for _, coordinate := range fc.Geometry.Coordinates[0] {
// 		latLong := s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
// 		points = append(points, latLong)
// 	}

// 	return points
// }

// func getChildren(cellID s2.CellID, level int) []s2.CellID {
// 	children := make([]s2.CellID, 0)

// 	if cellID.Level() >= level {
// 		return []s2.CellID{cellID.Parent(level)}
// 	}

// 	i := cellID.ChildBeginAtLevel(level)
// 	for i != cellID.ChildEndAtLevel(level) {
// 		children = append(children, i)
// 	}
// 	return children
// }

// func unique(intSlice []uint64) []uint64 {
// 	keys := make(map[uint64]bool)
// 	var list []uint64
// 	for _, entry := range intSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }

// func getS2Ids(points []s2.Point, level int, maxCells int) []uint64 {

// 	regionCoverer := s2.RegionCoverer{
// 		MinLevel: level,
// 		MaxLevel: level,
// 		MaxCells: maxCells}

// 	var loops []*s2.Loop
// 	loops = append(loops, s2.LoopFromPoints(points))

// 	var coverings []s2.CellUnion
// 	coverings = append(coverings, regionCoverer.Covering(s2.PolygonFromLoops(loops)))
// 	for _, point := range points {
// 		coverings = append(coverings, regionCoverer.CellUnion(point))
// 	}

// 	var coveringsUpdated []s2.CellUnion
// 	if level > 0 {
// 		for _, cells := range coverings {
// 			for _, cell := range cells {
// 				coveringsUpdated = append(coveringsUpdated, getChildren(cell, level))
// 			}
// 		}
// 	}

// 	var s2IDS []uint64
// 	for _, cells := range coveringsUpdated {
// 		for _, cell := range cells {
// 			s2IDS = append(s2IDS, uint64(cell))
// 		}
// 	}

// 	return unique(s2IDS)
// }

// func parseArgument() (fc *geojson.Feature, level int, delimiter string, err error) {
// 	level = 1

// 	raw, err := ioutil.ReadFile("geojson.geojson")
// 	if err != nil {
// 		return nil, -1, "", err
// 	}

// 	fc, err = geojson.UnmarshalFeature(raw)
// 	if err != nil {
// 		return nil, -1, "", err
// 	}

// 	if fc == nil {
// 		return nil, -1, "", errors.New("looks like bad GeoJson file. Make sure the root is a 'Feature'")
// 	}

// 	return fc, level, delimiter, nil
// }

// // Run as CLI entry point
// func Run() error {

// 	const MaxCells = 1
// 	const level = 1
// 	const delimiter = ","

// 	// fc, level, delimiter, err := parseArgument()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	geojson, _ := readGeoJSON("asd.geojson")

// 	points := getCoordinates(geojson.Features[0])
// 	s2IDs := getS2Ids(points, level, MaxCells)

// 	fmt.Printf("S2IDs Level %v :\n", level)
// 	fmt.Println(len(s2IDs))

// 	for i, s2id := range s2IDs {
// 		fmt.Print(s2id)
// 		if i < (len(s2IDs) - 1) {
// 			fmt.Print(fmt.Sprintf("%s ", delimiter))
// 		}
// 	}

// 	fmt.Println("")
// 	return nil
// }

// func readGeoJSON(filename string) (s2id.GeoJSON, error) {
// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return s2id.GeoJSON{}, err
// 	}

// 	// var geoJSON s2id.GeoJSON
// 	// if err := json.Unmarshal(data, &geoJSON); err != nil {
// 	// 	return nil, err
// 	// }

// 	// return &geoJSON, nil

// 	var asd s2id.GeoJSON
// 	if err := json.Unmarshal(data, &asd); err != nil {
// 		return s2id.GeoJSON{}, err
// 	}

// 	return asd, nil
// }
