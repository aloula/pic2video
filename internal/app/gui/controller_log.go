package gui

import "sync"

type LogController struct {
	mu    sync.Mutex
	store *LogStore
	view  *LogView
}

func NewLogController(store *LogStore, view *LogView) *LogController {
	return &LogController{store: store, view: view}
}

func (c *LogController) Append(stream, line string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Append(stream, line)
	c.view.SetText(c.store.Text())
}
