#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



Ext_Path(){
	case "${version}" in 
		'52')
		extFile="/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/ldap.so"
		;;
		'53')
		extFile="/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/ldap.so"
		;;
		'54')
		extFile="/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/ldap.so"
		;;
		'55')
		extFile="/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/ldap.so"
		;;
		'56')
		extFile="/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/ldap.so"
		;;
		'70')
		extFile="/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/ldap.so"
		;;
		'71')
		extFile="/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/ldap.so"
		;;
		'72')
		extFile="/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/ldap.so"
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/ldap.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/ldap.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/ldap.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/ldap.so'
		;;
		'82')
        extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/ldap.so'
        ;;

	esac
}

Install_ldap()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'ldap.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过ldap,请选择其它版本!"
		echo "php-$vphp is installed ldap, Plese select other version!"
		return
	fi
	

	if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ];then
		Pack="openldap openldap-devel"
		${PM} install ${Pack} -y
		ln -sf /usr/lib64/libldap* /usr/lib
	elif [ "${PM}" == "apt-get" ];then
		Pack="libldap2-dev"
		${PM} install ${Pack} -y
		ln -sf /usr/lib/x86_64-linux-gnu/libldap* /usr/lib
	fi
	
	if [ ! -d "/www/server/php/$version/src/ext/ldap" ];then
		mkdir -p /www/server/php/$version/src
		wget -O $version-ext.tar.gz $download_Url/extensions/php/install/lib/$version-ext.tar.gz
		tar -zxf $version-ext.tar.gz -C /www/server/php/$version/src/ 
		rm -f $version-ext.tar.gz
	fi
		
	if [ ! -f "${extFile}" ];then
		cd /www/server/php/$version/src/ext/ldap
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
	fi
	
	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	fi

	echo -e "extension = ldap.so" >> /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		echo -e "extension = ldap.so" >> /www/server/php/$version/etc/php-cli.ini
	fi
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

Uninstall_ldap()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'ldap.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装ldap,请选择其它版本!"
		echo "php-$vphp not install ldap, Plese select other version!"
		return
	fi

	sed -i '/ldap.so/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		sed -i '/ldap.so/d' /www/server/php/$version/etc/php-cli.ini
	fi
	rm -f ${extFile}

	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
Ext_Path
if [ "$actionType" == 'install' ];then
	Install_ldap
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_ldap
fi

