package watcher

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

const debounceDelay = 300 * time.Millisecond

// FileChangedMsg signals that a watched file has changed.
type FileChangedMsg struct{}

// ErrorMsg signals that an error occurred during file watching.
type ErrorMsg struct{ Err error }

// Watcher wraps an fsnotify-based file watcher.
type Watcher struct {
	fsw  *fsnotify.Watcher
	ch   chan tea.Msg
	done chan struct{}
}

// New creates a Watcher that monitors the given paths.
// Non-existent paths are silently ignored.
func New(paths []string) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue // skip non-existent paths
		}
		_ = fsw.Add(p)
	}

	w := &Watcher{
		fsw:  fsw,
		ch:   make(chan tea.Msg, 1),
		done: make(chan struct{}),
	}
	go w.loop()
	return w, nil
}

// loop receives fsnotify events, debounces them, and forwards to ch.
func (w *Watcher) loop() {
	var timer *time.Timer

	for {
		select {
		case <-w.done:
			if timer != nil {
				timer.Stop()
			}
			return

		case event, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			// Filter only relevant events
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}
			// Debounce: reset the timer
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounceDelay, func() {
				select {
				case w.ch <- FileChangedMsg{}:
				default: // non-blocking: skip if a message is already pending
				}
			})

		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			select {
			case w.ch <- ErrorMsg{Err: err}:
			default:
			}
		}
	}
}

// WaitForChange returns a tea.Cmd that waits for the next file change message.
func (w *Watcher) WaitForChange() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-w.ch
		if !ok {
			return nil
		}
		return msg
	}
}

// Close shuts down the watcher and releases its resources.
func (w *Watcher) Close() {
	close(w.done)
	w.fsw.Close()
}
