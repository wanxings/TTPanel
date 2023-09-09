#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
Ext_Path(){
    case "${version}" in 
        '54')
        extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/igbinary.so'
        ;;
        '55')
        extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/igbinary.so'
        ;;
        '56')
        extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/igbinary.so'
        ;;
        '70')
        extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/igbinary.so'
        ;;
        '71')
        extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/igbinary.so'
        ;;
        '72')
        extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/igbinary.so'
        ;;
        '73')
        extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/igbinary.so'
        ;;
        '74')
        extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/igbinary.so'
        ;;
        '80')
        extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/igbinary.so'
        ;;
        '81')
        extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/igbinary.so'
        ;;
        '82')
        extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/igbinary.so'
        ;;
    esac
}
Install_igbinary(){

    if [ ! -f "/www/server/php/$version/bin/php-config" ];then
        echo "php-$vphp 未安装,请选择其它版本!"
        echo "php-$vphp not install, Plese select other version!"
        return
    fi

    isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'igbinary.so'`
    if [ "${isInstall}" != "" ];then
        echo "php-$vphp 已安装过igbinary,请选择其它版本!"
        echo "php-$vphp is already installed igbinary, Plese select other version!"
        exit
    fi
    
    if [ ! -f "${extFile}" ];then
        if [ "${version}" -ge 70 ];then
            igbinaryVer="3.2.12"
        else 
            igbinaryVer="2.0.8"
        fi
        wget $download_Url/extensions/php/install/lib/igbinary-${igbinaryVer}.tgz -T 5
        tar zxvf igbinary-${igbinaryVer}.tgz
        rm -f igbinary-${igbinaryVer}.tgz
        cd igbinary-${igbinaryVer}
        /www/server/php/${version}/bin/phpize
        ./configure --with-php-config=/www/server/php/${version}/bin/php-config
        make && make install
        cd ../
        rm -rf igbinary-${igbinaryVer}
    fi

    if [ ! -f "${extFile}" ];then
        echo 'error';
        exit 0;
    fi
    echo -e "extension = ${extFile}\n" >> /www/server/php/$version/etc/php.ini
    if [[ -f /www/server/php/$version/etc/php-cli.ini ]];then
		echo -e "extension = " ${extFile} >> /www/server/php/$version/etc/php-cli.ini
	fi
    /etc/init.d/php-fpm-$version reload
    echo '==============================================='
    echo 'successful!'
}
Uninstall_igbinary(){
    if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'sysvmsg.so'`
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装sysvmsg,请选择其它版本!"
		echo "php-$vphp not install sysvmsg, Plese select other version!"
		return
	fi

	sed -i '/igbinary.so/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        sed -i '/igbinary.so/d' /www/server/php/$version/etc/php-cli.ini
    fi
    
    rm -f ${extFile}

    /etc/init.d/php-fpm-$version reload
    echo '==============================================='
    echo 'successful!'
}
if [ "$actionType" == 'install' ];then
    Ext_Path
    Install_igbinary
elif [ "$actionType" == 'uninstall' ];then
    Ext_Path
    Uninstall_igbinary
fi
