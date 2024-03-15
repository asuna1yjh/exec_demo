package pkg

import "fmt"

// 构造一个切片函数，根据传入的索引，从切片中取出对应值，并把切片的数据删除
func GetSlice(data *[][]string, index int) [][]string {
	if index >= len(*data) {
		index = len(*data)
	}
	res := (*data)[:index]
	*data = (*data)[index:]
	return res
}

func SliceNew() {
	a := []string{"a", "b", "c"}
	b := []string{"d", "e", "f"}
	a = append(a, b[:]...)
	fmt.Println(a)
}
