package matcher

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"sync"
)

// AsyncMatcher 异步的匹配器
type AsyncMatcher interface {
	// SetOnMatch 设置触发匹配
	SetOnMatch(func(userId, groupId int64, messages []*zero.Ctx))
	// Submit 提交事件
	Submit(event *zero.Ctx)
}

// SafeMatcher 可以线程安全地进行更换matcher
type SafeMatcher struct {
	rw sync.RWMutex
	AsyncMatcher
}

func (s *SafeMatcher) SetOnMatch(f func(userId int64, groupId int64, messages []*zero.Ctx)) {
	s.rw.RLock()
	matcher := s.AsyncMatcher
	s.rw.RUnlock()
	matcher.SetOnMatch(f)
}

func (s *SafeMatcher) Submit(event *zero.Ctx) {
	s.rw.RLock()
	matcher := s.AsyncMatcher
	s.rw.RUnlock()
	matcher.Submit(event)
}

func (s *SafeMatcher) Swap(m AsyncMatcher) (swapped AsyncMatcher) {
	s.rw.Lock()
	swapped = s.AsyncMatcher
	s.AsyncMatcher = m
	s.rw.Unlock()
	return
}
