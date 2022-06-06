package main

import "time"

type aContainer struct {
	baseContainer
	conditions
	itemList itemList
	retrievedXML
}

type itemList []aItem

type baseContainer struct {
	bl          string
	year        int
	etd         time.Time
	eta         time.Time
	totalWeight float64
}

func (r Row) NewContainer() aContainer {
	var result = aContainer{}

	result.eta = r.eta
	result.etd = r.etd
	result.bl = r.bl
	result.year = r.eta.Year()
	result.conditions = r.contConditions
	result.itemList = append(result.itemList, r.NewItem())
	result.totalWeight = r.weight

	return result
}

func (rows Rows) toContainers(ctx globalContext) Containers {
	result := make(Containers, 0)

	for _, v := range rows {
		if ok, cntr := result.haveThisCont(v); ok {
			cntr.itemList = append(cntr.itemList, v.NewItem())
			cntr.totalWeight += v.weight
		} else {
			result = append(result, v.NewContainer())
		}

	}
	return result
}

func (cntrs Containers) haveThisCont(r Row) (bool, *aContainer) {
	for i := range cntrs {
		if cntrs[i].bl == r.bl {
			return true, &cntrs[i]
		}
	}
	return false, nil
}

func (cntr *aContainer) checkConditions() {
	cntr.checkContCondition()
	for i := range cntr.itemList {
		cntr.itemList[i].checkItemCondition(cntr)
	}
}

func (cntr aContainer) getItemTotalWeight() float64 {
	var result float64

	for _, v := range cntr.itemList {
		result += v.weight
	}

	return result
}
