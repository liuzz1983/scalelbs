package geoindex

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

func ParseCell(id string) *Cell {
	v, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	cell := &Cell{
		X: v >> 32,
		Y: v & 0x00000000ffffffff,
	}
	return cell
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
