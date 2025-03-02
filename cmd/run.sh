#!/usr/bin/env bash

# 添加在脚本开始处，定义默认值
kafka=${kafka:-0}  # 默认不启用 kafka

echo "Updating and running the project..."

# 设置环境变量
export APP_ENV="test"
echo "Setting APP_ENV to: $APP_ENV and kafka: $kafka"

# 设置项目路径 (改用 HOME 目录)
PROJECT_ROOT="$HOME/game-fun-be"
ROOT_ENV_FILE="$HOME/.env.$APP_ENV"
PROJECT_ENV_FILE="$PROJECT_ROOT/.env.$APP_ENV"

# 设置 Go 的路径
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# 验证 Go 是否可用
if ! command -v go &> /dev/null; then
    echo "Go is not available. Please check your Go installation."
    exit 1
fi

# 输出 Go 版本，用于调试
go version

# 检查项目是否已经在运行，如果是则停止它
if pgrep -x "game-fun-be" > /dev/null; then
    echo "Stopping existing project..."
    pkill -f "game-fun-be"
    sleep 2
fi

# 清理项目目录
echo "Cleaning project directory..."
rm -rf "$PROJECT_ROOT"
mkdir -p "$PROJECT_ROOT"

# 根据环境选择分支
if [ "$APP_ENV" = "test" ]; then
    BRANCH="test"
else
    BRANCH="master"
fi

# 只克隆指定分支，使用 --single-branch 和 --depth 1 来最小化下载
echo "Cloning the repository (branch: $BRANCH)..."
git clone --depth 1 --single-branch --branch $BRANCH git@github.com:memeScan/game-fun-be.git "$PROJECT_ROOT"
if [ $? -ne 0 ]; then
    echo "Failed to clone the repository. Exiting..."
    exit 1
fi

# 切换到项目目录
cd "$PROJECT_ROOT" || exit

# 复制环境文件
if [ -f "$ROOT_ENV_FILE" ]; then
    echo "Copying environment file from $ROOT_ENV_FILE to $PROJECT_ENV_FILE..."
    cp "$ROOT_ENV_FILE" "$PROJECT_ENV_FILE"
else
    echo "Warning: Environment file $ROOT_ENV_FILE not found"
fi

# 更新依赖
echo "Updating dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "Failed to update dependencies. Exiting..."
    exit 1
fi

# 构建项目
echo "Building the project..."
go build -o game-fun-be main.go
if [ $? -ne 0 ]; then
    echo "Failed to build the project. Exiting..."
    exit 1
fi

# 确保日志目录存在 (同样改用 HOME 目录)
LOG_DIR="$HOME/logs/game-fun-be"
mkdir -p "$LOG_DIR"

# 在后台运行新的项目，重定向输出到指定日志文件
echo "Running the project in the background with APP_ENV=$APP_ENV and kafka=$kafka..."
nohup env APP_ENV=$APP_ENV kafka=$kafka LOG_DIR="$LOG_DIR" "$PROJECT_ROOT/game-fun-be" > "$LOG_DIR/nohup.log" 2>&1 &
echo $! > "$PROJECT_ROOT/game-fun-be.pid"

echo "Project is now running in the background."
echo "To stop the project, use 'kill $(cat $PROJECT_ROOT/game-fun-be.pid)'."

echo "Startup script completed."
