package manager

type Config struct {
	// 禁言相关配置
	SensitiveConfig `yaml:"sensitive"`
	// 入群相关配置
	JoinGroupConfig `yaml:"join_group"`
}

type SensitiveConfig struct {
	// 需要过滤的违禁词
	Words []string `yaml:"words"`
	// ban次数的重置cd，单位小时
	BanCd int64 `yaml:"ban_cd"`
	// 撤回消息后的提示词
	RecallTips string `yaml:"recall_tips"`
	// 禁言的提示词
	BanTips string `yaml:"ban_tips"`
}

type JoinGroupConfig struct {
	// 是否开启入群欢迎
	EnableHello bool `yaml:"enable_hello"`
	// 新人加群的欢迎词
	JoinGroup string `yaml:"join_group"`
	// 新人加群的额外提示词
	JoinGroupTips string `yaml:"join_group_tips"`
	// 入群的请求的答案
	RequestAnswers []string `yaml:"request_answers"`
	// 答案不对时是否直接拒绝入群
	Refuse bool `yaml:"refuse"`
	// 拒绝理由
	RefuseReason string `yaml:"refuse_reason"`
}
