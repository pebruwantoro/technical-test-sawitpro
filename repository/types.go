// This file contains types that are used in the repository layer.
package repository

type Estate struct {
	Id     string
	Width  int
	Length int
}

type EstateTree struct {
	Id       string
	EstateId string
	X        int
	Y        int
	Height   int
}

type StatsEstate struct {
	Count  int
	Max    int
	Min    int
	Median float64
}
