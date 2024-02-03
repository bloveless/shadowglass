package main

import "shadowglass/internal/gen/tradersv1"

type Traders struct{}

func (t Traders) GetMyShips(request *tradersv1.GetMyShipsRequest) tradersv1.Ships {
	return tradersv1.Ships{
		Ships: []*tradersv1.Ship{
			{Id: "1"},
			{Id: "2"},
			{Id: "3"},
		},
	}
}
