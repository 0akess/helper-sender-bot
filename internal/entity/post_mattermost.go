package entity

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

func (p *Post) IsExistReaction(name string) bool {
	for _, r := range p.Metadata.Reactions {
		if r.EmojiName == name {
			return true
		}
	}
	return false
}
