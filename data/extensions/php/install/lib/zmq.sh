#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL



zmq_Version='4.3.4'

expPath(){
  case "${version}" in 
    '53')
    extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/zmq.so'
    ;;
    '54')
    extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/zmq.so'
    ;;
    '55')
    extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/zmq.so'
    ;;
    '56')
    extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/zmq.so'
    ;;
    '70')
    extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/zmq.so'
    ;;
    '71')
    extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/zmq.so'
    ;;
    '72')
    extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/zmq.so'
    ;;
    '73')
    extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/zmq.so'
    ;;
    '74')
    extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/zmq.so'
    ;;
    '80')
    extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/zmq.so'
    ;;
    '81')
    extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/zmq.so'
    ;;
  esac
}

Install_zmq()
{

  if [ ! -f "/www/server/php/$version/bin/php-config" ];then
    echo "php-$vphp 未安装,请选择其它版本!"
    echo "php-$vphp not install, Plese select other version!"
    return
  fi
  
  isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'zmq.so'`
  if [ "${isInstall}" != "" ];then
    echo "php-$vphp 已安装过zmq,请选择其它版本!"
    echo "php-$vphp zmq has been installed, Plese select other version!"
    return
  fi

  runPath=/root
  expPath

  if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ];then
    Pack="zeromq-devel"
    ${PM} install ${Pack} -y
  elif [ "${PM}" == "apt-get" ];then
    Pack="libzmq3-dev"
    ${PM} install ${Pack} -y
  fi

  if [ ! -f "${extFile}" ];then
    wget $download_Url/extensions/php/install/lib/zmq-$zmq_Version.tgz
    tar -zxvf zmq-$zmq_Version.tgz
    cd zmq-$zmq_Version
    /www/server/php/$version/bin/phpize
    ./configure --with-php-config=/www/server/php/$version/bin/php-config --with-zmq
    make && make install
    cd ../
    rm -rf zmq*
   fi

   if [ ! -f "${extFile}" ];then
     echo 'error';
     exit 0;
   fi
     echo -e "\n[zmq]\nextension = zmq.so\n" >> /www/server/php/$version/etc/php.ini

    /etc/init.d/php-fpm-$version reload
}

Uninstall_zmq()
{
  expPath
  sed -i '/zmq/d' /www/server/php/$version/etc/php.ini
  rm -f ${extFile}
  /etc/init.d/php-fpm-$version reload
  echo '==============================================='
  echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
  Install_zmq
elif [ "$actionType" == 'uninstall' ];then
  Uninstall_zmq
fi