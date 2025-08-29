package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
)

const (
	PocketBaseURL = "http://localhost:8090"
	MCPPort       = "8093"
)

//go:embed static
var StaticFiles embed.FS

func main() {
	// 创建 PocketBase 实例
	app := pocketbase.New()

	// 设置数据目录
	dataDir := filepath.Join(".", "pb_data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// 创建 MCP 服务器
	mcpServer := server.NewMCPServer(
		"PocketBase MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// 添加 PocketBase 控制工具
	addPocketBaseTools(mcpServer, app)

	// 设置 HTTP 路由
	http.HandleFunc("/mcp/sse", func(w http.ResponseWriter, r *http.Request) {
		handleSSE(w, r, mcpServer)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "healthy",
			"service":   "pocketbase-mcp",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 静态文件服务
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/static/test.html"
		}
		http.FileServer(http.FS(StaticFiles)).ServeHTTP(w, r)
	})

	// 启动 PocketBase
	go func() {
		app.Bootstrap()
		serveCmd := cmd.NewServeCommand(app, false)
		if err := serveCmd.Execute(); err != nil {
			log.Fatal(err)
		}
	}()

	// 启动 MCP HTTP 服务器
	log.Printf("Starting MCP server with SSE support on port %s", MCPPort)
	log.Printf("SSE endpoint: http://localhost:%s/mcp/sse", MCPPort)
	log.Printf("Health check: http://localhost:%s/health", MCPPort)
	log.Printf("Test page: http://localhost:%s/", MCPPort)
	log.Printf("PocketBase admin UI: http://localhost:8090/_/")

	if err := http.ListenAndServe(":"+MCPPort, nil); err != nil {
		log.Fatalf("Failed to start MCP HTTP server: %v", err)
	}
}

func handleSSE(w http.ResponseWriter, r *http.Request, mcpServer *server.MCPServer) {
	// 设置 SSE 头部
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 创建 SSE 连接
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// 发送连接建立消息
	fmt.Fprintf(w, "data: %s\n\n", `{"type": "connection", "status": "connected"}`)
	flusher.Flush()

	// 保持连接活跃
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// 处理客户端断开连接
	notify := w.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-notify:
			log.Println("Client disconnected")
			return
		case <-ticker.C:
			// 发送心跳
			fmt.Fprintf(w, "data: %s\n\n", `{"type": "heartbeat", "timestamp": "`+time.Now().Format(time.RFC3339)+`"}`)
			flusher.Flush()
		}
	}
}

func addPocketBaseTools(mcpServer *server.MCPServer, app *pocketbase.PocketBase) {
	// 创建集合工具
	createCollectionTool := mcp.NewTool(
		"create_collection",
		mcp.WithDescription("Create a new collection in PocketBase"),
		mcp.WithString("name", mcp.Description("Collection name"), mcp.Required()),
		mcp.WithString("type", mcp.Description("Collection type (base or view)")),
		mcp.WithObject("schema", mcp.Description("Collection schema fields")),
	)

	mcpServer.AddTool(createCollectionTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := req.GetString("name", "")
		collectionType := req.GetString("type", "base")

		// 这里实现创建集合的逻辑
		// 注意：实际的 PocketBase API 需要更复杂的实现
		result := fmt.Sprintf("Collection '%s' of type '%s' would be created", name, collectionType)

		return mcp.NewToolResultText(result), nil
	})

	// 获取集合列表工具
	listCollectionsTool := mcp.NewTool(
		"list_collections",
		mcp.WithDescription("List all collections in PocketBase"),
	)

	mcpServer.AddTool(listCollectionsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 使用 Go SDK 直接获取集合，这应该有管理员权限
		collections, err := app.FindAllCollections("base")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get collections: %v", err)), nil
		}

		var result []string
		for _, collection := range collections {
			result = append(result, collection.Name)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Collections: %v", result)), nil
	})

	// 创建记录工具
	createRecordTool := mcp.NewTool(
		"create_record",
		mcp.WithDescription("Create a new record in a collection"),
		mcp.WithString("collection", mcp.Description("Collection name"), mcp.Required()),
		mcp.WithObject("data", mcp.Description("Record data")),
	)

	mcpServer.AddTool(createRecordTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionName := req.GetString("collection", "")

		collection, err := app.FindCollectionByNameOrId(collectionName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Collection not found: %v", err)), nil
		}

		record := core.NewRecord(collection)

		// 获取数据参数
		if data := req.GetRawArguments(); data != nil {
			if dataMap, ok := data.(map[string]any); ok {
				if dataObj, exists := dataMap["data"]; exists {
					if dataMap, ok := dataObj.(map[string]any); ok {
						for key, value := range dataMap {
							record.Set(key, value)
						}
					}
				}
			}
		}

		if err := app.Save(record); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to save record: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Record created with ID: %s", record.Id)), nil
	})

	// 查询记录工具
	queryRecordsTool := mcp.NewTool(
		"query_records",
		mcp.WithDescription("Query records from a collection"),
		mcp.WithString("collection", mcp.Description("Collection name"), mcp.Required()),
		mcp.WithString("filter", mcp.Description("Filter expression")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of records to return")),
	)

	mcpServer.AddTool(queryRecordsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionName := req.GetString("collection", "")
		filter := req.GetString("filter", "")
		limit := req.GetInt("limit", 0)

		// 使用 Go SDK 的 FindRecordsByFilter 方法
		records, err := app.FindRecordsByFilter(collectionName, filter, "", limit, 0)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to query records: %v", err)), nil
		}

		result := fmt.Sprintf("Found %d records", len(records))
		return mcp.NewToolResultText(result), nil
	})

	// 删除记录工具
	deleteRecordTool := mcp.NewTool(
		"delete_record",
		mcp.WithDescription("Delete a record from a collection"),
		mcp.WithString("collection", mcp.Description("Collection name"), mcp.Required()),
		mcp.WithString("id", mcp.Description("Record ID"), mcp.Required()),
	)

	mcpServer.AddTool(deleteRecordTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionName := req.GetString("collection", "")
		recordID := req.GetString("id", "")

		collection, err := app.FindCollectionByNameOrId(collectionName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Collection not found: %v", err)), nil
		}

		record, err := app.FindRecordById(collection, recordID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Record not found: %v", err)), nil
		}

		if err := app.Delete(record); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete record: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Record %s deleted successfully", recordID)), nil
	})

	// 服务器状态工具
	serverStatusTool := mcp.NewTool(
		"server_status",
		mcp.WithDescription("Get PocketBase server status"),
	)

	mcpServer.AddTool(serverStatusTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 使用 Go SDK 检查服务器状态
		collections, err := app.FindAllCollections()
		if err != nil {
			return mcp.NewToolResultText("PocketBase server is not responding: " + err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("PocketBase server is running and healthy. Found %d collections", len(collections))), nil
	})

	// 获取 PocketBase 信息工具
	getPocketBaseInfoTool := mcp.NewTool(
		"get_pocketbase_info",
		mcp.WithDescription("Get PocketBase server information"),
	)

	mcpServer.AddTool(getPocketBaseInfoTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collections, err := app.FindAllCollections()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get collections: %v", err)), nil
		}

		info := fmt.Sprintf("Admin UI: http://localhost:8090/_/, API: http://localhost:8090/api/, Status: running, Collections: %d", len(collections))
		return mcp.NewToolResultText(info), nil
	})
}
