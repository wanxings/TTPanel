#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


expFile(){
    case "${version}" in 
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/xdebug.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/xdebug.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/xdebug.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/xdebug.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/xdebug.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/xdebug.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/xdebug.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/xdebug.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/xdebug.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/xdebug.so'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/xdebug.so'
		;;
	esac
}

Install_xdebug(){
	expFile

	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装，请选择其他版本"
		echo "php-$vphp not install, Plese select other version!"
		exit 0
	fi

	isInstall=$(cat /www/server/php/$version/etc/php.ini | grep 'xdebug.so')
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过xdebug,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		exit 0
	fi

	if [ ! -f "${extFile}" ];then
		if [ "${version}" == "53" ] || [ "${version}" == "54" ] || [ "${version}" == "55" ] || [ "${version}" == "56" ];then
			xdebug_version="2.2.7"
		elif [ "$version" == "70" ] || [ "${version}" == "71" ];then
			xdebug_version='2.8.0'
		elif [ "$version" == "72" ] || [ "${version}" == "73" ] || [ "${version}" == "74" ];then
			xdebug_version='3.1.6'
		else
			xdebug_version='3.2.0'
		fi

		wget $download_Url/extensions/php/install/lib/xdebug-$xdebug_version.tgz -T 5
		tar -zxvf xdebug-$xdebug_version.tgz
		cd xdebug-$xdebug_version

		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ..
		rm -rf xdebug-*
	fi

	if [ ! -f "${extFile}" ];then
		echo '安装失败！';
		exit 0 
	fi

	echo "zend_extension=$extFile" >> /www/server/php/$version/etc/php.ini
	if [ -f "/www/server/php/$version/etc/php-cli.ini" ];then
		echo "zend_extension=$extFile" >> /www/server/php/$version/etc/php-cli.ini
	fi
	/etc/init.d/php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}

Uninstall_xdebug(){
	expFile
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return
	fi
	isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'xdebug.so')
	if [ "${isInstall}" == "" ];then
		echo "php-$vphp 未安装xdebug,请选择其它版本!"
		return
	fi

	sed -i '/xdebug.so/d' /www/server/php/$version/etc/php.ini
	if [ -f "/www/server/php/$version/etc/php-cli.ini" ];then
		sed -i '/xdebug.so/d' /www/server/php/$version/etc/php-cli.ini
	fi

	rm -f ${extFile}
	/etc/init.d/php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}

if [ "$actionType" == "install" ];then
	Install_xdebug
elif [[ "$actionType" == "uninstall" ]];then
	Uninstall_xdebug
fi
