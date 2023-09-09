#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH


header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


Swoole_Version='4.5.11'
runPath=/root

extPath()
{
 	case "${version}" in 
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/swoole.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/swoole.so'
		;;	
  		'72')
  		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/swoole.so'
  		;;
  		'73')
  		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/swoole.so'
  		;;
  		'74')
  		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/swoole.so'
  		;;
  		'80')
  		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/swoole.so'
  		;;
      '81')
      extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/swoole.so'
      ;;
	esac
}


Install_Swoole()
{
    if [ ! -f "/www/server/php/$version/bin/php-config" ];then
      echo "php-$vphp 未安装,请选择其它版本!"
      echo "php-$vphp not install, Plese select other version!"
      return
    fi
    
    isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'swoole.so'`
    if [ "${isInstall}" != "" ];then
      echo "php-$vphp 已安装过swoole,请选择其它版本!"
      echo "php-$vphp not install, Plese select other version!"
      return
    fi
  	extPath
	if [ ! -f "${extFile}" ];then 
		Swoole_Version="4.8.12"
		if [ "$version" -ge "80" ];then
			Swoole_Version="5.0.1"
		fi
		cd ${runPath}
  		wget $download_Url/extensions/php/install/lib/swoole-$Swoole_Version.tgz
  		tar -zxvf swoole-$Swoole_Version.tgz
  		cd swoole-$Swoole_Version
  		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config --enable-openssl --with-openssl-dir=/usr/local/openssl --enable-sockets
		make && make install
        cd ../
        rm -rf swoole*
	fi

 	if [ ! -f "${extFile}" ];then
 		echo 'error';
 		exit 0;
 	fi
 	
 	echo -e "\n[swoole]\nextension = swoole.so\n" >> /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        echo -e "\n[swoole]\nextension = swoole.so\n" >> /www/server/php/$version/etc/php-cli.ini
    fi
 	service php-fpm-$version reload
}



Uninstall_Swoole()
{
	extPath
	sed -i '/swoole/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        sed -i '/swoole/d' /www/server/php/$version/etc/php-cli.ini
    fi
  	rm -f ${extFile}
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_Swoole
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Swoole
fi
