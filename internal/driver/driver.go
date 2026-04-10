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

// Package driver 驱动相关功能
package driver

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// markBaseUpload 上传标记基址
const markBaseUpload = 0x0100

// markBaseDownload 下载标记基址
const markBaseDownload = 0x1000

// maxDevices 最大设备数
const maxDevices = markBaseDownload - markBaseUpload

// burstBytes 计算突发字节
// 入参: rateKbps 速率
// 返回: burst 突发字节
func burstBytes(rateKbps int) int {
	const burst = 32 * 1024
	_ = rateKbps
	return burst
}

// nftScript 执行脚本
// 入参: script 脚本内容
// 返回: error 错误信息
func nftScript(script string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(script)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nft script: %v (%s)", err, bytes.TrimSpace(stderr.Bytes()))
	}
	return nil
}

// Clean 清除规则
// 返回: error 错误信息
func Clean() error {
	return nftScript("table inet nftrl\ndelete table inet nftrl\n")
}

// DeviceRule 设备规则
type DeviceRule struct {
	MAC          string
	UploadKbps   int
	DownloadKbps int
}

// Apply 部署规则
// 入参: devices 设备列表
// 返回: error 错误信息
func Apply(devices []DeviceRule) error {
	if len(devices) > maxDevices {
		return fmt.Errorf("device count %d exceeds max %d", len(devices), maxDevices)
	}
	var sb strings.Builder
	sb.WriteString("table inet nftrl\n")
	sb.WriteString("delete table inet nftrl\n")
	sb.WriteString("table inet nftrl {\n")
	sb.WriteString("  chain prerouting {\n")
	sb.WriteString("    type filter hook prerouting priority -150; policy accept;\n")
	for i, d := range devices {
		upMark := markBaseUpload + i
		downMark := markBaseDownload + i
		fmt.Fprintf(&sb, "    ether saddr %s meta mark set 0x%04x ct mark set 0x%04x\n",
			d.MAC, upMark, downMark)
	}
	sb.WriteString("  }\n")
	sb.WriteString("  chain forward {\n")
	sb.WriteString("    type filter hook forward priority filter + 10; policy accept;\n")
	for i, d := range devices {
		upMark := markBaseUpload + i
		downMark := markBaseDownload + i
		fmt.Fprintf(&sb, "    meta mark != 0x%04x ct mark 0x%04x meta mark set 0x%04x\n",
			upMark, downMark, downMark)
		if d.UploadKbps > 0 {
			fmt.Fprintf(&sb,
				"    meta mark 0x%04x limit rate over %d bytes/second burst %d bytes counter drop\n",
				upMark, d.UploadKbps*125, burstBytes(d.UploadKbps))
		}
		if d.DownloadKbps > 0 {
			fmt.Fprintf(&sb,
				"    meta mark 0x%04x limit rate over %d bytes/second burst %d bytes counter drop\n",
				downMark, d.DownloadKbps*125, burstBytes(d.DownloadKbps))
		}
	}
	sb.WriteString("  }\n")
	sb.WriteString("}\n")
	return nftScript(sb.String())
}
