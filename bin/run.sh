#!/bin/bash


# 执行文件名称
FILE_NAME="dbProbe_linux"
# 当前执行文件和脚本所在路径，执行完deploy.sh文件后会自动修改改参数，所以可以不用手动修改改参数
LOCAL_PATH=/data/app/dbtz

#进入脚本所在目录，否则退出。避免在服务自启中起不来问题
cd "${LOCAL_PATH}" || exit 1

# 如果输入格式不对，给出提示！
tips() {
    echo ""
    echo "警告!!!......请使用命令：sh run.sh [start|stop|restart|status]。例如：sh run.sh start"
    echo ""
    exit 1
}

# 启动方法
start() {
    # 重新获取一下pid，因为其它操作如stop会导致pid的状态更新
    pid=$(ps -ef | grep "$FILE_NAME" | grep -v grep | awk '{print $2}')
    # -z 表示如果$pid为空时执行
    if [ -z "$pid" ]; then
		# 执行启动命令
        nohup ./$FILE_NAME > /dev/null 2>&1 &
        echo "服务 ${FILE_NAME} 正在启动!  请执行tail -f logfile 跟踪日志信息 ！！"
	# 暂停3s
	sleep 3
	pid=$(ps -ef | grep "$FILE_NAME" | grep -v grep | awk '{print $2}')
        echo ""
        echo "...........服务 ${FILE_NAME} 启动成功！pid为${pid}！........................."
    else
        echo ""
        echo "服务 ${FILE_NAME} 已经在运行中，其pid为 ${pid}。如有需要，请使用命令：sh run.sh restart。"
        echo ""
    fi
}

# 停止方法
stop() {
    # 重新获取一下pid，因为其它操作如start会导致pid的状态更新
    pid=$(ps -ef | grep "$FILE_NAME" | grep -v grep | awk '{print $2}')
    # -z 表示如果$pid为空时执行
    if [ -z "$pid" ]; then
        echo ""
        echo "服务 ${FILE_NAME} 未在运行中，无需停止！"
        echo ""
    else
        kill -9 "$pid"
        echo ""
        echo "服务停止成功！pid:${pid} 已被强制结束！"
        echo ""
    fi
}

# 输出运行状态方法
status() {
    echo "path：${FILE_NAME}"
    # 重新获取一下pid，因为其它操作如stop、restart、start等会导致pid的状态更新
    pid=`ps -ef | grep "$FILE_NAME" | grep -v grep | awk '{print $2}'`
    # -z 表示如果$pid为空时执行
    if [ -z "$pid" ]; then
        echo ""
        echo "服务 ${FILE_NAME} 未在运行中！"
        echo ""
    else
        echo ""
        echo "服务 ${FILE_NAME} 正在运行，其pid为 ${pid}"
        echo ""
    fi
}

# 重启方法
restart() {
    echo ""
    echo ".............................正在重启.............................."
    echo "....................................................................."
    # 重新获取一下pid，因为其它操作如start会导致pid的状态更新
    pid=$( ps -ef | grep "$FILE_NAME" | grep -v grep | awk '{print $2}' )
    # -z 表示如果$pid为空时执行
    if [ ! -z "$pid" ]; then
        kill -9 "$pid"
    fi
    start
    echo "....................重启成功！..........................."
}

# 根据输入参数执行对应方法，不输入则执行tips提示方法
case "$1" in
    "start")
        start
        ;;
    "stop")
        stop
        ;;
    "status")
        status
        ;;
    "restart")
        restart
        ;;
    *)
    tips
    ;;
esac
