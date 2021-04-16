#!/bin/bash
#
#there is need awscli loggin.pem ins.sh
#

time_format='+%F %T'

[ -d '/tmp/ENA' ] || mkdir -pv /tmp/ENA
logging_file="/tmp/ENA/cost_compression.log"
#old_type= c4 m4 r4 r3
#new_type c5 m5 r5 r4
work_dir=`pwd`

aws=`which aws`
[ -z ${aws} ] && loger "error" 'There is no install aws tool' && exit 1


loger(){
case $1 in
"error"|"1")
  echo -e "[`date "${time_format}"`]  \033[31m$2\033[0m" |tee -a ${logging_file}
  ;;
"success"|"2")
  echo -e "[`date "${time_format}"`]  \033[32m$2\033[0m" |tee -a ${logging_file}
  ;;
"info"|"3")
  echo "[`date "${time_format}"`]  $2" |tee -a ${logging_file}
  ;;
*)
  echo "Null"
esac
}

ip_parttern(){
 ip=$1
 parttern='^((25[0-5]|2[0-4][[:digit:]]|[01]?[[:digit:]][[:digit:]]?)\.){3}(25[0-5]|2[0-4][[:digit:]]|[01]?[[:digit:]][[:digit:]]?)$'
 if [[ "${ip}" =~ ${parttern} ]];then
   return $?
 else
   return 1
 fi
}

suffix(){
  ins_type=$1
  suf1=`echo ${ins_type} |awk -F'[[:digit:]]' '{print $2}'`
  suf2=`echo ${ins_type} |awk -F'.' '{print $1}'`
  case $suf2 in
c4)
  echo c5${suf1}
;;
m4)
  echo m5${suf1}
;;
r4)
  echo r5${suf1}
;;
r3)
  echo r4${suf1}
;;
t2)
  echo t3${suf1}
;;
*)
 echo 1
esac
}

#there are return a array for ec2 infomation
#first parameter is root device name
#second is instance type
#third is root device volume id
#the last is ENA check result
ec2_instance_info(){
  instance_id=$1
  ${aws} ec2 describe-instances --instance-id ${instance_id} | sed -n -e /RootDeviceName/p -e/InstanceType/p -e/EnaSupport/p -e/DeviceName/p -e /VolumeId/p |tr -d ' ' |tr -d '"'|tr -d ','>/tmp/ec2_info.log
  root_device_name=`cat /tmp/ec2_info.log|grep RootDeviceName|tail -1|cut -d':' -f2`
  instance_type=`cat /tmp/ec2_info.log|grep InstanceType |tail -1|cut -d':' -f2`
  root_volume_id=`cat /tmp/ec2_info.log|grep -A1 DeviceName:${root_device_name} | grep -vE "RootDeviceName|--"|tail -1|cut -d':' -f2`
  ena_support=`cat /tmp/ec2_info.log|grep EnaSupport|cut -d':' -f2`
  echo ${root_device_name} ${instance_type} ${root_volume_id} ${ena_support}
}

#wait snapshot completed
snapshot_stat_check(){
  snap_id=$1
  Flag=true
  while ${Flag};do
  snap_stat=`aws ec2 describe-snapshots --snapshot-ids ${snap_id}|grep State|tr -d ' '|tr -d '"'|tr -d ',' |cut -d':' -f2`
  if [ "${snap_stat}" == 'completed' ];then
    Flag=false
    loger 'success' "Snapshot ${snap_id} has successful created"
  fi
  sleep 10
  done
}


f_process(){
  for ip in ${ips[@]};do
   ip_parttern $ip
   if [ $? -eq 0 ];then
    instance_id=`ssh -i logging.pem -p22022  tier2@grafana "sudo -uroot ansible $ip -m shell -a'curl http://169.254.169.254/latest/meta-data/instance-id'" |grep 'i-*' |awk '{print $1}'`
    loger 'success' "${ip} instance id: ${instance_id}"
    if [[ "${instance_id}" =~ "i-" ]];then
      instance_info=(`ec2_instance_info ${instance_id}`)
    else
      loger 'error' "There is no instance id"
      continue
    fi
#backup instance
    echo -ne "$ip need \033[31mbackup\033[0m? y/n: "
    read  ami
    if [ ${ami} == 'y' ];then
     snapshot_id=`${aws} ec2 create-snapshot --volume-id ${instance_info[2]} --description 'Just for backup' --tag-specifications "ResourceType=snapshot,Tags=[{Key=old-root-name,Value=${instance_info[0]}},{Key=Name,Value=${ip}}]" |grep 'SnapshotId' |tr -d ' '|tr -d '"' |tr -d ','|cut -d':' -f2`
     snapshot_stat_check ${snapshot_id}
     echo "${ip} ${snapshot_id}" >>/tmp/ENA/snapshot.log
    fi
#get target id type
    new_type=`suffix ${instance_info[1]}`
    loger 'info' "${ip} has a new instance type: ${new_type}"
    if [ ${new_type} == 1 ];then
      loger "info" "$ip no need change instance type"
      continue
    fi
    if [ -z "${instance_info[3]}" ] || [ "${instance_info[3]}" == 'false' ];then
        loger "info" "Start ENA update process"
        ssh -i ${work_dir}/logging.pem -p22022 tier2@grafana "sudo -uroot ansible $ip -m script -a '/tmp/ins.sh'" 2>&1 >/dev/null
        update_stat=`ssh -i ${work_dir}/logging.pem  -p22022 tier2@grafana "sudo -uroot ansible $ip -m shell -a'cat /tmp/ENA_* '"|grep -i 'Successful'`
        if [ -n "${update_stat}" ];then
          ${aws} ec2   reboot-instances --instance-ids  ${instance_id}
        else
          loger "error" "${ip} There has Unkown error ENA module install failure"
          continue
        fi
        flag=true
        while $flag;do
          timeout 1 nc -zv $ip 22022
          [ $? -eq 0 ] && flag=false
          sleep 10
          echo -e "\033[32mConnecting $ip ...\033[0m"
        done
        echo -e "\033[32mInstance up and update ENA successfule\033[0m"
    fi

      flag=true
      while $flag;do
       instance_stat=`${aws} ec2   stop-instances --instance-ids  ${instance_id}|grep -A 2 "CurrentState"|tail -1|awk  '{print $2}'|tr -d '"'`
       [ ${instance_stat} == "stopped" ] && flag=false
       sleep 10
      done
      echo -e "\033[31mInstance stopped\033[0m"
      echo -e "Instance type change to \033[32m${new_type}\033[0m"
      ${aws} ec2  modify-instance-attribute --instance-id  ${instance_id} --instance-type  ${new_type} 2>&1 >/dev/null
      sleep 2
      ${aws} ec2 modify-instance-attribute --instance-id ${instance_id} --ena-support  2>&1 >/dev/null
      sleep 2
      ${aws} ec2  start-instances  --instance-ids  ${instance_id} 2>&1 >/dev/null

      flag=true
      while $flag;do
       ec2_stat=`${aws} ec2 describe-instance-status --instance-ids  ${instance_id}|grep -w 'Status'`
       check_stat=`echo ${ec2_stat}|grep -vE 'passed|ok'`
       [ -n "${ec2_stat}" ] && [ -z "${check_stat}" ] && flag=false
       echo -e "\033[32mInstance ${instance_id} init \033[0m..."
       sleep 10
      done
      echo -e "Instance \033[32mstarted\033[0m..."

      flag=true
     time_s=`date +%s`
      while $flag;do
        new_s=`date +%s`
        let run_t="${new_s} - ${time_s}"
        [ "${run_t}" -gt 300 ] && up_stat=1 && break
        timeout 1 nc -zv $ip 22022
        [ $? -eq 0 ] && flag=false
        sleep 10
        echo -e "Checking \033[32m$ip\033[0m network status \033[32m...\033[0m"
      done
      [ -z "${up_stat}" ] && echo "$ip :Success update" >>${logging_file} || echo "$ip :Timeout Error update ENA module, Pls have check"
  fi
  unset time_s
  unset new_s
  unset run_t
 done
echo
echo
}


[ -f "${work_dir}/logging.pem" ] ||`loger "error" "There is no authtication file" && exit 1`

echo -e "There have two mode you can choice:\033[31mf(file)|s(single)\033[0m"
echo -en "Please input your \033[32mchoice\033[0m: "
read choice

case ${choice} in
f|choice|F)
  echo -en "Please input your file name(\033[31mfull of path\033[0m): "
  read path
  if [ -f ${path} ];then
   ips=(`cat $path`)
  fi
  f_process
;;
s|S|single)
  echo -en "Please input your \033[31mip\033[0m: "
  read host
  ips=(`echo $host`)
  f_process
;;
*)
 echo -e "\033[31mError Input\033[0m"
esac

