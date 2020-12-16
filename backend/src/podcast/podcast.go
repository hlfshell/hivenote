package podcast

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Podcast struct {
	ID             string     `json:"id,omitempty" db:"id"`
	Name           string     `json:"name,omitempty" db:"name"`
	URL            string     `json:"url,omitempty" db:"url"`
	Episodes       []*Episode `json:"episodes"`
	CurrentEpisode *Episode   `json:"current_episode,omitempty"`
}

func (manager *Manager) CreatePodcast(podcast *Podcast) error {
	return nil
}

func (manager *Manager) GetPodcasts() ([]*Podcast, error) {
	rows, err := manager.db.Query(``)
	if err != nil {
		return nil, err
	}

	podcasts := []*Podcast{}
	for rows.Next() {

	}

	return podcasts, nil
}

func (podcast *Podcast) Save(db *sql.DB) error {
	return nil
}

func (podcast *Podcast) FetchLatestPodcast(basepath string) error {
	var err error
	//This deferred error func is to check for an existing err and,
	//if there is one, clean up after ourselves in terms of files
	defer func() {

	}()
	//First, prepare the file
	file, err := os.Create(fmt.Sprintf("%s/%s", basepath, podcast.LatestPodcast.GUID))

	response, err := http.Get(podcast.LatestPodcast.MediaURL)
	if err != nil {
		return err
	}

	//If it's not a 2XX, we consider it an error
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Status code is %d, not 2XX for subscription %s", response.StatusCode, podcast.LatestPodcast.GUID))
	}

	//Now begin copying the file over
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil
	}

	//Not done
	return nil
}
