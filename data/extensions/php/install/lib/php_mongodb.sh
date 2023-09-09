#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


runPath=/root
Ext_Path(){
	case "${version}" in
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/mongodb.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/mongodb.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/mongodb.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/mongodb.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/mongodb.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/mongodb.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/mongodb.so'
		;;
	esac
}
Install_mongodb()
{
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'mongodb.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装mongodb,请选择其它版本!"
		return
	fi
	
	if [ ! -f "$extFile" ];then
		if [ -z "${cpuCore}" ]; then
			cpuCore="1"
		fi
        
        cd ${runPath}
        if [ ${version} -ge 72 ];then
            mongodbVersion="1.12.0"
        else
		    mongodbVersion="1.9.0"
	    fi
		wget -O mongodb-${mongodbVersion}.tgz $download_Url/extensions/php/install/lib/mongodb-${mongodbVersion}.tgz
		tar xvf mongodb-${mongodbVersion}.tgz
		cd mongodb-${mongodbVersion}
		
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make -j${cpuCore}
		make install
		cd ..
		rm -rf mongodb-${mongodbVersion}*
	fi
	
	if [ ! -f "$extFile" ];then
		echo "ERROR!"
		return;
	fi
	echo "extension=$extFile" >> /www/server/php/$version/etc/php.ini
    if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        echo -e "extension = $extFile" >> /www/server/php/$version/etc/php-cli.ini
    fi
	
	service php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'
}


Uninstall_mongodb()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return
	fi
	
	if [ ! -f "$extFile" ];then
		echo "php-$vphp 未安装mongodb,请选择其它版本!"
		return
	fi
	
	sed -i '/mongodb.so/d'  /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        sed -i '/mongodb.so/d' /www/server/php/$version/etc/php-cli.ini
    fi
	rm -f $extFile

	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
Ext_Path
if [ "$actionType" == 'install' ];then
	Install_mongodb
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_mongodb
fi

