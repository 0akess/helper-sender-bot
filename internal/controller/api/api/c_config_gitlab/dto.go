package c_config_gitlab

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/entity"
	"strconv"
)

type ttlReviewItem struct {
	SLA        int    `json:"sla"`
	MRSize     int    `json:"mr_size"`
	MRSizeName string `json:"mr_size_name"`
}

func checkAndBuildTTLReview(req []ttlReviewItem) ([]entity.TTLReviewItem, error) {
	ttlReview := make([]entity.TTLReviewItem, len(req))
	for i, v := range req {
		if v.SLA <= 0 {
			return nil, fmt.Errorf("sla should be represented as element of object ttl_review")
		}
		if v.MRSize <= 0 {
			return nil, fmt.Errorf("mr_size should be represented as element of object ttl_review")
		}
		if v.MRSizeName == "" {
			return nil, fmt.Errorf("mr_size_name should be represented as element of object ttl_review")
		}
		ttlReview[i] = entity.TTLReviewItem{
			SLA:        v.SLA,
			MRSize:     v.MRSize,
			MRSizeName: v.MRSizeName,
		}
	}
	return ttlReview, nil
}

func toGetTTL(src []entity.TTLReviewItem) []ttlReviewItem {
	out := make([]ttlReviewItem, len(src))
	for i, v := range src {
		out[i] = ttlReviewItem{
			SLA:        v.SLA,
			MRSize:     v.MRSize,
			MRSizeName: v.MRSizeName,
		}
	}
	return out
}

func gitUrlAndIdQuery(e echo.Context) (string, int, error) {
	pidStr := e.QueryParam("project_id")
	gitURL := e.QueryParam("git_url")
	if pidStr == "" || gitURL == "" {
		return "", 0, fmt.Errorf("project_id and git_url are required")
	}
	projectID, err := strconv.Atoi(pidStr)
	if err != nil {
		return "", 0, fmt.Errorf("project_id must be an integer")
	}
	return gitURL, projectID, nil
}
