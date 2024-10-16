package controller

import "errors"

type Idx struct {
	value int
}

func FromInt(i int) (*Idx, error) {
	if i < 0 {
		return nil, errors.New("Index can't be less than zero")
	}

	return &Idx{value: i}, nil
}

func (idx *Idx) Value() int {
	return idx.value
}

func (idx *Idx) IsIn(max int) bool {
	return idx.value < max
}
