package models

import (
  "testing"
  "time"
  "fmt"
)

func TestPrevDate(t *testing.T) {
  type testDate struct {
    startDate time.Time
    prevDate time.Time
  }
  tests := []testDate{ {time.Date(2017, 11, 14, 0, 0, 0, 0, time.UTC), time.Date(2017, 11, 1, 0, 0, 0, 0, time.UTC)},
    {time.Date(2017, 11, 16, 0, 0, 0, 0, time.UTC), time.Date(2017, 11, 15, 0, 0, 0, 0, time.UTC)},
    {time.Date(2017, 12, 31, 0, 0, 0, 0, time.UTC), time.Date(2017, 12, 15, 0, 0, 0, 0, time.UTC)},
    {time.Date(2017, 12, 15, 0, 0, 0, 0, time.UTC), time.Date(2017, 12, 15, 0, 0, 0, 0, time.UTC)},
  }

  for _, test := range tests {
    resultDate := prevPayDate(test.startDate)
    fmt.Println("test")
    if resultDate != test.prevDate {
      t.Errorf("Incorrect date %v", resultDate  )
    }
  }
}
