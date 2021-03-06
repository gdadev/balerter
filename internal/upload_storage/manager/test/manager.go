package manager

import (
	"github.com/balerter/balerter/internal/config/storages/upload"
	moduleMock "github.com/balerter/balerter/internal/mock"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/upload_storage/provider/s3"
	"go.uber.org/zap"
)

type Provider interface {
}

type Manager struct {
	logger *zap.Logger

	modules map[string]modules.ModuleTest
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:  logger,
		modules: make(map[string]modules.ModuleTest),
	}

	return m
}

func (m *Manager) Init(cfg *upload.Upload) error {
	if cfg == nil {
		return nil
	}
	for _, c := range cfg.S3 {
		mod := moduleMock.New(s3.ModuleName(c.Name), s3.Methods(), m.logger)

		m.modules[mod.Name()] = mod
	}

	return nil
}

func (m *Manager) Get() []modules.ModuleTest {
	mm := make([]modules.ModuleTest, 0)

	for _, module := range m.modules {
		mm = append(mm, module)
	}

	return mm
}

func (m *Manager) Result() ([]modules.TestResult, error) {
	var result []modules.TestResult
	for _, m := range m.modules {
		results, err := m.Result()
		if err != nil {
			return nil, err
		}
		for _, r := range results {
			r.ModuleName = "storage." + r.ModuleName
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *Manager) Clean() {
	for _, m := range m.modules {
		m.Clean()
	}
}
