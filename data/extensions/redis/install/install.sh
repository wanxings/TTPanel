#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL

System_Lib(){
	if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ] ; then
		Pack="sudo"
		${PM} install ${Pack} -y
	elif [ "${PM}" == "apt-get" ]; then
		Pack="sudo"
		${PM} install ${Pack} -y
	fi

}
Install_Redis()
{
	groupadd redis
	useradd -g redis -s /sbin/nologin redis
	if [ ! -f '/www/server/redis/src/redis-server' ];then
		rm -rf /www/server/redis
		cd /www/server

		wget -O redis_$redis_version.tar.gz $download_Url/extensions/redis/redis_$redis_version.tar.gz
		tar zxvf redis_$redis_version.tar.gz
		cd redis
		make -j ${cpuCore}

		[ ! -f "/www/server/redis/src/redis-server" ] && Error_Msg

		VM_OVERCOMMIT_MEMORY=$(cat /etc/sysctl.conf|grep vm.overcommit_memory)
		NET_CORE_SOMAXCONN=$(cat /etc/sysctl.conf|grep net.core.somaxconn)
		if [ -z "${VM_OVERCOMMIT_MEMORY}" ] && [ -z "${NET_CORE_SOMAXCONN}" ];then
			echo "vm.overcommit_memory = 1" >> /etc/sysctl.conf
			echo "net.core.somaxconn = 1024" >> /etc/sysctl.conf
			sysctl -p
		fi


		ln -sf /www/server/redis/src/redis-cli /usr/bin/redis-cli
		chown -R redis.redis /www/server/redis
		sed -i 's/dir .\//dir \/www\/server\/redis\//g' /www/server/redis/redis.conf
		sed -i 's/daemonize no/daemonize yes/g' /www/server/redis/redis.conf
		sed -i 's#^pidfile .*#pidfile /www/server/redis/redis.pid#g' /www/server/redis/redis.conf

    # 重启redis服务
    systemctl restart redis

		wget -O /etc/init.d/redis ${download_Url}/extensions/redis/init_${redis_version}
#    wget -O /www/server/redis/redis.conf ${download_Url}/conf/redis_${redis_version}.conf

		ARM_CHECK=$(uname -a|grep aarch64)
		if [ "${ARM_CHECK}" ];then
			echo "ignore-warnings ARM64-COW-BUG" >> /www/server/redis/redis.conf
		fi

		chmod +x /etc/init.d/redis
		/etc/init.d/redis start
		rm -f /www/server/redis_$redis_version.tar.gz
	fi

	echo '==============================================='
	echo 'install redis successful!'
}
Uninstall_Redis()
{
	pkill -9 redis
	rm -f /var/run/redis_6379.pid
	Service_Del
	rm -f /usr/bin/redis-cli
	rm -f /etc/init.d/redis
	rm -rf /www/server/redis
	echo '==============================================='
	echo 'uninstall redis successful!'
}
Service_Del(){
	if [ -f "/usr/bin/yum" ];then
		chkconfig --level 2345 redis off
	elif [ -f "/usr/bin/apt" ]; then
		update-rc.d redis remove
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



actionType=$1
redis_version=$2


if [ "$actionType" == 'install' ];then
  Uninstall_Redis
	System_Lib
	Gcc_Version_Check
	Install_Redis
	Service_Add
elif [ "$actionType" == 'uninstall' ];then
	Uninstall_Redis
fi