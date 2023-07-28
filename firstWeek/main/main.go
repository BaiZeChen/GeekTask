package main

import "fmt"

func NewSlice[T any](capacity int) *Slice[T] {
	slice := &Slice[T]{}
	return slice.init(capacity)
}

type Slice[T any] struct {
	arr  []T
	size int
}

func (s *Slice[T]) init(capacity int) *Slice[T] {
	s.arr = make([]T, 0, capacity)
	return s
}

func (s *Slice[T]) Append(elems ...T) {
	capacity := cap(s.arr)
	addLength := len(elems)
	s.size = addLength + s.size

	// 得扩容了
	if s.size > capacity {
		s.expansion(capacity)
	}
	s.arr = append(s.arr, elems...)
}

func (s *Slice[T]) Insert(key int, value T) {
	if key < 0 || key >= s.size {
		panic("数组越界")
	}

	capacity := cap(s.arr)
	s.size++
	if s.size > capacity {
		s.expansion(capacity)
	}
	s.arr = append(s.arr, value)

	for i := s.size - 2; i >= key; i-- {
		s.arr[i+1] = s.arr[i]
	}
	s.arr[key] = value
}

func (s *Slice[T]) Update(key int, value T) {
	if key < 0 || key >= s.size {
		panic("数组越界")
	}
	s.arr[key] = value
}

func (s *Slice[T]) Get(key int) T {
	if key < 0 || key >= s.size {
		panic("数组越界")
	}
	return s.arr[key]
}

func (s *Slice[T]) RemoveByIndex(index int) {
	if index < 0 || index >= s.size {
		panic("数组越界")
	}

	for i := index + 1; i < s.size; i++ {
		s.arr[i-1] = s.arr[i]
	}
	s.size--
	s.arr = s.arr[:s.size]
	s.shrinkage()
}

func (s *Slice[T]) RemoveLast() {
	s.RemoveByIndex(s.size - 1)
}

func (s *Slice[T]) RemoveFirst() {
	s.RemoveByIndex(0)
}

func (s *Slice[T]) Len() int {
	return s.size
}

// 扩容方式本来向参考1.18版本之后append扩容源码的代码
// 但我看还涉及到内存对齐，所以这里统一就按照2倍处理了，这样也方便了缩容计算
func (s *Slice[T]) expansion(oldCap int) {
	var newCap = 0
	for newCap < s.size {
		if newCap == 0 {
			newCap = oldCap * 2
		} else {
			newCap = newCap * 2
		}
	}

	tmp := make([]T, 0, newCap)
	s.arr = append(tmp, s.arr...)
}

// 为了防止复杂度震荡，所以采取len=1/4cap时，新cap=旧cap1/2
// 而不是len=1/2cap时，新cap=旧cap1/2
func (s *Slice[T]) shrinkage() {
	length := s.size
	capacity := cap(s.arr)

	// 至少保留2个cap
	if length == (capacity/4) && (capacity > 2) {
		newCap := capacity / 2
		newSlice := make([]T, 0, newCap)
		for i := 0; i < length; i++ {
			newSlice = append(newSlice, s.arr[i])
		}
		s.arr = newSlice
	}
}

func (s *Slice[T]) String() string {
	res := fmt.Sprintf("Array: size = %d , realSize = %d ,capacity = %d\n", s.size, len(s.arr), cap(s.arr))
	res += "["
	for i := 0; i < s.size; i++ {
		res += fmt.Sprintf("%v", s.arr[i])
		if i != s.size-1 {
			res += ", "
		}
	}
	res += "]"
	return res
}

func main() {

}
