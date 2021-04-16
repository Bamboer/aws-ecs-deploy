#!/bin/bash
#

device_list=()
device(){
  meta=$1
  length=${#device_list[@]}
  device_list[$length]=$meta
  echo ${device_list[@]}
}


root_device=`df -h|grep -w "/"|awk '{print $1}'`
device_list=(`device ${root_device}`)

other_device=`df -h|grep "/dev"|grep -vw "/"|awk '{print $1}'`

[ -f "/root/fstab" ]|| mv /etc/fstab /root/
for i in ${other_device};do
  device_list=(`device $i`)
done    

for i in ${device_list[@]};do
  uid=`blkid|grep "$i"|awk '{print $2}'|tr -d '"'`
  sed -i s#$i#$uid# /etc/fstab
done
