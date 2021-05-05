package futures

import (
	"math"
	"time"
)

func ContractRollLastTradingDay(dtime time.Time, current, previous *Contract) bool {
	if dtime.Equal(previous.series.EndTime()) || dtime.Before(previous.series.EndTime()) {
		return true
	} else {
		return false
	}
}

func ContractRollLastNTradingDay(dtime time.Time, current, previous *Contract, days int) bool {
	if previous.series.EndTime().Sub(dtime).Hours()/24.0 >= float64(days) {
		return true
	} else {
		return false
	}
}

func ContractRollFirstOfMonth(
	dtime time.Time,
	current,
	previous *Contract,
) bool {

	pi := previous.series.IndexInFullTimes(dtime)
	if pi <= 0 {
		return false
	} else {
		pdt := previous.series.FullTimes()[pi-1]
		if dtime.Month() == previous.ContractMonth() && pdt.Month() == time.Month(previous.ContractMonth()-1) {
			return true
		} else {
			return false
		}
	}

}

func ContractRollOpenInterest(dtime time.Time, current, previous *Contract) bool {
	coi := current.series.ValueInFullTimes(dtime, "openinterest", math.NaN())
	poi := previous.series.ValueInFullTimes(dtime, "openinterest", math.NaN())

	cv := current.series.ValueInFullTimes(dtime, "volume", math.NaN())
	pv := previous.series.ValueInFullTimes(dtime, "volume", math.NaN())

	vok := true
	oiok := true
	if math.IsNaN(cv) || math.IsNaN(pv) || cv == 0 || pv == 0 {
		vok = false
	}

	if math.IsNaN(coi) || math.IsNaN(poi) || coi == 0 || poi == 0 {
		oiok = false
	}

	switch {
	case vok && !oiok:
		return cv < pv
	case !vok && oiok:
		return coi < poi
	case vok && oiok:
		return cv < pv || coi < poi
	default:
		return ContractRollLastNTradingDay(dtime, current, previous, 8)
	}
}
