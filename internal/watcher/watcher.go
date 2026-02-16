package watcher

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

const debounceDelay = 300 * time.Millisecond

// FileChangedMsg는 감시 중인 파일이 변경되었음을 알린다.
type FileChangedMsg struct{}

// ErrorMsg는 파일 감시 중 에러가 발생했음을 알린다.
type ErrorMsg struct{ Err error }

// Watcher는 fsnotify 기반 파일 감시자를 래핑한다.
type Watcher struct {
	fsw  *fsnotify.Watcher
	ch   chan tea.Msg
	done chan struct{}
}

// New는 주어진 경로들을 감시하는 Watcher를 생성한다.
// 존재하지 않는 경로는 무시한다.
func New(paths []string) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue // 존재하지 않는 경로 무시
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

// loop는 fsnotify 이벤트를 수신하고 디바운싱 후 ch로 전달한다.
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
			// 관심 있는 이벤트만 필터링
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}
			// 디바운스: 타이머 리셋
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounceDelay, func() {
				select {
				case w.ch <- FileChangedMsg{}:
				default: // non-blocking: 이미 메시지가 대기 중이면 스킵
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

// WaitForChange는 다음 파일 변경 메시지를 기다리는 tea.Cmd를 반환한다.
func (w *Watcher) WaitForChange() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-w.ch
		if !ok {
			return nil
		}
		return msg
	}
}

// Close는 watcher를 종료하고 리소스를 정리한다.
func (w *Watcher) Close() {
	close(w.done)
	w.fsw.Close()
}
