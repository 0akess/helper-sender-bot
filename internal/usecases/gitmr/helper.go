package gitmr

import (
	"context"
	"fmt"
	"helper-sender-bot/internal/entity"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

// getGitCfg метод для получения конфигурации git
func (gm *GitMR) getGitCfg(ctx context.Context, mr entity.MergeRequestPayload) (entity.GitlabConfig, string, error) {
	base := normalizeInstanceURL(mr.ProjectURL)
	gitCfg, err := gm.repo.GetGitlabConfig(ctx, mr.ProjectID, base)
	if err != nil {
		gm.log.Error("Не удалось получить конфигурацию к git", "project", mr.ProjectID, "base", base)
		return entity.GitlabConfig{}, "", err
	}
	return gitCfg, base, nil
}

// buildNewMRMessage собирает текст сообщения для отправки в мм при новом МР
func buildNewMRMessage(
	mr entity.MergeRequestPayload,
	cfg entity.GitlabConfig,
	mrInfo entity.MergeRequestInfo,
) (mrMsg string, ttlR entity.TTLReviewItem, reviewersLine string) {

	reviewersLine = getReviewersLine(mr, cfg)

	qaLine := getQaReviewerLine(cfg, mrInfo)

	slaLine, ttlRule := getSlaAndSizeMRLine(cfg, mrInfo)

	link := fmt.Sprintf("%s/-/merge_requests/%d", mr.ProjectURL, mr.MrID)

	return fmt.Sprintf(
			"**Новый МР в проекте**: %s\n**Название**: [%s](%s)\n**Автор**: @%s\n\n**Ревьюеры**: %s%s%s",
			cfg.ProjectName, escape(mr.MRTitle), link, mr.AuthorUsername, reviewersLine, qaLine, slaLine,
		),
		ttlRule, reviewersLine
}

func getSlaAndSizeMRLine(cfg entity.GitlabConfig, mrInfo entity.MergeRequestInfo) (string, entity.TTLReviewItem) {
	slaMsg := ""
	found := false
	var rule entity.TTLReviewItem
	if len(cfg.TTLReview) > 0 {
		sort.Slice(cfg.TTLReview, func(i, j int) bool {
			return cfg.TTLReview[i].MRSize < cfg.TTLReview[j].MRSize
		})

		for _, ttl := range cfg.TTLReview {
			if mrInfo.TotalLinesChanged > ttl.MRSize {
				continue
			}
			rule = ttl
			found = true
			break
		}
		if found {
			sla := rule.SLA
			hours := sla / 60
			minutes := sla % 60
			if hours == 0 {
				slaMsg = fmt.Sprintf(
					"\n**Размер МР и SLA**: %s | %d мин\n",
					rule.MRSizeName, minutes,
				)
			} else {
				slaMsg = fmt.Sprintf(
					"\n**Размер МР и SLA**: %s | %d часов %d мин\n",
					rule.MRSizeName, hours, minutes,
				)
			}
		}
	}
	return slaMsg, rule
}

func getQaReviewerLine(cfg entity.GitlabConfig, mrInfo entity.MergeRequestInfo) string {
	qaLine := ""
	if cfg.RequiresQaReview && cfg.QAReviewers != "" && mrInfo.HasTestChanges {
		qaLine = fmt.Sprintf("\n**Нуждается в QA**: %s", cfg.QAReviewers)
	}
	return qaLine
}

func getReviewersLine(mr entity.MergeRequestPayload, cfg entity.GitlabConfig) string {
	candidates := make([]string, 0, len(cfg.Reviewers))
	for _, r := range cfg.Reviewers {
		if !strings.EqualFold(r, mr.AuthorUsername) {
			candidates = append(candidates, r)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	revC := cfg.ReviewersCount
	if revC > len(candidates) {
		revC = len(candidates)
	}
	reviewersLine := strings.Join(candidates[:revC], ", ")

	return reviewersLine
}

func escape(s string) string {
	r := strings.NewReplacer(`*`, `\*`, "_", `\_`, "`", "\\`")
	return r.Replace(s)
}

func normalizeInstanceURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
}
