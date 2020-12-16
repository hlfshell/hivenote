package podcast

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"go.uber.org/atomic"
)

type Manager struct {
	db *sql.DB

	parser *gofeed.Parser

	lock sync.Mutex

	lastUpdated time.Time

	scan      *atomic.Bool
	scanEvery time.Duration
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db:          db,
		parser:      gofeed.NewParser(),
		lock:        sync.Mutex{},
		lastUpdated: time.Time{},
		scan:        atomic.NewBool(false),
		scanEvery:   5 * time.Minute,
	}
}

//GetSubscriptions gets the subscriptions and populates their podcasts
//as well
func (manager *Manager) GetSubscriptions() ([]Podcast, error) {
	return nil, nil
}

func (manager *Manager) Scan() {
	manager.scan.Store(true)

	for {
		//Check to see if we should stop (ie we were cancelled)
		if !manager.scan.Load() {
			break
		}

		// Ensure it's been long enough
		if !time.Now().After(manager.lastUpdated.Add(manager.scanEvery)) {
			// No need to react quickly, and since this
			// this is always threaded, we can get away
			// with an egregious sleep
			time.Sleep(1 * time.Second)
			continue
		}

		// Finally, we can start the process
		err := manager.checkSubscriptions()
		if err != nil {
			fmt.Println("An error has occured checking for subscriptions", err.Error())
		}
		manager.lastUpdated = time.Now()
	}
}

func (manager *Manager) StopScan() {
	manager.scan.Store(false)
}

func (manager *Manager) checkSubscriptions() error {
	// First we need to get the list of subscriptions
	subscriptions, err := manager.GetSubscriptions()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	//Results prep
	errors := []error{} //Switch this to a map
	errorLock := sync.Mutex{}
	queue := []Podcast{}
	queueLock := sync.Mutex{}

	for _, podcast := range subscriptions {
		wg.Add(1)
		go func() {
			defer wg.Done()

			feed, err := manager.parser.ParseURL(podcast.URL)
			if err != nil {
				errorLock.Lock()
				errors = append(errors, err)
				errorLock.Unlock()
			}

			// Only proceed if our RSS feed fits our expectations and
			// has items to deal with
			if len(feed.Items) > 0 && len(feed.Items[0].Enclosures) > 0 {
				if podcast.LatestPodcast == nil || feed.Items[0].GUID != podcast.LatestPodcast.GUID {
					// Queue a new episode to download
					queueLock.Lock()
					podcast.LatestPodcast = &Episode{
						GUID:     feed.Items[0].GUID,
						MediaURL: feed.Items[0].Enclosures[0].URL,
						Position: 0,
					}
					queue = append(queue, podcast)
					queueLock.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	//Handle errors here!

	//If there's nothing to do here, we're done.
	if len(queue) <= 0 {
		return nil
	}

	// Now that we have a list of subscriptions that we should
	// update, let's start that process.
	wg = sync.WaitGroup{}
	for _, subscription := range queue {
		func() {
			defer wg.Done()
			err := subscription.FetchLatestPodcast("")
			if err != nil {

			}
		}()
	}
	wg.Wait()

	return nil
}
