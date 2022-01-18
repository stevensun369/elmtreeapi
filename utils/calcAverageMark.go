package utils

import (
	"math"
)

func CalculateAverageMark(marks []int, finalMark int) int {
  finalMarkFloat := float64(finalMark)
  
  markValuesCounter := 0
  marksCounter := 0
  for _, mark := range marks {
    markValuesCounter += mark
    marksCounter += 1
  }

  marksValue := float64(markValuesCounter) / float64(marksCounter)
  marksValueFormat := math.Floor(marksValue * 100) / 100

  value := math.Round((marksValueFormat * 3 + finalMarkFloat) / 4)


  return int(value)
}