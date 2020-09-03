// Copyright 2014 Google Inc. All Rights Reserved.
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

package sysinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/google/cadvisor/utils/sysfs/fakesysfs"
	"github.com/stretchr/testify/assert"
)

func TestGetHugePagesInfo(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
		&fakesysfs.FileInfo{EntryName: "hugepages-1048576kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages":    "1",
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	hugePagesInfo, err := GetHugePagesInfo(&fakeSys, "/fakeSysfs/devices/system/node/node0/hugepages/")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(hugePagesInfo))
}

func TestGetHugePagesInfoWithHugePagesDirectory(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	hugePagesInfo, err := GetHugePagesInfo(&fakeSys, "/fakeSysfs/devices/system/node/node0/hugepages/")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(hugePagesInfo))
}

func TestGetHugePagesInfoWithWrongDirName(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-abckB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePagesInfo, err := GetHugePagesInfo(&fakeSys, "/fakeSysfs/devices/system/node/node0/hugepages/")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hugePagesInfo))
}

func TestGetHugePagesInfoWithReadingNrHugePagesError(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
		&fakesysfs.FileInfo{EntryName: "hugepages-1048576kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages":    "1",
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, fmt.Errorf("Error in reading nr_hugepages"))

	hugePagesInfo, err := GetHugePagesInfo(&fakeSys, "/fakeSysfs/devices/system/node/node0/hugepages/")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hugePagesInfo))
}

func TestGetHugePagesInfoWithWrongNrHugePageValue(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
		&fakesysfs.FileInfo{EntryName: "hugepages-1048576kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages":    "*****",
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-1048576kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	hugePagesInfo, err := GetHugePagesInfo(&fakeSys, "/fakeSysfs/devices/system/node/node0/hugepages/")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hugePagesInfo))
}

func TestGetNodesInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	c := sysfs.CacheInfo{
		Size:  32 * 1024,
		Type:  "unified",
		Level: 3,
		Cpus:  2,
	}
	fakeSys.SetCacheInfo(c)

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
		"/fakeSysfs/devices/system/node/node1",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
		"/fakeSysfs/devices/system/node/node1": {
			"/fakeSysfs/devices/system/node/node0/cpu2",
			"/fakeSysfs/devices/system/node/node0/cpu3",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetCoreThreads(coreThread, nil)

	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages": "1",
		"/fakeSysfs/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	physicalPackageIDs := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, nil)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, 4, cores)

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)
	expectedNodes := `
	[
      {
        "node_id": 0,
        "memory": 33604804608,
        "hugepages": [
          {
            "page_size": 2048,
            "num_pages": 1
          }
        ],
        "cores": [
          {
            "core_id": 0,
            "thread_ids": [
              0,
              1
            ],
            "caches": null,
	    "socket_id": 0
          }
        ],
        "caches": [
          {
            "size": 32768,
            "type": "unified",
            "level": 3
          }
        ]
      },
      {
        "node_id": 1,
        "memory": 33604804608,
        "hugepages": [
          {
            "page_size": 2048,
            "num_pages": 1
          }
        ],
        "cores": [
          {
            "core_id": 1,
            "thread_ids": [
              2,
              3
            ],
            "caches": null,
	    "socket_id": 1
          }
        ],
        "caches": [
          {
            "size": 32768,
            "type": "unified",
            "level": 3
          }
        ]
      }
    ]
    `
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesInfoWithOfflineCPUs(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	c := sysfs.CacheInfo{
		Size:  32 * 1024,
		Type:  "unified",
		Level: 3,
		Cpus:  1,
	}
	fakeSys.SetCacheInfo(c)

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
		"/fakeSysfs/devices/system/node/node1",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
		"/fakeSysfs/devices/system/node/node1": {
			"/fakeSysfs/devices/system/node/node0/cpu2",
			"/fakeSysfs/devices/system/node/node0/cpu3",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetCoreThreads(coreThread, nil)
	fakeSys.SetOnlineCPUs(map[string]interface{}{
		"/fakeSysfs/devices/system/node/node0/cpu0": nil,
		"/fakeSysfs/devices/system/node/node0/cpu2": nil,
	})

	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages": "1",
		"/fakeSysfs/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	physicalPackageIDs := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, nil)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, 2, cores)

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)
	expectedNodes := `
	[
      {
        "node_id": 0,
        "memory": 33604804608,
        "hugepages": [
          {
            "page_size": 2048,
            "num_pages": 1
          }
        ],
        "cores": [
          {
            "core_id": 0,
            "thread_ids": [
              0
            ],
            "caches": null,
	    "socket_id": 0
          }
        ],
        "caches": [
          {
            "size": 32768,
            "type": "unified",
            "level": 3
          }
        ]
      },
      {
        "node_id": 1,
        "memory": 33604804608,
        "hugepages": [
          {
            "page_size": 2048,
            "num_pages": 1
          }
        ],
        "cores": [
          {
            "core_id": 1,
            "thread_ids": [
              2
            ],
            "caches": null,
	    "socket_id": 1
          }
        ],
        "caches": [
          {
            "size": 32768,
            "type": "unified",
            "level": 3
          }
        ]
      }
    ]
    `
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesWithoutMemoryInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	c := sysfs.CacheInfo{
		Size:  32 * 1024,
		Type:  "unified",
		Level: 3,
		Cpus:  2,
	}
	fakeSys.SetCacheInfo(c)

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
		"/fakeSysfs/devices/system/node/node1",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
		"/fakeSysfs/devices/system/node/node1": {
			"/fakeSysfs/devices/system/node/node0/cpu2",
			"/fakeSysfs/devices/system/node/node0/cpu3",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetCoreThreads(coreThread, nil)

	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages": "1",
		"/fakeSysfs/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.NotNil(t, err)
	assert.Equal(t, []info.Node([]info.Node(nil)), nodes)
	assert.Equal(t, 0, cores)
}

func TestGetNodesInfoWithoutCacheInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
		"/fakeSysfs/devices/system/node/node1",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
		"/fakeSysfs/devices/system/node/node1": {
			"/fakeSysfs/devices/system/node/node0/cpu2",
			"/fakeSysfs/devices/system/node/node0/cpu3",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetCoreThreads(coreThread, nil)

	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages": "1",
		"/fakeSysfs/devices/system/node/node1/hugepages/hugepages-2048kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	physicalPackageIDs := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, nil)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, 4, cores)

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)
	expectedNodes := `
	[
      {
        "node_id": 0,
        "memory": 33604804608,
        "hugepages": [
          {
            "page_size": 2048,
            "num_pages": 1
          }
        ],
        "cores": [
	  {
            "core_id": 0,
            "thread_ids": [
              0,
              1
            ],
            "caches": null,
	    "socket_id": 0
          }
        ],
        "caches": null
      },
      {
        "node_id": 1,
        "memory": 33604804608,
        "hugepages": [
          {
            "page_size": 2048,
            "num_pages": 1
          }
        ],
        "cores": [
          {
            "core_id": 1,
            "thread_ids": [
              2,
              3
            ],
            "caches": null,
	    "socket_id": 1
          }
        ],
        "caches": null
      }
    ]`
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesInfoWithoutHugePagesInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	c := sysfs.CacheInfo{
		Size:  32 * 1024,
		Type:  "unified",
		Level: 2,
		Cpus:  2,
	}
	fakeSys.SetCacheInfo(c)

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
		"/fakeSysfs/devices/system/node/node1",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
		"/fakeSysfs/devices/system/node/node1": {
			"/fakeSysfs/devices/system/node/node0/cpu2",
			"/fakeSysfs/devices/system/node/node0/cpu3",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetCoreThreads(coreThread, nil)

	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	physicalPackageIDs := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
		"/fakeSysfs/devices/system/node/node0/cpu2": "1",
		"/fakeSysfs/devices/system/node/node0/cpu3": "1",
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, nil)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, 4, cores)

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)
	expectedNodes := `
	[
      {
        "node_id": 0,
        "memory": 33604804608,
        "hugepages": null,
        "cores": [
          {
            "core_id": 0,
            "thread_ids": [
              0,
              1
            ],
            "caches": [
              {
                "size": 32768,
                "type": "unified",
                "level": 2
              }
            ],
	    "socket_id": 0
          }
        ],
        "caches": null
      },
      {
        "node_id": 1,
        "memory": 33604804608,
        "hugepages": null,
        "cores": [
          {
            "core_id": 1,
            "thread_ids": [
              2,
              3
            ],
            "caches": [
              {
                "size": 32768,
                "type": "unified",
                "level": 2
              }
            ],
	    "socket_id": 1
          }
        ],
        "caches": null
      }
    ]`
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesInfoWithoutNodes(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}

	c := sysfs.CacheInfo{
		Size:  32 * 1024,
		Type:  "unified",
		Level: 1,
		Cpus:  2,
	}
	fakeSys.SetCacheInfo(c)

	nodesPaths := []string{}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		cpusPath: {
			cpusPath + "/cpu0",
			cpusPath + "/cpu1",
			cpusPath + "/cpu2",
			cpusPath + "/cpu3",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		cpusPath + "/cpu0": "0",
		cpusPath + "/cpu1": "0",
		cpusPath + "/cpu2": "1",
		cpusPath + "/cpu3": "1",
	}
	fakeSys.SetCoreThreads(coreThread, nil)

	physicalPackageIDs := map[string]string{
		"/sys/devices/system/cpu/cpu0": "0",
		"/sys/devices/system/cpu/cpu1": "0",
		"/sys/devices/system/cpu/cpu2": "1",
		"/sys/devices/system/cpu/cpu3": "1",
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, nil)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, 4, cores)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id < nodes[j].Id
	})

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)

	expectedNodes := `[
		{
			"node_id":0,
			"memory":0,
			"hugepages":null,
			"cores":[
			   {
				  "core_id":0,
				  "thread_ids":[
					 0,
					 1
				  ],
				  "caches":[
					 {
						"size":32768,
						"type":"unified",
						"level":1
					 }
				  ],
				  "socket_id": 0
			   }
			],
			"caches":null
		 },
		 {
			"node_id":1,
			"memory":0,
			"hugepages":null,
			"cores":[
			   {
				  "core_id":1,
				  "thread_ids":[
					 2,
					 3
				  ],
				  "caches":[
					 {
						"size":32768,
						"type":"unified",
						"level":1
					 }
				  ],
				  "socket_id": 1
			   }
			],
			"caches":null
		 }
	]`
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesInfoWithoutNodesWhenPhysicalPackageIDMissingForOneCPU(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}

	nodesPaths := []string{}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		cpusPath: {
			cpusPath + "/cpu0",
			cpusPath + "/cpu1",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		cpusPath + "/cpu0": "0",
		cpusPath + "/cpu1": "0",
	}

	coreThreadErrors := map[string]error{
		cpusPath + "/cpu0": nil,
		cpusPath + "/cpu1": nil,
	}
	fakeSys.SetCoreThreads(coreThread, coreThreadErrors)

	physicalPackageIDs := map[string]string{
		cpusPath + "/cpu0": "0",
		cpusPath + "/cpu1": "0",
	}

	physicalPackageIDErrors := map[string]error{
		cpusPath + "/cpu0": nil,
		cpusPath + "/cpu1": os.ErrNotExist,
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, physicalPackageIDErrors)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, 2, cores)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id < nodes[j].Id
	})

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)

	fmt.Println(string(nodesJSON))

	expectedNodes := `[
		{
			"node_id":0,
			"memory":0,
			"hugepages":null,
			"cores":[
			   {
				  "core_id":0,
				  "thread_ids":[
					 0
				  ],
				  "caches": null,
				  "socket_id": 0
			   }
			],
			"caches":null
		}
	]`
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesInfoWithoutNodesWhenPhysicalPackageIDMissing(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}

	nodesPaths := []string{}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		cpusPath: {
			cpusPath + "/cpu0",
			cpusPath + "/cpu1",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		cpusPath + "/cpu0": "0",
		cpusPath + "/cpu1": "0",
	}

	coreThreadErrors := map[string]error{
		cpusPath + "/cpu0": nil,
		cpusPath + "/cpu1": nil,
	}
	fakeSys.SetCoreThreads(coreThread, coreThreadErrors)

	physicalPackageIDs := map[string]string{
		cpusPath + "/cpu0": "0",
		cpusPath + "/cpu1": "0",
	}

	physicalPackageIDErrors := map[string]error{
		cpusPath + "/cpu0": os.ErrNotExist,
		cpusPath + "/cpu1": os.ErrNotExist,
	}
	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, physicalPackageIDErrors)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(nodes))
	assert.Equal(t, 2, cores)
}

func TestGetNodesWhenTopologyDirMissingForOneCPU(t *testing.T) {
	/*
		Unit test for case in which:
		- there are two cpus (cpu0 and cpu1) in /sys/devices/system/node/node0/ and /sys/devices/system/cpu
		- topology directory is missing for cpu1 but it exists for cpu0
	*/
	fakeSys := &fakesysfs.FakeSysFs{}

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
	}

	coreThreadErrors := map[string]error{
		"/fakeSysfs/devices/system/node/node0/cpu0": nil,
		"/fakeSysfs/devices/system/node/node0/cpu1": os.ErrNotExist,
	}
	fakeSys.SetCoreThreads(coreThread, coreThreadErrors)

	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	physicalPackageIDs := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
	}

	physicalPackageIDErrors := map[string]error{
		"/fakeSysfs/devices/system/node/node0/cpu0": nil,
		"/fakeSysfs/devices/system/node/node0/cpu1": os.ErrNotExist,
	}

	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, physicalPackageIDErrors)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, 1, cores)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id < nodes[j].Id
	})

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)

	expectedNodes := `[
		{
		   "node_id":0,
		   "memory":33604804608,
		   "hugepages":[
			  {
				 "page_size":2048,
				 "num_pages":1
			  }
		   ],
		   "cores":[
			  {
				 "core_id":0,
				 "thread_ids":[
					0
				 ],
				 "caches":null,
				 "socket_id":0
			  }
		   ],
		   "caches": null
		}
	 ]`
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodesWhenPhysicalPackageIDMissingForOneCPU(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}

	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
	}
	fakeSys.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
	}
	fakeSys.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
	}

	coreThreadErrors := map[string]error{
		"/fakeSysfs/devices/system/node/node0/cpu0": nil,
		"/fakeSysfs/devices/system/node/node0/cpu1": nil,
	}
	fakeSys.SetCoreThreads(coreThread, coreThreadErrors)

	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	hugePages := []os.FileInfo{
		&fakesysfs.FileInfo{EntryName: "hugepages-2048kB"},
	}
	fakeSys.SetHugePages(hugePages, nil)

	hugePageNr := map[string]string{
		"/fakeSysfs/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages": "1",
	}
	fakeSys.SetHugePagesNr(hugePageNr, nil)

	physicalPackageIDs := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
	}

	physicalPackageIDErrors := map[string]error{
		"/fakeSysfs/devices/system/node/node0/cpu0": nil,
		"/fakeSysfs/devices/system/node/node0/cpu1": os.ErrNotExist,
	}

	fakeSys.SetPhysicalPackageIDs(physicalPackageIDs, physicalPackageIDErrors)

	nodes, cores, err := GetNodesInfo(fakeSys)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, 2, cores)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id < nodes[j].Id
	})

	nodesJSON, err := json.Marshal(nodes)
	assert.Nil(t, err)

	expectedNodes := `[
		{
		   "node_id":0,
		   "memory":33604804608,
		   "hugepages":[
			  {
				 "page_size":2048,
				 "num_pages":1
			  }
		   ],
		   "cores":[
			  {
				 "core_id":0,
				 "thread_ids":[
					0, 1
				 ],
				 "caches":null,
				 "socket_id":0
			  }
		   ],
		   "caches": null
		}
	 ]`
	assert.JSONEq(t, expectedNodes, string(nodesJSON))
}

func TestGetNodeMemInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	memTotal := "MemTotal:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	mem, err := getNodeMemInfo(fakeSys, "/fakeSysfs/devices/system/node/node0")
	assert.Nil(t, err)
	assert.Equal(t, uint64(32817192*1024), mem)
}

func TestGetNodeMemInfoWithMissingMemTotaInMemInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	memTotal := "MemXXX:       32817192 kB"
	fakeSys.SetMemory(memTotal, nil)

	mem, err := getNodeMemInfo(fakeSys, "/fakeSysfs/devices/system/node/node0")
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), mem)
}

func TestGetNodeMemInfoWhenMemInfoMissing(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	memTotal := ""
	fakeSys.SetMemory(memTotal, fmt.Errorf("Cannot read meminfo file"))

	mem, err := getNodeMemInfo(fakeSys, "/fakeSysfs/devices/system/node/node0")
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), mem)
}

func TestGetCoresInfoWhenCoreIDIsNotDigit(t *testing.T) {
	sysFs := &fakesysfs.FakeSysFs{}
	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
	}
	sysFs.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
		},
	}
	sysFs.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "abc",
	}
	sysFs.SetCoreThreads(coreThread, nil)

	cores, err := getCoresInfo(sysFs, []string{"/fakeSysfs/devices/system/node/node0/cpu0"})
	assert.NotNil(t, err)
	assert.Equal(t, []info.Core(nil), cores)
}

func TestGetCoresInfoWithOnlineOfflineFile(t *testing.T) {
	sysFs := &fakesysfs.FakeSysFs{}
	nodesPaths := []string{
		"/fakeSysfs/devices/system/node/node0",
	}
	sysFs.SetNodesPaths(nodesPaths, nil)

	cpusPaths := map[string][]string{
		"/fakeSysfs/devices/system/node/node0": {
			"/fakeSysfs/devices/system/node/node0/cpu0",
			"/fakeSysfs/devices/system/node/node0/cpu1",
		},
	}
	sysFs.SetCPUsPaths(cpusPaths, nil)

	coreThread := map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
	}
	sysFs.SetCoreThreads(coreThread, nil)
	sysFs.SetOnlineCPUs(map[string]interface{}{"/fakeSysfs/devices/system/node/node0/cpu0": nil})
	sysFs.SetPhysicalPackageIDs(map[string]string{
		"/fakeSysfs/devices/system/node/node0/cpu0": "0",
		"/fakeSysfs/devices/system/node/node0/cpu1": "0",
	}, nil)

	cores, err := getCoresInfo(
		sysFs,
		[]string{"/fakeSysfs/devices/system/node/node0/cpu0", "/fakeSysfs/devices/system/node/node0/cpu1"},
	)
	assert.NoError(t, err)
	expected := []info.Core{
		{
			Id:       0,
			Threads:  []int{0},
			Caches:   nil,
			SocketID: 0,
		},
	}
	assert.Equal(t, expected, cores)
}

func TestGetBlockDeviceInfo(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	disks, err := GetBlockDeviceInfo(&fakeSys)
	if err != nil {
		t.Errorf("expected call to GetBlockDeviceInfo() to succeed. Failed with %s", err)
	}
	if len(disks) != 1 {
		t.Errorf("expected to get one disk entry. Got %d", len(disks))
	}
	key := "8:0"
	disk, ok := disks[key]
	if !ok {
		t.Fatalf("expected key 8:0 to exist in the disk map.")
	}
	if disk.Name != "sda" {
		t.Errorf("expected to get disk named sda. Got %q", disk.Name)
	}
	size := uint64(1234567 * 512)
	if disk.Size != size {
		t.Errorf("expected to get disk size of %d. Got %d", size, disk.Size)
	}
	if disk.Scheduler != "cfq" {
		t.Errorf("expected to get scheduler type of cfq. Got %q", disk.Scheduler)
	}
}

func TestGetNetworkDevices(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	fakeSys.SetEntryName("eth0")
	devs, err := GetNetworkDevices(&fakeSys)
	if err != nil {
		t.Errorf("expected call to GetNetworkDevices() to succeed. Failed with %s", err)
	}
	if len(devs) != 1 {
		t.Errorf("expected to get one network device. Got %d", len(devs))
	}
	eth := devs[0]
	if eth.Name != "eth0" {
		t.Errorf("expected to find device with name eth0. Found name %q", eth.Name)
	}
	if eth.Mtu != 1024 {
		t.Errorf("expected mtu to be set to 1024. Found %d", eth.Mtu)
	}
	if eth.Speed != 1000 {
		t.Errorf("expected device speed to be set to 1000. Found %d", eth.Speed)
	}
	if eth.MacAddress != "42:01:02:03:04:f4" {
		t.Errorf("expected mac address to be '42:01:02:03:04:f4'. Found %q", eth.MacAddress)
	}
}

func TestIgnoredNetworkDevices(t *testing.T) {
	fakeSys := fakesysfs.FakeSysFs{}
	ignoredDevices := []string{"veth1234", "lo", "docker0"}
	for _, name := range ignoredDevices {
		fakeSys.SetEntryName(name)
		devs, err := GetNetworkDevices(&fakeSys)
		if err != nil {
			t.Errorf("expected call to GetNetworkDevices() to succeed. Failed with %s", err)
		}
		if len(devs) != 0 {
			t.Errorf("expected dev %s to be ignored, but got info %+v", name, devs)
		}
	}
}

func TestGetCacheInfo(t *testing.T) {
	fakeSys := &fakesysfs.FakeSysFs{}
	cacheInfo := sysfs.CacheInfo{
		Size:  1024,
		Type:  "Data",
		Level: 3,
		Cpus:  16,
	}
	fakeSys.SetCacheInfo(cacheInfo)
	caches, err := GetCacheInfo(fakeSys, 0)
	if err != nil {
		t.Errorf("expected call to GetCacheInfo() to succeed. Failed with %s", err)
	}
	if len(caches) != 1 {
		t.Errorf("expected to get one cache. Got %d", len(caches))
	}
	if caches[0] != cacheInfo {
		t.Errorf("expected to find cacheinfo %+v. Got %+v", cacheInfo, caches[0])
	}
}

func TestGetNetworkStats(t *testing.T) {
	expectedStats := info.InterfaceStats{
		Name:      "eth0",
		RxBytes:   1024,
		RxPackets: 1024,
		RxErrors:  1024,
		RxDropped: 1024,
		TxBytes:   1024,
		TxPackets: 1024,
		TxErrors:  1024,
		TxDropped: 1024,
	}
	fakeSys := &fakesysfs.FakeSysFs{}
	netStats, err := getNetworkStats("eth0", fakeSys)
	if err != nil {
		t.Errorf("call to getNetworkStats() failed with %s", err)
	}
	if expectedStats != netStats {
		t.Errorf("expected to get stats %+v, got %+v", expectedStats, netStats)
	}
}

func TestGetSocketFromCPU(t *testing.T) {
	topology := []info.Node{
		{
			Id:        0,
			Memory:    0,
			HugePages: nil,
			Cores: []info.Core{
				{
					Id:       0,
					Threads:  []int{0, 1},
					Caches:   nil,
					SocketID: 0,
				},
				{
					Id:       1,
					Threads:  []int{2, 3},
					Caches:   nil,
					SocketID: 0,
				},
			},
			Caches: nil,
		},
		{
			Id:        1,
			Memory:    0,
			HugePages: nil,
			Cores: []info.Core{
				{
					Id:       0,
					Threads:  []int{4, 5},
					Caches:   nil,
					SocketID: 1,
				},
				{
					Id:       1,
					Threads:  []int{6, 7},
					Caches:   nil,
					SocketID: 1,
				},
			},
			Caches: nil,
		},
	}
	socket := GetSocketFromCPU(topology, 6)
	assert.Equal(t, socket, 1)

	// Check if return "-1" when there is no data about passed CPU.
	socket = GetSocketFromCPU(topology, 8)
	assert.Equal(t, socket, -1)
}


func TestGetVMStatsGetNuma(t *testing.T) {
	vmstatFile := "testdata/vmstat_data"
	re := "numa.*"
	result, err := GetVMStats(&re, vmstatFile)
	expected := map[string]int{
		"numa_foreign":           0,
		"numa_hint_faults":       0,
		"numa_hint_faults_local": 0,
		"numa_hit":               120883140,
		"numa_huge_pte_updates":  0,
		"numa_interleave":        23061,
		"numa_local":             120883140,
		"numa_miss":              0,
		"numa_other":             0,
		"numa_pages_migrated":    0,
		"numa_pte_updates":       0,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)

}
func TestGetVMStatsGetFaults(t *testing.T) {
	vmstatFile := "testdata/vmstat_data"
	re := ".*(faults|failed).*"
	result, err := GetVMStats(&re, vmstatFile)
	expected := map[string]int{
		"numa_hint_faults":           0,
		"numa_hint_faults_local":     0,
		"thp_collapse_alloc_failed":  0,
		"thp_split_page_failed":      0,
		"thp_zero_page_alloc_failed": 0,
		"zone_reclaim_failed":        0,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)

}
func TestGetVMStatsGetAll(t *testing.T) {
	vmstatFile := "testdata/vmstat_data"
	re := ".*"
	result, err := GetVMStats(&re, vmstatFile)
	expected := map[string]int{"allocstall_dma": 0,
		"allocstall_dma32":               0,
		"allocstall_movable":             0,
		"allocstall_normal":              0,
		"balloon_deflate":                0,
		"balloon_inflate":                0,
		"balloon_migrate":                0,
		"compact_daemon_free_scanned":    0,
		"compact_daemon_migrate_scanned": 0,
		"compact_daemon_wake":            0,
		"compact_fail":                   0,
		"compact_free_scanned":           0,
		"compact_isolated":               0,
		"compact_migrate_scanned":        0,
		"compact_stall":                  0,
		"compact_success":                0,
		"drop_pagecache":                 0,
		"drop_slab":                      0,
		"htlb_buddy_alloc_fail":          0,
		"htlb_buddy_alloc_success":       0,
		"kswapd_high_wmark_hit_quickly":  0,
		"kswapd_inodesteal":              0,
		"kswapd_low_wmark_hit_quickly":   0,
		"nr_active_anon":                 1164948,
		"nr_active_file":                 567385,
		"nr_anon_pages":                  1262075,
		"nr_anon_transparent_hugepages":  19,
		"nr_bounce":                      0,
		"nr_dirtied":                     1329549,
		"nr_dirty":                       1145,
		"nr_dirty_background_threshold":  644547,
		"nr_dirty_threshold":             1290670,
		"nr_file_hugepages":              0,
		"nr_file_pages":                  1721185,
		"nr_file_pmdmapped":              0,
		"nr_foll_pin_acquired":           0,
		"nr_foll_pin_released":           0,
		"nr_free_cma":                    0,
		"nr_free_pages":                  4981979,
		"nr_inactive_anon":               56987,
		"nr_inactive_file":               972664,
		"nr_isolated_anon":               0,
		"nr_isolated_file":               0,
		"nr_kernel_misc_reclaimable":     0,
		"nr_kernel_stack":                23952,
		"nr_mapped":                      517738,
		"nr_mlock":                       12,
		"nr_page_table_pages":            14081,
		"nr_shmem":                       278886,
		"nr_shmem_hugepages":             0,
		"nr_shmem_pmdmapped":             0,
		"nr_slab_reclaimable":            111629,
		"nr_slab_unreclaimable":          69328,
		"nr_unevictable":                 221278,
		"nr_unstable":                    0,
		"nr_vmscan_immediate_reclaim":    0,
		"nr_vmscan_write":                0,
		"nr_writeback":                   0,
		"nr_writeback_temp":              0,
		"nr_written":                     1325274,
		"nr_zone_active_anon":            1164948,
		"nr_zone_active_file":            567385,
		"nr_zone_inactive_anon":          56987,
		"nr_zone_inactive_file":          972664,
		"nr_zone_unevictable":            221278,
		"nr_zone_write_pending":          1139,
		"nr_zspages":                     0,
		"numa_foreign":                   0,
		"numa_hint_faults":               0,
		"numa_hint_faults_local":         0,
		"numa_hit":                       120883140,
		"numa_huge_pte_updates":          0,
		"numa_interleave":                23061,
		"numa_local":                     120883140,
		"numa_miss":                      0,
		"numa_other":                     0,
		"numa_pages_migrated":            0,
		"numa_pte_updates":               0,
		"oom_kill":                       0,
		"pageoutrun":                     0,
		"pgactivate":                     5149662,
		"pgalloc_dma":                    0,
		"pgalloc_dma32":                  1,
		"pgalloc_movable":                0,
		"pgalloc_normal":                 126730966,
		"pgdeactivate":                   0,
		"pgfault":                        108683933,
		"pgfree":                         131715713,
		"pginodesteal":                   0,
		"pglazyfree":                     828926,
		"pglazyfreed":                    0,
		"pgmajfault":                     20994,
		"pgmigrate_fail":                 0,
		"pgmigrate_success":              0,
		"pgpgin":                         4079261,
		"pgpgout":                        6255350,
		"pgrefill":                       0,
		"pgrotated":                      11,
		"pgscan_direct":                  0,
		"pgscan_direct_throttle":         0,
		"pgscan_kswapd":                  0,
		"pgskip_dma":                     0,
		"pgskip_dma32":                   0,
		"pgskip_movable":                 0,
		"pgskip_normal":                  0,
		"pgsteal_direct":                 0,
		"pgsteal_kswapd":                 0,
		"pswpin":                         0,
		"pswpout":                        0,
		"slabs_scanned":                  0,
		"swap_ra":                        0,
		"swap_ra_hit":                    0,
		"thp_collapse_alloc":             918,
		"thp_collapse_alloc_failed":      0,
		"thp_deferred_split_page":        891,
		"thp_fault_alloc":                450,
		"thp_fault_fallback":             0,
		"thp_fault_fallback_charge":      0,
		"thp_file_alloc":                 0,
		"thp_file_fallback":              0,
		"thp_file_fallback_charge":       0,
		"thp_file_mapped":                0,
		"thp_split_page":                 891,
		"thp_split_page_failed":          0,
		"thp_split_pmd":                  891,
		"thp_split_pud":                  0,
		"thp_swpout":                     0,
		"thp_swpout_fallback":            0,
		"thp_zero_page_alloc":            1,
		"thp_zero_page_alloc_failed":     0,
		"unevictable_pgs_cleared":        0,
		"unevictable_pgs_culled":         9542080,
		"unevictable_pgs_mlocked":        270870,
		"unevictable_pgs_munlocked":      270858,
		"unevictable_pgs_rescued":        9300587,
		"unevictable_pgs_scanned":        11752725,
		"unevictable_pgs_stranded":       0,
		"workingset_activate":            0,
		"workingset_nodereclaim":         0,
		"workingset_nodes":               0,
		"workingset_refault":             0,
		"workingset_restore":             0,
		"zone_reclaim_failed":            0,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
