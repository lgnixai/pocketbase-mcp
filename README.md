# PocketBase MCP 服务

一个集成了 Model Context Protocol (MCP) 的 PocketBase 服务器，允许 AI 助手通过 MCP 协议直接操作 PocketBase 数据库。

## 🚀 功能特性

### 核心功能
- **PocketBase 集成**: 嵌入式 PocketBase 数据库服务器
- **MCP 协议支持**: 完整的 Model Context Protocol 实现
- **HTTP 服务器**: 提供 Web 界面和 API 端点
- **SSE 支持**: Server-Sent Events 实时通信
- **Web 测试界面**: 内置测试页面，方便调试和演示

### MCP 工具
- `create_collection` - 创建新的集合
- `list_collections` - 列出所有集合
- `create_record` - 在集合中创建记录
- `query_records` - 查询集合中的记录
- `delete_record` - 删除指定记录
- `server_status` - 获取服务器状态
- `get_pocketbase_info` - 获取 PocketBase 信息

## 📋 系统要求

- Go 1.23.6 或更高版本
- 现代浏览器（支持 SSE）

## 🛠️ 安装和运行

### 1. 克隆项目
```bash
git clone <repository-url>
cd ok
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 编译项目
```bash
go build -o aicms main.go
```

### 4. 运行服务
```bash
./aicms
```

## 🌐 访问地址

启动后，你可以访问以下地址：

- **PocketBase 管理界面**: http://localhost:8090/_/
- **MCP 测试页面**: http://localhost:8093/
- **健康检查**: http://localhost:8093/health
- **SSE 端点**: http://localhost:8093/mcp/sse

## 📖 使用指南

### 1. 创建集合和记录

#### 通过 PocketBase 管理界面
1. 访问 http://localhost:8090/_/
2. 点击 "Collections" → "New collection"
3. 设置集合名称和字段
4. 创建记录

#### 通过 MCP 工具
使用 AI 助手调用 MCP 工具：
```
请创建一个名为 "posts" 的集合，包含 title、content、author 字段
```

### 2. 测试 MCP 功能

访问 http://localhost:8093/ 测试页面：

1. **连接状态**: 检查 MCP 服务器和 PocketBase 连接状态
2. **工具测试**: 测试各种 MCP 工具功能
3. **实时日志**: 查看操作日志和结果

### 3. 与 AI 助手集成

MCP 服务器支持与支持 MCP 协议的 AI 助手集成，如：
- Claude Desktop
- 其他支持 MCP 的客户端

## 🔧 配置说明

### 端口配置
- PocketBase: 8090
- MCP HTTP 服务器: 8093

### 数据存储
- 数据目录: `./pb_data/`
- 自动创建和管理

## 🏗️ 项目结构

```
ok/
├── main.go              # 主程序入口
├── static/
│   └── test.html        # 测试页面
├── setup_demo_data.go   # 演示数据设置脚本
├── go.mod               # Go 模块文件
├── go.sum               # 依赖校验文件
└── README.md           # 项目文档
```

## 🔍 开发说明

### 添加新的 MCP 工具

1. 在 `addPocketBaseTools` 函数中定义新工具
2. 使用 `mcp.NewTool` 创建工具
3. 实现工具处理函数
4. 使用 `mcpServer.AddTool` 注册工具

示例：
```go
newTool := mcp.NewTool(
    "tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param", mcp.Description("Parameter description")),
)

mcpServer.AddTool(newTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // 工具实现逻辑
    return mcp.NewToolResultText("Result"), nil
})
```

### 扩展功能

- **认证支持**: 可以添加用户认证和权限控制
- **更多工具**: 根据需求添加更多 MCP 工具
- **数据验证**: 增强数据验证和错误处理
- **监控日志**: 添加详细的日志记录和监控

## 🐛 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   lsof -i :8090  # 检查 PocketBase 端口
   lsof -i :8093  # 检查 MCP 端口
   kill <PID>     # 停止占用进程
   ```

2. **权限错误**
   - 确保数据目录有写入权限
   - 检查文件系统权限

3. **连接失败**
   - 检查防火墙设置
   - 确认端口未被占用

### 日志查看
```bash
./aicms 2>&1 | tee app.log
```

## 📄 许可证

本项目基于 MIT 许可证开源。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如有问题，请：
1. 查看本文档
2. 检查故障排除部分
3. 提交 Issue

---

**注意**: 这是一个演示项目，生产环境使用前请进行充分的安全评估和配置。
