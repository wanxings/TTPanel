#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL


Install_Exif()
{
    isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'exif.so')
    if [ "${isInstall}" != "" ];then
        echo "php-$vphp 已安装exif,请选择其它版本!"
        return
    fi
    
    if [ ! -d "/www/server/php/$version/src/ext/exif" ];then
        mkdir -p /www/server/php/$version/src
		wget -O $version-ext.tar.gz $download_Url/extensions/php/install/lib/$version-ext.tar.gz
		tar -zxf $version-ext.tar.gz -C /www/server/php/$version/src/ 
		rm -f $version-ext.tar.gz
    fi
    
    case "${version}" in 
        '53')
        extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/exif.so'
        ;;
        '54')
        extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/exif.so'
        ;;
        '55')
        extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/exif.so'
        ;;
        '56')
        extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/exif.so'
        ;;
        '70')
        extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/exif.so'
        ;;
        '71')
        extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/exif.so'
        ;;
        '72')
        extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/exif.so'
        ;;
        '73')
        extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/exif.so'
        ;;
        '74')
        extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/exif.so'
        ;;
        '80')
        extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/exif.so'
        ;;
        '81')
        extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/exif.so'
        ;;
        '82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/exif.so'
		;;
    esac

    Centos7Check=$(cat /etc/redhat-release|grep ' 7.'|grep -i centos)
    if [ "${Centos7Check}" ] && [ "${version}" == "80" ];then
        yum install centos-release-scl-rh -y
        yum install devtoolset-7-gcc devtoolset-7-gcc-c++ -y
        yum install cmake3 -y
        cmakeV="cmake3"
        export CC=/opt/rh/devtoolset-7/root/usr/bin/gcc
        export CXX=/opt/rh/devtoolset-7/root/usr/bin/g++
    fi
    
    if [ ! -f "$extFile" ];then
        cd /www/server/php/$version/src/ext/exif
        /www/server/php/$version/bin/phpize
        ./configure --with-php-config=/www/server/php/$version/bin/php-config
        if [ "${version}" -ge "80" ];then
            sed -i "s#CFLAGS = -g -O2#CFLAGS = -std=c99 -g -O2#g" Makefile
        fi
        make && make install
    fi
    
    if [ ! -f "$extFile" ];then
        echo "ERROR!"
        return;
    fi
    
    echo -e "extension = $extFile" >> /www/server/php/$version/etc/php.ini
    if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        echo -e "extension = $extFile" >> /www/server/php/$version/etc/php-cli.ini
    fi
    service php-fpm-$version reload
    echo '==========================================================='
    echo 'successful!'
}


Uninstall_Exif()
{
    if [ ! -f "/www/server/php/$version/bin/php-config" ];then
        echo "php-$vphp 未安装,请选择其它版本!"
        return
    fi
    isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'exif.so')
    if [ "${isInstall}" == "" ];then
        echo "php-$vphp 未安装exif,请选择其它版本!"
        return
    fi
    
    sed -i '/exif.so/d' /www/server/php/$version/etc/php.ini
    if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        sed -i '/exif.so/d' /www/server/php/$version/etc/php-cli.ini
    fi
    service php-fpm-$version reload
    echo '==============================================='
    echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
    Install_Exif
elif [ "$actionType" == 'uninstall' ];then
    Uninstall_Exif
fi
