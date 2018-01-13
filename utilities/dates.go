package utilities

import (
"fmt"
"time"
)

type DefaultDate struct {
  Today time.Time
  CurrentMonth int
  CurrentYear int
}

func (defaultDate DefaultDate) GetEst() (*time.Location) {
    utc, err := time.LoadLocation("America/New_York")
    if err != nil {
      fmt.Println(err)
    }
    return utc
}

func InitDates() DefaultDate {
  defaultDate := DefaultDate{}
  defaultDate.Today = time.Now().In(defaultDate.GetEst())
  defaultDate.CurrentMonth = int(time.Now().In(defaultDate.GetEst()).Month())
  defaultDate.CurrentYear = time.Now().In(defaultDate.GetEst()).Year()
  return defaultDate
}
