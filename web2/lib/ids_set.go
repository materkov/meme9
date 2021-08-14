package lib

type IdsSet map[int]bool

func (i IdsSet) Add(id int) {
	i[id] = true
}

func (i IdsSet) Get() []int {
	result := make([]int, len(i))

	idx := 0
	for id := range i {
		result[idx] = id
		idx++
	}

	return result
}
