#!/bin/bash
#
#This script not add dkms 
#if you need dkms to keep ena module included in future kernel update,pls run dkms_ins
host=`curl http://169.254.169.254/latest/meta-data/local-ipv4`
time_format='+%F %T'
logging_file="/tmp/ENA_${host}.log"
file="/etc/centos-release"
version=`cat /etc/centos-release|awk '{split($(NF-1),a,".");print a[1]}'`
md5_code="ed928fb1c4732b2ae5472b8009d10b50"
judge_url="www.baidu.com"
error_repositry=''

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

yum_repositry(){
#Use CN centos respositry
cd /etc/yum.repos.d
 
 wget http://mirrors.163.com/.help/CentOS6-Base-163.repo
 repositry_code="8940c01e18087820cc7d0bf1bcd7389c"
 new_repositry_code=`md5sum CentOS-Base.repo|cut -d' ' -f1`
 [ ${new_repositry_code} == ${repositry_code} ] && rm -f CentOS6-Base-163.repo || `mv CentOS-Base.repo CentOS-Base.repo_bak && mv CentOS6-Base-163.repo CentOS-Base.repo`
 yum clean all
}

yum_judge(){
#internet access judge
  `timeout 1 curl ${judge_url} >/dev/null` && `timeout 1 curl ${judge_url}:443 >/dev/null`
  judge=$?

  if [ ${judge} -ne 0 ];then
   loger "error" "This host no internet access"
   exit 1
  fi
# yum judge
  yum repolist 2>&1|grep "cloudfarms" >/dev/null
  yum_status=$?
  if [ "${yum_status}" -eq 0 ];then
    mv /etc/yum.repos.d/${error_repositry} /etc/yum.repos.d/${error_repositry}_bak
    yum_repositry
    yum repolist 2>&1|grep "Errno" >/dev/null
    yum_status=$?
    [ ${yum_status} -eq 0 ] && loger "error" "yum repositry Error" && exit 1
  fi
   yum list all|grep $(uname -r|awk -F'.x86_64' '{print $1}')|grep kernel-devel
  yum_status=$?
  [ ${yum_status} -ne 0 ] && loger "error" "No  kernel-devel-$(uname -r) package in yum repositry" && \
  exit 1
}

centos6(){
#Install Compile tools
 yum -y install kernel-devel-$(uname -r) gcc git patch rpm-build wget 

#download ENA source package
 workdir="/tmp"
 [ -d "${workdir}/amzn-drivers-master" ] && rm -rf ${workdir}/amzn-drivers-master
 cd ${workdir}
 if [ -f "${workdir}/master.zip" ];then
    newmd5_code=`md5sum ${workdir}/master.zip |cut -d' ' -f1`
    [ ${newmd5_code} == ${md5_code} ] || `rm -f ${workdir}/master.zip && \
    wget https://github.com/amzn/amzn-drivers/archive/master.zip`
 else
   wget https://github.com/amzn/amzn-drivers/archive/master.zip
 fi
 
 unzip master.zip
 cd ${workdir}/amzn-drivers-master/kernel/linux/ena
 [ -d "/lib/modules/$(uname -r)/build" ] && make 
 compile_status=$?
 
 if [ ${compile_status} -ne 0 ];then
   loger "error" "Make error"
   exit 1
 fi
 cp ena.ko /lib/modules/$(uname -r)/  
 insmod ena.ko  
 depmod 

 echo 'add_drivers+=" ena "' >> /etc/dracut.conf.d/ena.conf
 dracut -f -v
 lsinitrd /boot/initramfs-$(uname -r).img | grep ena.ko >/dev/null
 init_ram=$?
 [ ${init_ram} -ne 0 ] && loger "error" " Failure Dracut install"
 rm -f ${workdir}/master.zip
 rm -rf ${workdir}/amzn-drivers-master
 [ -f "/etc/yum.repos.d/${error_repositry}_bak" ] && mv -f /etc/yum.repos.d/${error_repositry}_bak /etc/yum.repos.d/${error_repositry}
 [ -f "/etc/yum.repos.d/CentOS-Base.repo_bak" ] && mv -f CentOS-Base.repo_bak Centos-Base.repo  
 >/etc/dracut.conf.d/ena.conf
 
}

Centos7(){
  yum -y update
}


###################################################
dkms_ins(){
  yum -y install http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm
  yum install dkms

  VER=$( grep ^VERSION /root/amzn-drivers-master/kernel/linux/rpm/Makefile | cut -d' ' -f2 )   # Detect current version

  sudo cp -a /root/amzn-drivers-master /usr/src/amzn-drivers-${VER}   # Copy source into the source directory.

  cat > /usr/src/amzn-drivers-${VER}/dkms.conf <<EOM                  # Generate the dkms config file.
PACKAGE_NAME="ena"
PACKAGE_VERSION="$VER"
CLEAN="make -C kernel/linux/ena clean"
MAKE="make -C kernel/linux/ena/ BUILD_KERNEL=\${kernelver}"
BUILT_MODULE_NAME[0]="ena"
BUILT_MODULE_LOCATION="kernel/linux/ena"
DEST_MODULE_LOCATION[0]="/updates"
DEST_MODULE_NAME[0]="ena"
AUTOINSTALL="yes"
EOM

  dkms add -m amzn-drivers -v $VER
  dkms build -m amzn-drivers -v $VER
  dkms install -m amzn-drivers -v $VER
  cp -a /boot/grup/menu.lst /boot/grup/menu.lst_bak
  sed -i -e "/vmlinuz-$(uname -r).*/s/$/ net.ifname=0/" /boot/grup/menu.lst 
}


if [ -f ${file} ];then
  echo
else
  loger "info" "No Centos system"
  exit 1
fi

[ -f "/etc/udev/rules.d/70-persistent-net.rules" ] && mv -f /etc/udev/rules.d/70-persistent-net.rules /root/

yum_repositry
yum_judge

case ${version} in
6)
centos6 
;;
7)
centos7 
;;
*)
loger "info" "None"
esac

modinfo ena >/dev/null
mod_stat=$?
if [ ${mod_stat} -eq 0 ];then
 loger "success" "Successful Installed"
else
 loger "error" "Failure Installed"
fi
