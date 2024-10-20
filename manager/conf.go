package manager

type Config struct {
	// 禁言相关配置
	SensitiveConfig `mapstructure:"sensitive"`
	// 入群相关配置
	JoinGroupConfig `mapstructure:"join_group"`
}

type SensitiveConfig struct {
	// 需要过滤的违禁词
	words []string `mapstructure:"words"`
	// ban次数的重置cd，单位小时
	banCd int64 `mapstructure:"ban_cd"`
	// 撤回消息后的提示词
	recallTips string `mapstructure:"recall_tips"`
	// 禁言的提示词
	banTips string `mapstructure:"ban_tips"`
}

type JoinGroupConfig struct {
	// 新人加群的欢迎词
	joinGroup string `mapstructure:"join_group"`
	// 新人加群的额外提示词
	joinGroupTips string `mapstructure:"join_group_tips"`
	// 入群的请求的答案
	requestAnswers []string `mapstructure:"request_answers"`
	// 答案不对时是否直接拒绝入群
	refuse bool `mapstructure:"refuse"`
	// 拒绝理由
	refuseReason string `mapstructure:"refuse_reason"`
}
