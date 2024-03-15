package pkg

import (
	"fmt"
	"testing"
)

func Test_GetSlice(t *testing.T) {
	tt := [][]string{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}}
	s1 := GetSlice(&tt, 2)
	s2 := GetSlice(&tt, 1)
	//fmt.Println(tt)
	s3 := GetSlice(&tt, 20)
	s4 := GetSlice(&tt, 2)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)
}

func TestSliceNew(t *testing.T) {
	SliceNew()
}
