#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


Install_yac()
{
	case "${version}" in 
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/yac.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/yac.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/yac.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/yac.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/yac.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/yac.so'
		;;
	esac
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'yac.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装yac,请选择其它版本!"
		return
	fi

	if [ ! -f "$extFile" ];then

		yacVersion='2.2.1'
		wget -O yac-$yacVersion.tgz $download_Url/extensions/php/install/lib/yac-$yacVersion.tgz
		tar -zxf yac-$yacVersion.tgz
		cd yac-$yacVersion

		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install

		cd ..
		rm -rf yac-*
	fi

	if [ ! -f "$extFile" ];then
		echo "ERROR!"
		return;
	fi

	echo "extension=$extFile" >> /www/server/php/$version/etc/php.ini
	
	/etc/init.d/php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}
Uninstall_yac()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return
	fi

	case "${version}" in 
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/yac.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/yac.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/yac.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/yac.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/yac.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/yac.so'
		;;
	esac
	
	if [ ! -f "$extFile" ];then
		echo "php-$vphp 未安装yac,请选择其它版本!"
		return
	fi

	sed -i '/yac.so/d'  /www/server/php/$version/etc/php.ini
		
	rm -f $extFile
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'

}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_yac
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_yac
fi
