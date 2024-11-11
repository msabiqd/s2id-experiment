package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"s2id-experiment/rewind"
	"s2id-experiment/s2id"
	"strconv"
	"strings"

	"github.com/golang/geo/s2"
)

type Hasil struct {
	BlockID      string
	BlockName    string
	S2ids        []s2.CellID
	S2idsUint64  []uint64
	S2idsToken   []string
	S2idsStrings []string
	S2idsString  string
	Query        string
}

type ReqDadakan struct {
	BlockName string
	Lat       float64
	Long      float64
}

func main() {
	queries := ""
	query := "update public.block \n"
	querySet := "set s2ids = '{%s}', updated_by='sabiq', updated_at_utc0=EXTRACT(EPOCH FROM NOW()) \n"
	queryWhere := "where business_id = 'a2886ccd-5283-4c33-a408-87386f155b81' and name = '%s';"

	// read file
	geoJSON, err := readGeoJSON("geojsongunawan.json")
	if err != nil {
		log.Fatalf("Failed to read GeoJSON file: %v", err)
	}

	original := geoJSON

	// read file
	geoJSONInter, err := readGeoJSONInter("jakarta.geojson")
	if err != nil {
		log.Fatalf("Failed to read GeoJSON file: %v", err)
	}

	rewinded := rewind.Rewind(geoJSONInter, true)

	rewindJSON, err := json.Marshal(rewinded)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("rewinded.geojson", rewindJSON, 0644)

	var result []Hasil
	var pointz []s2.Point
	var loopz *s2.Loop
	// var reqDadakans []ReqDadakan
	var tokens []string
	// var sorted s2id.GeoJSON

	// sorted = geoJSON
	// sorted.Features = []s2id.GeoJSONFeature{}

	// for _, feature := range geoJSON.Features {
	// 	if feature.Properties.BlockName != "B4" && feature.Properties.BlockName != "No Block" {
	// 		geoJSON, err := readGeoJSON(fmt.Sprintf("sorted/%s.json", feature.Properties.BlockName))
	// 		if err != nil {
	// 			log.Fatalf("Failed to read GeoJSON file: %v", err)
	// 		}
	// 		sorted.Features = append(sorted.Features, geoJSON.Features...)
	// 	}
	// }

	for _, feature := range geoJSON.Features {
		// read file

		// if feature.Properties.BlockName != "B4" && feature.Properties.BlockName != "No Block" {
		// 	geoJSON, err := readGeoJSON(fmt.Sprintf("sorted/%s.json", feature.Properties.BlockName))
		// 	if err != nil {
		// 		log.Fatalf("Failed to read GeoJSON file: %v", err)
		// 	}
		// 	sorted.Features = append(sorted.Features, geoJSON.Features...)
		// }

		geometryType := feature.Geometry.Type

		//split purpose
		split := original
		split.Features = []s2id.GeoJSONFeature{feature}
		jsonBytes, err := json.Marshal(split)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(feature.Properties.BlockName)
		fmt.Println(jsonBytes)
		ioutil.WriteFile(fmt.Sprintf("splitgunawan/%s.txt", feature.Properties.BlockName), jsonBytes, 0644)
		fmt.Println("nulis")
		temp := Hasil{}
		temp.BlockID = feature.Properties.Id
		temp.BlockName = feature.Properties.BlockName
		cellIds := []s2.CellID{}

		switch geometryType {
		case "MultiPolygon":
			multiPolygon := feature.Geometry.Coordinates
			loop, zxc := s2id.CreateLoopFromMultiPolygon(multiPolygon)
			loopz = loop
			cellIDs := s2id.LoopCovering(loop)
			cellIds = append(cellIds, cellIDs...)
			pointz = zxc
		default:
			fmt.Printf("Unsupported geometry type: %s\n", geometryType)
		}

		if feature.Properties.BlockName == "E10" {
			fmt.Print("asd")
		}

		temp.S2ids = cellIds
		tokenz := []string{}
		for _, s2id := range cellIds {
			tokenz = append(tokenz, s2id.ToToken())
		}
		temp.S2idsToken = tokenz
		result = append(result, temp)
		if feature.Properties.BlockName == "E10" {
			for _, cellId := range cellIds {
				tokens = append(tokens, cellId.ToToken())
			}
		}
	}

	for i, res := range result {
		s2iduint64 := []uint64{}
		s2idstring := []string{}
		for _, s2id := range res.S2ids {
			s2iduint64 = append(s2iduint64, uint64(s2id))
			s2idstring = append(s2idstring, strconv.FormatUint(uint64(s2id), 10))
		}

		// s2idstring = res.S2idsToken

		if res.BlockName == "E10" {
			fmt.Print("asd")
		}
		result[i].S2idsUint64 = s2iduint64
		result[i].S2idsStrings = s2idstring
		result[i].S2idsString = strings.Join(s2idstring, ",")

		result[i].Query = query
		result[i].Query += " " + fmt.Sprintf(querySet, result[i].S2idsString)
		result[i].Query += " " + fmt.Sprintf(queryWhere, result[i].BlockName)
		queries += result[i].Query
		queries += "\n"
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ioutil.WriteFile("result.geojson", jsonBytes, 0644)

	jsonBytes, err = json.Marshal(queries)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ioutil.WriteFile("queries.txt", jsonBytes, 0644)

	// jsonBytes, err = json.Marshal(sorted)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// ioutil.WriteFile("sortedddd.json", jsonBytes, 0644)

	jsonBytes, err = json.Marshal(tokens)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ioutil.WriteFile("tokens.json", jsonBytes, 0644)

	// asd := s2.PointFromLatLng(s2.LatLngFromDegrees(-0.3689108, 101.695171))
	asd := s2.PointFromLatLng(s2.LatLngFromDegrees(-6.319850, 107.129511))
	qwe := asd.Contains(pointz[0])
	contain := loopz.ContainsPoint(asd)
	// result[0].S2ids[0].
	fmt.Println(asd)
	fmt.Println(qwe)
	fmt.Println(contain)
	// g2s2.Run()
}

func readGeoJSON(filename string) (s2id.GeoJSON, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return s2id.GeoJSON{}, err
	}

	var asd s2id.GeoJSON
	if err := json.Unmarshal(data, &asd); err != nil {
		return s2id.GeoJSON{}, err
	}

	return asd, nil
}

func readGeoJSONInter(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var asd map[string]interface{}
	if err := json.Unmarshal(data, &asd); err != nil {
		return map[string]interface{}{}, err
	}

	return asd, nil
}
