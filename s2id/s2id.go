package s2id

import (
	"github.com/golang/geo/s2"
)

type GeoJSON struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

type GeoJSONFeature struct {
	Type       string          `json:"type"`
	Geometry   GeoJSONGeometry `json:"geometry"`
	Properties GeoJSONProperty `json:"properties"`
}

type GeoJSONProperty struct {
	Id        string `json:"ID_Block"`
	BlockName string `json:"NO_BLOK"`
}

type GeoJSONGeometry struct {
	Type        string       `json:"type"`
	Coordinates MultiPolygon `json:"coordinates"`
}

type MultiPolygon [][][][]float64
type Polygon [][][]float64

func CreateLoopFromMultiPolygon(multiPolygon MultiPolygon) (*s2.Loop, []s2.Point) {
	var points []s2.Point

	for _, polygon := range multiPolygon {
		for _, ring := range polygon {
			for _, coord := range ring {
				lat := coord[1]
				lng := coord[0]
				points = append(points, s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)))
			}
		}
	}

	return s2.LoopFromPoints(points), points
}

func CreateLoopFromPolygon(multiPolygon Polygon) (*s2.Loop, []s2.Point) {
	var points []s2.Point

	for _, polygon := range multiPolygon {
		for _, coord := range polygon {
			lat := coord[1]
			lng := coord[0]
			points = append(points, s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)))
		}
	}

	return s2.LoopFromPoints(points), points
}

func LoopCovering(loop *s2.Loop) []s2.CellID {
	loops := []*s2.Loop{loop}
	poly := s2.PolygonFromLoops(loops)

	rc := &s2.RegionCoverer{MinLevel: 1, MaxLevel: 16, MaxCells: 100}
	asd := s2.Region(poly)
	cellUnion := rc.Covering(asd)

	var cellIDs []s2.CellID
	for _, cellID := range cellUnion {
		cellIDs = append(cellIDs, cellID)
	}

	return cellIDs
}
