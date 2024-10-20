package manager

import (
	"fmt"
	"github.com/kohmebot/manager/manager/matcher/textmatcher"
	"github.com/kohmebot/plugin/pkg/chain"
	"github.com/kohmebot/plugin/pkg/gopool"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

func (s *managerPlugin) SetOnWord(engine *zero.Engine) {
	engine.OnCommand("word", s.env.SuperUser().Rule()).Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		var err error
		defer func() {
			if err != nil {
				s.env.Error(ctx, err)
				return
			}
		}()
		err = ctx.Parse(&cmd)
		if err != nil {
			return
		}
		content := cmd.Args

		err = s.appendToFile(s.dictPath, content)
		if err != nil {
			gopool.Go(func() {
				var msgChain chain.MessageChain
				ctx.Send(msgChain.Join(message.Text(fmt.Sprintf("%s 写入成功！", content))))
			})
		}
		words := s.tryRead(s.dictPath)
		dwords := make([]string, len(words)+len(s.conf.words))
		copy(dwords, words)
		copy(dwords[len(words):], s.conf.words)
		s.matcher.Swap(textmatcher.NewTextMatcher(dwords...))
	})
}

func (s *managerPlugin) SetOnReload(engine *zero.Engine) {
	engine.OnCommand("reload", s.env.SuperUser().Rule()).Handle(func(ctx *zero.Ctx) {
		var err error
		defer func() {
			if err != nil {
				s.env.Error(ctx, err)
				return
			}
		}()
		words := s.tryRead(s.dictPath)
		dwords := make([]string, len(words)+len(s.conf.words))
		copy(dwords, words)
		copy(dwords[len(words):], s.conf.words)
		s.matcher.Swap(textmatcher.NewTextMatcher(dwords...))
		if err == nil {
			gopool.Go(func() {
				var msgChain chain.MessageChain
				ctx.Send(msgChain.Join(message.Text("reload成功！")))
			})
		}
	})
}

// SetOnJoinRequest 设置处理加群申请
func (s *managerPlugin) SetOnJoinRequest(engine *zero.Engine) {
	engine.OnRequest(s.env.Groups().Rule()).Handle(func(ctx *zero.Ctx) {
		if ctx.Event.RequestType != "group" {
			return
		}
		comment := strings.TrimSpace(ctx.Event.Comment)
		var pass bool
		for _, answer := range s.conf.requestAnswers {
			if answer == comment {
				pass = true
				break
			}
		}
		gopool.Go(func() {
			if pass {
				ctx.SetGroupAddRequest(ctx.Event.Flag, ctx.Event.SubType, true, "")
				return
			}
			if s.conf.refuse {
				ctx.SetGroupAddRequest(ctx.Event.Flag, ctx.Event.SubType, false, s.conf.refuseReason)
				return
			}
		})

	})
}

// SetOnJoinGroup 设置有新人加群
func (s *managerPlugin) SetOnJoinGroup(engine *zero.Engine) {
	engine.OnNotice(s.env.Groups().Rule()).Handle(func(ctx *zero.Ctx) {
		if ctx.Event.NoticeType != "group_increase" {
			return
		}
		var msgChain chain.MessageChain
		msgChain.Line(message.Text(s.conf.joinGroup), message.At(ctx.Event.UserID))
		msgChain.Join(message.Text(s.conf.joinGroupTips))
		gopool.Go(func() {
			ctx.Send(msgChain)
		})

	})
}
