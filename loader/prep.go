package loader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path"
)

// prepare runs any special hacks for prepping the data before passing it onto
// the loader
func prepare(url, dir string) error {
	switch url {
	case "http://www.nyc.gov/html/dot/downloads/misc/siferry-gtfs.zip":
		return siFerry(dir)
	}

	return nil
}

// siFerry preps the Staten Island Ferry download by adding a direction id
// to files and fixing the header name of trip_headsign
func siFerry(dir string) error {
	var directionID string
	var err error

	// Get the incoming file
	tripFile := path.Join(dir, "trips.txt")
	inFH, err := os.Open(tripFile)
	if err != nil {
		return err
	}
	r := csv.NewReader(inFH)
	r.LazyQuotes = true

	// Create an outgoing csv file for transformed data
	w, outFH := writecsvtmp(dir)
	defer outFH.Close()
	defer os.Remove(outFH.Name())

	header, err := r.Read()
	if err != nil {
		return err
	}

	// Fix incorrectly named header
	for k, v := range header {
		if v == "headsign" {
			header[k] = "trip_headsign"
		}
	}

	headsignIdx := find(header, "trip_headsign")

	// Add a direction id header
	header = append(header, "direction_id")
	err = w.Write(header)
	if err != nil {
		return err
	}

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		headsign := rec[headsignIdx]

		// SI Ferry has two destinations, add a 0/1 direction ID
		if headsign == "To Whitehall" {
			directionID = "0"
		} else {
			directionID = "1"
		}

		rec = append(rec, directionID)

		err = w.Write(rec)
		if err != nil {
			return err
		}

	}

	w.Flush()
	err = outFH.Close()
	if err != nil {
		return err
	}

	log.Println("what is the name?", outFH.Name())

	err = os.Rename(outFH.Name(), tripFile)
	if err != nil {
		return err
	}

	return nil
}