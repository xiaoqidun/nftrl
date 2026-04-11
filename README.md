# NFTRL
一个基于 nftables 的零侵入 OpenWrt 限速工具

# 源码安装

## 1. 构建程序
```batch
:: 1. 克隆仓库
git clone https://github.com/xiaoqidun/nftrl.git
cd nftrl
:: 2. 交叉编译（以从 Windows 编译 Linux amd64 为例）
set GOOS=linux
set GOARCH=amd64
go build -o ./openwrt-dist/usr/sbin/nftrl -trimpath -ldflags "-s -w -buildid=" ./cmd/main.go
```

## 2. 打包上传
```batch
tar -cf openwrt-dist.tar openwrt-dist
scp openwrt-dist.tar root@ip:/root
```

## 3. 配置服务
```shell
ssh root@ip
tar -xf openwrt-dist.tar
cp -r openwrt-dist/* /
chmod +x /usr/sbin/nftrl /etc/init.d/nftrl
/etc/init.d/nftrl enable
/etc/init.d/nftrl start
rm -r openwrt-dist openwrt-dist.tar
```

# 授权协议
本项目使用 [Apache License 2.0](https://github.com/xiaoqidun/nftrl/blob/main/LICENSE) 授权协议