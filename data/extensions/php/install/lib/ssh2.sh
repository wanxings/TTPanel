#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


Install_ssh2()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'ssh2.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过ssh2,请选择其它版本!"
		echo "php-$vphp is installed ssh2, Plese select other version!"
		return
	fi
	

	if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ];then
		Pack="libssh2-devel"
	elif [ "${PM}" == "apt-get" ];then
		Pack="libssh2-1-dev"
	fi
	${PM} install ${Pack} -y
	
	case "${version}" in 
		'70')
		extFile="/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/ssh2.so"
		;;
		'71')
		extFile="/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/ssh2.so"
		;;
		'72')
		extFile="/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/ssh2.so"
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/ssh2.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/ssh2.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/ssh2.so'
		;;
	esac
	
	if [ ! -f "${extFile}" ];then
		ssh2Ver="1.3.1"
		wget -O ssh2-${ssh2Ver}.tgz $download_Url/extensions/php/install/lib/ssh2-${ssh2Ver}.tgz
		tar -xvf ssh2-${ssh2Ver}.tgz
		cd ssh2-${ssh2Ver}
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ..
		rm -rf ssh2-*
	fi
	
	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	fi

	echo -e "extension = " ${extFile} >> /www/server/php/$version/etc/php.ini
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

Uninstall_ssh2()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'ssh2.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装ssh2,请选择其它版本!"
		echo "php-$vphp not install ssh2, Plese select other version!"
		return
	fi

	sed -i '/ssh2.so/d' /www/server/php/$version/etc/php.ini

	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_ssh2
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_ssh2
fi
