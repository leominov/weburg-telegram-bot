package watcher

import "sync"

func (r *RssWatcher) StartWatch() {
	var w sync.WaitGroup
	var totalAgents int

	totalAgents = len(RssAgentsCollection)
	w.Add(totalAgents)

	for i := 0; i <= totalAgents-1; i++ {
		go func(i int) {
			RssAgentsCollection[i].Start(r.Sender)
			w.Done()
		}(i)
	}

	w.Wait()
}
