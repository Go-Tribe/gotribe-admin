#!/bin/bash

# 加载颜色配置
RED_COLOR='\E[1;31m'
GREEN_COLOR='\E[1;32m'
YELLOW_COLOR='\E[1;33m'
BLUE_COLOR='\E[1;34m'
PINK='\E[1;35m'
RESET_COLOR='\E[0m'



echo -e  "${COLOR_SUCCESS}首次手动部署: ${COLOR_SUCCESS}"
echo -e  "${COLOR_PINK}pm2 start/data/web/gotribe/admin/gotribe-admin${COLOR_PINK}"

echo -e  "${COLOR_WARNING}****开始执行自动化部署****${RES}\n\n"

# 拉取最新代码
git pull || {
    echo -e "${COLOR_ERROR}---step1:拉取最新代码失败---${RES}"
    echo "请检查是否存在合并冲突或网络问题。"
}
echo -e "${COLOR_BLUE}---step1:拉取最新代码成功---${RES}\n"

# 编译
make build || {
    echo -e "${COLOR_ERROR}---step2:编译失败---${RES}"
    echo "请检查编译工具是否安装正确。"
    exit 1
}
echo -e "${COLOR_BLUE}---step2:编译完成---${RES}\n"

# 启动
pm2 restart gotribe-admin || {
    echo -e "${COLOR_ERROR}---step3:启动失败---${RES}"
    echo "请检查应用程序配置是否正确。"
    exit 1
}
echo -e "${COLOR_BLUE}---step3:启动完成---${RES}\n"

echo -e "${COLOR_WARNING}查看日志:${COLOR_WARNING} ${COLOR_RED} pm2 log gotribe-admin${COLOR_RED}"

echo -e "${COLOR_SUCCESS}****部署成功****${RES}\n"
