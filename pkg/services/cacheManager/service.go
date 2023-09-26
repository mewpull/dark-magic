package cacheManager

import (
	"github.com/gravestench/runtime"
	"github.com/rs/zerolog"

	"github.com/gravestench/dark-magic/pkg/cache"
)

type cacheLookup = map[string]any

type Service struct {
	runtime       runtime.Runtime
	logger        *zerolog.Logger
	boundServices cacheLookup
}

func (s *Service) Init(rt runtime.Runtime) {
	s.runtime = rt

	s.FlushAllCaches()
}

func (s *Service) Name() string {
	return "Cache Manager"
}

func (s *Service) BindLogger(logger *zerolog.Logger) {
	s.logger = logger
}

func (s *Service) Logger() *zerolog.Logger {
	return s.logger
}

func (s *Service) FlushAllCaches() {
	s.boundServices = make(cacheLookup)

	for _, service := range s.runtime.Services() {
		s.tryToFlushCacheForService(service)
	}
}

func (s *Service) tryToFlushAllCaches(rt runtime.Runtime) {
	for _, service := range rt.Services() {
		s.tryToFlushCacheForService(service)
	}
}

func (s *Service) tryToFlushCacheForService(service runtime.Service) {
	if _, exists := s.boundServices[service.Name()]; exists {
		return
	}

	candidate, ok := service.(HasCache)
	if !ok {
		return
	}

	s.logger.Info().Msgf("flushing cache for service: %v", service.Name())

	candidate.FlushCache(cache.New(candidate.CacheBudget()))

	s.boundServices[service.Name()] = nil
}
