package geoindex

type GeoIndex struct {
	resolution Meters
	index      map[Cell]interface{}
	newEntry   func() interface{}
}

// Creates new geo index with resolution a function that returns a new entry that is stored in each cell.
func NewGeoIndex(resolution Meters, newEntry func() interface{}) *GeoIndex {
	return &GeoIndex{resolution, make(map[Cell]interface{}), newEntry}
}

func (i *GeoIndex) Clone() *GeoIndex {
	clone := &GeoIndex{
		resolution: i.resolution,
		index:      make(map[Cell]interface{}, len(i.index)),
		newEntry:   i.newEntry,
	}
	for k, v := range i.index {
		set, ok := v.(set)
		if !ok {
			panic("Cannot cast value to set")
		}
		clone.index[k] = set.Clone()
	}

	return clone
}

// AddEntryAt adds an entry if missing, returns the entry at specific position.
func (geoIndex *GeoIndex) AddEntryAt(point Point) interface{} {
	square := CellOf(point, geoIndex.resolution)

	if _, ok := geoIndex.index[square]; !ok {
		geoIndex.index[square] = geoIndex.newEntry()
	}

	return geoIndex.index[square]
}

// GetEntryAt gets an entry from the geoindex, if missing returns an empty entry without changing the index.
func (geoIndex *GeoIndex) GetEntryAt(point Point) interface{} {
	square := CellOf(point, geoIndex.resolution)

	entries, ok := geoIndex.index[square]
	if !ok {
		return geoIndex.newEntry()
	}

	return entries
}

// Range returns the index entries within lat, lng range.
func (geoIndex *GeoIndex) Range(topLeft Point, bottomRight Point) []interface{} {
	topLeftIndex := CellOf(topLeft, geoIndex.resolution)
	bottomRightIndex := CellOf(bottomRight, geoIndex.resolution)

	return geoIndex.get(bottomRightIndex.X, topLeftIndex.X, topLeftIndex.Y, bottomRightIndex.Y)
}

func (geoIndex *GeoIndex) get(minx int, maxx int, miny int, maxy int) []interface{} {
	entries := make([]interface{}, 0, 0)

	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			if indexEntry, ok := geoIndex.index[Cell{x, y}]; ok {
				entries = append(entries, indexEntry)
			}
		}
	}

	return entries
}

func (g *GeoIndex) getCells(minx int, maxx int, miny int, maxy int) []Cell {
	indices := make([]Cell, 0)

	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			indices = append(indices, Cell{x, y})
		}
	}

	return indices
}
