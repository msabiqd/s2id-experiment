package rewind

// var gj map[string]interface{}
// 	err := json.Unmarshal([]byte(geoJSON), &gj)

func rewindRings(rings []interface{}, outer bool) {
	if len(rings) == 0 {
		return
	}

	rewindRing(rings[0].([]interface{}), outer)
	for i := 1; i < len(rings); i++ {
		rewindRing(rings[i].([]interface{}), !outer)
	}
}

func Rewind(gj map[string]interface{}, outer bool) map[string]interface{} {
	switch gj["type"] {
	case "FeatureCollection":
		features := gj["features"].([]interface{})
		for i := range features {
			features[i] = Rewind(features[i].(map[string]interface{}), outer)
		}

	case "GeometryCollection":
		geometries := gj["geometries"].([]interface{})
		for i := range geometries {
			geometries[i] = Rewind(geometries[i].(map[string]interface{}), outer)
		}

	case "Feature":
		gj["geometry"] = Rewind(gj["geometry"].(map[string]interface{}), outer)

	case "Polygon":
		rewindRings(gj["coordinates"].([]interface{}), outer)

	case "MultiPolygon":
		coordinates := gj["coordinates"].([]interface{})
		for i := range coordinates {
			rewindRings(coordinates[i].([]interface{}), outer)
		}
	}

	return gj
}

func rewindRing(ring []interface{}, dir bool) {
	area := 0.0
	err := 0.0
	for i, len := 0, len(ring); i < len; i++ {
		j := (i + len - 1) % len
		pointI := ring[i].([]interface{})
		pointJ := ring[j].([]interface{})
		k := (pointI[0].(float64) - pointJ[0].(float64)) * (pointJ[1].(float64) + pointI[1].(float64))
		m := area + k
		if abs(area) >= abs(k) {
			err += area - m + k
		} else {
			err += k - m + area
		}
		area = m
	}
	if (area+err >= 0) != dir {
		reverseRing(ring)
	}
}

func reverseRing(ring []interface{}) {
	for i, j := 0, len(ring)-1; i < j; i, j = i+1, j-1 {
		ring[i], ring[j] = ring[j], ring[i]
	}
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
