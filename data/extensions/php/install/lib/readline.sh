#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL


Install_Readline() {
    if [ ! -f "/www/server/php/$version/bin/php-config" ]; then
        echo "php-$vphp 未安装,请选择其它版本!"
        echo "php-$vphp not install, Plese select other version!"
        return
    fi

    isInstall=$(cat /www/server/php/$version/etc/php.ini | grep 'readline.so')
    if [ "${isInstall}" != "" ]; then
        echo "php-$vphp 已安装过readline,请选择其它版本!"
        echo "php-$vphp is installed readline, Plese select other version!"
        return
    fi
    
    if [ ! -d "/www/server/php/$version/src/ext/readline" ]; then
        mkdir -p /www/server/php/$version/src
		wget -O $version-ext.tar.gz $download_Url/extensions/php/install/lib/$version-ext.tar.gz
		tar -zxf $version-ext.tar.gz -C /www/server/php/$version/src/ 
		rm -f $version-ext.tar.gz
    fi

    case "${version}" in
    '52')
        extFile="/www/server/php/52/lib/php/extensions/no-debug-non-zts-20060613/readline.so"
        ;;
    '53')
        extFile="/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/readline.so"
        ;;
    '54')
        extFile="/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/readline.so"
        ;;
    '55')
        extFile="/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/readline.so"
        ;;
    '56')
        extFile="/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/readline.so"
        ;;
    '70')
        extFile="/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/readline.so"
        ;;
    '71')
        extFile="/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/readline.so"
        ;;
    '72')
        extFile="/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/readline.so"
        ;;
    '73')
        extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/readline.so'
        ;;
    '74')
        extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/readline.so'
        ;;
    '80')
        extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/readline.so'
        ;;
    '81')
        extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/readline.so'
        ;;
    esac

    if [ ! -f "${extFile}" ]; then
        if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ]; then
            Pack="libedit-devel"
            ${PM} install ${Pack} -y
        elif [ "${PM}" == "apt-get" ]; then
            Pack="libedit-dev"
            ${PM} install ${Pack} -y
        fi
        cd /www/server/php/$version/src/ext/readline
        /www/server/php/$version/bin/phpize
        ./configure --with-php-config=/www/server/php/$version/bin/php-config
        make && make install
    fi

    if [ ! -f "${extFile}" ]; then
        echo 'error'
        exit 0
    fi

    echo -e "extension = " ${extFile} >>/www/server/php/$version/etc/php.ini
    if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        echo -e "extension = $extFile" >> /www/server/php/$version/etc/php-cli.ini
    fi

    service php-fpm-$version reload
    echo '==============================================='
    echo 'successful!'
}

Uninstall_Readline() {
    if [ ! -f "/www/server/php/$version/bin/php-config" ]; then
        echo "php-$vphp 未安装,请选择其它版本!"
        echo "php-$vphp not install, Plese select other version!"
        return
    fi

    isInstall=$(cat /www/server/php/$version/etc/php.ini | grep 'readline.so')
    if [ "${isInstall}" = "" ]; then
        echo "php-$vphp 未安装readline,请选择其它版本!"
        echo "php-$vphp not install readline, Plese select other version!"
        return
    fi

    sed -i '/readline.so/d' /www/server/php/$version/etc/php.ini
    if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        sed -i '/readline.so/d' /www/server/php/$version/etc/php-cli.ini
    fi

    service php-fpm-$version reload
    echo '==============================================='
    echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ]; then
    Install_Readline
elif [ "$actionType" == 'uninstall' ]; then
    Uninstall_Readline
fi
