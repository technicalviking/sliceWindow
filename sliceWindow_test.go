package sliceWindow

import (
	"testing"
)

func TestAdd(t *testing.T) {
	vals := []float64{
		//extra values
		0.0,
		25.6,
		100.4,
		-64.3,
		2.0,
		//values the test cares about.
		11.78125,
		13.208984375, //max
		-1.400390625,
		2.849609375,
		-3.814453125,
		0.888671875,
		7.208984375,
		-15.650390625,
		-9.869140625,
		-35.416015625, //min
		-23.126953125,
		-19.189453125,
		-26.314453125,
		-22.158203125,
		-22.158203125,
		-22.158203125,
		-22.158203125,
		-17.001953125,
		-29.884765625,
		-23.009765625,
	}

	window := New(20)

	for _, curVal := range vals {
		window.PushBack(curVal)
	}

	t.Errorf("min: %+v \n max: %+v", window.min, window.max)

}

func TestGetNormalizedSlice(t *testing.T) {

	window := New(5)

	for startingValue := 1.0; startingValue <= 10.0; startingValue += 1.0 {
		window.PushBack(startingValue)
	}

	if window.min != 6.0 {
		t.Errorf("invalid min value %+v\n", window.min)
	}

	if window.max != 10.0 {
		t.Errorf("invalid max value %+v\n", window.max)
	}

	params := [][]int{
		{0, 3},
		{0, window.Len()},
		{1, 4},
		{-1, 1},
		{-1, 3},
		{-3, 2},
	}

	for index, args := range params {
		slice, err := window.GetNormalizedSlice(args[0], args[1])

		if err != nil {
			t.Errorf("err getting slice %s", err.Error())
		}
		t.Errorf("what am I %d? %v \n", index, slice)
	}
}
