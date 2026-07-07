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

// Package config 配置文件解析
package config

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Path 配置文件路径
const Path = "/etc/config/nftrl"

// Config 配置信息
type Config struct {
	Global  Global
	Devices []Device
}

// Global 全局配置
type Global struct {
	Enabled bool
}

// Device 设备配置
type Device struct {
	Enabled      bool
	MAC          string
	EgressLimit  int
	IngressLimit int
	Comment      string
}

// Parse 解析配置文件
// 入参: path 配置文件路径
// 返回: cfg 配置信息, err 错误信息
func Parse(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()
	cfg := &Config{}
	scanner := bufio.NewScanner(file)
	var section string
	var device *Device
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "config":
			if len(parts) < 2 {
				continue
			}
			if device != nil {
				cfg.Devices = append(cfg.Devices, *device)
				device = nil
			}
			switch parts[1] {
			case "nftrl":
				section = "global"
			case "device":
				section = "device"
				device = &Device{}
			default:
				section = ""
			}
		case "option":
			if len(parts) < 3 {
				continue
			}
			key := parts[1]
			val := strings.TrimSpace(strings.Join(parts[2:], " "))
			val = strings.Trim(val, "'\"")
			switch section {
			case "global":
				if key == "enabled" {
					cfg.Global.Enabled = val == "1"
				}
			case "device":
				if device == nil {
					continue
				}
				switch key {
				case "enabled":
					device.Enabled = val == "1"
				case "mac":
					mac, err := net.ParseMAC(val)
					if err != nil {
						return nil, fmt.Errorf("invalid mac %q: %w", val, err)
					}
					if len(mac) != 6 {
						return nil, fmt.Errorf("invalid mac %q: expected 6 bytes", val)
					}
					device.MAC = mac.String()
				case "egress_limit":
					v, err := strconv.Atoi(val)
					if err != nil {
						return nil, fmt.Errorf("invalid egress_limit %q: %w", val, err)
					}
					if v < 0 {
						return nil, fmt.Errorf("egress_limit must be non-negative: %d", v)
					}
					device.EgressLimit = v
				case "ingress_limit":
					v, err := strconv.Atoi(val)
					if err != nil {
						return nil, fmt.Errorf("invalid ingress_limit %q: %w", val, err)
					}
					if v < 0 {
						return nil, fmt.Errorf("ingress_limit must be non-negative: %d", v)
					}
					device.IngressLimit = v
				case "comment":
					device.Comment = val
				}
			}
		}
	}
	if device != nil {
		cfg.Devices = append(cfg.Devices, *device)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan config: %w", err)
	}
	seen := make(map[string]bool)
	for _, c := range cfg.Devices {
		if !c.Enabled || c.MAC == "" {
			continue
		}
		if seen[c.MAC] {
			return nil, fmt.Errorf("duplicate mac: %s", c.MAC)
		}
		seen[c.MAC] = true
	}
	return cfg, nil
}
