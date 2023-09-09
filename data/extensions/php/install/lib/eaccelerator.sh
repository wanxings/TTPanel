#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL


extensionName='eaccelerator'
extensionVersion='0.9.6.1'

Install_eaccelerator()
{
	

    runPath=/root

    case "${version}" in 
    	'52')
    	extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/eaccelerator.so'
    	;;
    	'53')
    	extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/eaccelerator.so'
    	;;
    	'54')
    	extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/eaccelerator.so'
    	;;
    esac

    if [ ! -f "${extFile}" ];then
    	wget $download_Url/extensions/php/install/lib/$extensionName-$extensionVersion.tar.gz
    	tar -zxvf $extensionName-$extensionVersion.tar.gz
    	cd eaccelerator-eaccelerator-42067ac
    	/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ../
		rm -rf $extensionName*
   	fi

   	if [ ! -f "${extFile}" ];then
   		echo ${extFile};
   		echo 'error';
   		exit 0;
   	fi
   		
   	cat >> /www/server/php/$version/etc/php.ini<<EOF
[eaccelerator]
extension=eaccelerator.so
eaccelerator.cache_dir="/tmp/eaccelerator"
eaccelerator.shm_size="8" 
eaccelerator.enable="1"
eaccelerator.optimizer="1"
eaccelerator.check_mtime="1"
eaccelerator.debug="0" 
eaccelerator.filter=""
eaccelerator.shm_max="0"
eaccelerator.shm_ttl="0"
eaccelerator.shm_prune_period="0"
eaccelerator.shm_only="0"
eaccelerator.compress="1"
eaccelerator.compress_level="9"
eaccelerator.keys = "disk_only"
eaccelerator.sessions = "disk_only"
eaccelerator.content = "disk_only"
EOF
   	service php-fpm-$version reload
}



Uninstall_eaccelerator()
{
	sed -i '/eaccelerator/d' /www/server/php/$version/etc/php.ini	
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_$extensionName
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_$extensionName
fi
