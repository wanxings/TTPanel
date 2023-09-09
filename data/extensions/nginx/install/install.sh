#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8

header_file=/www/panel/data/shell/install_header.sh
. $header_file
download_Url=$NODE_URL

Root_Path=/www
Setup_Path=$Root_Path/server/nginx
run_path="/root"
Is_64bit=$(getconf LONG_BIT)

ARM_CHECK=$(uname -a | grep -E 'aarch64|arm|ARM')
LUAJIT_VER="2.0.4"
LUAJIT_INC_PATH="luajit-2.0"

if [ "${ARM_CHECK}" ]; then
  LUAJIT_VER="2.1.0"
  LUAJIT_INC_PATH="luajit-2.1"
fi

System_Lib() {
  if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ]; then
    Pack="gcc gcc-c++ curl curl-devel libtermcap-devel ncurses-devel libevent-devel readline-devel libuuid-devel"
    ${PM} install ${Pack} -y
  elif [ "${PM}" == "apt-get" ]; then
    LIBCURL_VER=$(dpkg -l | grep libx11-6 | awk '{print $3}')
    if [ "${LIBCURL_VER}" == "2:1.6.9-2ubuntu1.3" ]; then
      apt remove libx11* -y
      apt install libx11-6 libx11-dev libx11-data -y
    fi
    Pack="gcc g++ libgd3 libgd-dev libevent-dev libncurses5-dev libreadline-dev uuid-dev"
    ${PM} install ${Pack} -y
  fi

}

Service_Add() {
  if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ]; then
    chkconfig --add nginx
    chkconfig --level 2345 nginx on
  elif [ "${PM}" == "apt-get" ]; then
    update-rc.d nginx defaults
  fi
  if [ "$?" == "127" ]; then
    wget -O /usr/lib/systemd/system/nginx.service ${download_Url}/extensions/nginx/nginx.service
    systemctl enable nginx.service
  fi
}
Service_Del() {
  if [ "${PM}" == "yum" ] || [ "${PM}" == "dnf" ]; then
    chkconfig --del nginx
    chkconfig --level 2345 nginx off
  elif [ "${PM}" == "apt-get" ]; then
    update-rc.d nginx remove
  fi
}
Set_Time() {
  BASH_DATE=$(stat nginx.sh | grep Modify | awk '{print $2}' | tr -d '-')
  SYS_DATE=$(date +%Y%m%d)
  [ "${SYS_DATE}" -lt "${BASH_DATE}" ] && date -s "$(curl https://panel.cixing.io/api/get_date.php)"
}
Install_Jemalloc() {
  if [ ! -f '/usr/local/lib/libjemalloc.so' ]; then
    wget -O jemalloc-5.0.1.tar.bz2 ${download_Url}/install/src/jemalloc-5.0.1.tar.bz2
    tar -xvf jemalloc-5.0.1.tar.bz2
    cd jemalloc-5.0.1
    ./configure
    make && make install
    ldconfig
    cd ..
    rm -rf jemalloc*
  fi
}
Install_LuaJIT() {
  if [ ! -f '/usr/local/lib/libluajit-5.1.so' ] || [ ! -f "/usr/local/include/${LUAJIT_INC_PATH}/luajit.h" ]; then
    wget -c -O LuaJIT-${LUAJIT_VER}.tar.gz ${download_Url}/install/src/LuaJIT-${LUAJIT_VER}.tar.gz -T 10
    tar xvf LuaJIT-${LUAJIT_VER}.tar.gz
    cd LuaJIT-${LUAJIT_VER}
    make linux
    make install
    cd ..
    rm -rf LuaJIT-*
    export LUAJIT_LIB=/usr/local/lib
    export LUAJIT_INC=/usr/local/include/${LUAJIT_INC_PATH}/
    ln -sf /usr/local/lib/libluajit-5.1.so.2 /usr/local/lib64/libluajit-5.1.so.2
    echo "/usr/local/lib" >>/etc/ld.so.conf
    ldconfig
  fi
}
Install_cjson() {
  if [ ! -f /usr/local/lib/lua/5.1/cjson.so ]; then
    wget -O lua-cjson-2.1.0.tar.gz $download_Url/install/src/lua-cjson-2.1.0.tar.gz -T 20
    tar xvf lua-cjson-2.1.0.tar.gz
    rm -f lua-cjson-2.1.0.tar.gz
    cd lua-cjson-2.1.0
    make
    make install
    cd ..
    rm -rf lua-cjson-2.1.0
  fi
}
Download_Src() {
  mkdir -p ${Setup_Path}
  cd ${Setup_Path}
  rm -rf ${Setup_Path}/src
  if [ "${version}" == "tengine" ] || [ "${version}" == "openresty" ]; then
    wget -O ${Setup_Path}/src.tar.gz ${download_Url}/extensions/nginx/src/${version}-${nginxVersion}.tar.gz -T20
    tar -xvf src.tar.gz
    mv ${version}-${nginxVersion} src
  else
    wget -O ${Setup_Path}/src.tar.gz ${download_Url}/extensions/nginx/src/nginx-${nginxVersion}.tar.gz -T20
    tar -xvf src.tar.gz
    tar -xvf src.tar.gz
    mv nginx-${nginxVersion} src
  fi

  cd src

  TLSv13_NGINX=$(echo ${nginxVersion} | tr -d '.' | cut -c 1-3)
  if [ "${TLSv13_NGINX}" -ge "115" ] && [ "${TLSv13_NGINX}" != "181" ]; then
    opensslVer="1.1.1q"
  else
    opensslVer="1.0.2u"
  fi

  wget -O openssl.tar.gz ${download_Url}/install/src/openssl-${opensslVer}.tar.gz
  tar -xvf openssl.tar.gz
  mv openssl-${opensslVer} openssl
  rm -f openssl.tar.gz

  pcre_version="8.43"
  wget -O pcre-$pcre_version.tar.gz ${download_Url}/install/src/pcre-$pcre_version.tar.gz
  tar zxf pcre-$pcre_version.tar.gz

  wget -O ngx_cache_purge.tar.gz ${download_Url}/install/src/ngx_cache_purge-2.3.tar.gz
  tar -zxvf ngx_cache_purge.tar.gz
  mv ngx_cache_purge-2.3 ngx_cache_purge
  rm -f ngx_cache_purge.tar.gz

  wget -O nginx-sticky-module.zip ${download_Url}/install/src/nginx-sticky-module.zip
  unzip -o nginx-sticky-module.zip
  rm -f nginx-sticky-module.zip

  wget -O nginx-http-concat.zip ${download_Url}/install/src/nginx-http-concat-1.2.2.zip
  unzip -o nginx-http-concat.zip
  mv nginx-http-concat-1.2.2 nginx-http-concat
  rm -f nginx-http-concat.zip

  #lua_nginx_module
  LuaModVer="0.10.13"
  wget -c -O lua-nginx-module-${LuaModVer}.zip ${download_Url}/install/src/lua-nginx-module-${LuaModVer}.zip
  unzip -o lua-nginx-module-${LuaModVer}.zip
  mv lua-nginx-module-${LuaModVer} lua_nginx_module
  rm -f lua-nginx-module-${LuaModVer}.zip

  #ngx_devel_kit
  NgxDevelKitVer="0.3.1"
  wget -c -O ngx_devel_kit-${NgxDevelKitVer}.zip ${download_Url}/install/src/ngx_devel_kit-${NgxDevelKitVer}.zip
  unzip -o ngx_devel_kit-${NgxDevelKitVer}.zip
  mv ngx_devel_kit-${NgxDevelKitVer} ngx_devel_kit
  rm -f ngx_devel_kit-${NgxDevelKitVer}.zip

  #nginx-dav-ext-module
  NgxDavVer="3.0.0"
  wget -c -O nginx-dav-ext-module-${NgxDavVer}.tar.gz ${download_Url}/install/src/nginx-dav-ext-module-${NgxDavVer}.tar.gz
  tar -xvf nginx-dav-ext-module-${NgxDavVer}.tar.gz
  mv nginx-dav-ext-module-${NgxDavVer} nginx-dav-ext-module
  rm -f nginx-dav-ext-module-${NgxDavVer}.tar.gz

  if [ "${Is_64bit}" = "64" ]; then
    if [ "${version}" == "tengine" ]; then
      NGX_PAGESPEED_VAR="1.13.35.2"
      wget -O ngx-pagespeed-${NGX_PAGESPEED_VAR}.tar.gz ${download_Url}/install/src/ngx-pagespeed-${NGX_PAGESPEED_VAR}.tar.gz
      tar -xvf ngx-pagespeed-${NGX_PAGESPEED_VAR}.tar.gz
      mv ngx-pagespeed-${NGX_PAGESPEED_VAR} ngx-pagespeed
      rm -f ngx-pagespeed-${NGX_PAGESPEED_VAR}.tar.gz
    fi
  fi
}
Install_Configure() {
  Run_User="www"
  wwwUser=$(cat /etc/passwd | grep www)
  if [ "${wwwUser}" == "" ]; then
    groupadd ${Run_User}
    useradd -s /sbin/nologin -g ${Run_User} ${Run_User}
  fi

  [ -f "/www/panel/data/shell/nginx_prepare.sh" ] && . /www/server/panel/install/nginx_prepare.sh
  [ -f "/www/panel/data/shell/nginx_configure.pl" ] && ADD_EXTENSION=$(cat /www/panel/data/shell/nginx_configure.pl)
  if [ -f "/usr/local/lib/libjemalloc.so" ] && [ -z "${ARM_CHECK}" ]; then
    jemallocLD="--with-ld-opt="-ljemalloc""
  fi

  if [ "${version}" == "1.8" ]; then
    ENABLE_HTTP2="--with-http_spdy_module"
  else
    ENABLE_HTTP2="--with-http_v2_module --with-stream --with-stream_ssl_module --with-stream_ssl_preread_module"
  fi

  WebDav_NGINX=$(echo ${nginxVersion} | tr -d '.' | cut -c 1-3)
  if [ "${WebDav_NGINX}" -ge "114" ] && [ "${WebDav_NGINX}" != "181" ]; then
    ENABLE_WEBDAV="--with-http_dav_module --add-module=${Setup_Path}/src/nginx-dav-ext-module"
  fi

  if [ "${version}" == "openresty" ]; then
    ENABLE_LUA="--with-luajit"
  elif [ -z "${ARM_CHECK}" ] && [ -f "/usr/local/include/${LUAJIT_INC_PATH}/luajit.h" ]; then
    ENABLE_LUA="--add-module=${Setup_Path}/src/ngx_devel_kit --add-module=${Setup_Path}/src/lua_nginx_module"
  fi

  ENABLE_STICKY="--add-module=${Setup_Path}/src/nginx-sticky-module"
  if [ "$version" == "1.23" ]; then
    ENABLE_LUA=""
    ENABLE_STICKY=""
  fi

  if [ "${ARM_CHECK}" ]; then
    ARM_LUA="--add-module=/www/server/nginx/src/ngx_devel_kit --add-module=/www/server/nginx/src/lua_nginx_module"
  else
    ARM_LUA=""
  fi


  #    name=nginx
  #    i_path=/www/server/panel/install/$name
  #
  #    i_args=$(cat $i_path/config.pl | xargs)
  #    i_make_args=""
  #    for i_name in $i_args; do
  #        init_file=$i_path/$i_name/init.sh
  #        if [ -f $init_file ]; then
  #            bash $init_file
  #        fi
  #
  #        args_file=$i_path/$i_name/args.pl
  #        if [ -f $args_file ]; then
  #            args_string=$(cat $args_file)
  #            i_make_args="$i_make_args $args_string"
  #        fi
  #    done

  cd ${Setup_Path}/src

  export LUAJIT_LIB=/usr/local/lib
  export LUAJIT_INC=/usr/local/include/${LUAJIT_INC_PATH}/
  export LD_LIBRARY_PATH=/usr/local/lib/:$LD_LIBRARY_PATH

  ./configure --user=www --group=www --prefix=${Setup_Path} ${ENABLE_LUA} ${ARM_LUA} --add-module=${Setup_Path}/src/ngx_cache_purge ${ENABLE_STICKY} --with-openssl=${Setup_Path}/src/openssl --with-pcre=pcre-${pcre_version} ${ENABLE_HTTP2} --with-http_stub_status_module --with-http_ssl_module --with-http_image_filter_module --with-http_gzip_static_module --with-http_gunzip_module --with-ipv6 --with-http_sub_module --with-http_flv_module --with-http_addition_module --with-http_realip_module --with-http_mp4_module --with-ld-opt="-Wl,-E" --with-cc-opt="-Wno-error" ${jemallocLD} ${ENABLE_WEBDAV} ${ENABLE_NGX_PAGESPEED} ${ADD_EXTENSION} ${i_make_args}
  make -j${cpuCore}
}
Install_Nginx() {
  make install
  if [ "${version}" == "openresty" ]; then
    ln -sf /www/server/nginx/nginx/html /www/server/nginx/html
    ln -sf /www/server/nginx/nginx/conf /www/server/nginx/conf
    ln -sf /www/server/nginx/nginx/logs /www/server/nginx/logs
    ln -sf /www/server/nginx/nginx/sbin /www/server/nginx/sbin

  fi

  if [ ! -f "${Setup_Path}/sbin/nginx" ]; then
    echo '========================================================'
    GetSysInfo
    echo -e "ERROR: nginx-${nginxVersion} installation failed."
    if [ -z "${SYS_VERSION}" ]; then
      echo -e "============================================"
      echo -e "检测到为非常用系统安装,请尝试安装其他Nginx版本看是否正常"
      echo -e "如无法正常安装，建议更换至Centos-7或Debian-10+或Ubuntu-20+系统安装面板"
      echo -e "特殊情况可通过以下联系方式寻求安装协助情况"
      echo -e "============================================"
    fi
    echo -e "安装失败，请保存以上报错信息"
    echo -e "============================================"
    echo -e "联系邮箱：qianqianwanxingsu@gmail.com"
    echo -e "============================================"
    rm -rf ${Setup_Path}
    exit 1
  fi

  \cp -rpa ${Setup_Path}/sbin/nginx /www/backup/nginxBak
  chmod -x /www/backup/nginxBak
  md5sum ${Setup_Path}/sbin/nginx >/www/panel/data/extensions/nginx/nginx_md5.pl
  ln -sf ${Setup_Path}/sbin/nginx /usr/bin/nginx
  rm -f ${Setup_Path}/conf/nginx.conf

  cd ${Setup_Path}
  rm -f src.tar.gz
}
Update_Nginx() {
  if [ "${nginxVersion}" = "openresty" ]; then
    make install
    echo -e "done"
    nginx -v
    echo "${nginxVersion}" >${Setup_Path}/version.pl
    rm -f ${Setup_Path}/version_check.pl
    exit
  fi
  if [ ! -f ${Setup_Path}/src/objs/nginx ]; then
    echo '========================================================'
    GetSysInfo
    echo -e "ERROR: nginx-${nginxVersion} installation failed."
    echo -e "升级失败，请截图以上报错信息发送邮件至qianqianwanxingsu@gmail.com求助"
    exit 1
  fi
  sleep 1
  /etc/init.d/nginx stop
  mv -f ${Setup_Path}/sbin/nginx ${Setup_Path}/sbin/nginxBak
  \cp -rfp ${Setup_Path}/src/objs/nginx ${Setup_Path}/sbin/
  sleep 1
  /etc/init.d/nginx start
  rm -rf ${Setup_Path}/src
  nginx -v

  echo "${nginxVersion}" >${Setup_Path}/version.pl
  rm -f ${Setup_Path}/version_check.pl
  if [ "${version}" == "tengine" ]; then
    echo "2.2.4(${tengine})" >${Setup_Path}/version_check.pl
  fi
  exit
}
Set_Conf() {
  Default_Website_Dir=$Root_Path'/wwwroot/default'
  mkdir -p ${Default_Website_Dir}
  mkdir -p ${Root_Path}/wwwlogs
  mkdir -p ${Setup_Path}/conf/vhost
  mkdir -p /usr/local/nginx/logs
  mkdir -p ${Setup_Path}/conf/rewrite

  mkdir -p /www/wwwlogs/load_balancing/tcp
  mkdir -p /www/panel/data/extensions/nginx/vhost/tcp

  wget -O ${Setup_Path}/conf/nginx.conf ${download_Url}/extensions/nginx/conf/nginx1.conf -T20
  wget -O ${Setup_Path}/conf/pathinfo.conf ${download_Url}/extensions/nginx/conf/pathinfo.conf -T20
  wget -O ${Setup_Path}/conf/enable-php.conf ${download_Url}/extensions/nginx/conf/enable-php.conf -T20
  wget -O ${Setup_Path}/html/index.html ${download_Url}/extensions/nginx/error/index.html -T20

  chmod 755 /www/server/nginx/
  chmod 755 /www/server/nginx/html/
  chmod 755 /www/wwwroot/
  chmod 644 /www/server/nginx/html/*

  cat >${Root_Path}/panel/data/extensions/nginx/vhost/main/phpfpm_status.conf <<EOF
server {
    listen 80;
    server_name 127.0.0.1;
    allow 127.0.0.1;
    location /nginx_status {
        stub_status on;
        access_log off;
    }
EOF
  echo "" >/www/server/nginx/conf/enable-php-00.conf
  for phpV in 52 53 54 55 56 70 71 72 73 74 75 80 81 82; do
    cat >${Setup_Path}/conf/enable-php-${phpV}.conf <<EOF
    location ~ [^/]\.php(/|$)
    {
        try_files \$uri =404;
        fastcgi_pass  unix:/tmp/php-cgi-${phpV}.sock;
        fastcgi_index index.php;
        include fastcgi.conf;
        include pathinfo.conf;
    }
EOF
    cat >>${Root_Path}/panel/data/extensions/nginx/vhost/main/phpfpm_status.conf <<EOF
    location /phpfpm_${phpV}_status {
        fastcgi_pass unix:/tmp/php-cgi-${phpV}.sock;
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME \$fastcgi_script_name;
    }
EOF
  done
  echo \} >>${Root_Path}/panel/data/extensions/nginx/vhost/main/phpfpm_status.conf

  cat >${Setup_Path}/conf/proxy.conf <<EOF
proxy_temp_path ${Setup_Path}/proxy_temp_dir;
proxy_cache_path ${Setup_Path}/proxy_cache_dir levels=1:2 keys_zone=cache_one:20m inactive=1d max_size=5g;
client_body_buffer_size 512k;
proxy_connect_timeout 60;
proxy_read_timeout 60;
proxy_send_timeout 60;
proxy_buffer_size 32k;
proxy_buffers 4 64k;
proxy_busy_buffers_size 128k;
proxy_temp_file_write_size 128k;
proxy_next_upstream error timeout invalid_header http_500 http_503 http_404;
proxy_cache cache_one;
EOF

  cat >${Setup_Path}/conf/luawaf.conf <<EOF
lua_shared_dict limit 10m;
lua_package_path "/www/server/nginx/waf/?.lua";
init_by_lua_file  /www/server/nginx/waf/init.lua;
access_by_lua_file /www/server/nginx/waf/waf.lua;
EOF

  mkdir -p /www/wwwlogs/waf
  chown www.www /www/wwwlogs/waf
  chmod 744 /www/wwwlogs/waf
  mkdir -p /www/panel/data/vhost

  sed -i "s#include vhost/\*.conf;#include /www/panel/data/extensions/nginx/vhost/main/\*.conf;#" ${Setup_Path}/conf/nginx.conf
  sed -i "s#/www/wwwroot/default#/www/server/phpmyadmin#" ${Setup_Path}/conf/nginx.conf
  sed -i "/pathinfo/d" ${Setup_Path}/conf/enable-php.conf
  sed -i "s/#limit_conn_zone.*/limit_conn_zone \$binary_remote_addr zone=perip:10m;\n\tlimit_conn_zone \$server_name zone=perserver:10m;/" ${Setup_Path}/conf/nginx.conf
  sed -i "s/mime.types;/mime.types;\n\t\tinclude proxy.conf;\n/" ${Setup_Path}/conf/nginx.conf
  #if [ "${nginx_version}" == "1.12.2" ] || [ "${nginx_version}" == "openresty" ] || [ "${nginx_version}" == "1.14.2" ];then
#  sed -i "s/mime.types;/mime.types;\n\t\t#include luawaf.conf;\n/" ${Setup_Path}/conf/nginx.conf
  #fi

  PHPVersion=""
  for phpVer in 52 53 54 55 56 70 71 72 73 74 80; do
    if [ -d "/www/server/php/${phpVer}/bin" ]; then
      PHPVersion=${phpVer}
    fi
  done

  if [ "${PHPVersion}" ]; then
    \cp -r -a ${Setup_Path}/conf/enable-php-${PHPVersion}.conf ${Setup_Path}/conf/enable-php.conf
  fi

  wget -O /etc/init.d/nginx ${download_Url}/extensions/nginx/nginx.init -T 5
  chmod +x /etc/init.d/nginx
}
Set_Version() {
  if [ "${version}" == "tengine" ]; then
    echo "-Tengine2.2.3" >${Setup_Path}/version.pl
    echo "2.2.4(${tengine})" >${Setup_Path}/version_check.pl
  elif [ "${version}" == "openresty" ]; then
    echo "openresty" >${Setup_Path}/version.pl
    echo "openresty-${openresty}" >${Setup_Path}/version_check.pl
  else
    echo "${nginxVersion}" >${Setup_Path}/version.pl
  fi

}

Uninstall_Nginx() {
  if [ -f "/etc/init.d/nginx" ]; then
    Service_Del
    /etc/init.d/nginx stop
    rm -f /etc/init.d/nginx
  fi
  pkill -9 nginx
  rm -rf ${Setup_Path}
}

actionType=$1
version=$2




if [ "${actionType}" == "uninstall" ]; then
  Service_Del
  Uninstall_Nginx
else
  if [[ $version == *-* ]]; then
    # 版本号包含 "-" 符号
    nginxVersion="${version%%-*}"   # 提取 "-" 符号前的部分
    version="${version#*-}"    # 提取 "-" 符号后的部分
  else
    # 版本号不包含 "-" 符号
    nginxVersion=$version
    version="nginx"
  fi
  if [ "${actionType}" == "install" ]; then
    if [ -f "/www/server/nginx/sbin/nginx" ]; then
      Uninstall_Nginx
    fi
    System_Lib

    #arm架构有问题
    Install_Jemalloc
    Install_LuaJIT
    Install_cjson

    Download_Src
    Install_Configure
    Install_Nginx
    Set_Conf
    Set_Version
    Service_Add
    /etc/init.d/nginx start
  elif [ "${actionType}" == "update" ]; then
    Download_Src
    Install_Configure
    Update_Nginx
  fi
fi
