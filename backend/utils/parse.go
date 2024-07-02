package utils

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"vorto/backend/db"
)

func ParseFile(fileName string, problemID int64) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		loadNumber, _ := strconv.Atoi(parts[0])
		pickup := parseCoordinates(parts[1])
		dropoff := parseCoordinates(parts[2])

		load := db.Load{
			LoadNumber: loadNumber,
			PickupX:    pickup[0],
			PickupY:    pickup[1],
			DropoffX:   dropoff[0],
			DropoffY:   dropoff[1],
			FileName:   fileName,
			ProblemID:  problemID,
		}

		db.InsertLoad(load)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseCoordinates(coord string) [2]float64 {
	coord = strings.Trim(coord, "()")
	parts := strings.Split(coord, ",")
	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)
	return [2]float64{x, y}
}
