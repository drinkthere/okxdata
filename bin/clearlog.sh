#!/bin/bash

# 指定日志文件所在的目录
log_directory="/root/code/okxdata/bin/logs"

# 保留最近的两个okxmm.log文件
ls -t $log_directory/okxdata.log* | tail -n +5 | xargs rm -f

