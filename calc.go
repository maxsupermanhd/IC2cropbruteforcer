package main

import (
	"log"
)

func calcWithRNG(rng int, stats crop, params cropEnvironment) int {
	base := 3 + rng + stats.Growth
	need := (stats.Tier-1)*4 + stats.Growth + stats.Gain + stats.Resist
	if need < 0 {
		need = 0
	}

	have := int(float64(params.Humidity)*stats.Humidity+float64(params.Nutrients)*stats.Nutrients+float64(params.Air)*stats.Air) * 5
	if have >= need {
		base = base * (100 + (have - need)) / 100
	} else {
		neg := (need - have) * 4
		if neg > 100 {
			base = -1
		} else {
			base = base * (100 - neg) / 100
			if base < 0 {
				base = 0
			}
		}
	}
	return base
}

func calcByParams(stats crop, params cropEnvironment) float64 {
	avgs := 0
	avgc := 0
	for rng := 0; rng < 7; rng++ {
		curr := calcWithRNG(rng, stats, params)
		if curr == -1 {
			return 0.0
		}
		avgs += curr
		avgc++
	}
	return float64(avgs) / float64(avgc)
}

func getCellFromPosOffset(field *cropField, ctop, cleft, otop, oleft int) int {
	if ctop < 0 || ctop > field.Height || cleft < 0 || cleft > field.Width {
		log.Fatal("bad lookup params")
	}
	top := ctop + otop
	left := cleft + oleft
	for top < 0 {
		top = field.Height + top
	}
	for left < 0 {
		left = field.Height + left
	}
	for top >= field.Height {
		top = top - field.Height
	}
	for left >= field.Height {
		left = left - field.Height
	}
	return top*field.Height + left
}

func calcFieldEnvironment(field *cropField) {
	for cr := 0; cr < field.Height; cr++ {
		for cc := 0; cc < field.Width; cc++ {
			i := cr*field.Height + cc
			if field.Cells[i].What != "crop" {
				continue
			}

			// air core/crop/TileEntityCrop.java updateAirQuality
			air := (field.CropYlevel - 64) / 15
			if air > 4 {
				air = 4
			}
			if air < 0 {
				air = 0
			}
			nearby := 9
			for nr := -1; nr < 2; nr++ {
				for nc := -1; nc < 2; nc++ {
					n := field.Cells[getCellFromPosOffset(field, cr, cc, nr, nc)].What
					if n == "crop" || n == "block" {
						nearby--
					}
				}
			}
			air += nearby / 2
			if field.CropSkyAccess {
				air += 2
			}
			field.Cells[i].Environment.Air = air

			// nutrients core/crop/TileEntityCrop.java updateNutrients
			nutrients := field.CropBiomeNutrientsBonus + field.CropDirtBelow
			fert := 0
			if field.CropFertilizing {
				fert = 90 // assume automatic
			}
			nutrients += (fert + 19) / 20
			field.Cells[i].Environment.Nutrients = nutrients

			// humidity core/crop/TileEntityCrop.java updateHumidity
			humidity := field.CropBiomeHumidityBonus + 2 // assume water is present otherwise it is not a field
			waterstorage := 0
			if field.CropWatering {
				humidity += 2
				waterstorage = 200
			}
			humidity += (waterstorage + 24) / 25
			field.Cells[i].Environment.Humidity = humidity
		}
	}
}

func calcFieldAverageGrouth(field *cropField) float64 {
	avgc := 0.0
	avgs := 0.0
	for cr := 0; cr < field.Height; cr++ {
		for cc := 0; cc < field.Width; cc++ {
			i := cr*field.Height + cc
			if field.Cells[i].What != "crop" {
				continue
			}
			avgs += calcByParams(field.CropStat, field.Cells[i].Environment)
			avgc++
		}
	}
	return avgs / avgc
}

func calcFieldGrowthTime(field *cropField) (map[float64]int, int) {
	ret := map[float64]int{}
	c := 0
	for cr := 0; cr < field.Height; cr++ {
		for cc := 0; cc < field.Width; cc++ {
			i := cr*field.Height + cc
			if field.Cells[i].What != "crop" {
				continue
			}
			c++
			need := float64(field.CropStat.GrowthDuration) / float64(calcByParams(field.CropStat, field.Cells[i].Environment))
			if v, ok := ret[need]; ok {
				ret[need] = v + 1
			} else {
				ret[need] = 1
			}
		}
	}
	return ret, c
}

// func simulateFiledGrowth(field *cropField) int {

// 	 /
// }
