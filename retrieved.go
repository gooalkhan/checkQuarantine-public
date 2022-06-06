package main

import (
	"encoding/xml"
	"regexp"
	"time"
)

type retrievedXML struct {
	XMLName                  xml.Name          `xml:"cargCsclPrgsInfoQryRtnVo"`
	NtceInfo                 string            `xml:"ntceInfo"`
	TCnt                     int               `xml:"tCnt"`
	CargCsclPrgsInfoQryVo    cargoInfo         `xml:"cargCsclPrgsInfoQryVo"`
	CargCsclPrgsInfoDtlQryVo []cargoDetailInfo `xml:"cargCsclPrgsInfoDtlQryVo"`
}
type cargoInfo struct {
	CsclPrgsStts        string `xml:"csclPrgsStts"`
	Vydf                string `xml:"vydf"`
	RlseDtyPridPassTpcd string `xml:"rlseDtyPridPassTpcd"`
	Prnm                string `xml:"prnm"`
	LdprCd              string `xml:"ldprCd"`
	ShipNat             string `xml:"shipNat"`
	BlPt                string `xml:"blPt"`
	DsprNm              string `xml:"dsprNm"`
	EtprDt              string `xml:"etprDt"`
	PrgsStCd            string `xml:"prgsStCd"`
	Msrm                string `xml:"msrm"`
	WghtUt              string `xml:"wghtUt"`
	DsprCd              string `xml:"dsprCd"`
	CntrGcnt            string `xml:"cntrGcnt"`
	CargTp              string `xml:"cargTp"`
	ShcoFlcoSgn         string `xml:"shcoFlcoSgn"`
	PckGcnt             string `xml:"pckGcnt"`
	EtprCstm            string `xml:"etprCstm"`
	ShipNm              string `xml:"shipNm"`
	HblNo               string `xml:"hblNo"`
	PrcsDttm            string `xml:"prcsDttm"`
	FrwrSgn             string `xml:"frwrSgn"`
	SpcnCargCd          string `xml:"spcnCargCd"`
	Ttwg                string `xml:"ttwg"`
	LdprNm              string `xml:"ldprNm"`
	FrwrEntsConm        string `xml:"frwrEntsConm"`
	DclrDelyAdtxYn      string `xml:"dclrDelyAdtxYn"`
	MtTrgtCargYnNm      string `xml:"mtTrgtCargYnNm"`
	CargMtNo            string `xml:"cargMtNo"`
	CntrNo              string `xml:"cntrNo"`
	MblNo               string `xml:"mblNo"`
	BlPtNm              string `xml:"blPtNm"`
	LodCntyCd           string `xml:"lodCntyCd"`
	PrgsStts            string `xml:"prgsStts"`
	ShcoFlco            string `xml:"shcoFlco"`
	PckUt               string `xml:"pckUt"`
	ShipNatNm           string `xml:"shipNatNm"`
	Agnc                string `xml:"agnc"`
}

type cargoDetailInfo struct {
	ShedNm               string      `xml:"shedNm"`
	PrcsDttm             processDate `xml:"prcsDttm"`
	DclrNo               string      `xml:"dclrNo"`
	RlbrDttm             string      `xml:"rlbrDttm"`
	Wght                 float64     `xml:"wght"`
	RlbrBssNo            string      `xml:"rlbrBssNo"`
	BfhnGdncCn           string      `xml:"bfhnGdncCn"`
	WghtUt               string      `xml:"wghtUt"`
	PckGcnt              string      `xml:"pckGcnt"`
	CargTrcnRelaBsopTpcd string      `xml:"cargTrcnRelaBsopTpcd"`
	PckUt                string      `xml:"pckUt"`
	RlbrCn               string      `xml:"rlbrCn"`
	ShedSgn              string      `xml:"shedSgn"`
}

type processDate time.Time

type dclrInfos []cargoDetailInfo

func (c *processDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "20060102150405" // yyyymmdd date format
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = processDate(parse)
	return nil
}

func (dcinfo dclrInfos) Weight() float64 {
	var result float64
	for _, v := range dcinfo {
		result += v.Wght
	}
	return result
}

func (dcinfo dclrInfos) LatestDate() time.Time {
	var result time.Time
	var temp time.Time

	for _, v := range dcinfo {

		temp = time.Time(v.PrcsDttm)
		if result.Before(temp) {
			result = temp
		}
	}

	return result
}

func (dcinfo dclrInfos) ExtractByContent(reg string) dclrInfos {
	result := dclrInfos{}

	re := regexp.MustCompile(reg)

	for _, v := range dcinfo {
		if re.Match([]byte(v.RlbrCn)) {
			result = append(result, v)
		}
	}

	return result
}

func (dcinfo dclrInfos) ExtractByDclr(reg string) dclrInfos {
	result := dclrInfos{}

	re := regexp.MustCompile(reg)

	for _, v := range dcinfo {
		if re.Match([]byte(v.DclrNo)) {
			result = append(result, v)
		}
	}

	return result
}

func (dcinfo dclrInfos) BundleByContent(reg string) []dclrInfos {
	result := make([]dclrInfos, 0)

	re := regexp.MustCompile(reg)

	for _, v := range dcinfo {
		findString := re.FindString(v.CargTrcnRelaBsopTpcd)

		for i := range result {

			if findString != "" {
				result[i] = append(result[i], v)
			}
			break
		}
		if len(result) == 0 && findString != "" {
			result = append(result, dclrInfos{v})
		}
	}
	return result
}

func (dcinfo dclrInfos) BundleByDclrNum(reg string) []dclrInfos {
	result := make([]dclrInfos, 0)

	re := regexp.MustCompile(reg)

	for _, v := range dcinfo {
		matches := re.FindStringSubmatch(v.DclrNo)

		if matches != nil {
			if len(result) == 0 {
				result = append(result, dclrInfos{v})
				continue
			}
			for i := range result {

				if len(matches) > 1 {
					if matches[1] == re.FindStringSubmatch(result[i][0].DclrNo)[1] {
						result[i] = append(result[i], v)
						break
					} else if i == len(result)-1 {
						result = append(result, dclrInfos{v})
					}

				} else if len(matches) == 1 {
					result[i] = append(result[i], v)
					break
				}
			}
		}
	}
	return result
}
