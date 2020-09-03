// Copyright 2020 Google Inc. All Rights Reserved.
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

package machine

import (
	v1 "github.com/google/cadvisor/info/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetVmStatPerNumaGetWorkingset(t *testing.T) {
	numaVmStatFiles = "testdata/node[0-9]*/vmstat"
	re := "workingset.*"
	result, err := GetVmStatPerNuma(&re)
	expected := []v1.VmStatNuma{
		{
			Node: "0",
			Stats: map[string]int{
				"workingset_activate":    0,
				"workingset_nodereclaim": 0,
				"workingset_nodes":       0,
				"workingset_refault":     0,
				"workingset_restore":     0,
			}},
		{
			Node: "11",
			Stats: map[string]int{
				"workingset_activate":    1,
				"workingset_nodereclaim": 1,
				"workingset_nodes":       1,
				"workingset_refault":     1,
				"workingset_restore":     1,
			}}}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
func TestGetVmStatPerNumaGetAll(t *testing.T) {
	numaVmStatFiles = "testdata/node[0-9]*/vmstat"
	re := ".*"
	result, err := GetVmStatPerNuma(&re)
	expected := []v1.VmStatNuma{
		{
			Node: "0",
			Stats: map[string]int{
				"nr_active_anon":                1192699,
				"nr_active_file":                655116,
				"nr_anon_pages":                 1251468,
				"nr_anon_transparent_hugepages": 14,
				"nr_bounce":                     0,
				"nr_dirtied":                    1772217,
				"nr_dirty":                      175,
				"nr_file_hugepages":             0,
				"nr_file_pages":                 2092573,
				"nr_file_pmdmapped":             0,
				"nr_foll_pin_acquired":          0,
				"nr_foll_pin_released":          0,
				"nr_free_cma":                   0,
				"nr_free_pages":                 4596753,
				"nr_inactive_anon":              74864,
				"nr_inactive_file":              1203848,
				"nr_isolated_anon":              0,
				"nr_isolated_file":              0,
				"nr_kernel_misc_reclaimable":    0,
				"nr_kernel_stack":               25888,
				"nr_mapped":                     492913,
				"nr_mlock":                      12,
				"nr_page_table_pages":           15699,
				"nr_shmem":                      292950,
				"nr_shmem_hugepages":            0,
				"nr_shmem_pmdmapped":            0,
				"nr_slab_reclaimable":           130991,
				"nr_slab_unreclaimable":         70818,
				"nr_unevictable":                217508,
				"nr_unstable":                   0,
				"nr_vmscan_immediate_reclaim":   0,
				"nr_vmscan_write":               0,
				"nr_writeback":                  0,
				"nr_writeback_temp":             0,
				"nr_written":                    1641298,
				"nr_zone_active_anon":           1192699,
				"nr_zone_active_file":           655116,
				"nr_zone_inactive_anon":         74864,
				"nr_zone_inactive_file":         1203848,
				"nr_zone_unevictable":           217508,
				"nr_zone_write_pending":         175,
				"nr_zspages":                    0,
				"numa_foreign":                  0,
				"numa_hit":                      83241101,
				"numa_interleave":               23066,
				"numa_local":                    83241101,
				"numa_miss":                     0,
				"numa_other":                    0,
				"workingset_activate":           0,
				"workingset_nodereclaim":        0,
				"workingset_nodes":              0,
				"workingset_refault":            0,
				"workingset_restore":            0,
			}},
		{
			Node: "11",
			Stats: map[string]int{
				"nr_active_anon":                1192699,
				"nr_active_file":                655116,
				"nr_anon_pages":                 1251468,
				"nr_anon_transparent_hugepages": 14,
				"nr_bounce":                     0,
				"nr_dirtied":                    1772217,
				"nr_dirty":                      175,
				"nr_file_hugepages":             0,
				"nr_file_pages":                 2092573,
				"nr_file_pmdmapped":             0,
				"nr_foll_pin_acquired":          0,
				"nr_foll_pin_released":          0,
				"nr_free_cma":                   0,
				"nr_free_pages":                 4596753,
				"nr_inactive_anon":              74864,
				"nr_inactive_file":              1203848,
				"nr_isolated_anon":              0,
				"nr_isolated_file":              0,
				"nr_kernel_misc_reclaimable":    0,
				"nr_kernel_stack":               25888,
				"nr_mapped":                     492913,
				"nr_mlock":                      12,
				"nr_page_table_pages":           15699,
				"nr_shmem":                      292950,
				"nr_shmem_hugepages":            0,
				"nr_shmem_pmdmapped":            0,
				"nr_slab_reclaimable":           130991,
				"nr_slab_unreclaimable":         70818,
				"nr_unevictable":                217508,
				"nr_unstable":                   0,
				"nr_vmscan_immediate_reclaim":   0,
				"nr_vmscan_write":               0,
				"nr_writeback":                  0,
				"nr_writeback_temp":             0,
				"nr_written":                    1641298,
				"nr_zone_active_anon":           1192699,
				"nr_zone_active_file":           655116,
				"nr_zone_inactive_anon":         74864,
				"nr_zone_inactive_file":         1203848,
				"nr_zone_unevictable":           217508,
				"nr_zone_write_pending":         175,
				"nr_zspages":                    0,
				"numa_foreign":                  0,
				"numa_hit":                      83241101,
				"numa_interleave":               23066,
				"numa_local":                    83241101,
				"numa_miss":                     0,
				"numa_other":                    0,
				"workingset_activate":           1,
				"workingset_nodereclaim":        1,
				"workingset_nodes":              1,
				"workingset_refault":            1,
				"workingset_restore":            1,
			}}}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
