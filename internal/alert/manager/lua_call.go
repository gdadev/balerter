package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

type options struct {
	Fields   []string
	Channels []string
	Quiet    bool
	Repeat   int
	Image    string
}

func defaultOptions() options {
	return options{
		Fields: nil,
		Quiet:  false,
		Repeat: 0,
	}
}

func (m *Manager) getAlertData(L *lua.LState) (alertName string, alertText string, alertOptions options, err error) {
	alertOptions = defaultOptions()

	alertNameLua := L.Get(1)
	if alertNameLua.Type() == lua.LTNil {
		err = fmt.Errorf("alert name must be provided")
		return
	}

	alertName = strings.TrimSpace(alertNameLua.String())
	if alertName == "" {
		err = fmt.Errorf("alert name must be not empty")
		return
	}

	alertTextLua := L.Get(2)
	if alertTextLua.Type() == lua.LTNil {
		return
	}

	alertText = alertTextLua.String()

	alertOptionsLua := L.Get(3)
	if alertOptionsLua.Type() == lua.LTNil {
		return
	}

	if alertOptionsLua.Type() != lua.LTTable {
		err = fmt.Errorf("options must be a table")
		return
	}

	err = gluamapper.Map(alertOptionsLua.(*lua.LTable), &alertOptions)
	if err != nil {
		err = fmt.Errorf("wrong options format: %v", err)
		return
	}

	return
}

func (m *Manager) luaCall(s *script.Script, alertLevel alert.Level) lua.LGFunction {
	return func(L *lua.LState) int {
		alertName, alertText, alertOptions, err := m.getAlertData(L)
		if err != nil {
			m.logger.Error("error get args", zap.Error(err))
			L.Push(lua.LString("error get arguments: " + err.Error()))
			return 1
		}

		metrics.SetAlertLevel(alertName, alertLevel)

		if len(alertOptions.Channels) == 0 {
			alertOptions.Channels = s.Channels
		}

		m.logger.Debug("call alert luaCall", zap.String("alertName", alertName), zap.String("scriptName", s.Name), zap.String("alertText", alertText), zap.Int("alertLevel", int(alertLevel)), zap.Any("alertOptions", alertOptions))

		a, err := m.engine.Alert().GetOrNew(alertName)
		if err != nil {
			m.logger.Error("error get alert from storage", zap.Error(err))
			L.Push(lua.LString("internal error get alert from storage: " + err.Error()))
			return 1
		}

		if a.Level() == alertLevel {
			a.Inc()

			if !alertOptions.Quiet && alertOptions.Repeat > 0 && a.Count()%alertOptions.Repeat == 0 {
				m.Send(alertLevel.String(), alertName, alertText, alertOptions.Channels, alertOptions.Fields, alertOptions.Image)
			}

			return 0
		}

		a.UpdateLevel(alertLevel)

		if !alertOptions.Quiet {
			m.Send(alertLevel.String(), alertName, alertText, alertOptions.Channels, alertOptions.Fields, alertOptions.Image)
		}

		m.engine.Alert().Release(a)

		return 0
	}
}
