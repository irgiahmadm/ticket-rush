package domain

import "time"
type Event struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Date time.Time `json:"date"`
}
