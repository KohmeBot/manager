package textmatcher

import (
	"github.com/kohmebot/manager/manager/ahocorasick"
	"github.com/kohmebot/manager/manager/matcher"
	"github.com/kohmebot/plugin/pkg/gopool"
	zero "github.com/wdvxdr1123/ZeroBot"
	"sync/atomic"
	"unsafe"
)

type callback func(userId, groupId int64, messages []*zero.Ctx)

type TextMatcher struct {
	matcher *ahocorasick.Matcher
	on      atomic.Value
}

func NewTextMatcher(texts ...string) matcher.AsyncMatcher {
	m := ahocorasick.NewStringMatcher(texts)
	return &TextMatcher{matcher: m}
}

func (t *TextMatcher) SetOnMatch(f func(userId, groupId int64, messages []*zero.Ctx)) {
	t.on.Store(callback(f))
}

func (t *TextMatcher) Submit(ctx *zero.Ctx) {
	event := ctx.Event
	if event.GroupID == 0 {
		return
	}
	msgs := event.Message
	var texts []string
	for _, msg := range msgs {
		if msg.Type == "text" {
			v, ok := msg.Data["text"]
			if ok {
				texts = append(texts, v)
			}
		}
	}
	if len(texts) <= 0 {
		return
	}
	for _, text := range texts {
		if t.matcher.Contains(unsafe.Slice(unsafe.StringData(text), len(text))) {
			cb, ok := t.on.Load().(callback)
			if ok {
				gopool.Go(func() {
					cb(event.UserID, event.GroupID, []*zero.Ctx{ctx})
				})
				break
			}
		}
	}

}
