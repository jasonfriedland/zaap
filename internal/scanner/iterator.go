package scanner

type Iterator interface {
	Next() bool
	Item() Item
	Err() error
}

type sliceIterator struct {
	items []Item
	index int
}

func NewIterator(items []Item) Iterator {
	return &sliceIterator{items: items, index: -1}
}

func (it *sliceIterator) Next() bool {
	it.index++
	return it.index < len(it.items)
}

func (it *sliceIterator) Item() Item {
	if it.index < 0 || it.index >= len(it.items) {
		return Item{}
	}
	return it.items[it.index]
}

func (it *sliceIterator) Err() error {
	return nil
}

func Collect(it Iterator) ([]Item, error) {
	var items []Item
	for it.Next() {
		items = append(items, it.Item())
	}
	return items, it.Err()
}
