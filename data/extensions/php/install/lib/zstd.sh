#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL




zstd_Version='0.11.0'

expPath(){
  case "${version}" in 
    '53')
    extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/zstd.so'
    ;;
    '54')
    extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/zstd.so'
    ;;
    '55')
    extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/zstd.so'
    ;;
    '56')
    extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/zstd.so'
    ;;
    '70')
    extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/zstd.so'
    ;;
    '71')
    extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/zstd.so'
    ;;
    '72')
    extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/zstd.so'
    ;;
    '73')
    extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/zstd.so'
    ;;
    '74')
    extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/zstd.so'
    ;;
    '80')
    extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/zstd.so'
    ;;
    '81')
    extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/zstd.so'
    ;;
  esac
}

Install_zstd()
{

  if [ ! -f "/www/server/php/$version/bin/php-config" ];then
    echo "php-$vphp 未安装,请选择其它版本!"
    echo "php-$vphp not install, Plese select other version!"
    return
  fi
  
  isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'zstd.so'`
  if [ "${isInstall}" != "" ];then
    echo "php-$vphp 已安装过zstd,请选择其它版本!"
    echo "php-$vphp zstd has been installed, Plese select other version!"
    return
  fi

  runPath=/root
  expPath

  if [ ! -f "${extFile}" ];then
    wget $download_Url/extensions/php/install/lib/zstd-$zstd_Version.tgz
    tar -zxvf zstd-$zstd_Version.tgz
    cd zstd-$zstd_Version
    /www/server/php/$version/bin/phpize
    ./configure --with-php-config=/www/server/php/$version/bin/php-config --with-zstd
    make && make install
    cd ../
    rm -rf zstd-0*
   fi

   if [ ! -f "${extFile}" ];then
     echo 'error';
     exit 0;
   fi
     echo -e "\n[zstd]\nextension = zstd.so\n" >> /www/server/php/$version/etc/php.ini

    /etc/init.d/php-fpm-$version reload
}

Uninstall_zstd()
{
  expPath
  sed -i '/zstd/d' /www/server/php/$version/etc/php.ini
  rm -f ${extFile}
  /etc/init.d/php-fpm-$version reload
  echo '==============================================='
  echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
  Install_zstd
elif [ "$actionType" == 'uninstall' ];then
  Uninstall_zstd
fi