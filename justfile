# 显示所有可用的 just 命令
default:
    just --list

# 安装前端依赖包
[working-directory: 'client']
frontend-init:
    bun install

# 启动前端开发服务器
[working-directory: 'client']
frontend:
    bun run dev

# 下载并整理 Go 模块依赖
[working-directory: 'server']
backend-init:
    go mod tidy

# 启动后端服务
[working-directory: 'server']
backend:
    go run main.go

# 生成 JWT 令牌（用于测试或开发）
[working-directory: 'server']
get-jwt-token:
    go run scripts/generate_jwt.go

# 运行后端单元测试
[working-directory: 'server']
backend-test:
    go test -v ./...

# 运行 Go 代码静态检查
[working-directory: 'server']
backend-lint:
    go vet ./...

# 创建 tmux 开发环境窗口
tmux-dev:
    #!/usr/bin/env bash
    # 获取当前 tmux session 名称
    SESSION=$(tmux display-message -p '#S')

    # 创建 deamon 窗口（2个 pane）
    tmux rename-window deamon
    tmux split-window -t $SESSION:deamon -h

    # 创建 nvim 窗口
    tmux new-window -t $SESSION -n nvim

    # 创建 opencode 窗口
    tmux new-window -t $SESSION -n opencode

# 后台启动前后端服务（在 deamon 窗口的左右 pane 中运行）
dev:
    #!/usr/bin/env bash
    SESSION=$(tmux display-message -p '#S')

    # 在 deamon 窗口的左侧 pane 启动后端
    tmux send-keys -t $SESSION:deamon.1 'just backend' Enter

    # 在 deamon 窗口的右侧 pane 启动前端
    tmux send-keys -t $SESSION:deamon.2 'just frontend' Enter

    echo "✓ Backend started in deamon.1"
    echo "✓ Frontend started in deamon.2"
