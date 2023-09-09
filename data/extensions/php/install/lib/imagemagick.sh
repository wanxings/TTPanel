#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8

header_file=/www/panel/data/shell/install_header.sh
. $header_file

download_Url=$NODE_URL



Centos8Check=$(cat /etc/redhat-release | grep ' 8.' | grep -i centos)
if [ "${Centos8Check}" ];then
	dnf config-manager --set-enabled PowerTools
fi
extPath(){
	case "${version}" in 
		'54')
		extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/imagick.so'
		;;
		'55')
		extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/imagick.so'
		;;
		'56')
		extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/imagick.so'
		;;
		'70')
		extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/imagick.so'
		;;
		'71')
		extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/imagick.so'
		;;
		'72')
		extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/imagick.so'
		;;
		'73')
		extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/imagick.so'
		;;
		'74')
		extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/imagick.so'
		;;
		'80')
		extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/imagick.so'
		;;
		'81')
        extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/imagick.so'
        ;;
		'82')
		extFile='/www/server/php/82/lib/php/extensions/no-debug-non-zts-20220829/imagick.so'
		;;
	esac
}
Install_imagemagick()
{
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'imagick.so')
	if [ "${isInstall}" != "" ];then
		echo "php-$vphp 已安装过imagemagick,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ];then
		Pack="ImageMagick ImageMagick-devel"
	elif [ "${PM}" == "apt-get" ];then
		Pack="imagemagick libmagickwand-dev libmagick++-dev"
	fi
	${PM} install ${Pack} -y

	if [ "${PM}" == "yum" ];then
		centos7Check=$(cat /etc/redhat-release | grep ' 7.' | grep -i centos)
		if [ "${centos7Check}" ] || [ ! -f "/usr/bin/MagickWand-config" ];then
			if [ ! -f "/usr/local/ImageMagick-7.1.0/bin/MagickWand-config" ];then
				ImageMagickVer="7.1.0-16"
				wget -O ImageMagick-${ImageMagickVer}.tar.gz $download_Url/extensions/php/install/lib/ImageMagick-${ImageMagickVer}.tar.gz
				tar -xvf ImageMagick-${ImageMagickVer}.tar.gz
				cd ImageMagick-${ImageMagickVer}
				./configure --prefix=/usr/local/ImageMagick-7.1.0
				make
				make install
				cd ..
				rm -rf ImageMagick-${ImageMagickVer}*
			fi
			ImageMagick_DIR="--with-imagick=/usr/local/ImageMagick-7.1.0"
		fi
	fi

	if [ ! -f "$extFile" ];then
		rm -rf imagick*
		if [ "${version}" -ge "80" ];then
			wget $download_Url/extensions/php/install/lib/imagick-3.7.0.tgz
			tar -xvf imagick-3.7.0.tgz
			cd imagick-3.7.0 || exit
		else
			wget $download_Url/extensions/php/install/lib/imagick-3.4.4.tgz -T 5
			tar -zxf imagick-3.4.4.tgz
			cd imagick-3.4.4
		fi
		/www/server/php/$version/bin/phpize
		./configure --with-php-config=/www/server/php/$version/bin/php-config ${ImageMagick_DIR}
		make && make install
	fi
	
	if [ ! -f "$extFile" ];then
		echo 'error';
		exit 0;
	fi
	
	
	echo -e "\n[ImageMagick]\nextension = \"imagick.so\"\n" >> /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		echo -e "\n[ImageMagick]\nextension = \"imagick.so\"\n" >> /www/server/php/$version/etc/php-cli.ini
	fi

	cd ../
	rm -rf imagick-*
	/etc/init.d/php-fpm-$version reload
}


Uninstall_imagemagick()
{
	extPath
	if [ ! -f "/www/server/php/$version/bin/php-config" ];then
		echo "php-$vphp 未安装,请选择其它版本!"
		echo "php-$vphp not install, Plese select other version!"
		return
	fi
	
	isInstall=$(cat /www/server/php/$version/etc/php.ini|grep 'imagick.so')
	if [ "${isInstall}" = "" ];then
		echo "php-$vphp 未安装imagemagick,请选择其它版本!"
		echo "php-$vphp not install imagemagick, Plese select other version!"
		return
	fi
	
	sed -i '/imagick.so/d' /www/server/php/$version/etc/php.ini
	sed -i '/ImageMagick/d' /www/server/php/$version/etc/php.ini
	if [ -f /www/server/php/$version/etc/php-cli.ini ];then
		sed -i '/imagick.so/d' /www/server/php/$version/etc/php-cli.ini
		sed -i '/ImageMagick/d' /www/server/php/$version/etc/php-cli.ini
	fi
	rm -f ${extFile}
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}


actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
extPath
if [ "$actionType" == 'install' ];then
	Install_imagemagick
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_imagemagick
fi

