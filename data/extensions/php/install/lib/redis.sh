#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
redis_version=7.0.5
runPath=/root
header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL

if [ -z "${cpuCore}" ]; then
	cpuCore="1"
fi

Error_Msg(){
	if [ "${actionType}" == "install" ];then
		AC_TYPE="安装"
	elif [ "${actionType}" == "update" ]; then
		AC_TYPE="升级"
	fi

	echo '========================================================'
	GetSysInfo
	echo -e "ERROR: redis-${redis_version} ${actionType} failed.";
	exit 1;
}

System_Lib(){
	if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ] ; then
		Pack="sudo"
		${PM} install ${Pack} -y
	elif [ "${PM}" == "apt-get" ]; then
		Pack="sudo"
		${PM} install ${Pack} -y
	fi

}
Service_Add(){
	if [ -f "/usr/bin/yum" ];then
		chkconfig --add redis
		chkconfig --level 2345 redis on
	elif [ -f "/usr/bin/apt" ]; then
		apt-get install sudo -y	
		update-rc.d redis defaults
	fi
}
Service_Del(){
	if [ -f "/usr/bin/yum" ];then
		chkconfig --level 2345 redis off
	elif [ -f "/usr/bin/apt" ]; then
		update-rc.d redis remove
	fi
}
Gcc_Version_Check(){
	if [ "${PM}" == "yum" ];then
		Centos7Check=$(cat /etc/redhat-release | grep ' 7.' | grep -iE 'centos|Red Hat')
		gccV=$(gcc -v 2>&1|grep "gcc version"|awk '{printf("%d",$3)}')
		sysType=$(uname -a|grep x86_64)
		armType=$(uname -a|grep aarch64)
		if [ "${Centos7Check}" ];then
			yum install centos-release-scl-rh -y
			yum install devtoolset-7-gcc devtoolset-7-gcc-c++ -y
			if [ "${armType}" ];then
				yum install devtoolset-7-libatomic-devel -y
			fi
			if [ -f "/opt/rh/devtoolset-7/root/usr/bin/gcc" ] && [ "${sysType}" != "${armType}" ];then
				export CC=/opt/rh/devtoolset-7/root/usr/bin/gcc
			else
				redis_version="5.0.8"
			fi
		elif [ "${gccV}" -le "5" ];then
			redis_version="5.0.8"
		fi
	fi
}
ext_Path(){
	case "${version}" in 
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/redis.so'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/redis.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/redis.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/redis.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/redis.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/redis.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/redis.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/redis.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/redis.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/redis.so'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/redis.so'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/redis.so'
		;;
	esac
}
Install_Redis()
{
	
	if [ ! -d /www/server/php/$version ];then
		return;
	fi
	
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'redis.so'`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过Redis,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	

	if [ ! -f "${extFile}" ];then		
		if [ "${version}" == "52" ];then
			rVersion='2.2.7'
		elif [ "${version}" -ge "70" ];then
			rVersion='5.3.7'
		else
			rVersion='4.3.0'
		fi
		
		wget $download_Url/extensions/php/install/lib/src/redis-$rVersion.tgz -T 5
		tar -xzvf redis-$rVersion.tgz
		rm -f redis-$rVersion.tgz
		cd redis-$rVersion
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config
		make && make install
		cd ../
		rm -rf redis-$rVersion*
	fi
	
	if [ ! -f "${extFile}" ];then
		echo 'error';
		exit 0;
	fi
	
	echo -e "\n[redis]\nextension = ${extFile}\n" >> /www/server/php/$version/etc/php.ini
	if [ -f "/www/server/php/$version/etc/php-cli.ini" ]; then
		echo -e "\n[redis]\nextension = ${extFile}\n" >> /www/server/php/$version/etc/php-cli.ini
	fi

	/etc/init.d/php-fpm-$version restart
	echo '==============================================='
	echo 'successful!'
}

Uninstall_Redis()
{
	pkill -9 redis
	rm -f /var/run/redis_6379.pid
	Service_Del
	rm -f /usr/bin/redis-cli
	rm -f /etc/init.d/redis
	rm -rf /www/server/redis
	rm -rf /www/server/panel/plugin/redis
	
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'redis.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装Redis,请选择其它版本!"
		echo "php-$vphp not install Redis, Plese select other version!"
		return
	fi
	
	sed -i '/redis.so/d' /www/server/php/$version/etc/php.ini
	sed -i '/\[redis\]/d' /www/server/php/$version/etc/php.ini
	if [ -f "/www/server/php/$version/etc/php-cli.ini" ]; then
		sed -i '/redis.so/d' /www/server/php/$version/etc/php-cli.ini
		sed -i '/\[redis\]/d' /www/server/php/$version/etc/php-cli.ini
	fi
	
	/etc/init.d/php-fpm-$version restart
	echo '==============================================='
	echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}

if [ "$version" == "7.0" ];then
	redis_version="7.0.5"
elif [ "$version" == "6.2" ]; then
	redis_version="6.2.7"
fi

if [ "$actionType" == 'install' ];then
	System_Lib
	ext_Path
	Gcc_Version_Check
	Install_Redis
	Service_Add
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Redis
fi


