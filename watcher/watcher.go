package watcher

import "sync"

type Watcher struct {
	Telegram Telegram
}

func (w *Watcher) Start() {
	var wg sync.WaitGroup
	var totalAgents int

	totalAgents = len(AgentsCollection)
	wg.Add(totalAgents)

	for i := 0; i <= totalAgents-1; i++ {
		go func(i int) {
			AgentsCollection[i].Start(w.Telegram)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
