#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8


header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


runPath="/root"
Is_64bit=`getconf LONG_BIT`

extFile(){
	case "${version}" in 
		'52')
		extFile='/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/ixed.lin'
		;;
		'53')
		extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/ixed.lin'
		;;
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/ixed.lin'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/ixed.lin'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/ixed.lin'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/ixed.lin'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/ixed.lin'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/ixed.lin'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/ixed.lin'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/ixed.lin'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/ixed.lin'
		;;
		'81')
		extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/ixed.lin'
		;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/ixed.lin'
		;;
	esac;
}
Install_sg11()
{
	extFile
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep ${extFile}`
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装sg11,请选择其它版本!"
		return
	fi

	if [ ! -f "$extFile" ];then

		SysType=$(uname -m)
		sed -i '/ixed/d'  /www/server/php/$version/etc/php.ini
		
		case "${version}" in 
			'52')
			mkdir -p /www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'53')
			mkdir -p /www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'54')
			mkdir -p /www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'55')
			mkdir -p /www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'56')
			mkdir -p /www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'70')
			mkdir -p /www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'71')
			mkdir -p /www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'72')
			mkdir -p /www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'73')
			mkdir -p /www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'74')
			mkdir -p /www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'80')
			mkdir -p /www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin 
			;;
			'81')
			mkdir -p /www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin
			;;
			'82')
			mkdir -p /www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/
			wget -O ${extFile} https://download.bt.cn/src/sg11/${Is_64bit}/${SysType}/ixed.${vphp}.lin
			;;
		esac;

	fi

	if [ ! -f "$extFile" ];then
		echo "ERROR!"
		return;
	fi

	echo "extension=$extFile" >> /www/server/php/$version/etc/php.ini
	service php-fpm-$version reload
	echo '==========================================================='
	echo 'successful!'

}
Uninstall_sg11()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		return;
	fi

	extFile

	if [ ! -f "$extFile" ];then
		echo "php-$vphp 未安装sg11,请选择其它版本!"
		return
	fi

	sed -i '/ixed/d'  /www/server/php/$version/etc/php.ini
		
	rm -f $extFile
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}
actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_sg11
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_sg11
fi


