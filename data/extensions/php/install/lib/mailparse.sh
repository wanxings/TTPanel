#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



mailparse_version="3.1.4"
runPath=/root


expFile(){
	case "${version}" in 
		'70')
			extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/mailparse.so'
		;;
		'71')
			extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/mailparse.so'
		;;
		'72')
			extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/mailparse.so'
		;;
		'73')
			extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/mailparse.so'
		;;
		'74')
			extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/mailparse.so'
		;;
		'80')
			extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/mailparse.so'
		;;
		'81')
			extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/mailparse.so'
		;;
		'82')
    		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/mailparse.so'
    	;;
	esac
}

Install_mailparse()
{
	expFile

	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'mailparse.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过mailparse,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	if [ ! -f "${extFile}" ];then
		cd ${runPath}
		wget -O mailparse-${mailparse_version}.tgz $download_Url/extensions/php/install/lib/mailparse-${mailparse_version}.tgz
		tar -xvf mailparse-${mailparse_version}.tgz
		cd mailparse-${mailparse_version}
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
	    make && make install
		cd ../
		rm -rf ${runPath}/mailparse*
	fi
	
	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	else
		echo -e "extension = $extFile" >> /www/server/php/$version/etc/php.ini
		if [ -f /www/server/php/$version/etc/php-cli.ini ];then
			echo -e "extension = $extFile" >> /www/server/php/$version/etc/php-cli.ini
		fi
	fi
	service php-fpm-$version restart
	echo '==============================================='
	echo 'successful!'
}

Uninstall_mailparse()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'mailparse.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装mailparse,请选择其它版本!"
		echo "php-$vphp not install mailparse, Plese select other version!"
		return
	fi

	expFile

	sed -i '/mailparse.so/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		sed -i '/mailparse.so/d' /www/server/php/$version/etc/php-cli.ini
	fi
	
	rm -f ${extFile}
	service php-fpm-$version restart
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_mailparse
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_mailparse
fi
