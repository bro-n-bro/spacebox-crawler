package modules

import (
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

type (
	ModuleLoader struct {
		logger  *zerolog.Logger
		modules []types.Module
	}
)

func NewModuleLoader() *ModuleLoader {
	return &ModuleLoader{
		modules: make([]types.Module, 0),
	}
}

func (m *ModuleLoader) Len() int { return len(m.modules) }

func (m *ModuleLoader) WithLogger(log *zerolog.Logger) *ModuleLoader {
	m.logger = log
	return m
}

func (m *ModuleLoader) Add(mod types.Module) {
	m.modules = append(m.modules, mod)
	if m.logger != nil {
		m.logger.Info().Msgf("module %s registered", mod.Name())
	}
}

func (m *ModuleLoader) Build() []types.Module {
	if m.logger != nil {
		m.logger.Info().Int("count", m.Len()).Msgf("all modules registered")
	}

	return m.modules
}
