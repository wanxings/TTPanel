# TTPanel - Linux管理面板

### 简单的Linux管理面板，支持CentOS、Ubuntu、Debian

[English README](README_en.md) [截图](/screenshots)

- ⚠️ 在 1.0.0 正式版之前请勿用于生产环境
- ⚠️ 由于前端使用了商业框架，无法开源前端代码


## 


## web环境

- ✅ nginx
- ✅ php
- ✅ mysql
- ✅ phpmyadmin
- ✅ redis
- ✅ nodejs
- ❌ nginx防火墙




## 安装

### centos7

```
yum install -y wget && wget -O install_panel.sh https://download.ttpanel.org/install/src/install_panel_0.1.0.sh && sh install_panel.sh
```

### Ubuntu18+

```
wget -O install_panel.sh https://download.ttpanel.org/install/src/install_panel_0.1.0.sh && sudo bash install_panel.sh
```

### Debian10+

```
wget -O install_panel.sh https://download.ttpanel.org/install/src/install_panel_0.1.0.sh && bash install_panel.sh
```

## 命令

### 菜单

```
tt
```
### 启动面板

```
tt start
```
### 停止面板

```
tt stop
```
### 重启面板

```
tt restart
```
### 面板状态

```
tt status
```
