#!/bin/bash

start() {
    echo "Starting the process..."
    nohup taskset -c 0 ./okxdata ../config/config.json > /data/dc/okxdata/nohup.log 2>&1 &
    echo "Process started."
}

stop() {
    echo "Stopping the process..."
    pid=$(pgrep -f "okxdata ../config/config.json")
    if [ -n "$pid" ]; then
        kill -SIGINT $pid
        echo "Process stopping"
        sleep 5
        kill $pid
        echo "Process stopped."
    else
        echo "Process is not running."
    fi
}

restart() {
    stop
    sleep 5
    start
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
        ;;
esac

exit 0