package main

import (
	"github.com/liuzz1983/scalelbs/geoindex"
)

type GeoIndexer struct {
	resolution geoindex.Meters
	indexes    map[string]map[string]*Pos
}

func NewGeoIndexer(resolution geoindex.Meters) *GeoIndexer {
	return &GeoIndexer{
		indexes:    make(map[string]map[string]*Pos),
		resolution: resolution,
	}
}

func (indexer *GeoIndexer) Cells(lat float64, lng float64) []geoindex.Cell {
	cell := geoindex.CellOf2(lat, lng, indexer.resolution)
	cells := make([]geoindex.Cell, 0)
	indexes := []int{0, 1, 2}
	for _, x := range indexes {
		for _, y := range indexes {
			cells = append(cells, geoindex.Cell{
				X: cell.X - 1 + x,
				Y: cell.Y - 1 + y,
			})
		}
	}
	return cells
}

func (indexer *GeoIndexer) Cell(lat float64, lng float64) geoindex.Cell {
	return geoindex.CellOf2(lat, lng, indexer.resolution)
}

func (indexer *GeoIndexer) PosCell(pos *Pos) geoindex.Cell {
	return geoindex.CellOf2(pos.Lat, pos.Lng, indexer.resolution)
}

func (indexer *GeoIndexer) AddPos(pos *Pos) {
	cell := indexer.PosCell(pos)
	values, ok := indexer.indexes[cell.Id()]
	if !ok {
		values = make(map[string]*Pos, 0)
	}
	values[cell.Id()] = pos
	indexer.indexes[cell.Id()] = values
}

func (indexer *GeoIndexer) Get(cellId string) map[string]*Pos {
	v, _ := indexer.indexes[cellId]
	return v
}

type GeoLevelIndexer struct {
}
