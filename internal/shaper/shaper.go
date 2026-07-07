// Copyright 2026 肖其顿 (XIAO QI DUN)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package shaper 业务编排逻辑
package shaper

import (
	"github.com/xiaoqidun/nftrl/internal/config"
	"github.com/xiaoqidun/nftrl/internal/driver"
)

// Apply 部署规则
// 返回: error 错误信息
func Apply() error {
	cfg, err := config.Parse(config.Path)
	if err != nil {
		return err
	}
	if !cfg.Global.Enabled {
		return Clean()
	}
	var rules []driver.DeviceRule
	for _, d := range cfg.Devices {
		if !d.Enabled || d.MAC == "" {
			continue
		}
		rules = append(rules, driver.DeviceRule{
			MAC:          d.MAC,
			UploadKbps:   d.EgressLimit,
			DownloadKbps: d.IngressLimit,
		})
	}
	if len(rules) == 0 {
		return Clean()
	}
	return driver.Apply(rules)
}

// Clean 清理规则
// 返回: error 错误信息
func Clean() error {
	return driver.Clean()
}
