#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


Install_Pgsql()
{	
	
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep ' pgsql.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过pgsql,请选择其它版本!"
		echo "php-$vphp is alerady install pgsql, Plese select other version!"
		return
	fi

	if [ ! -d "/www/server/php/$version/src/ext/pgsql" ];then
		mkdir -p /www/server/php/$version/src
		wget -O $version-ext.tar.gz $download_Url/extensions/php/install/lib/$version-ext.tar.gz
		tar -zxf $version-ext.tar.gz -C /www/server/php/$version/src/ 
		rm -f $version-ext.tar.gz
	fi
	case "${version}" in 
		'52')
		extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/pgsql.so'
		;;
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/pgsql.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/pgsql.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/pgsql.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/pgsql.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/pgsql.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/pgsql.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/pgsql.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/pgsql.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/pgsql.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/pgsql.so'
		;;
	esac

	if [ ! -f "${extFile}" ];then
		cd /www/server/php/$version/src/ext/pgsql
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config --with-pgsql=/www/server/pgsql
		make && make install
	fi

	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	fi

	echo -e "extension = pgsql.so" >> /www/server/php/$version/etc/php.ini
	/www/server/php/${version}/bin/php -m |grep pgsql
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'

}
Uninstall_Pgsql()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi

	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'pgsql.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装Fileinfo,请选择其它版本!"
		echo "php-$vphp not install Fileinfo, Plese select other version!"
		return
	fi

	sed -i '/ pgsql.so/d' /www/server/php/$version/etc/php.ini

	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}
actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_Pgsql
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Pgsql
fi
