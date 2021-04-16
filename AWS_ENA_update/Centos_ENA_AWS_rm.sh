#!/bin/bash
#SRCG: booo.wang

host=`curl http://169.254.169.254/latest/meta-data/local-ipv4`
time_format='+%F %T'
logging_file="/tmp/ENA_${host}.log"
file="/etc/centos-release"
version=`cat /etc/centos-release|awk '{split($(NF-1),a,".");print a[1]}'`

loger(){
case $1 in
"error"|"1")
  echo -e "[`date "${time_format}"`] ${host} \033[31m$2\033[0m" >>${logging_file}
  ;;
"success"|"2")
  echo -e "[`date "${time_format}"`] ${host} \033[32m$2\033[0m" >>${logging_file}
  ;;
"info"|"3")
  echo "[`date "${time_format}"`] ${host} $2" >>${logging_file}  
  ;;
*)  
  echo "Null"
esac
}



if [ ${version} == 6 ];then
#omit ena module from initramfs
   lsinitrd /boot/initramfs-$(uname -r).img | grep ena.ko >/dev/null
   init_ram=$?
   if [ ${init_ram} -ne 0 ];then
     loger "error" "No ena module in initramfs"
   else
     echo 'omit_drivers+=" ena "' >/etc/dracut.conf.d/ena.conf
     dracut -f -v
     lsinitrd /boot/initramfs-$(uname -r).img | grep ena.ko >/dev/null
     init_ram=$?
     [ ${init_ram} -eq 0 ] && loger "error" " Failure Dracut remove" && exit 1
     loger "success" "remove ena form initramfs success"
   fi
   
#remove ena module from running kernel   
   mod_status=`modprobe -rf ena 2>&1`
   if [ -z ${mod_status} ];then
     loger "info" "Remove ena module from running kernel"
   else  
     loger "error" "${mod_status}" 
     exit 1
   fi
   sed -i /ena\.ko/d /lib/modules/$(uname -r)/modules.dep |grep ena
   sed_stat=$?
   [ ${sed_stat} -ne 0 ] && \
   [ -f "/lib/modules/$(uname -r)/ena.ko" ] && \
   rm -f /lib/modules/$(uname -r)/ena.ko
   loger "success" "rollback successful"
 fi 
