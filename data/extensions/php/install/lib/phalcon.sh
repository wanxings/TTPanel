#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL




Is_64bit=`getconf LONG_BIT`
run_path="/root"


Ext_Path(){
	case "${version}" in 
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/phalcon.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/phalcon.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/phalcon.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/phalcon.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/phalcon.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/phalcon.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/phalcon.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/phalcon.so'
		;;
	esac
}
Install_Psr(){
	case "${version}" in 
		'72')
		PsrFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/psr.so'
		;;
		'73')
		PsrFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/psr.so'
		;;
		'74')
		PsrFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/psr.so'
		;;
	esac

	if [ ! -f "${PsrFile}" ];then
		psrVersion="1.0.0"
		wget -O psr-${psrVersion}.tgz $download_Url/extensions/php/install/lib/psr-${psrVersion}.tgz
		tar -xvf psr-${psrVersion}.tgz
		cd psr-${psrVersion}
		/www/server/php/${version}/bin/phpize
		./configure --with-php-config=/www/server/php/${version}/bin/php-config
		make -j${cpuCore}
		make install
		if [ ! -f "${PsrFile}" ];then
			echo '========================================================'
			GetSysInfo
			echo -e "psr installation failed.";
			echo -e "安装失败，请截图以上报错信息发帖至论坛www.bt.cn/bbs求助"
			exit 1;
		else
			local isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'psr.so')
			if [ -z "${isInstall}" ];then
				echo "extension = psr.so" >> /www/server/php/$version/etc/php.ini
			fi
		fi
		cd ..
		rm -rf psr-${psrVersion}
	fi 
}
Install_Phalcon(){

	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi

	local isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'phalcon.so'`
	if [ "${isInstall}" ];then
		echo "php-$vphp 已安装过phalcon,请选择其它版本!"
		echo "php-$vphp is already installed phpalcon, Plese select other version!"
		return
	fi

	if [ ! -f "${extFile}" ];then
		cd ${run_path}
		if [ "${version}" -ge "72" ]; then
			Install_Psr
			phalconVer="4.0.5"
			wget -O phalcon-${phalconVer}.tgz $download_Url/extensions/php/install/lib/phalcon-${phalconVer}.tgz
			tar zxvf phalcon-${phalconVer}.tgz
			rm -f phalcon-${phalconVer}.tgz
			cd phalcon-${phalconVer}
		else
			phalconVer="3.4.4"
			wget $download_Url/extensions/php/install/lib/cphalcon-${phalconVer}.tar.gz
			tar zxvf cphalcon-${phalconVer}.tar.gz
			rm -f cphalcon-${phalconVer}.tar.gz
			if [ "${version:0:1}" == "5" ]; then
				phpV="php5"
			else
				phpV="php7"
			fi
			cd cphalcon-${phalconVer}/build/${phpV}/${Is_64bit}bits
		fi
		/www/server/php/${version}/bin/phpize
		./configure --with-php-config=/www/server/php/${version}/bin/php-config --enable-phalcon
		make -j${cpuCore}
		make install
		cd ${run_path}
		rm -rf cphalcon-${phalconVer}
		rm -rf phalcon-${phalconVer}
	fi

	if [ ! -f "${extFile}" ];then
		echo '========================================================'
		GetSysInfo
		echo -e "Phalcon installation failed.";
		echo -e "安装失败，请截图以上报错信息发帖至论坛www.bt.cn/bbs求助"
		exit 1;
	fi
	echo -e "extension = ${extFile}\n" >> /www/server/php/$version/etc/php.ini
	service php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}
Uninstall_Phalcon(){
	sed -i '/phalcon.so/d' /www/server/php/$version/etc/php.ini
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
	Install_Phalcon
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Phalcon
fi

