package entity

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	Team                  string
	ChannelID             string
	DutyTtlInMinute       int
	DutyRepeatTtlInMinute int
	EmojiStart            string
	EmojiDone             string
	WorkdayStart          int
	WorkdayEnd            int
}

type Post struct {
	ID       string `json:"id"`
	RootId   string `json:"root_id"`
	Message  string `json:"message"`
	CreateAt int64  `json:"create_at"`
	Type     string `json:"type"`
	Metadata struct {
		Reactions []struct {
			EmojiName string `json:"emoji_name"`
		} `json:"reactions"`
	} `json:"metadata"`
}

type Team struct {
	Name    string
	Token   uuid.UUID
	LeadEID string
}

type PostsInfoDuty struct {
	ChannelID  string
	PostID     string
	CreateAt   time.Time
	LastPushAt time.Time
	InProgress bool
}

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
	MRIID          int
	MRTitle        string
	SourceBranch   string
	TargetBranch   string
	ProjectURL     string
	AuthorID       int
	AuthorUsername string
	IsDraft        bool
	MRState        string
}

type AuthMeta struct {
	Team  string
	Token uuid.UUID
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
	PushedReview bool
	PostID       string
	TTLReview    TTLReviewItem
	Reviewers    string
}
