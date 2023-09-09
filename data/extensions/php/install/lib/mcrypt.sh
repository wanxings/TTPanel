#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



mcrypt_version="1.0.5"
runPath=/root

if [ -z "${cpuCore}" ]; then
	cpuCore="1"
fi

ext_Path(){
	case "${version}" in 
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/mcrypt.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/mcrypt.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/mcrypt.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/mcrypt.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/mcrypt.so'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/mcrypt.so'
		;;
	esac
}

Install_Mcrypt(){
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'mcrypt.so'`
	if [ "${isInstall}" ];then
		echo "php-$vphp 已安装过mcrypt,请选择其它版本!"
		echo "php-$vphp mcrypt has been installed. Please select another version"
		exit 0
		return
	fi
	if [ ! -f "${extFile}" ];then 
		cd ${runPath}
		wget -O mcrypt-${mcrypt_version}.tgz $download_Url/extensions/php/install/lib/mcrypt-${mcrypt_version}.tgz
		tar -xvf mcrypt-${mcrypt_version}.tgz
		cd mcrypt-${mcrypt_version}
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ../
		rm -rf mcrypt-${mcrypt_version}*
	fi

	if [ ! -f "${extFile}" ];then
		GetSysInfo
		echo 'error';
		exit 1;
	fi

	echo -e "extension = mcrypt.so\n" >> /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		echo -e "extension = mcrypt.so\n" >> /www/server/php/$version/etc/php-cli.ini
	fi
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

Uninstall_Mcrypt(){
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp Not installed, please select another version!"
		return
	fi
	isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'mcrypt.so')
	if [ "${isInstall}" == "" ];then
		echo "php-$vphp 未安装mcrypt,请选择其它版本!"
		echo "php-$vphp mcrypt is not installed, please choose another version!"
		return
	fi

	rm -f ${extFile}
	sed -i '/mcrypt.so/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		sed -i '/mcrypt.so/d' /www/server/php/$version/etc/php-cli.ini
	fi
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
ext_Path
if [ "$actionType" == 'install' ];then
	Install_Mcrypt
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Mcrypt
fi

