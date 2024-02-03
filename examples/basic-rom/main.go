package main

import (
	"fmt"

	"shadowglass/internal/gen/tradersv1"
)

func main() {
	ships, err := tradersv1.GetMyShips(tradersv1.GetMyShipsRequest{})
	fmt.Printf("My ships from go: %+v\n", ships.Ships)
	if err != nil {
		panic(err)
	}
}
