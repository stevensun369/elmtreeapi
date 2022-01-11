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

  // new algo
  if dateMonthInt >= env.TermOneMonth {
    if dateMonthInt == env.TermOneMonth {
      if dateDayInt >= env.TermOneDay {
        isTermOne = true
      } else {
        isTermOne = false
      }
    } else {
      isTermOne = true
    }
  }
  if dateMonthInt <= env.TermTwoMonth {
    if dateMonthInt == env.TermTwoMonth {
      if dateDayInt <= env.TermTwoDay {
        isTermOne = true
      } else {
        isTermOne = false
      }
    } else {
      isTermOne = true
    }
  }

  return isTermOne
}