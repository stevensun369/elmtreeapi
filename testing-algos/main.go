package main

import (
	"fmt"
	"math"
)

func main() {
  marks := []int{9, 10, 10}
  finalMark := 9

  markValuesCounter := 0
  marksCounter := 0
  for _, mark := range marks {
    markValuesCounter += mark
    marksCounter += 1
  }

  marksValue := float64(markValuesCounter) / float64(marksCounter)
  marksValueFormat := math.Floor(marksValue * 100) / 100

  value := math.Round((marksValueFormat * 3 + float64(finalMark)) / 4)

  fmt.Printf("%v \n", marksValueFormat)
  fmt.Printf("%v", value)
}