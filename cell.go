package main

import "strconv"

var (
	minLon          = -180.0
	minLat          = -90.0
	latDegreeLength = Km(111.0)
	lonDegreeLength = Km(85.0)
)

type Meters float64

func Km(km float64) Meters {
	return Meters(km * 1000)
}

func Meter(meters float64) Meters {
	return Meters(meters)
}

type Cell struct {
	X int
	Y int
}

func (cell *Cell) Id() string {
	return strconv.Itoa(cell.X<<32 + cell.Y)
}

func ParseCell(id string) (Cell, error) {
	v, err := strconv.Atoi(id)
	if err != nil {
		return Cell{}, err
	}
	cell := Cell{
		X: v >> 32,
		Y: v & 0x00000000ffffffff,
	}
	return cell, nil
}

func CellOf(point Point, resolution Meters) Cell {
	x := int((-minLat + point.Lat()) * float64(latDegreeLength) / float64(resolution))
	y := int((-minLon + point.Lon()) * float64(lonDegreeLength) / float64(resolution))

	return Cell{x, y}
}

func CellOf2(lat float64, lng float64, resolution Meters) Cell {
	x := int((-minLat + lat) * float64(latDegreeLength) / float64(resolution))
	y := int((-minLon + lng) * float64(lonDegreeLength) / float64(resolution))
	return Cell{x, y}
}

func CellsOf(lat float64, lng float64, resolution Meters) []Cell {
	cell := CellOf2(lat, lng, resolution)
	cells := make([]Cell, 0, 9)
	indexes := []int{0, 1, 2}
	for _, x := range indexes {
		for _, y := range indexes {
			cells = append(cells, Cell{
				X: cell.X - 1 + x,
				Y: cell.Y - 1 + y,
			})
		}
	}
	return cells
}
