package utils

import (
	"math"
	"sort"
	"time"
	"vorto/backend/db"

	"github.com/gin-gonic/gin"
)

const (
	fixedCostPerDriver = 500
)

type Driver struct {
	Loads     []db.Load
	TotalTime float64
	Cost      float64
}

func CalculateResults(problemID int64) (gin.H, error) {
	start := time.Now()

	rows, err := db.Db.Query("SELECT loadNumber, pickupX, pickupY, dropoffX, dropoffY, fileName, problem_id FROM loads WHERE problem_id = ?", problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loadsByFile := make(map[string][]db.Load)
	for rows.Next() {
		var load db.Load
		if err := rows.Scan(&load.LoadNumber, &load.PickupX, &load.PickupY, &load.DropoffX, &load.DropoffY, &load.FileName, &load.ProblemID); err != nil {
			return nil, err
		}
		loadsByFile[load.FileName] = append(loadsByFile[load.FileName], load)
	}

	var drivers []Driver
	driversPerFile := make(map[string]int)
	fileDurations := make(map[string]float64)
	fileCosts := make(map[string]float64)

	for fileName, loads := range loadsByFile {
		fileStart := time.Now()
		fileDrivers, fileDriverCount := solveVRP(loads)
		fileDuration := time.Since(fileStart).Seconds()

		drivers = append(drivers, fileDrivers...)
		driversPerFile[fileName] = fileDriverCount

		totalFileTime := 0.0
		for _, driver := range fileDrivers {
			totalFileTime += driver.TotalTime
		}
		fileDurations[fileName] = fileDuration
		fileCosts[fileName] = (float64(fixedCostPerDriver) * float64(fileDriverCount)) + totalFileTime
	}

	totalDuration := time.Since(start).Seconds()

	totalSets := len(loadsByFile) // Each distinct fileName represents a problem set
	totalDrivers := len(drivers)
	totalCost := 0.0

	for _, cost := range fileCosts {
		totalCost += cost
	}
	meanCost := totalCost / float64(totalSets)

	kpis := gin.H{
		"number_of_drivers": totalDrivers,
		"total_cost":        totalCost,
		"mean_cost":         meanCost,
		"time_taken":        totalDuration,
		"drivers_per_file":  driversPerFile,
		"file_durations":    fileDurations,
		"file_costs":        fileCosts,
	}

	return gin.H{"drivers": drivers, "kpis": kpis}, nil
}

func solveVRP(loads []db.Load) ([]Driver, int) {
	var drivers []Driver
	var currentDriver Driver
	driverCount := 0

	// Sort loads by proximity to form clusters.
	sort.Slice(loads, func(i, j int) bool {
		return euclideanDistance(loads[i].PickupX, loads[i].PickupY, loads[i].DropoffX, loads[i].DropoffY) <
			euclideanDistance(loads[j].PickupX, loads[j].PickupY, loads[j].DropoffX, loads[j].DropoffY)
	})

	for len(loads) > 0 {
		if len(currentDriver.Loads) == 0 {
			// Start with the first load in the list if the driver has no loads yet.
			currentDriver.Loads = append(currentDriver.Loads, loads[0])
			currentDriver.TotalTime += euclideanDistance(0, 0, loads[0].PickupX, loads[0].PickupY) +
				euclideanDistance(loads[0].PickupX, loads[0].PickupY, loads[0].DropoffX, loads[0].DropoffY)
			loads = loads[1:]
		} else {
			lastLoad := currentDriver.Loads[len(currentDriver.Loads)-1]
			lastDropoffX, lastDropoffY := lastLoad.DropoffX, lastLoad.DropoffY

			// Find the next load within the cluster that minimizes the distance.
			closestLoadIndex := -1
			minDistance := math.MaxFloat64

			for i, load := range loads {
				distanceToPickup := euclideanDistance(lastDropoffX, lastDropoffY, load.PickupX, load.PickupY)
				if distanceToPickup < minDistance {
					minDistance = distanceToPickup
					closestLoadIndex = i
				}
			}

			if closestLoadIndex != -1 {
				load := loads[closestLoadIndex]
				distanceToPickup := euclideanDistance(lastDropoffX, lastDropoffY, load.PickupX, load.PickupY)
				distanceToDropoff := euclideanDistance(load.PickupX, load.PickupY, load.DropoffX, load.DropoffY)
				newTotalTime := currentDriver.TotalTime + distanceToPickup + distanceToDropoff

				// Check if adding this load would exceed the 12-hour limit, including return to depot.
				returnDistance := euclideanDistance(load.DropoffX, load.DropoffY, 0, 0)
				if newTotalTime+returnDistance > 720 {
					// If it exceeds, finalize the current driver and start a new one.
					if len(currentDriver.Loads) > 0 {
						currentDriver.TotalTime += euclideanDistance(lastDropoffX, lastDropoffY, 0, 0)
						currentDriver.Cost = currentDriver.TotalTime
						drivers = append(drivers, currentDriver)
						driverCount++
					}
					currentDriver = Driver{}
					continue
				}

				// Add the current load to the driver's list of loads.
				currentDriver.Loads = append(currentDriver.Loads, load)
				currentDriver.TotalTime = newTotalTime

				// Remove the load from the remaining loads.
				loads = append(loads[:closestLoadIndex], loads[closestLoadIndex+1:]...)
			}
		}
	}

	// Ensure the last driver also returns to the depot.
	if len(currentDriver.Loads) > 0 {
		lastLoad := currentDriver.Loads[len(currentDriver.Loads)-1]
		currentDriver.TotalTime += euclideanDistance(lastLoad.DropoffX, lastLoad.DropoffY, 0, 0)
		currentDriver.Cost = currentDriver.TotalTime
		drivers = append(drivers, currentDriver)
		driverCount++
	}

	return drivers, driverCount
}

func euclideanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
