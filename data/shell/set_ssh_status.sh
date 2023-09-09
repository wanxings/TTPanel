#!/bin/bash

# Function to start the SSH service
start_ssh() {
  # 检查是否已经运行了SSH服务（使用 ssh 或 sshd 两种服务名）
  if systemctl is-active --quiet ssh || systemctl is-active --quiet sshd; then
    echo "SSH service is already running"
  else
    echo "Starting SSH service..."
    if systemctl start ssh; then
      echo "ssh successfully"
    elif systemctl start sshd; then
      echo "sshd successfully"
    else
      echo "Unable to find SSH service name. Please start the SSH service manually."
    fi
  fi
}

# Function to stop the SSH service
stop_ssh() {
  local ssh_service=""

  if [ -x "$(command -v systemctl)" ]; then
    # Check for common SSH service names using systemctl
    if systemctl status ssh &>/dev/null; then
      ssh_service="ssh"
    elif systemctl status sshd &>/dev/null; then
      ssh_service="sshd"
    else
      echo "Unable to find SSH service name. Please stop the SSH service manually."
      exit 1
    fi
  else
    echo "Systemctl not found. Please stop the SSH service manually."
    exit 1
  fi

  echo "Stopping SSH service..."
  systemctl stop "$ssh_service"
  echo "successfully"
}

case "$1" in
start)
  start_ssh
  ;;
stop)
  stop_ssh
  ;;
restart)
  stop_ssh
  sleep 1
  start_ssh
  ;;
*)
  echo "Usage: $0 {start|stop|restart}"
  exit 1
  ;;
esac

exit 0
