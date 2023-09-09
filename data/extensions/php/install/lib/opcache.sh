#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL



Install_Opcache()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'opcache.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过Opcache,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	case "${version}" in 
		'52')
		extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/opcache.so'
		;;
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/opcache.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/opcache.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/opcache.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/opcache.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/opcache.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/opcache.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/opcache.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/opcache.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/opcache.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/opcache.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/opcache.so'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/opcache.so'
	esac
	
	if [ ! -f "$extFile" ];then
		wget $download_Url/extensions/php/install/lib/zendopcache-7.0.5.tgz
		tar -zxf zendopcache-7.0.5.tgz
		rm -f zendopcache-7.0.5.tgz
		cd zendopcache-7.0.5
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make
		make install
		rm -rf zendopcache-7.0.5
	fi
	
	sed -i '/;opcache./d' /www/server/php/$version/etc/php.ini	
	if [ "$version" -ge '80' ];then
		sed -i "s#;opcache#;opcache\n[Zend Opcache]\nzend_extension=${extFile}\nopcache.enable = 1\nopcache.memory_consumption=128\nopcache.interned_strings_buffer=32\nopcache.max_accelerated_files=80000\nopcache.revalidate_freq=3\nopcache.fast_shutdown=1\nopcache.enable_cli=1\nopcache.jit_buffer_size=128m\nopcache.jit=1205#" /www/server/php/$version/etc/php.ini
	else
		sed -i "s#;opcache#;opcache\n[Zend Opcache]\nzend_extension=${extFile}\nopcache.enable = 1\nopcache.memory_consumption=128\nopcache.interned_strings_buffer=32\nopcache.max_accelerated_files=80000\nopcache.revalidate_freq=3\nopcache.fast_shutdown=1\nopcache.enable_cli=1#" /www/server/php/$version/etc/php.ini
	fi
	
	/etc/init.d/php-fpm-$version reload
}

Uninstall_Opcache()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'opcache.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装Opcache,请选择其它版本!"
		echo "php-$vphp not install Opcache, Plese select other version!"
		return
	fi
	
	sed -i '/opcache./d' /www/server/php/$version/etc/php.ini
	sed -i '/Opcache/d' /www/server/php/$version/etc/php.ini
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'						
}



actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_Opcache
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Opcache
fi
