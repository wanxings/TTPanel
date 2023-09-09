#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin

export PATH
export LANG=en_US.UTF-8
export LANGUAGE=en_US:en

PanelVersion="0.0.1"

get_node_url(){
  local nodes=("https://download.ttpanel.org" "https://download.ttpanel.org")
  local tmp_file1=$(mktemp /dev/shm/net_test1.XXXXXX)
  local tmp_file2=$(mktemp /dev/shm/net_test2.XXXXXX)
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
        echo "${RES} ${node}" >> "${tmp_file1}"
      fi
      if [ "${TIME_TOTAL}" -lt 100 ] && [ "${RES}" -ge 1500 ]; then
        echo "${TIME_TOTAL} ${node}" >> "${tmp_file2}"
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
  NODE_URL="${fastest_node}"
  rm -f "${tmp_file1}" "${tmp_file2}"
}

GetCpuStat(){
	time1=$(cat /proc/stat |grep 'cpu ')
	sleep 1
	time2=$(cat /proc/stat |grep 'cpu ')
	cpuTime1=$(echo ${time1}|awk '{print $2+$3+$4+$5+$6+$7+$8}')
	cpuTime2=$(echo ${time2}|awk '{print $2+$3+$4+$5+$6+$7+$8}')
	runTime=$((${cpuTime2}-${cpuTime1}))
	idelTime1=$(echo ${time1}|awk '{print $5}')
	idelTime2=$(echo ${time2}|awk '{print $5}')
	idelTime=$((${idelTime2}-${idelTime1}))
	useTime=$(((${runTime}-${idelTime})*3))
	[ ${useTime} -gt ${runTime} ] && cpuBusy="true"
	if [ "${cpuBusy}" == "true" ]; then
		cpuCore=$((${cpuInfo}/2))
	else
		cpuCore=$((${cpuInfo}-1))
	fi
}
GetPackManager(){
	if [ -f "/usr/bin/yum" ] && [ -f "/etc/yum.conf" ]; then
		PM="yum"
	elif [ -f "/usr/bin/apt-get" ] && [ -f "/usr/bin/dpkg" ]; then
		PM="apt-get"
	fi
}


GetSysInfo(){
	if [ "${PM}" = "yum" ]; then
		SYS_VERSION=$(cat /etc/redhat-release)
	elif [ "${PM}" = "apt-get" ]; then
		SYS_VERSION=$(cat /etc/issue)
	fi
	SYS_INFO=$(uname -msr)
	SYS_BIT=$(getconf LONG_BIT)
	MEM_TOTAL=$(free -m|grep Mem|awk '{print $2}')
	CPU_INFO=$(getconf _NPROCESSORS_ONLN)
	GCC_VER=$(gcc -v 2>&1|grep "gcc version"|awk '{print $3}')
	CMAKE_VER=$(cmake --version|grep version|awk '{print $3}')

	echo -e ${SYS_VERSION}
	echo -e Bit:${SYS_BIT} Mem:${MEM_TOTAL}M Core:${CPU_INFO} gcc:${GCC_VER} cmake:${CMAKE_VER}
	echo -e ${SYS_INFO}
}
cpuInfo=$(getconf _NPROCESSORS_ONLN)
if [ "${cpuInfo}" -ge "4" ];then
	GetCpuStat
else
	cpuCore="1"
fi
GetPackManager

if [ ! $NODE_URL ];then
	echo '正在选择下载节点...';
	echo "selecting download node...";
	get_node_url
fi


