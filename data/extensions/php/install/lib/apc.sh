#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file



download_Url=$NODE_URL

Install_apc()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'apc.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过apc,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	
	case "${version}" in 
		'52')
		extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/apc.so'
		;;
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/apc.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/apc.so'
		;;
	esac

	if [ ! -f "${extFile}" ];then
		wget -c wget ${download_Url}/extensions/php/install/lib/apc-2b75460.tar.gz -T 5
		tar -zxvf apc-2b75460.tar.gz
		cd apc-2b75460
		/www/server/php/$version/bin/phpize
		./configure  --with-php-config=/www/server/php/$version/bin/php-config --enable-apc
		make && make install
	fi
	
	if [ ! -f "${extFile}" ];then
		echo 'error'
		exit 0;
	fi

	echo -e "\n[APC]\nextension=apc.so\napc.enabled = 1\napc.shm_segments = 1\napc.shm_size = 64M\napc.optimization = 1\napc.num_files_hint = 0\napc.ttl = 0\napc.gc_ttl = 3600\napc.cache_by_default = on" >> /www/server/php/$version/etc/php.ini
	service php-fpm-$version reload
	cd ..
	rm -rf apc*
	echo '==============================================='
	echo 'successful!'
}

Uninstall_apc()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'apc.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装apc,请选择其它版本!"
		echo "php-$vphp not install apc, Plese select other version!"
		return
	fi
		
	sed -i '/apc/d'  /www/server/php/$version/etc/php.ini
	sed -i '/APC/d'  /www/server/php/$version/etc/php.ini
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_apc
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_apc
fi