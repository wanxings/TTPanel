#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



geoip_Version='1.1.1'

expPath(){
  case "${version}" in 
    '52')
    extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/geoip.so'
    ;;
    '53')
    extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/geoip.so'
    ;;
    '54')
    extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/geoip.so'
    ;;
    '55')
    extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/geoip.so'
    ;;
    '56')
    extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/geoip.so'
    ;;
    '70')
    extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/geoip.so'
    ;;
    '71')
    extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/geoip.so'
    ;;
    '72')
    extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/geoip.so'
    ;;
    '73')
    extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/geoip.so'
    ;;
	  '74')
    extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/geoip.so'
    ;;
  esac
}

Install_geoip()
{

  if [ ! -f "/www/server/php/$version/bin/php-config" ];then
    echo "php-$vphp 未安装,请选择其它版本!"
    echo "php-$vphp not install, Plese select other version!"
    return
  fi
  
  isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'geoip.so'`
  if [ "${isInstall}" != "" ];then
    echo "php-$vphp 已安装过geoip,请选择其它版本!"
    echo "php-$vphp geoip has been installed, Plese select other version!"
    return
  fi


  runPath=/root
  expPath

  if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ];then
		Pack="geoip geoip-devel"
		${PM} install ${Pack} -y
	elif [ "${PM}" == "apt-get" ];then
		Pack="libgeoip-dev"
		${PM} install ${Pack} -y
	fi

  if [ ! -f "${extFile}" ];then
  	wget $download_Url/extensions/php/install/lib/geoip-$geoip_Version.tgz
  	tar -zxvf geoip-$geoip_Version.tgz
  	cd geoip-$geoip_Version
  	/www/server/php/$version/bin/phpize
  	./configure --with-php-config=/www/server/php/$version/bin/php-config --with-geoip
  	make && make install
  	cd ../
  	rm -rf geoip*
 	fi

 	if [ ! -f "${extFile}" ];then
 		echo 'error';
 		exit 0;
 	fi
   	echo -e "\n[geoip]\nextension = geoip.so\n" >> /www/server/php/$version/etc/php.ini

    /etc/init.d/php-fpm-$version reload
}

Uninstall_geoip()
{
  expPath
	sed -i '/geoip/d' /www/server/php/$version/etc/php.ini
  rm -f ${extFile}
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_geoip
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_geoip
fi