package manager

import (
	"github.com/kohmebot/manager/manager/matcher"
	"github.com/kohmebot/manager/manager/matcher/textmatcher"
	"github.com/kohmebot/plugin"
	"github.com/kohmebot/plugin/pkg/command"
	"github.com/kohmebot/plugin/pkg/version"
	"github.com/wdvxdr1123/ZeroBot"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type dictWords struct {
	words []string
	rw    sync.RWMutex
}

func (w *dictWords) Load() []string {
	w.rw.RLock()
	defer w.rw.RUnlock()
	return w.words
}
func (w *dictWords) Set(words []string) {
	w.rw.Lock()
	defer w.rw.Unlock()
	w.words = words
}

type managerPlugin struct {
	matcher   matcher.SafeMatcher
	conf      Config
	env       plugin.Env
	dictPath  string
	dictWords dictWords
}

func NewPlugin() plugin.Plugin {
	return &managerPlugin{}
}

func (s *managerPlugin) Init(engine *zero.Engine, env plugin.Env) error {
	s.env = env
	err := env.GetConf(&s.conf)
	if err != nil {
		return err
	}

	r := &BanRecord{}
	db, err := env.GetDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&r)
	if err != nil {
		return err
	}
	path, err := env.FilePath()
	if err != nil {
		return err
	}
	s.dictPath = filepath.Join(path, "words.txt")
	words := s.tryRead(s.dictPath)
	s.dictWords.Set(words)

	dwords := make([]string, len(words)+len(s.conf.Words))
	copy(dwords, words)
	copy(dwords[len(words):], s.conf.Words)

	s.matcher.Swap(textmatcher.NewTextMatcher(dwords...))
	s.matcher.SetOnMatch(s.onMatch)
	g := env.Groups()

	engine.OnMessage(g.Rule()).Handle(func(ctx *zero.Ctx) {
		s.matcher.Submit(ctx)
	})
	s.SetOnWord(engine)
	s.SetOnReload(engine)
	s.SetOnJoinRequest(engine)
	s.SetOnJoinGroup(engine)
	return nil
}

func (s *managerPlugin) tryRead(path string) (words []string) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(contentBytes)
	return strings.Fields(content)
}

func (s *managerPlugin) appendToFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

func (s *managerPlugin) Name() string {
	return "manager"
}

func (s *managerPlugin) Description() string {
	return "群管理插件"
}

func (s *managerPlugin) Commands() command.Commands {
	return command.NewCommands()
}

func (s *managerPlugin) Version() version.Version {
	return version.NewVersion(0, 0, 21)
}
