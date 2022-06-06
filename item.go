package main

type aItem struct {
	baseItem
	conditions
}

type baseItem struct {
	row    int
	code   string
	name   string
	weight float64
}

func (r Row) NewItem() aItem {
	result := aItem{}

	result.row = r.rowNum
	result.code = r.code
	result.name = r.name
	result.weight = r.weight * 1000

	result.conditions = r.itemConditions

	return result
}
