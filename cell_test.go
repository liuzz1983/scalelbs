package main

import (
	"fmt"
	"testing"
)

func TestCell(t *testing.T) {

	v := CellOf2(23.4, 45.5, Km(0.1))
	fmt.Println("-----:" + v.Id())
	id := v.Id()
	cell, _ := ParseCell(id)
	fmt.Printf("%v %v %v %v\n", v.X, v.Y, cell.X, cell.Y)

}
