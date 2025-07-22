package wh_gitlab

import (
	"github.com/labstack/echo/v4"
	"helper-sender-bot/internal/entity"
	"log/slog"
	"net/http"
)

type mergeRequestPayload struct {
	ProjectID      int    `json:"project_id"`
	ProjectName    string `json:"project_name"`
	MRIID          int    `json:"mr_iid"`
	MRTitle        string `json:"mr_title"`
	SourceBranch   string `json:"source_branch"`
	TargetBranch   string `json:"target_branch"`
	ProjectURL     string `json:"project_url"`
	AuthorID       int    `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	IsDraft        bool   `json:"is_draft"`
	MRState        string `json:"mr_state"`
}

func (gc *GitlabController) handleGitlabWebhook(e echo.Context) error {
	token := e.Request().Header.Get("X-Gitlab-Token")
	if token != gc.token {
		slog.Warn("GitLab webhook: invalid secret", "provided", token)
		return e.NoContent(http.StatusUnauthorized)
	}

	var req mergeRequestPayload
	if err := e.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body := entity.MergeRequestPayload{
		ProjectID:      req.ProjectID,
		ProjectName:    req.ProjectName,
		MRIID:          req.MRIID,
		MRTitle:        req.MRTitle,
		SourceBranch:   req.SourceBranch,
		TargetBranch:   req.TargetBranch,
		ProjectURL:     req.ProjectURL,
		AuthorID:       req.AuthorID,
		AuthorUsername: req.AuthorUsername,
		IsDraft:        req.IsDraft,
		MRState:        req.MRState,
	}
	if body.IsDraft {
		return e.NoContent(http.StatusOK)
	}

	switch body.MRState {
	case "opened":
		go func() {
			gc.pushMR.SendPushNewMR(gc.ctx, body)
		}()
	case "merged":
		go func() {
			gc.pushMR.SendPushMergedMR(gc.ctx, body)
		}()
	case "closed":
		go func() {
			gc.pushMR.SendPushClosedMR(gc.ctx, body)
		}()
	}

	return e.NoContent(http.StatusOK)
}
