package watcher

import "sync"

type RssWatcher struct {
	Telegram Telegram
}

func (r *RssWatcher) Start() {
	var w sync.WaitGroup
	var totalAgents int

	totalAgents = len(RssAgentsCollection)
	w.Add(totalAgents)

	for i := 0; i <= totalAgents-1; i++ {
		go func(i int) {
			RssAgentsCollection[i].Start(r.Telegram)
			w.Done()
		}(i)
	}

	w.Wait()
}
