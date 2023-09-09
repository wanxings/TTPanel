#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

Action=$1
Name=$2
Version=$3

header_file=/www/panel/data/shell/install_header.sh
. $header_file

serverUrl=$NODE_URL/extensions/php/install/lib

wget -N --no-check-certificate -O lib_$Name.sh $serverUrl/$Name.sh
echo '|-Start Install Lib---'
bash lib_$Name.sh $Action $Version
rm -rf lib_$Name.sh
echo '|-Successify --- 命令已执行! ---'