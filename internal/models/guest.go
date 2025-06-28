package models

import "time"

type GuestUser struct {
	UUID   		string 		`json:"uuid" db:"uuid"`
	Banned 		bool 			`json:"banned" db:"banned"`
	CreatedAt 	time.Time	`json:"created_at" db:"created_at"`
}