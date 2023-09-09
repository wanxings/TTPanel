#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL

Root_Path="/www"
Setup_Path=$Root_Path/server/phpmyadmin
webserver="nginx"



Install_phpMyAdmin()
{

	wget -O phpMyAdmin.zip $download_Url/extensions/phpmyadmin/phpMyAdmin-${1}.zip -T20
	mkdir -p $Setup_Path

	unzip -o phpMyAdmin.zip -d $Setup_Path/ > /dev/null
	rm -f phpMyAdmin.zip
	rm -rf $Root_Path/server/phpmyadmin/phpmyadmin*


	phpmyadminExt=`cat /dev/urandom | head -n 32 | md5sum | head -c 16`;
	mv $Setup_Path/databaseAdmin $Setup_Path/phpmyadmin_$phpmyadminExt
	chmod -R 755 $Setup_Path/phpmyadmin_$phpmyadminExt
	chown -R www.www $Setup_Path/phpmyadmin_$phpmyadminExt

	secret=`cat /dev/urandom | head -n 32 | md5sum | head -c 32`;
	\cp -a -r $Setup_Path/phpmyadmin_$phpmyadminExt/config.sample.inc.php  $Setup_Path/phpmyadmin_$phpmyadminExt/config.inc.php
	sed -i "s#^\$cfg\['blowfish_secret'\].*#\$cfg\['blowfish_secret'\] = '${secret}';#" $Setup_Path/phpmyadmin_$phpmyadminExt/config.inc.php
	sed -i "s#^\$cfg\['blowfish_secret'\].*#\$cfg\['blowfish_secret'\] = '${secret}';#" $Setup_Path/phpmyadmin_$phpmyadminExt/libraries/config.default.php

	echo $1 > $Setup_Path/version.pl

	PHPVersion=""
	for phpVer in 52 53 54 55 56 70 71 72 73 74 80 81;
	do
		if [ -d "/www/server/php/${phpVer}/bin" ]; then
			PHPVersion=${phpVer}
		fi
	done


	sed -i "s#$Root_Path/wwwroot/default#$Root_Path/server/phpmyadmin#" $Root_Path/server/nginx/conf/nginx.conf
	rm -f $Root_Path/server/nginx/conf/enable-php.conf
	\cp $Root_Path/server/nginx/conf/enable-php-$PHPVersion.conf $Root_Path/server/nginx/conf/enable-php.conf
	sed -i "/pathinfo/d" $Root_Path/server/nginx/conf/enable-php.conf
	/etc/init.d/nginx reload

}

Uninstall_phpMyAdmin()
{
	rm -rf $Root_Path/server/phpmyadmin/phpmyadmin*
	rm -f $Root_Path/server/phpmyadmin/version.pl
	rm -f $Root_Path/server/phpmyadmin/version_check.pl
}

actionType=$1
version=$2

if [ "$actionType" == 'install' ];then
	Install_phpMyAdmin $version
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_phpMyAdmin
fi
