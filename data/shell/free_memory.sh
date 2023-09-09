#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

endDate=`date +"%Y-%m-%d %H:%M:%S"`
log="释放内存!"
echo "★[$endDate] $log"
echo '----------------------------------------------------------------------------'

if [ -f "/etc/init.d/php-fpm-52" ];then
	/etc/init.d/php-fpm-52 reload
fi

if [ -f "/etc/init.d/php-fpm-53" ];then
	/etc/init.d/php-fpm-53 reload
fi

if [ -f "/etc/init.d/php-fpm-54" ];then
	/etc/init.d/php-fpm-54 reload
fi

if [ -f "/etc/init.d/php-fpm-55" ];then
	/etc/init.d/php-fpm-55 reload
fi

if [ -f "/etc/init.d/php-fpm-56" ];then
	/etc/init.d/php-fpm-56 reload
fi

if [ -f "/etc/init.d/php-fpm-70" ];then
	/etc/init.d/php-fpm-70 reload
fi

if [ -f "/etc/init.d/php-fpm-71" ];then
	/etc/init.d/php-fpm-71 reload
fi

if [ -f "/etc/init.d/php-fpm-72" ];then
	/etc/init.d/php-fpm-72 reload
fi

if [ -f "/etc/init.d/php-fpm-73" ];then
	/etc/init.d/php-fpm-73 reload
fi

if [ -f "/etc/init.d/php-fpm-74" ];then
	/etc/init.d/php-fpm-74 reload
fi

if [ -f "/etc/init.d/mysqld" ];then
	/etc/init.d/mysqld reload
fi

if [ -f "/etc/init.d/nginx" ];then
	/etc/init.d/nginx reload
fi

sync
sleep 2
sync
echo 3 > /proc/sys/vm/drop_caches

echo '----------------------------------------------------------------------------'