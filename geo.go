package main

import (
	"sync"
	"time"
	"sort"
)

type IndexEntry struct {
	entries      map[string]*Pos
	lastInserted map[string]time.Time
	sync.RWMutex
}

func NewIndexEntry() *IndexEntry {
	return &IndexEntry{
		entries:      make(map[string]*Pos),
		lastInserted: make(map[string]time.Time),
	}
}

func (entry *IndexEntry) Add(pos *Pos) {
	entry.Lock()
	defer entry.Unlock()

	entry.lastInserted[pos.Id] = getNow()
	entry.entries[pos.Id] = pos
}

func (entry *IndexEntry) Remove(id string) {
	entry.Lock()
	defer entry.Unlock()
	delete(entry.entries, id)
}

func (entry *IndexEntry) Entries() []*Pos {
	entry.RLock()
	defer entry.RUnlock()

	points := make([]*Pos, 0, len(entry.entries))
	for _, v := range entry.entries {
		points = append(points, v)
	}
	return points
}

type GeoIndexer struct {
	resolution Meters
	expiredTimed time.Duration
	indexes    map[string]*IndexEntry
	sync.RWMutex
}

func NewGeoIndexer(resolution Meters, expiredTimed time.Duration) *GeoIndexer {
	return &GeoIndexer{
		indexes:    make(map[string]*IndexEntry),
		resolution: resolution,
		expiredTimed: expiredTimed,
	}
}

//expire the entries
func (indexer *GeoIndexer) expire() {

}

func (indexer *GeoIndexer) Cells(lat float64, lng float64) []Cell {
	cells := CellsOf(lat, lng, indexer.resolution)
	return cells
}

func (indexer *GeoIndexer) Cell(lat float64, lng float64) Cell {
	return CellOf2(lat, lng, indexer.resolution)
}

func (indexer *GeoIndexer) Add(pos *Pos) {

	cell := indexer.Cell(pos.Lat, pos.Lng)
	indexer.Lock()
	defer indexer.Unlock()

	entry, ok := indexer.indexes[cell.Id()]
	if !ok {
		entry = NewIndexEntry()
	}
	entry.Add(pos)
	indexer.indexes[cell.Id()] = entry
}

func (indexer *GeoIndexer) Get(cellId string) []*Pos {
	indexer.RLock()
	v, ok := indexer.indexes[cellId]
	indexer.RUnlock()

	if !ok {
		return nil
	}
	return v.Entries()
}

func (indexer *GeoIndexer) GetByLatLng( lat float64, lng float64 ) []*Pos {
	cell := indexer.Cell(lat, lng)
	return indexer.Get(cell.Id())
}

type MetersVector []Meters

// These methods exist to satisfy the sort.Interface for sorting.
func (s *MetersVector) Len() int           { return len(*s) }
func (s *MetersVector) Less(i, j int) bool { return (*s)[i] < (*s)[j] }
func (s *MetersVector) Swap(i, j int)      { (*s)[i], (*s)[j] = (*s)[j], (*s)[i] }


type GeoMultiIndexer struct {
	indexes []*GeoIndexer
}

func NewGeoMultiIndexer( resolutions MetersVector, expiredTimed time.Duration) *GeoMultiIndexer{
	sort.Sort(&resolutions)
	indexes := make([]*GeoIndexer, 0, len(resolutions))
	for _, r := range resolutions {
		index := NewGeoIndexer( r, expiredTimed)
		indexes = append(indexes, index)
	}

	return &GeoMultiIndexer{
		indexes: indexes,
	}
}



