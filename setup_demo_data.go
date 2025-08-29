package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const POCKETBASE_URL = "http://localhost:8090"

func main() {
	fmt.Println("开始设置演示数据...")

	// 1. 创建 posts 集合
	fmt.Println("1. 创建 posts 集合...")
	postsCollection := map[string]interface{}{
		"name": "posts",
		"type": "base",
		"schema": []map[string]interface{}{
			{
				"name":     "title",
				"type":     "text",
				"required": true,
			},
			{
				"name":     "content",
				"type":     "editor",
				"required": true,
			},
			{
				"name":     "author",
				"type":     "text",
				"required": true,
			},
			{
				"name":     "published",
				"type":     "bool",
				"required": false,
			},
			{
				"name":     "tags",
				"type":     "text",
				"required": false,
			},
		},
	}

	createCollection(postsCollection)

	// 2. 创建示例记录
	fmt.Println("2. 创建示例记录...")
	samplePosts := []map[string]interface{}{
		{
			"title":     "欢迎使用 PocketBase",
			"content":   "这是一个使用 PocketBase 创建的示例文章。PocketBase 是一个开源的实时后端和数据库。",
			"author":    "系统管理员",
			"published": true,
			"tags":      "pocketbase,backend,database",
		},
		{
			"title":     "MCP 服务器集成",
			"content":   "我们成功集成了 Model Context Protocol (MCP) 服务器，可以通过 AI 助手来操作数据库。",
			"author":    "开发者",
			"published": true,
			"tags":      "mcp,ai,integration",
		},
		{
			"title":     "测试文章",
			"content":   "这是一篇测试文章，用于验证我们的系统功能。",
			"author":    "测试用户",
			"published": false,
			"tags":      "test,demo",
		},
	}

	for _, post := range samplePosts {
		createRecord("posts", post)
	}

	fmt.Println("🎉 演示数据设置完成！")
	fmt.Println("📊 现在你可以：")
	fmt.Println("   - 访问 http://localhost:8090/_/ 查看 PocketBase 管理界面")
	fmt.Println("   - 访问 http://localhost:8093/ 测试 MCP 功能")
	fmt.Println("   - 使用测试页面上的按钮来测试各种功能")
}

func createCollection(collection map[string]interface{}) {
	jsonData, _ := json.Marshal(collection)
	resp, err := http.Post(POCKETBASE_URL+"/api/collections", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 创建集合失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Printf("✅ %s 集合创建成功\n", collection["name"])
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("⚠️  %s 集合可能已存在或创建失败: %s\n", collection["name"], string(body))
	}
}

func createRecord(collectionName string, record map[string]interface{}) {
	jsonData, _ := json.Marshal(record)
	resp, err := http.Post(POCKETBASE_URL+"/api/collections/"+collectionName+"/records", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 创建记录失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		if id, ok := result["id"].(string); ok {
			fmt.Printf("✅ 记录创建成功: %s (ID: %s)\n", record["title"], id)
		} else {
			fmt.Printf("✅ 记录创建成功: %s\n", record["title"])
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ 记录创建失败: %s - %s\n", record["title"], string(body))
	}
}
