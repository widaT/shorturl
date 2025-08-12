#!/bin/bash

# 初始化shint key 脚本
# 使用4位随机数 * (shard_index+1) * 10000最为初始值

# 默认配置

RIDIS_HOST="127.0.0.1"
RIDIS_PORT=6379
SHARD_COUNT=128
# RIDIS_PASSWORD="123456"
SHINT_KEY="link:shint:"

REDIS_CMD="redis-cli -h ${RIDIS_HOST} -p ${RIDIS_PORT}"

if [ -n "$RIDIS_PASSWORD" ]; then
    REDIS_CMD="${REDIS_CMD} -a ${RIDIS_PASSWORD}"
fi

echo "init shint key: ${SHINT_KEY}" 
echo "shard count: ${SHARD_COUNT}"
echo "redis cmd: ${REDIS_CMD}"
# 解析命令行

for ((i=0;i<${SHARD_COUNT};i++)); do
    random_num=$(shuf -i 1000-9999 -n 1)
    intial_value=$((random_num * (i+1) * 10000))
    key="${SHINT_KEY}${i}"
    echo "set ${key} ${intial_value} random_num: ${random_num}"
    ${REDIS_CMD} set "${key}" "${intial_value}" > /dev/null
    if [ $? -eq 0 ]; then
        echo "set ${key} success"
    else
        echo "set ${key} failed"
    fi
done

echo "init shint key done"
