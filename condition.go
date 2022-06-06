package main

import (
	"errors"
	"log"
	"math"
	"regexp"
	"time"
)

const (
	arriveCondition conditionCat = iota
	warehouseCondition
	quarantineCondition
	foodCondition
)

const (
	foodConditionTolerance float64 = 0.5
	foodDclrReg            string  = `11-([A-Z\d]{16})-\d{2}`
	qurantineDclrReg       string  = `12-([A-Z\d]{14})-01`
	warehouseContentReg    string  = "보세운송 반입"
	arrivedContentReg      string  = "입항 반입"
)

var itemConditions = []conditionCat{foodCondition}

var containerConditions = []conditionCat{arriveCondition, warehouseCondition, quarantineCondition}

var mustConditions = []conditionCat{quarantineCondition, foodCondition}

type conditionCat uint

type condition struct {
	condName         conditionCat
	isCondApproved   bool
	dateCondApproved time.Time
	errorMsg         string
}

type conditions map[conditionCat]condition

func (cc conditionCat) New() *condition {
	return &condition{
		condName: cc,
	}
}

func indexList(element conditionCat, data []conditionCat) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func (cc conditionCat) Context(ic ioCandidate) (map[string]string, error) {
	var result map[string]string
	switch ic {
	case excelSource:
		switch cc {
		case foodCondition:
			result = map[string]string{
				"name":       "식품검역",
				"header":     excelOrder,
				"Identifier": excelFood,
			}
		case quarantineCondition:
			result = map[string]string{
				"name":       "동물검역",
				"header":     excelOrder,
				"Identifier": excelQuarantine,
			}
		case warehouseCondition:
			result = map[string]string{
				"name":       "보세창고 입고",
				"header":     excelOrder,
				"Identifier": excelWarehouse,
			}
		case arriveCondition:
			result = map[string]string{
				"name":       "입항",
				"header":     excelOrder,
				"Identifier": excelArrive,
			}
		default:
			return result, errors.New("condition not implemented")
		}
	default:
		return result, errors.New("IO Source not implemented")
	}
	return result, nil
}

// func (cc conditionCat) checkCondition(cntr *aContainer) {

// }

func ParseCondtions(r map[string]string, cc []conditionCat, ic ioCandidate) conditions {
	result := conditions{}
	var newCondition condition

	for _, cond := range cc {
		v, err := cond.Context(ic)
		if err != nil {
			log.Fatal(err.Error())
		}
		re := regexp.MustCompile(v["Identifier"])

		if re.Match([]byte(r[v["header"]])) {
			newCondition = *cond.New()
			result[cond] = newCondition
		}
	}

	return result
}

func checkFood(cntr *aContainer) {
	items := dclrInfos(cntr.retrievedXML.CargCsclPrgsInfoDtlQryVo)
	var conditionToChange condition
	var isChecked bool

	bundles := items.BundleByDclrNum(foodDclrReg)
	var bundleWeight float64
	for _, v := range bundles {
		bundleWeight += v.Weight()
	}

	if len(bundles) != len(cntr.itemList) && math.Abs(cntr.totalWeight*1000-bundleWeight) <= foodConditionTolerance {
		log.Printf("%s-유니패스 검역신청수(%d)와 엑셀의 아이템수(%d)가 일치하지 않습니다\n", cntr.bl, len(bundles), len(cntr.itemList))
	}

	for _, v := range bundles {
		isChecked = false
		for i := range cntr.itemList {
			if _, ok := cntr.itemList[i].conditions[foodCondition]; math.Abs(v.Weight()-cntr.itemList[i].weight) <= foodConditionTolerance && ok {
				if isChecked == false {
					conditionToChange = cntr.itemList[i].conditions[foodCondition]
					conditionToChange.dateCondApproved = v.LatestDate()
					conditionToChange.isCondApproved = true
					cntr.itemList[i].conditions[foodCondition] = conditionToChange
					isChecked = true
				} else if isChecked == true {
					log.Println("duplicate found! weight should be unique")
				}
			}
		}
	}

}

func (cntr *aContainer) checkContCondition() {

	for _, v := range cntr.conditions {
		switch v.condName {
		case quarantineCondition:
			checkQuarantine(cntr)
		case warehouseCondition:
			checkWarehouse(cntr)
		case arriveCondition:
			checkArrive(cntr)
		}
	}
}

func (cs conditions) checkItemCondition(cntr *aContainer) {

	for _, v := range cs {
		switch v.condName {
		case foodCondition:
			checkFood(cntr)
		}
	}

}

func checkQuarantine(cntr *aContainer) {
	var conditionToChange condition
	items := dclrInfos(cntr.retrievedXML.CargCsclPrgsInfoDtlQryVo)
	quarantineDclr := items.ExtractByDclr(qurantineDclrReg)
	if len(quarantineDclr) == 1 {
		if _, ok := cntr.conditions[quarantineCondition]; ok {
			conditionToChange = cntr.conditions[quarantineCondition]
			conditionToChange.isCondApproved = true
			conditionToChange.dateCondApproved = quarantineDclr.LatestDate()
			cntr.conditions[quarantineCondition] = conditionToChange
		}
	} else if len(quarantineDclr) > 1 {
		log.Println(cntr.bl, "두 개 이상의 동물검역 결과가 있습니다")
		if _, ok := cntr.conditions[quarantineCondition]; math.Abs(cntr.totalWeight-quarantineDclr.Weight()) <= foodConditionTolerance && ok {
			conditionToChange = cntr.conditions[quarantineCondition]
			conditionToChange.isCondApproved = true
			conditionToChange.dateCondApproved = quarantineDclr.LatestDate()
			conditionToChange.errorMsg = "두 개 이상의 동물 검역 결과가 있지만 합계가 콘테이너 총 중량과 일치"
			cntr.conditions[quarantineCondition] = conditionToChange
		} else if _, ok := cntr.conditions[quarantineCondition]; math.Abs(cntr.totalWeight-quarantineDclr.Weight()) > foodConditionTolerance && ok {
			log.Println(cntr.bl, "엑셀상의 아이템 합계중량과 유니패스상의 아이템 합계중량이 일치하지 않습니다.", "엑셀에는 한 콘테이너에 담긴 모든 아이템이 기재되어있어야 합니다")
			if _, ok := cntr.conditions[quarantineCondition]; ok {
				conditionToChange = cntr.conditions[quarantineCondition]
				conditionToChange.errorMsg = "엑셀상의 아이템 합계중량과 유니패스상의 아이템 합계중량이 일치하지 않습니다./n 엑셀에는 한 콘테이너에 담긴 모든 아이템이 기재되어있어야 합니다"
				cntr.conditions[quarantineCondition] = conditionToChange
			}
		}
	}
}

func checkWarehouse(cntr *aContainer) {
	err := checkByContentBase(cntr, warehouseContentReg, warehouseCondition)
	if err != nil {
		log.Println("warehouseCondition", err)
	}
}

func checkArrive(cntr *aContainer) {
	err := checkByContentBase(cntr, arrivedContentReg, arriveCondition)
	if err != nil {
		log.Println("arrivedCondition", err)
	}
}

func checkByContentBase(cntr *aContainer, regString string, conditionID conditionCat) error {
	var conditionToChange condition
	items := dclrInfos(cntr.retrievedXML.CargCsclPrgsInfoDtlQryVo)
	condDclr := items.ExtractByContent(regString)
	if len(condDclr) == 1 {
		if _, ok := cntr.conditions[conditionID]; ok {
			conditionToChange = cntr.conditions[conditionID]
			conditionToChange.isCondApproved = true
			conditionToChange.dateCondApproved = condDclr.LatestDate()
			cntr.conditions[conditionID] = conditionToChange
			return nil
		}
	} else if _, ok := cntr.conditions[conditionID]; ok && len(condDclr) > 1 {
		conditionToChange = cntr.conditions[conditionID]
		conditionToChange.errorMsg = "조건의 결과가 두개 이상입니다"
		cntr.conditions[conditionID] = conditionToChange
		return errors.New("scondition has 2 more result")
	}
	return errors.New("condition has 0 result")
}

func (row Row) mustConditionReducer(ic ioCandidate) (bool, time.Time, []string) {
	var mustTime time.Time
	var notime time.Time
	var failedConditions []string

	for _, j := range mustConditions {
		for _, k := range row.contConditions {
			if j == k.condName {
				if condContext, err := k.condName.Context(ic); k.isCondApproved == false && err == nil {
					failedConditions = append(failedConditions, condContext["name"])
				} else if mustTime.Before(k.dateCondApproved) && true == k.isCondApproved {
					mustTime = k.dateCondApproved
				}
			}
		}

		for _, k := range row.itemConditions {
			if j == k.condName {
				if condContext, err := k.condName.Context(ic); k.isCondApproved == false && err == nil {
					failedConditions = append(failedConditions, condContext["name"])
				} else if mustTime.Before(k.dateCondApproved) && true == k.isCondApproved {
					mustTime = k.dateCondApproved
				}
			}
		}
	}

	if 0 == len(failedConditions) {
		return true, mustTime, failedConditions
	}
	return false, notime, failedConditions
}
