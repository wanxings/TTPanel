#!/bin/bash
PATH=/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
export HOME=/root

Install(){
  $npm_bin_path install "${modules_name}" -g
}
Uninstall(){
  $npm_bin_path uninstall "${modules_name}" -g
}
Upgrade(){
  $npm_bin_path update -global "${modules_name}"
}

action=$1
node_version=$2
modules_name=$3
registry_url=$4


node_path="/www/server/nodejs/${node_version}"
node_bin_path="${node_path}/bin"
npm_bin_path="${node_path}/bin/npm"
yarn_bin_path="${node_path}/bin/yarn"
if [[ -z "$node_bin_path" ]]; then
  echo "node_bin_path is empty"
  exit 1
fi
export PATH=$PATH:$node_bin_path
export NODE_PATH="${node_path}/etc/node_modules"
$npm_bin_path config set registry "${registry_url}"
$npm_bin_path config set prefix "${node_path}/"
$npm_bin_path config set cache "${node_path}/cache"
if [ -d "$yarn_bin_path" ]; then
    $yarn_bin_path config set registry "${registry_url}"
fi
if [ "$action" == 'install' ] ;then
  Install
elif [ "$action" == 'uninstall' ];then
	Uninstall
elif [ "$action" == 'upgrade' ];then
	Upgrade
fi

echo '----- successful -----'