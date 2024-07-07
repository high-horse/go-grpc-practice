package main

import (
	"log"

	"github.com/shirou/gopsutil/v3/process"
)

func getProcessInfo() ([]*process.Process, error) {

	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error retriving processes : %v", err)
	}

	return processes, nil
}
