package manager

type Config struct {
	// 禁言相关配置
	SensitiveConfig `mapstructure:"sensitive"`

	// 新人加群的欢迎词
	joinGroup string `mapstructure:"join_group"`
	// 新人加群的额外提示词
	joinGroupTips string `mapstructure:"join_group_tips"`
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
