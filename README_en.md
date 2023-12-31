# TTPanel - Linux Management Panel

### Simple Linux control panel，supports CentOS、Ubuntu、Debian

[简体中文](README.md)  [Screenshots](/screenshots)

- ⚠️ Please do not use it in production environment
- ⚠️ The frontend uses a commercial framework and the frontend code cannot be open sourced



## Web Environment

- ✅ nginx
- ✅ php
- ✅ mysql
- ✅ phpmyadmin
- ✅ redis
- ✅ nodejs
- ❌ nginx Waf




## Installation

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

## Commands

### Menu

```
tt
```
### Start Panel

```
tt start
```
### Stop Panel

```
tt stop
```
### Restart Panel

```
tt restart
```
### Panel Status

```
tt status
```
