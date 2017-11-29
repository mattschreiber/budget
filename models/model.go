package models

import (
  "time"
)
// Generic model for describing a budget or ledger
type Model struct {
  Id int `json:"id"`
  Credit int `json:"credit"`
  Debit int `json:"debit"`
  Trans_date time.Time `json:"trans_date"`
  // Store_name string `json:"store_name"`
  // Category_name string `json:"category_name"`
  // Store_id int `json:"store_id"`
  // Category_id int `json:"category_id"`
  Cat Category `json:"category"`
  St Store `json:"store"`
}
