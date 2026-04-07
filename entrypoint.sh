#!/bin/sh
case "$1" in
  backup)  ./white_archive -mode=backup -dir=./data ;;
  restore) ./white_archive -mode=restore -dir=./data ;;
  cron)    
    echo "$2 su runner -c '/arch/white_archive -mode=backup -dir=/arch/data'" > /etc/crontabs/root
    crond -f -c /etc/crontabs ;;
esac