#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



Install_Intl()
{
	case "${version}" in 
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/intl.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/intl.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/intl.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/intl.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/intl.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/intl.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/intl.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/intl.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/intl.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/intl.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/intl.so'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/intl.so'
		;;
	esac
	
	isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'intl.so')
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装intl,请选择其它版本!"
		return
	fi
	
	if [ ! -d "/www/server/php/$version/src/ext/intl" ];then
		mkdir -p /www/server/php/$version/src
		wget -O $version-ext.tar.gz $download_Url/extensions/php/install/lib/$version-ext.tar.gz
		tar -zxf $version-ext.tar.gz -C /www/server/php/$version/src/ 
		rm -f $version-ext.tar.gz
	fi
	
	
	if [ ! -f "$extFile" ];then
		yum install icu -y
		cd /www/server/php/$version/src/ext/intl
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
	fi
	
	if [ ! -f "$extFile" ];then
		echo "ERROR!"
		return;
	fi
	if [ "$version" = '53' ];then
		echo "extension=$extFile" >> /www/server/php/$version/etc/php.ini
	else
		echo ";extension=$extFile" >> /www/server/php/$version/etc/php.ini
	fi
	service php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}


Uninstall_Intl()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return
	fi
	
	case "${version}" in 
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/intl.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/intl.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/intl.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/intl.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/intl.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/intl.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/intl.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/intl.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/intl.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/intl.so'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/intl.so'
		;;
	esac
	if [ ! -f "$extFile" ];then
		echo "php-$vphp 未安装intl,请选择其它版本!"
		return
	fi
	
	sed -i '/intl.so/d'  /www/server/php/$version/etc/php.ini
		
	rm -f $extFile
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_Intl
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Intl
fi
