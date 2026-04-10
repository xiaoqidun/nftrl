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

// Package main 一个基于 nftables 的零侵入 OpenWrt 限速工具
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xiaoqidun/nftrl/internal/shaper"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <command>\n\ncommands:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  apply    read config and deploy rules")
		fmt.Fprintln(os.Stderr, "  clean    remove nftrl rules")
		fmt.Fprintln(os.Stderr, "  license  print license information")
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	switch flag.Arg(0) {
	case "apply":
		if err := shaper.Apply(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "clean":
		if err := shaper.Clean(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "license":
		fmt.Println("NFTRL")
		fmt.Println("Copyright 2026 肖其顿 (XIAO QI DUN)")
		fmt.Println("")
		fmt.Println("This product includes software developed by")
		fmt.Println("肖其顿 (XIAO QI DUN) (https://github.com/xiaoqidun/nftrl).")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", flag.Arg(0))
		flag.Usage()
		os.Exit(1)
	}
}
