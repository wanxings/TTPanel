#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL



smbclient_Version='1.0.6'

expPath(){
  case "${version}" in 
    '53')
    extFile='/www/server/php/53/lib/php/extensions/no-debug-non-zts-20090626/smbclient.so'
    ;;
    '54')
    extFile='/www/server/php/54/lib/php/extensions/no-debug-non-zts-20100525/smbclient.so'
    ;;
    '55')
    extFile='/www/server/php/55/lib/php/extensions/no-debug-non-zts-20121212/smbclient.so'
    ;;
    '56')
    extFile='/www/server/php/56/lib/php/extensions/no-debug-non-zts-20131226/smbclient.so'
    ;;
    '70')
    extFile='/www/server/php/70/lib/php/extensions/no-debug-non-zts-20151012/smbclient.so'
    ;;
    '71')
    extFile='/www/server/php/71/lib/php/extensions/no-debug-non-zts-20160303/smbclient.so'
    ;;
    '72')
    extFile='/www/server/php/72/lib/php/extensions/no-debug-non-zts-20170718/smbclient.so'
    ;;
    '73')
    extFile='/www/server/php/73/lib/php/extensions/no-debug-non-zts-20180731/smbclient.so'
    ;;
	  '74')
    extFile='/www/server/php/74/lib/php/extensions/no-debug-non-zts-20190902/smbclient.so'
    ;;
    '80')
    extFile='/www/server/php/80/lib/php/extensions/no-debug-non-zts-20200930/smbclient.so'
    ;;
    '81')
    extFile='/www/server/php/81/lib/php/extensions/no-debug-non-zts-20210902/smbclient.so'
    ;;
  esac
}

Install_smbclient()
{

  if [ ! -f "/www/server/php/$version/bin/php-config" ];then
    echo "php-$vphp 未安装,请选择其它版本!"
    echo "php-$vphp not install, Plese select other version!"
    return
  fi
  
  isInstall=`cat /www/server/php/$version/etc/php.ini|grep 'smbclient.so'`
  if [ "${isInstall}" != "" ];then
    echo "php-$vphp 已安装过smbclient,请选择其它版本!"
    echo "php-$vphp smbclient has been installed, Plese select other version!"
    return
  fi

  runPath=/root
  expPath

  if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ];then
    redhat_version_file="/etc/redhat-release"
    os_version=$(cat $redhat_version_file|grep Red|grep -Eo '([0-9]+\.)+[0-9]+'|grep -Eo '^[0-9]')
    is_oracle=`cat /etc/os-release|awk '{FS=":"}{print $3}'|grep oracle`
    is_aarch64=$(uname -a|grep aarch64)

    if [ "$os_version" == "8" ] && [ "$is_oracle" == "oracle" ];then
      el_ver="el8_5"
      lib_ver="4.14.5-7"
      if [ "$is_aarch64" != "" ];then
        os_ver="aarch64"
        wget -O libsmbclient.rpm $download_Url/extensions/php/install/lib/libsmbclient-$lib_ver.$el_ver.$os_ver.rpm
        wget -O libsmbclient-devel.rpm $download_Url/extensions/php/install/lib/libsmbclient-devel-$lib_ver.$el_ver.$os_ver.rpm
        rpm -ivh libsmbclient.rpm --nodeps
        rpm -ivh libsmbclient-devel.rpm --nodeps
      else
        os_ver="x86_64"
        wget -O libsmbclient.rpm $download_Url/extensions/php/install/lib/libsmbclient-$lib_ver.$el_ver.$os_ver.rpm
        wget -O libsmbclient-devel.rpm $download_Url/extensions/php/install/lib/libsmbclient-devel-$lib_ver.$el_ver.$os_ver.rpm
        rpm -ivh libsmbclient.rpm --nodeps
        rpm -ivh libsmbclient-devel.rpm --nodeps
      fi
      rm -rf libsmbclient*
    elif [ "$os_version" == "7" ] && [ "$is_oracle" == "oracle" ];then
      el_ver="el7"
      lib_ver="4.10.16-5"
      if [ "$is_aarch64" != "" ];then
        os_ver="aarch64"
        wget -O libsmbclient.rpm $download_Url/extensions/php/install/lib/libsmbclient-$lib_ver.$el_ver.$os_ver.rpm
        wget -O libsmbclient-devel.rpm $download_Url/extensions/php/install/lib/libsmbclient-devel-$lib_ver.$el_ver.$os_ver.rpm
        rpm -ivh libsmbclient.rpm --nodeps
        rpm -ivh libsmbclient-devel.rpm --nodeps
      else
        os_ver="x86_64"
        wget -O libsmbclient.rpm $download_Url/extensions/php/install/lib/libsmbclient-$lib_ver.$el_ver.$os_ver.rpm
        wget -O libsmbclient-devel.rpm $download_Url/extensions/php/install/lib/libsmbclient-devel-$lib_ver.$el_ver.$os_ver.rpm
        rpm -ivh libsmbclient.rpm --nodeps
        rpm -ivh libsmbclient-devel.rpm --nodeps
      fi
      rm -rf libsmbclient*
    else
      Pack="libsmbclient libsmbclient-devel"
		  ${PM} install ${Pack} -y
    fi
	elif [ "${PM}" == "apt-get" ];then
		Pack="libsmbclient-dev"
		${PM} install ${Pack} -y
	fi

  if [ ! -f "${extFile}" ];then
  	wget $download_Url/extensions/php/install/lib/smbclient-$smbclient_Version.tgz
  	tar -zxvf smbclient-$smbclient_Version.tgz
  	cd smbclient-$smbclient_Version
  	/www/server/php/$version/bin/phpize
  	./configure --with-php-config=/www/server/php/$version/bin/php-config --with-smbclient
  	make && make install
  	cd ../
  	rm -rf smbclient*
 	fi

 	if [ ! -f "${extFile}" ];then
 		echo 'error';
 		exit 0;
 	fi
   	echo -e "\n[smbclient]\nextension = smbclient.so\n" >> /www/server/php/$version/etc/php.ini
    if [ -f /www/server/php/$version/etc/php-cli.ini ];then
        echo -e "\n[smbclient]\nextension = smbclient.so\n" >> /www/server/php/$version/etc/php-cli.ini
    fi

    /etc/init.d/php-fpm-$version reload
}

Uninstall_smbclient()
{
  expPath
	sed -i '/smbclient/d' /www/server/php/$version/etc/php.ini
  if [ -f /www/server/php/$version/etc/php-cli.ini ];then
    sed -i '/smbclient.so/d' /www/server/php/$version/etc/php-cli.ini
  fi
  rm -f ${extFile}
	/etc/init.d/php-fpm-$version reload
	echo '==============================================='
	echo 'successful!'
}

actionType=$1
version=$2
vphp=${version:0:1}.${version:1:1}
if [ "$actionType" == 'install' ];then
	Install_smbclient
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_smbclient
fi