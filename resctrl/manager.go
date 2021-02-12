// +build linux

// Copyright 2021 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Manager of resctrl for containers.
package resctrl

import (
	"errors"
	"time"

	"github.com/google/cadvisor/stats"
)

type Manager interface {
	Destroy()
	GetCollector(containerName string, getContainerPids func() ([]string, error)) (stats.Collector, error)
}

type manager struct {
	stats.NoopDestroy
	interval time.Duration
}

func (m *manager) GetCollector(containerName string, getContainerPids func() ([]string, error)) (stats.Collector, error) {
	collector := newCollector(containerName, getContainerPids, m.interval)
	err := collector.setup()
	if err != nil {
		return &stats.NoopCollector{}, err
	}

	return collector, nil
}

func NewManager(interval time.Duration, setup func() error) (Manager, error) {
	err := setup()
	if err != nil {
		return &NoopManager{}, err
	}

	if !isResctrlInitialized {
		return &NoopManager{}, errors.New("the resctrl isn't initialized")
	}
	if !(enabledCMT || enabledMBM) {
		return &NoopManager{}, errors.New("there are no monitoring features available")
	}

	return &manager{interval: interval}, nil
}

type NoopManager struct {
	stats.NoopDestroy
}

func (np *NoopManager) GetCollector(_ string, _ func() ([]string, error)) (stats.Collector, error) {
	return &stats.NoopCollector{}, nil
}
