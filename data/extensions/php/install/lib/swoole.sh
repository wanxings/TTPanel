#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL




Swoole_Version='1.10.1'

expPath(){
  case "${version}" in 
    '53')
    extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/swoole.so'
    ;;
    '54')
    extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/swoole.so'
    ;;
    '55')
    extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/swoole.so'
    ;;
    '56')
    extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/swoole.so'
    ;;
    '70')
    extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/swoole.so'
    ;;
    '71')
    extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/swoole.so'
    ;;
    '72')
    extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/swoole.so'
    ;;
  esac
}

Install_Swoole()
{

  if [ ! -f "/www/server/php/$version/bin/php-config" ];then
    echo "php-$vphp 未安装,请选择其它版本!"
    echo "php-$vphp not install, Plese select other version!"
    return
  fi
  
  isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'swoole.so'`
  if [ "${isInstall}" != "" ];then
    echo "php-$vphp 已安装过swoole,请选择其它版本!"
    echo "php-$vphp not install, Plese select other version!"
    return
  fi

  runPath=/root
  expPath
  if [ ! -f "${extFile}" ];then
    if [ "${version}" == "70" ] || [ "${version}" == "71" ] || [ "${version}" == "72" ]; then
      Swoole_Version='2.2.0'
    fi
  	wget $download_Url/extensions/php/install/lib/swoole-$Swoole_Version.tgz
  	tar -zxvf swoole-$Swoole_Version.tgz
  	cd swoole-$Swoole_Version
  	/www/server/php/$version/bin/phpize
  	./configure --with-php-config=/www/server/php/$version/bin/php-config --enable-openssl --with-openssl-dir=/usr/local/openssl --enable-sockets
  	make && make install
  	cd ../
  	rm -rf swoole*
 	fi

 	if [ ! -f "${extFile}" ];then
 		echo 'error';
 		exit 0;
 	fi
   	
   	echo -e "\n[swoole]\nextension = swoole.so\n" >> /www/server/php/$version/etc/php.ini

   	service php-fpm-$version reload
}



Uninstall_Swoole()
{
  expPath
	sed -i '/swoole/d' /www/server/php/$version/etc/php.ini
  rm -f ${extFile}
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_Swoole
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Swoole
fi
