#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8

#判断是否是root用户
if [ "$(whoami)" != "root" ]; then
  echo "当前非root用户,请使用root权限执行安装命令！"
  echo "The current non-root user, please execute the installation command with root authority!"
  exit 1
fi

Version="0.0.8"
PanelName="TTPanel"
setup_path="/www"
GetSysInfo() {
  if [ -s "/etc/redhat-release" ]; then
    SYS_VERSION=$(cat /etc/redhat-release)
  elif [ -s "/etc/issue" ]; then
    SYS_VERSION=$(cat /etc/issue)
  fi
  SYS_INFO=$(uname -a)
  SYS_BIT=$(getconf LONG_BIT)
  MEM_TOTAL=$(free -m | grep Mem | awk '{print $2}')
  CPU_INFO=$(getconf _NPROCESSORS_ONLN)

  echo -e "${SYS_VERSION}"
  echo -e Bit:"${SYS_BIT}" Mem:"${MEM_TOTAL}"M cpu-Core:"${CPU_INFO}"
  echo -e "${SYS_INFO}"

  echo -e "============================================"
  echo -e "安装失败"
  echo -e "installation failed"
  echo -e "============================================"
}
Red_Error() {
  echo '================================================='
  printf '\033[1;31;40m%b\033[0m\n' "$@"
  GetSysInfo
  exit 1
}
Init() {

  #检查是否已安装
  if [ -e "/etc/init.d/tt" ]; then
    echo "检测到已安装${PanelName},请勿重复安装！"
    echo "It has been detected that ${PanelName} has been installed, please do not install it again!"
    #提示是否安装
    read -r -p "输入yes强制安装/Enter yes to force the installation(yes/no): " go
    if [ "$go" != 'yes' ]; then
      echo -e "------------"
      echo "取消安装"
      echo "cancel installation"
      exit
    fi
  fi

  cd ~ || cd /

  echo "
  +----------------------------------------------------------------------
  | ${PanelName} FOR CentOS/Ubuntu/Debian
  +----------------------------------------------------------------------
  | Copyright © 2023 wanxing All rights reserved.
  +----------------------------------------------------------------------
  | 安装成功后，${PanelName}的URL将为http://IP:8888
  | The ${PanelName} URL will be http://IP:8888 when installed.
  +----------------------------------------------------------------------
  | 为了您的正常使用，请确保使用全新或纯净的系统安装Linux管理面板，不支持已部署项目/环境的系统安装
  | To ensure your normal use, please make sure to install the Linux management panel on a brand new or clean system. Installing on systems that have deployed projects/environments is not supported.
  +----------------------------------------------------------------------
  "
  #提示是否安装
  while [ "$go" != 'y' ] && [ "$go" != 'n' ]; do
    read -r -p "Do you want to install $PanelName to the $setup_path directory now?(y/n): " go
  done
  if [ "$go" == 'n' ]; then
    exit
  fi

  #检查是否存在其他web/mysql环境
  Software_Check

  #记录开始时间
  startTime=$(date +%s)
}

#检查环境
Software_Check() {
  MYSQL_CHECK=$(pgrep -f mysqld)
  PHP_CHECK=$(pgrep -f php-fpm | grep master | grep -v "$setup_path/server/php")
  NGINX_CHECK=$(pgrep -f nginx | grep master | grep -v "$setup_path/server/nginx")
  HTTPD_CHECK=$(pgrep -f 'httpd|apache' | grep -v grep)
  if [ "${PHP_CHECK}" ] || [ "${MYSQL_CHECK}" ] || [ "${NGINX_CHECK}" ] || [ "${HTTPD_CHECK}" ]; then
    #提示是否强制安装
    echo -e "----------------------------------------------------"
    echo -e "检查已有其他Web/mysql环境，安装${PanelName}可能影响现有站点及数据"
    echo -e "Check other existing Web/mysql environments, installing ${PanelName} may affect existing sites and data"
    echo -e "----------------------------------------------------"
    echo -e "已知风险/Enter yes to force installation"
    read -r -p "输入yes强制安装,其他操作退出安装: " forcedInstall
    if [ "$forcedInstall" != 'yes' ]; then
      echo -e "------------"
      echo "取消安装"
      echo "cancel installation"
      exit
    fi
  fi
}

#获取软件包管理器
Get_Package_Manager() {
  if [ -f "/usr/bin/yum" ]; then
    Package="yum"
  elif [ -f "/usr/bin/apt-get" ]; then
    Package="apt-get"
  elif [ -f "/usr/bin/pacman" ]; then
    Package="pacman"
  fi
}

#获取下载节点
Get_Download_Node() {
  if [ -n "$nodeURL" ]; then
    Download_Node=nodeURL
    return
  fi

  local nodes=("https://download.ttpanel.org")
  local tmp_file1
  tmp_file1=$(mktemp /dev/shm/net_test1.XXXXXX)
  local tmp_file2
  tmp_file2=$(mktemp /dev/shm/net_test2.XXXXXX)
  local default_node="https://download.ttpanel.org"
  local fastest_node=""
  # 测试每个节点的网络性能
  for node in "${nodes[@]}"; do
    NODE_CHECK=$(curl --connect-timeout 3 -m 3 -s -w "%{http_code} %{time_total}" "${node}/net_test" | xargs)
    RES=$(echo "${NODE_CHECK}" | awk '{print $1}')
    NODE_STATUS=$(echo "${NODE_CHECK}" | awk '{print $2}')
    TIME_TOTAL=$(echo "${NODE_CHECK}" | awk '{print $3 * 1000 - 500 }' | cut -d '.' -f 1)
    if [ "${NODE_STATUS}" == "200" ]; then
      if [ "${RES}" -ge 1500 ]; then
        echo "${RES} ${node}" >>"${tmp_file1}"
      fi
      if [ "${TIME_TOTAL}" -lt 100 ] && [ "${RES}" -ge 1500 ]; then
        echo "${TIME_TOTAL} ${node}" >>"${tmp_file2}"
      fi
    fi
  done
  # 筛选出请求最快的节点
  if [ -s "${tmp_file1}" ]; then
    fastest_node=$(sort -n "${tmp_file1}" | head -n 1 | awk '{print $2}')
  elif [ -s "${tmp_file2}" ]; then
    fastest_node=$(sort -n "${tmp_file2}" | head -n 1 | awk '{print $2}')
  fi
  # 如果没有可用节点，则使用默认节点
  if [ -z "${fastest_node}" ]; then
    echo "All nodes are unreachable. Using default node: ${default_node}"
    fastest_node="${default_node}"
  fi
  # 输出最快节点的URL地址和响应时间
  echo "Fastest node: ${fastest_node}"
  echo "当前下载节点: ${fastest_node}"
  Download_Node="${fastest_node}"
  rm -f "${tmp_file1}" "${tmp_file2}"
  sleep 3
  echo '---------------------------------------------'
}

#Arch Linux安装必要服务
Install_Arch_Pack() {
  pacman -Sy
  pacman -S curl wget tar freetype2 unzip firewalld openssl pkg-config make gcc cmake libxml2 libxslt libvpx gd libsodium oniguruma sqlite libzip autoconf inetutils sudo --noconfirm
}

#获取IP
Get_Ip_Address() {
  Extranet=""
  Extranet=$(curl -sS --connect-timeout 10 -m 60 ${Download_Node}/ask_me_ip.php)
  echo "服务器IP地址：${Extranet}"
  sleep 3
  #设置面板入口
  ${setup_path}/panel/TTPanel tools setIP -i "${Extranet}"
  LOCAL_IP=$(ip addr | grep -E -o '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}' | grep -E -v "^127\.|^255\.|^0\." | head -n 1)
}

Auto_Swap() {
  swap=$(free | grep Swap | awk '{print $2}')
  if [ "${swap}" -gt 1 ]; then
    echo "Swap total size: $swap"
    return
  fi
  swapFile="${setup_path}/swap"
  dd if=/dev/zero of=$swapFile bs=1M count=1024
  mkswap -f $swapFile
  swapon $swapFile
  echo "$swapFile    swap    swap    defaults    0 0" >>/etc/fstab
  swap=$(free | grep Swap | awk '{print $2}')
  if [ "$swap" -gt 1 ]; then
    echo "Swap total size: $swap"
    return
  fi

  sed -i "/\\${setup_path}\/swap/d" /etc/fstab
  rm -f $swapFile
}

#memCheck
MemCheck() {
  MEM_TOTAL=$(free -g | grep Mem | awk '{print $2}')
  if [ "${MEM_TOTAL}" -le "2" ]; then
    echo "内存低于2G，自动创建Swap"
    Auto_Swap
  fi
}
#创建目录
Create_Dir() {
  if [ ! -d $setup_path ]; then
    mkdir $setup_path
  fi
}
#安装deb包
Install_RPM_Pack() {
  #禁止部分软件更新
  yumPath=/etc/yum.conf
  isExc=$(grep httpd /etc/yum.conf)
  if [ "$isExc" = "" ]; then
    echo "exclude=httpd nginx php mysql mairadb" >>$yumPath
  fi
  #关闭selinux,介绍http://c.biancheng.net/view/3906.html
  sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config

  #安装依赖包
  yumPacks="libcurl-devel freetype-devel epel-release firewalld ntp  wget tar gcc make openssl openssl-devel gcc libxml2 libxml2-devel libxslt* zlib zlib-devel freetype freetype-devel lsof pcre pcre-devel vixie-cron crontabs icu libicu-devel c-ares libffi-devel ncurses-devel readline-devel gdbm-devel libpcap-devel"
  yum install -y "${yumPacks}"
  for yumPack in ${yumPacks}; do
    rpmPack=$(rpm -q "${yumPack}")
    packCheck=$(echo "${rpmPack}" | grep not)
    if [ "${packCheck}" ]; then
      yum install "${yumPack}" -y
    fi
  done
  if [ -f "/usr/bin/dnf" ]; then
    dnf install -y redhat-rpm-config
  fi
}

#安装rpm包
Install_Deb_Pack() {
  #
  ln -sf bash /bin/sh
  #待验证
  UBUNTU_22=$(grep "Ubuntu 22" /etc/issue)
  if [ "${UBUNTU_22}" ]; then
    apt-get remove needrestart -y
  fi
  ALIYUN_CHECK=$(grep "Alibaba Cloud " /etc/motd)
  if [ "${ALIYUN_CHECK}" ] && [ "${UBUNTU_22}" ]; then
    apt-get remove libicu70 -y
  fi

  apt-get update -y
  apt-get install bash -y
  if [ -f "/usr/bin/bash" ]; then
    ln -sf /usr/bin/bash /bin/sh
  fi
  apt-get install ruby -y
  apt-get install lsb-release -y

  LIBCURL_VER=$(dpkg -l | grep libcurl4 | awk '{print $3}')
  if [ "${LIBCURL_VER}" == "7.68.0-1ubuntu2.8" ]; then
    apt-get remove libcurl4 -y
    apt-get install curl -y
  fi

  #安装依赖包
  debPacks="wget curl ufw libcurl4-openssl-dev libfreetype6-dev gcc cmake libxslt-dev make tar openssl libssl-dev gcc libxml2 libxml2-dev zlib1g zlib1g-dev libjpeg-dev libpng-dev lsof libpcre3 libpcre3-dev cron net-tools swig build-essential libffi-dev libbz2-dev libncurses-dev libsqlite3-dev libreadline-dev tk-dev libgdbm-dev libdb-dev libdb++-dev libpcap-dev xz-utils git qrencode"
  apt-get install -y "$debPacks" --force-yes

  for debPack in ${debPacks}; do
    if dpkg -l | grep -q "${debPack}"; then
      apt-get install -y "$debPack"
    fi
  done

  if [ ! -d '/etc/letsencrypt' ]; then
    mkdir -p /etc/letsencryp
    mkdir -p /var/spool/cron
    if [ ! -f '/var/spool/cron/crontabs/root' ]; then
      echo '' >/var/spool/cron/crontabs/root
      chmod 600 /var/spool/cron/crontabs/root
    fi
  fi

}

#重置防火墙规则
Set_Firewall() {
  #  sshPort=$(cat /etc/ssh/sshd_config | grep 'Port '|awk '{print $2}')
  if [ -f "/usr/sbin/ufw" ]; then
    echo y | ufw enable
    ufw default deny
    ufw reload
    #关闭防火墙
    ufw disable
  fi
  if [ -f "/etc/init.d/iptables" ]; then
    /etc/init.d/iptables stop
    iptables -P INPUT ACCEPT
    iptables -P FORWARD ACCEPT
    iptables -P OUTPUT ACCEPT
    iptables -F
    chkconfig iptables off
    service iptables restart
    service iptables stop
  fi
  if [ -f "/etc/firewalld" ]; then
    rm -rf /usr/etc/firewalld/zones
    rm -rf /etc/firewalld/zones
    rm -f /etc/firewalld/direct.xml
    firewall-cmd --set-default-zone=public >/dev/null 2>&1
    service firewalld stop
    systemctl stop firewalld
  fi
}

#添加系统服务,centos9可能要加东西
Service_Add() {
  if [ "${Package}" == "yum" ] || [ "${Package}" == "dnf" ]; then
    chkconfig --add tt
    chkconfig --level 2345 tt on
  elif [ "${Package}" == "apt-get" ]; then
    update-rc.d tt defaults
  fi
}
turn_off_system_firewall() {
  if command -v ufw >/dev/null 2>&1; then
    echo "UFW detected"
    sudo ufw disable
  elif command -v firewall-cmd >/dev/null 2>&1; then
    echo "Firewalld detected"
    sudo systemctl stop firewalld
    sudo systemctl disable firewalld
  elif command -v iptables >/dev/null 2>&1; then
    echo "Iptables detected"
    service iptables stop
    chkconfig iptables off
  else
    echo "No supported firewall detected"
  fi
}
Set_Panel() {
  if [ -z "$panelPort" ]; then
    panelPort=8888
  fi
  if [ -z "$panelPwd" ]; then
    panelPwd=$(head -c 16 /dev/urandom | md5sum | head -c 8)
  fi
  if [ -z "$panelUserName" ]; then
    panelUserName=$(head -c 16 /dev/urandom | md5sum | head -c 8)
  fi

  if [ -z "$panelEntrance" ]; then
    panelEntrance=$(head -c 16 /dev/urandom | md5sum | head -c 8)
  fi

  panelSecret=$(head -c 16 /dev/urandom | md5sum | head -c 32)

  sleep 1

  Run_User="www"
  wwwUser=$(cut -d ":" -f 1 </etc/passwd | grep ^www$)
  if [ "${wwwUser}" != "www" ]; then
    groupadd ${Run_User}
    useradd -s /sbin/nologin -g ${Run_User} ${Run_User}
  fi
  cd ${setup_path}/panel/ || Red_Error "ERROR: ${setup_path}/panel/ not found." "ERROR: path does not exist"
  #设置面板端口
  ${setup_path}/panel/TTPanel tools setPort -p ${panelPort}
  sleep 1
  #关闭防火墙
  turn_off_system_firewall
  sleep 1
  #创建管理员账户
  ${setup_path}/panel/TTPanel tools createUser -u "${panelUserName}" -p "${panelPwd}"
  sleep 1
  #设置面板入口
  ${setup_path}/panel/TTPanel tools setEntrance -e "/${panelEntrance}"
  sleep 1
  #设置panelSecret
  ${setup_path}/panel/TTPanel tools setSecret -s "${panelSecret}"
  sleep 1

  #启动面板
  /etc/init.d/tt start

  sleep 3
  isStart=$(pgrep -f 'TTPanel')
  if [ -z "${isStart}" ]; then
    ls -lh /etc/init.d/tt
    Red_Error "ERROR: The TTPanel service startup failed." "ERROR: TTPanel failed to start"
  fi
}

Install_Panel() {

  mkdir -p ${setup_path}/panel
  mkdir -p ${setup_path}/server

  mkdir -p ${setup_path}/wwwroot
  mkdir -p ${setup_path}/wwwroot/default

  mkdir -p ${setup_path}/backup/database
  mkdir -p ${setup_path}/backup/project
  mkdir -p ${setup_path}/backup/panel

  if [ ! -d "/etc/init.d" ]; then
    mkdir -p /etc/init.d
  fi
  #需要修改
  if [ -f "/etc/init.d/tt" ]; then
    /etc/init.d/tt stop
    sleep 5
  fi

  wget -O /etc/init.d/tt ${Download_Node}/install/src/tt_${Version}.init -T 10
  #下载主程序
  if [ "$(uname -m)" == "aarch64" ]; then
    wget -O panel.tar.gz ${Download_Node}/install/src/panel_arm64_${Version}.tar.gz -T 10
  else
    wget -O panel.tar.gz ${Download_Node}/install/src/panel_amd64_${Version}.tar.gz -T 10
  fi

  #下载ttwaf防火墙
  mkdir -p /www/server/ttwaf
  wget -O ttwaf.tar.gz ${Download_Node}/install/src/ttwaf_${Version}.tar.gz -T 10

  chmod 755 ttwaf.tar.gz
  chmod 755 panel.tar.gz
  chmod 644 /etc/init.d/tt
  #如果有旧版本，备份旧版本
  if [ -f "${setup_path}/panel/TTPanel" ]; then
    backup_date="$(date +%Y%m%d)"
    if [ -d "${setup_path}/backup/panel/${backup_date}.tar.gz" ]; then
      rm -rf "${setup_path}/backup/panel/${backup_date}.tar.gz"
    fi
    tar -czf "${setup_path}/backup/panel/${backup_date}.tar.gz" "${setup_path}/panel/"
  fi

  tar -xzf panel.tar.gz -C ${setup_path}/panel
  tar -xzf ttwaf.tar.gz -C /www/server/ttwaf

  if [ ! -f ${setup_path}/panel/TTPanel ]; then
    ls -lh panel.tar.gz
    Red_Error "ERROR: Failed to download, please try install again!" "ERROR: 下载主程序失败，请尝试重新安装！"
  fi

  rm -f panel.tar.gz
  rm -f ttwaf.tar.gz

  chmod +x /etc/init.d/tt
  chmod -R 755 "${setup_path}/panel"
  ln -sf /etc/init.d/tt /usr/bin/tt

  ttwafAcessKey=$(head -c 16 /dev/urandom | md5sum | head -c 32)
  echo -n "$ttwafAcessKey" > /www/server/ttwaf/config/access_key
}
#安装结束
Install_Over() {
  HTTP_S="http"

  echo -e "=================================================================="
  echo -e "\033[32m安装成功！\033[0m"
  echo -e "=================================================================="
  echo "外网面板地址: ${HTTP_S}://${Extranet}:${panelPort}/${panelEntrance}"
  echo "内网面板地址: ${HTTP_S}://${LOCAL_IP}:${panelPort}/${panelEntrance}"
  echo -e "用户名: $panelUserName"
  echo -e "密  码: $panelPwd"
  echo -e "\033[33m安装过程默认关闭Linux防火墙，若无法访问面板，请检查服务器商的 防火墙/安全组 是否有放行[${panelPort}]端口\033[0m"
  echo -e "\033[33m如果未正常显示IP地址，在面板地址中手动加入IP地址即可\033[0m"
  echo -e "\033[33m面板密码仅显示一次，请做好保存，后续无法获得密码，只能通过 tt 命令修改密码\033[0m"
  echo -e "=================================================================="

  echo -e "=================================================================="
  echo -e "\033[32mCongratulations! Installed successfully!\033[0m"
  echo -e "=================================================================="
  echo "External panel address: ${HTTP_S}://IP:${panelPort}/${panelEntrance}"
  echo "Internal panel address: ${HTTP_S}://${LOCAL_IP}:${panelPort}/${panelEntrance}"
  echo -e "username: $panelUserName"
  echo -e "password: $panelPwd"
  echo -e "\033[33mIf you cannot access the panel,\033[0m"
  echo -e "\033[33mDuring the installation process, the Linux firewall is disabled by default. If you are unable to access the panel, please check if the server provider's firewall/security group has allowed access to the [${panelPort}] port.\033[0m"
  echo -e "\033[33mIf the IP address is not displayed correctly, you can manually add the IP address to the panel address.\033[0m"
  echo -e "\033[33mThe panel password is displayed only once. Please make sure to save it. You will not be able to obtain the password later. You can only modify the password using the tt command.\033[0m"
  echo -e "=================================================================="

  endTime=$(date +%s)
  outTime=$((endTime - startTime))
  echo -e "安装耗时:\033[32m $outTime \033[0m秒!"
  echo -e "Time consumed:\033[32m $outTime \033[0mSecond!"
}
Check_release_version() {
  #检查是否是64位系统
  is64bit=$(getconf LONG_BIT)
  if [ "${is64bit}" != '64' ]; then
    Red_Error "面板面板不支持32位系统,当前系统不是64位系统,请更换64位系统后再安装！\n" "The panel does not support 32-bit systems. The current system is not a 64-bit system. Please change to a 64-bit system before installing!"
  fi
}
Install_Pack() {
  if [ "${Package}" = "yum" ]; then
    Install_RPM_Pack
  elif [ "${Package}" = "apt-get" ]; then
    Install_Deb_Pack
  elif [ "${Package}" = "pacman" ]; then
    Install_Arch_Pack
  fi
}

Install() {
  Get_Package_Manager
  Check_release_version
  Init
  Install_Pack
  if [ -z "$nodeURL" ]; then
    Get_Download_Node
  else
    Download_Node=nodeURL
  fi
  Create_Dir
  MemCheck
  Install_Panel
  Set_Panel
  Service_Add
  Set_Firewall
  Get_Ip_Address
  Install_Over
}
while getopts "n:u:a:p:e:" opt; do
  case $opt in
  n)
    nodeURL="$OPTARG"
    ;;
  u)
    panelUserName="$OPTARG"
    ;;
  a)
    panelPwd="$OPTARG"
    ;;
  p)
    panelPort="$OPTARG"
    ;;
  e)
    panelEntrance="$OPTARG"
    ;;
  \?)
    echo "无效的选项: -$OPTARG" >&2
    exit 1
    ;;
  esac
done

Install
