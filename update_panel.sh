#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL
u0010() {
  mkdir -p /www/wwwlogs/analytics
  chown www.www /www/wwwlogs/analytics
  chmod 744 /www/wwwlogs/analytics
}
Update() {
  #清理临时目录
  rm -rf /tmp/update_panel
  mkdir -p /tmp/update_panel
  #下载更新包
  if [ "$(uname -m)" == "aarch64" ]; then
    wget -O update_panel.tar.gz "${download_Url}/update/${version}/TTPanel_arm64_${pre_version}.tar.gz" -T 10
  else
    wget -O update_panel.tar.gz "${download_Url}/update/${version}/TTPanel_amd64_${pre_version}.tar.gz" -T 10
  fi

  chmod 755 update_panel.tar.gz
  tar -xzf update_panel.tar.gz -C /tmp/update_panel
  if [ ! -f /tmp/update_panel/panel/TTPanel ]; then
    ls -lh update_panel.tar.gz
    Red_Error "ERROR: Failed to download, please try install again!" "ERROR: 下载主程序失败，请尝试重新更新！"
  fi
  if [ -f /tmp/update_panel/panel/tt ]; then
    mv -f /tmp/update_panel/panel/tt /etc/init.d/tt
    chmod +x /etc/init.d/tt
  fi

  #版本处理
  u010

  cp -rf /tmp/update_panel/* /www/
  rm -f update_panel.tar.gz
  rm -rf /tmp/update_panel

  echo "Update Successfully! "
  sleep 3
  nohup tt restart >/dev/null 2>/dev/null &
}
Red_Error() {
  echo '================================================='
  printf '\033[1;31;40m%b\033[0m\n' "$@"
  exit 1
}
version=$1
pre_version=$2

echo "Update panel.... version: ${version} pre:${pre_version}  "
Update
