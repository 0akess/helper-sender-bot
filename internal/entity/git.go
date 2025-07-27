package entity

import (
	"time"
)

type TTLReviewItem struct {
	SLA        int
	MRSize     int
	MRSizeName string
}

type GitlabConfig struct {
	Team              string
	GitlabURL         string
	ProjectName       string
	ProjectID         int
	ChannelID         string
	Reviewers         []string
	ReviewersCount    int
	TTLReview         []TTLReviewItem
	QAReviewers       string
	RequiresQaReview  bool
	PushQaAfterReview bool
}

type MergeRequestPayload struct {
	ProjectID      int
	ProjectName    string
	MrID           int
	MRTitle        string
	SourceBranch   string
	TargetBranch   string
	ProjectURL     string
	AuthorID       int
	AuthorUsername string
	IsDraft        bool
	MRState        string
}

type MergeRequestInfo struct {
	TotalLinesChanged int
	HasTestChanges    bool
}

type PostGitMR struct {
	TeamName     string
	ChannelID    string
	GitURL       string
	GitProjectID int
	GitMRID      int
	CreateAt     time.Time
	UpdateAT     time.Time
	PushedReview bool
	PostID       string
	TTLReview    TTLReviewItem
	IsDraft      bool
	Reviewers    string
}
