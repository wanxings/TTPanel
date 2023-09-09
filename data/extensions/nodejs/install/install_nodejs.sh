#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

run_path='/www/server/nodejs'

Install(){
  if [ ! -d $run_path ]; then
    mkdir -m 755 $run_path
  fi
  filename=$(basename "$download_url")
  pathname=$(basename "$download_url" .tar.gz)
  cd "$run_path" || exit
  if ! wget "$download_url" -T 20; then
      echo "Download failed: $download_url"
      exit 1
  fi
  tar -zxf "$filename"
  rm -f "$filename"
  mv "$pathname" "$version"
  chown -R root:root "$version"

  node_bin=$run_path/$version/bin/node
  npm_js=$run_path/$version/lib/node_modules/npm/bin/npm-cli.js
  npx_js=$run_path/$version/lib/node_modules/npm/bin/npx-cli.js
  sed -i "s|#!/usr/bin/env node|#!${node_bin}|g" "$npm_js"
  sed -i "s|#!/usr/bin/env node|#!${node_bin}|g" "$npx_js"
  echo '|--- The installation is complete ---'

}
Uninstall(){
  if [ -d "${run_path:?}/${version:?}" ]; then
    rm -rf "${run_path:?}/${version:?}"
  fi
  echo '|--- Uninstall completed ---'
}
action=$1
version=$2
download_url=$3
if [ "$action" == 'install' ] ;then
  Uninstall
  Install
elif [ "$action" == 'uninstall' ];then
	Uninstall
fi

echo '----- successful -----'