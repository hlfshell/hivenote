package podcast

type Episode struct {
	GUID       string `json:"guid,omitempty" db:"guid"`
	Podcast    string `json:"podcast",omitempty" db:"podcast"`
	MediaURL   string `json:"media_url,omitempty" db:"media_url"`
	Filepath   string `json:"filepath,omitempty" db:"filepath"`
	InProgress bool   `json:"in_progress,omitempty" db:"in_progress"`
	Position   int64  `json:"position,omitempty" db:"position"`
}

func (manager *Manager) GetEpisodeFromPodcast(podcast string) ([]*Episode, error) {
	return nil, nil
}
