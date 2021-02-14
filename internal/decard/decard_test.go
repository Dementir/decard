package decard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecard(t *testing.T) {
	input := []Point{
		{
			100, 100,
		},
		{
			400, 400,
		},
		{
			100, 400,
		},
		{
			400, 100,
		},
	}

	testResult := []Point{
		{
			100, 100,
		},
		{
			100, 400,
		},
		{
			400, 100,
		},
		{
			400, 400,
		},
	}

	result := Decard(input)

	assert.EqualValues(t, testResult, result)
}
