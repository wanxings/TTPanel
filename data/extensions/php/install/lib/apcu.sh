#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file



download_Url=$NODE_URL


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}

Ext_Path(){
	case "${version}" in 
		'52')
		extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/apcu.so'
		;;
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/apcu.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/apcu.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/apcu.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/apcu.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/apcu.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/apcu.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/apcu.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/apcu.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/apcu.so'
		;;
        '80')
        extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/apcu.so'
        ;;
        '81')
        extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/apcu.so'
        ;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/apcu.so'
		;;
	esac
}
Install_Apcu(){

	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi

	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'apcu.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过apcu,请选择其它版本!"
		echo "php-$vphp is already installed apcu, Plese select other version!"
		exit
	fi
	
	if [ ! -f "${extFile}" ];then
		if [ "${version}" -ge 80 ];then
			apcuVer="5.1.22"
		elif [ "${version}" -ge 70 ];then
			apcuVer="5.1.21"
		else 
			apcuVer="4.0.10"
		fi
		wget ${download_Url}/extensions/php/install/lib/apcu-${apcuVer}.tgz -T 5
		tar zxvf apcu-${apcuVer}.tgz
		rm -f apcu-${apcuVer}.tgz
		cd apcu-${apcuVer}
		/www/server/php/${version}/bin/phpize
		./configure --with-php-config=/www/server/php/${version}/bin/php-config
		make && make install
		cd ../
		rm -rf apcu-${apcuVer}
	fi

	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	fi
	echo -e "extension = ${extFile}\n" >> /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		echo -e "\extension = ${extFile}\n" >> /www/server/php/$version/etc/php-cli.ini
	fi
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}
Uninstall_Apcu(){
	sed -i '/apcu.so/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/$version/etc/php-cli.ini ];then
		sed -i '/apcud/d' /www/server/$version/etc/php-cli.ini
	fi

	rm -f $extFile
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}
if [ "$actionType" == 'install' ];then
	Ext_Path
	Install_Apcu
elif [ "$actionType" == 'uninstall' ];then
	Ext_Path
	Uninstall_Apcu
fi
