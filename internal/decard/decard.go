package decard

import "sort"

type Points []string

type Point struct {
	X int
	Y int
}

func Decard(input []Point) []Point {
	sort.Slice(input, func(i, j int) bool {
		if input[i].X < input[j].X {
			return true
		}

		if input[i].X == input[j].X {
			if input[i].Y < input[j].Y {
				return true
			}
		}
		return false
	})

	return input
}
