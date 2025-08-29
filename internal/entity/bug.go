package entity

import "time"

type BugSLA struct {
	Priority string
	SLA      int
}

type BugConfig struct {
	ID        int
	TeamName  string
	ChannelId string
	TrackURL  string
	TrackName string
	BugSLA    []BugSLA
	CreateAT  time.Time
	UpdateAT  time.Time
}

type Bug struct {
	ID       int
	CfgBugID int
	BugSLA   BugSLA
	CreateAT time.Time
	UpdateAT time.Time
}
