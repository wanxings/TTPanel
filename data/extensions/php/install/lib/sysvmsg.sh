#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


Install_sysvmsg()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'sysvmsg.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过sysvmsg,请选择其它版本!"
		echo "php-$vphp is installed sysvmsg, Plese select other version!"
		return
	fi
	
	if [ ! -d "/www/server/php/$version/src/ext/sysvmsg" ];then
		mkdir -p /www/server/php/$version/src
		wget -O $version-ext.tar.gz $download_Url/extensions/php/install/lib/$version-ext.tar.gz
		tar -zxf $version-ext.tar.gz -C /www/server/php/$version/src/ 
		rm -f $version-ext.tar.gz
	fi
	
	case "${version}" in 
		'52')
		extFile="/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/sysvmsg.so"
		;;
		'53')
		extFile="/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/sysvmsg.so"
		;;
		'54')
		extFile="/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/sysvmsg.so"
		;;
		'55')
		extFile="/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/sysvmsg.so"
		;;
		'56')
		extFile="/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/sysvmsg.so"
		;;
		'70')
		extFile="/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/sysvmsg.so"
		;;
		'71')
		extFile="/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/sysvmsg.so"
		;;
		'72')
		extFile="/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/sysvmsg.so"
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/sysvmsg.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/sysvmsg.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/sysvmsg.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/sysvmsg.so'
		;;
	esac
	
	if [ ! -f "${extFile}" ];then
		cd /www/server/php/$version/src/ext/sysvmsg
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
	fi
	
	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	fi

	echo -e "extension = " ${extFile} >> /www/server/php/$version/etc/php.ini
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

Uninstall_sysvmsg()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'sysvmsg.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装sysvmsg,请选择其它版本!"
		echo "php-$vphp not install sysvmsg, Plese select other version!"
		return
	fi

	sed -i '/sysvmsg.so/d' /www/server/php/$version/etc/php.ini

	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_sysvmsg
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_sysvmsg
fi
