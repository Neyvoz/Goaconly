package domain

import (
	"context"
	"time"
)

type CmdType int

const (
	CmdAdd CmdType = iota
	CmdRemove
	CmdUpdateInterval
)

type SchedulerCmd struct {
	Type              CmdType
	CmdAdd            string
	CmdRemove         time.Duration
	CmdUpdateInterval *Target
}

type TickerEntry struct {
	ticker *time.Ticker
	cancel context.CancelFunc
}
