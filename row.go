package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

type Row struct {
	rowNum         int
	bl             string
	etd            time.Time
	eta            time.Time
	weight         float64
	code           string
	name           string
	contConditions conditions
	itemConditions conditions
}

type Rows []Row

type Containers []aContainer

func NewRow(rowNum int, bl string, etd time.Time, eta time.Time, weight float64, code string, name string, contCond conditions, itemCond conditions) Row {
	result := Row{}
	result.rowNum = rowNum
	result.bl = bl
	result.etd = etd
	result.eta = eta
	result.weight = weight
	result.code = code
	result.name = name
	result.contConditions = contCond
	result.itemConditions = itemCond

	return result
}

func mapsToRows(rows []map[string]string, ctx globalContext) []Row {
	var result []Row

	for _, r := range rows {
		newRow, err := mapToRow(r, ctx)
		if err != nil {
			continue
		} else {
			result = append(result, newRow)
		}
	}

	return result
}

func mapToRow(r map[string]string, ctx globalContext) (Row, error) {
	var (
		contextMap = ctx.input.context
		err        error
	)

	var (
		rowNum int
		bl     string
		code   string
		name   string
		weight float64
		etd    time.Time
		eta    time.Time
	)

	if v, ok := r["Row"]; ok {
		rowNum, err = strconv.Atoi(v)
		if err != nil {
			log.Println(err)
		}
	}

	if v, ok := r[contextMap["bl"]]; ok {
		bl = strings.TrimSpace(v)
		if len(bl) > 7 && bl[len(bl)-7:] == "_x000D_" {
			bl = bl[0 : len(bl)-7]
		}
	} else {
		log.Println("no BL! this row will be disregarded")
		return Row{}, errors.New("No BL!")
	}

	if v, ok := r[contextMap["code"]]; ok {
		code = v
	}

	if v, ok := r[contextMap["etd"]]; ok {
		etd, err = time.Parse("2006-01-02", v)
		if err != nil {
			log.Println(err)
		}
	}

	if v, ok := r[contextMap["eta"]]; ok {
		eta, err = time.Parse("2006-01-02", v)
		if err != nil {
			log.Println(err)
			return Row{}, err
		}
	} else {
		log.Println("No ETA!, this row will be disregarded")
		return Row{}, errors.New("No ETA!")
	}

	if v, ok := r[contextMap["name"]]; ok {
		name = v
	}

	if v, ok := r[contextMap["qty"]]; ok {
		weight, err = strconv.ParseFloat(v, 32)
		if err != nil {
			log.Println(err, r[contextMap["qty"]])
		}
	} else {
		log.Println("No weight! this row will be disregarded")
		return Row{}, errors.New("No weight!")
	}

	itemconditions := ParseCondtions(r, itemConditions, excelSource)
	contconditions := ParseCondtions(r, containerConditions, excelSource)

	result := NewRow(rowNum, bl, etd, eta, weight, code, name, contconditions, itemconditions)

	return result, nil
}

func (cntrs Containers) checkAllConditions() {
	for i := range cntrs {
		cntrs[i].checkConditions()
	}
}

func (cntr aContainer) toRow() []Row {
	result := make([]Row, 0)
	var temp Row
	for _, item := range cntr.itemList {
		temp = Row{}
		temp.bl = cntr.bl
		temp.weight = item.weight
		temp.contConditions = cntr.conditions
		temp.itemConditions = item.conditions
		temp.rowNum = item.row
		temp.code = item.code
		temp.name = item.name
		temp.etd = cntr.etd
		temp.eta = cntr.eta
		result = append(result, temp)
	}
	return result
}

func (cntrs Containers) toRows() Rows {
	result := make(Rows, 0)
	var temps []Row
	for _, cntr := range cntrs {
		temps = cntr.toRow()
		for _, temp := range temps {
			result = append(result, temp)
		}
	}
	return result
}
