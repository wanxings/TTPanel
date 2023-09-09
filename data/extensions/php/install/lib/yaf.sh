#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL

Ext_Path(){
	case "${version}" in 
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/yaf.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/yaf.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/yaf.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/yaf.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/yaf.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/yaf.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/yaf.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/yaf.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/yaf.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/yaf.so'
		;;
	esac
}

Install_yaf()
{
	
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi

	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'yaf.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装yaf,请选择其它版本!"
		return
	fi
	
	if [ ! -f "$extFile" ];then
		wafV='2.3.5';
		if [ "$version" -ge "70" ];then
			wafV='3.3.5';
		fi
		wget -O yaf-$wafV.tgz $download_Url/extensions/php/install/lib/yaf-$wafV.tgz
		tar xvf yaf-$wafV.tgz
		cd yaf-$wafV
		
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ..
		rm -rf yaf-*
		rm -f package.xml
	fi
	
	if [ ! -f "$extFile" ];then
		echo "install failed!"
		return;
	fi
	echo "extension=$extFile" >> /www/server/php/$version/etc/php.ini
	
	/etc/init.d/php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}


Uninstall_yaf()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return
	fi
	
	if [ ! -f "$extFile" ];then
		echo "php-$vphp 未安装yaf,请选择其它版本!"
		return
	fi
	
	sed -i '/yaf.so/d'  /www/server/php/$version/etc/php.ini
	rm -f $extFile
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Ext_Path
	Install_yaf
elif [ "$actionType" == 'uninstall' ];then
	Ext_Path
	Uninstall_yaf
fi
