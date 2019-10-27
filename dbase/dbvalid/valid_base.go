package dbvalid

import (
	"time"
)

/*
CREATE TABLE `userinfo` (
	`uid` INTEGER PRIMARY KEY AUTOINCREMENT,
	`username` VARCHAR(64) NULL,
	`departname` VARCHAR(64) NULL,
	`created` DATE NULL
);

CREATE TABLE note (
	device VARCHAR(64) NOT NULL,
    currency INTEGER NOT NULL DEFAULT 0,
    nominal INTEGER NOT NULL DEFAULT 0,
    count INTEGER NOT NULL DEFAULT 0,
    amount INTEGER NOT NULL DEFAULT 0,
    UNIQUE (device, currency, nominal)
);
 */

const (
	StateUndefined uint16 = iota
	StateActive
	StateClosed
)
type ObjBatch struct {
	Id       int
	State    uint16
	Device   string
	Count    uint16
	Opened   time.Time
	Closed   time.Time
}

type ObjDeposit struct {
	Id       int
	BatchId  int
	ExtraId  int
	Device   string
	Currency uint16
	Nominal  float32
	Count    uint16
	Amount   float32
	Created  time.Time
}

type ObjBalance struct {
	Id       int
	BatchId  int
	Device   string
	Currency uint16
	Nominal  float32
	Count    uint16
	Amount   float32
	Created  time.Time
}



