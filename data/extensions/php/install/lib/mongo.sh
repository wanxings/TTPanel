#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL


Install_mongo()
{
	case "${version}" in
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/mongo.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/mongo.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/mongo.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/mongo.so'
		;;
	esac
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'mongo.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装mongo,请选择其它版本!"
		return
	fi
	
	
	if [ ! -f "$extFile" ];then
		wget -O mongo-1.6.14.tgz $download_Url/extensions/php/install/lib/mongo-1.6.14.tgz
		tar xvf mongo-1.6.14.tgz
		cd mongo-1.6.14
		
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ..
		rm -rf mongo-1.6.14*
		rm -f package.xml
	fi
	
	if [ ! -f "$extFile" ];then
		echo "ERROR!"
		return;
	fi
	echo "extension=$extFile" >> /www/server/php/$version/etc/php.ini
	
	service php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}


Uninstall_mongo()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return
	fi
	
	case "${version}" in
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/mongo.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/mongo.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/mongo.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/mongo.so'
		;;
	esac
	if [ ! -f "$extFile" ];then
		echo "php-$vphp 未安装mongo,请选择其它版本!"
		return
	fi
	
	sed -i '/mongo.so/d'  /www/server/php/$version/etc/php.ini
		
	rm -f $extFile
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_mongo
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_mongo
fi