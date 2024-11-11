package geojsonsort

import (
	"math"
	"s2id-experiment/s2id"
	"sort"
)

// Point represents a point with x and y coordinates
type Point struct {
	X, Y float64
}

func (p Point) Less(p2 Point) bool {
	if p.Y < p2.Y || (p.Y == p2.Y && p.X < p2.X) {
		return true
	}
	return false
}

// Angle calculates the angle of the point relative to the centroid
func (p Point) Angle(centroid Point) float64 {
	return math.Atan2(p.Y-centroid.Y, p.X-centroid.X)
}

// Centroid calculates the centroid of a polygon
func Centroid(polygon [][]float64) Point {
	var centroid Point
	var signedArea float64

	for i := range polygon {
		x0, y0 := polygon[i][0], polygon[i][1]
		x1, y1 := polygon[(i+1)%len(polygon)][0], polygon[(i+1)%len(polygon)][1]
		a := x0*y1 - x1*y0
		signedArea += a
		centroid.X += (x0 + x1) * a
		centroid.Y += (y0 + y1) * a
	}

	signedArea *= 0.5
	centroid.X /= (6 * signedArea)
	centroid.Y /= (6 * signedArea)

	return centroid
}

func SortGeoJSONCounterClockwise(geoJSON s2id.GeoJSONGeometry) s2id.GeoJSONGeometry {
	// for i, polygon := range geoJSON.Coordinates {
	// 	centroid := Centroid(polygon)
	// 	points := make([]Point, len(polygon))

	// 	for k, coord := range polygon {
	// 		points[k] = Point{X: coord[0], Y: coord[1]}
	// 	}

	// 	sort.Slice(points, func(a, b int) bool {
	// 		return points[a].Angle(centroid) < points[b].Angle(centroid)
	// 	})

	// 	for k, point := range points {
	// 		geoJSON.Coordinates[i][k] = []float64{point.X, point.Y}
	// 	}
	// }
	// return geoJSON
	// for i, multipolygon := range geoJSON.Coordinates {
	// 	for j, polygon := range multipolygon {
	// 		centroid := Centroid(polygon)
	// 		points := make([]Point, len(polygon))

	// 		for k, coord := range polygon {
	// 			points[k] = Point{X: coord[0], Y: coord[1]}
	// 		}

	// 		sort.Slice(points, func(a, b int) bool {
	// 			return points[a].Angle(centroid) < points[b].Angle(centroid)
	// 		})

	// 		for k, point := range points {
	// 			geoJSON.Coordinates[i][j][k] = []float64{point.X, point.Y}
	// 		}
	// 	}
	// }

	for i, multipolygon := range geoJSON.Coordinates {
		for j, polygon := range multipolygon {
			points := make([]Point, len(polygon))

			for k, coord := range polygon {
				points[k] = Point{X: coord[0], Y: coord[1]}
			}

			sortedPoints := SortPoints(points)

			for k, point := range sortedPoints {
				geoJSON.Coordinates[i][j][k] = []float64{point.X, point.Y}
			}
		}
	}

	return geoJSON
}

// SortPoints sorts points in the order specified
func SortPoints(points []Point) []Point {
	n := len(points)
	if n < 2 {
		return points
	}

	leftMost := points[0]
	leftMostIndex := 0

	// Find the leftmost (or bottom-leftmost) point
	for i := 1; i < n; i++ {
		if points[i].Less(leftMost) {
			leftMost = points[i]
			leftMostIndex = i
		}
	}

	// Place the leftmost point at the beginning
	points[0], points[leftMostIndex] = points[leftMostIndex], points[0]

	// Sort the rest of the points based on polar angle relative to the leftmost point
	sort.Slice(points[1:], func(i, j int) bool {
		pi, pj := points[i+1], points[j+1]
		d1x, d1y := pi.X-leftMost.X, pi.Y-leftMost.Y
		d2x, d2y := pj.X-leftMost.X, pj.Y-leftMost.Y
		return math.Atan2(d1y, d1x) < math.Atan2(d2y, d2x)
	})

	return points
}

// func main() {
// 	// Sample GeoJSON multipolygon
// 	geoJSONStr := `{
// 		"type": "MultiPolygon",
// 		"coordinates": [
// 			[
// 				[
// 					[1.0, 1.0],
// 					[4.0, 1.0],
// 					[4.0, 4.0],
// 					[1.0, 4.0],
// 					[1.0, 1.0]
// 				]
// 			],
// 			[
// 				[
// 					[10.0, 10.0],
// 					[14.0, 10.0],
// 					[14.0, 14.0],
// 					[10.0, 14.0],
// 					[10.0, 10.0]
// 				]
// 			]
// 		]
// 	}`

// 	var geoJSON GeoJSON
// 	err := json.Unmarshal([]byte(geoJSONStr), &geoJSON)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i, multipolygon := range geoJSON.Coordinates {
// 		for j, polygon := range multipolygon {
// 			centroid := Centroid(polygon)
// 			points := make([]Point, len(polygon))

// 			for k, coord := range polygon {
// 				points[k] = Point{X: coord[0], Y: coord[1]}
// 			}

// 			sort.Slice(points, func(a, b int) bool {
// 				return points[a].Angle(centroid) < points[b].Angle(centroid)
// 			})

// 			for k, point := range points {
// 				geoJSON.Coordinates[i][j][k] = []float64{point.X, point.Y}
// 			}
// 		}
// 	}

// 	sortedGeoJSON, err := json.MarshalIndent(geoJSON, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(string(sortedGeoJSON))
// }
