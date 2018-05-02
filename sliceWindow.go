package sliceWindow

import (
	"container/list"
	"fmt"
)

//SliceWindow wrapper for container/list type specifically designed to normalize the values within when performing certain operations.
type SliceWindow struct {
	list        *list.List
	min         float64
	max         float64
	minPosition int //index of element containing min val
	maxPosition int //index of element containing max val
	sum         float64
	maxLength   int
}

//Len get the length of the underlying list.
func (nw SliceWindow) Len() int {
	return nw.list.Len()
}

//Max get max value of the slice window.
func (nw SliceWindow) Max() float64 {
	return nw.max
}

//MaxPosition get the index of the maximum value in the list.
func (nw SliceWindow) MaxPosition() int {
	return nw.maxPosition
}

//Min get min value of the slice window.
func (nw SliceWindow) Min() float64 {
	return nw.min
}

//MinPosition get the index of the minimum value in the list.
func (nw SliceWindow) MinPosition() int {
	return nw.minPosition
}

//Mean get the mean of the values in the window.
func (nw SliceWindow) Mean() float64 {
	len := float64(nw.list.Len())
	return nw.sum / len
}

//Init reset underlying list.
func (nw *SliceWindow) Init() {
	nw.list = nw.list.Init()
}

//PushBack Add a value to the end of the list.
func (nw *SliceWindow) PushBack(val float64) {

	newEntry := val
	nw.list.PushBack(newEntry)
	nw.sum += newEntry

	for nw.list.Len() > nw.maxLength {
		removed := nw.list.Remove(nw.list.Front())
		nw.sum -= removed.(float64)

		if removed == nw.min {
			nw.min = 0
		}

		if removed == nw.max {
			nw.max = 0
		}
	}

	nw.setMinMax()
}

func (nw *SliceWindow) setMinMax() {
	i := 0

	for e := nw.list.Front(); e != nil; e = e.Next() {
		if (nw.min == 0) || nw.min >= e.Value.(float64) {
			nw.minPosition = i
			nw.min = e.Value.(float64)
		}

		if (nw.max == 0) || nw.max <= e.Value.(float64) {
			nw.maxPosition = i
			nw.max = e.Value.(float64)
		}

		i++
	}
}

//GetSliceFromList convert the list type into an equivalent []float64
func (nw SliceWindow) GetSliceFromList() []float64 {
	result := make([]float64, nw.list.Len())

	index := 0
	for e := nw.list.Front(); e != nil; e = e.Next() {
		result[index] = e.Value.(float64)
		index++
	}

	return result
}

//GetNormalizedSlice retrieve a chunk of the list specified by startIndex with a len = numElements.
//specifying a negative startindex implies to start from the end of the list, where -1 is the last element.
//Example: calling this for foo := [0,1,2,3,4] with startIndex = -2 and numElements = 2 should return
// [-3,4]
func (nw SliceWindow) GetNormalizedSlice(startIndex, numElements int) ([]float64, error) {

	if numElements > nw.list.Len() {
		return nil, fmt.Errorf("Cannot create slice window of length %d greater than list length %d", numElements, nw.list.Len())
	}

	result := make([]float64, 0)

	if nw.max == nw.min {
		for i := 0; i < numElements; i++ {
			result = append(result, float64(1.0))
		}

		return result, nil
	}

	if startIndex < 0 {
		startIndex *= -1

		startIndex = nw.list.Len() - startIndex
	}

	curIndex := 0
	for e := nw.list.Front(); curIndex < startIndex+numElements && e != nil; e = e.Next() {
		if curIndex < startIndex {
			curIndex++
			continue
		}

		entry := e.Value.(float64)

		numerator := entry - nw.min
		denominator := nw.max - nw.min

		normalizedEntry := numerator / denominator

		result = append(result, normalizedEntry)
		curIndex++
	}

	return result, nil
}

type sliceMapFunc = func(float64) float64

//Map perform an operation on every element of the list, returning a new list.
func (nw SliceWindow) Map(cb sliceMapFunc) *SliceWindow {
	output := New(nw.Len())

	for curWindowEl := nw.list.Front(); curWindowEl != nil; curWindowEl = curWindowEl.Next() {
		input := curWindowEl.Value.(float64)
		outputVal := cb(input)
		output.PushBack(outputVal)
	}

	return output
}

//New constructor for a normalizer window
func New(length int) *SliceWindow {
	return &SliceWindow{
		list:      list.New(),
		maxLength: length,
	}
}
