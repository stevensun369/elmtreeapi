package utils

import (
	"backend-go/env"
	"strconv"
)

func IsTermOne(dateDay string, dateMonth string) bool {
  var isTermOne bool

  dateDayInt, _ := strconv.Atoi(dateDay)
  dateMonthInt, _ := strconv.Atoi(dateMonth)

  if (dateMonthInt >= env.TermOneMonth && dateDayInt >= env.TermOneDay) || 
    (dateMonthInt <= env.TermTwoMonth && dateDayInt <= env.TermTwoDay) {
    isTermOne = true
  } else  {
    isTermOne = false
  }

  return isTermOne
}