package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	"go.uber.org/zap"
)

func (m *Manager) Send(level alert.Level, alertName, text string, channels []string, fields []string, chartData *chartModule.Data) {
	chs := make(map[string]alertChannel)

	if len(channels) > 0 {
		for _, channelName := range channels {
			ch, ok := m.channels[channelName]
			if !ok {
				m.logger.Error("channel not found", zap.String("channel name", channelName))
				continue
			}
			chs[channelName] = ch
		}
	} else {
		chs = m.channels
	}

	if len(chs) == 0 {
		m.logger.Error("empty channels")
		return
	}

	for name, module := range chs {
		if err := module.Send(level, message.New(alertName, text, fields), chartData); err != nil {
			m.logger.Error("error send message to channel", zap.String("channel name", name), zap.Error(err))
		}
	}
}
