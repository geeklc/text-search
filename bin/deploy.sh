#!/bin/bash

#定义自启服务的文件地址
AUTO_FILE=/etc/systemd/system/dbprobe.service
#定义执行文件名称
RUN_FILE=run.sh
#获取当前的路径信息
LOCAL_PATH=`pwd`


# -------------------------------修改脚本参数-----------------------------------
# 替换脚本中的路径
sed -i 's#^LOCAL_PATH.*#LOCAL_PATH='${LOCAL_PATH}'#g' ${LOCAL_PATH}/${RUN_FILE}




# --------------------------------定义自启服务------------------------------------
#先删除再创建
rm -rf ${AUTO_FILE}
#创建自启服务的文件
touch ${AUTO_FILE}

#写入相关执行命令到自启服务文件中
cat >> ${AUTO_FILE} <<EOF
[Unit]
#服务的说明
Description=数据库探针抓取服务
#服务的描述
After=network.target


#服务运行参数的设置
[Service]
Type=simple
#该命令作用的用户
#User=root
#该命令作用的用户组
#Group=root
#运行的命令
ExecStart=/bin/bash -c '${LOCAL_PATH}/${RUN_FILE} start'
#停止的命令
ExecStop=/bin/bash -c '${LOCAL_PATH}/${RUN_FILE} stop'
#为重启命令
ExecReload=/bin/bash -c '${LOCAL_PATH}/${RUN_FILE} restart'
# 防止执行systemctl start时继续执行stop命令
RemainAfterExit=yes
#表示给服务分配独立的临时空间
PrivateTmp=True
#只要不是通过systemctl stop来停止服务，任何情况下都必须要重启服务，默认值为no
Restart=always
#重启间隔，比如某次异常后，等待5(s)再进行启动，默认值0.1(s)
RestartSec=5
#无限次重启，默认是10秒内如果重启超过5次则不再重启，设置为0表示不限次数重启
StartLimitInterval=30


[Install]
WantedBy=multi-user.target
EOF


#重新加载配置
systemctl daemon-reload

# 加上系统自启
systemctl enbale dbprobe
