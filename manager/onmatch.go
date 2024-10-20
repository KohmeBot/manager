package manager

import (
	"github.com/kohmebot/plugin/pkg/chain"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

func (s *managerPlugin) onMatch(userId, groupId int64, ctxs []*zero.Ctx) {

	var err error
	var msgChain chain.MessageChain
	defer func() {
		if err != nil {
			s.env.Error(ctxs[0], err)
		}
	}()

	msgChain.Split(
		message.At(userId),
		message.Text(s.conf.RecallTips),
	)

	record := BanRecord{}
	db, err := s.env.GetDB()
	if err != nil {
		return
	}
	need, err := record.NeedBan(time.Duration(s.conf.BanCd)*time.Hour, userId, groupId, db)
	if err != nil {
		return
	}
	count := record.Count
	if need && count > 1 {
		// 禁言时间相当于 2的count次方分钟
		duration := int64((1 << count) * 60)
		ctxs[0].SetGroupBan(groupId, userId, duration)
		msgChain.Line()
		msgChain.Join(message.Text(s.conf.BanTips))
	}
	ctxs[0].Send(msgChain)

	for _, ctx := range ctxs {
		ctx.DeleteMessage(ctx.Event.MessageID)
	}

}
