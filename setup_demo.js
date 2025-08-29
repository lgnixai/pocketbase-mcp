// 演示数据设置脚本

const POCKETBASE_URL = 'http://localhost:8090';

async function setupDemo() {
    console.log('开始设置演示数据...');

    try {
        // 1. 创建 posts 集合
        console.log('1. 创建 posts 集合...');
        const postsCollection = {
            name: "posts",
            type: "base",
            schema: [
                {
                    name: "title",
                    type: "text",
                    required: true
                },
                {
                    name: "content",
                    type: "editor",
                    required: true
                },
                {
                    name: "author",
                    type: "text",
                    required: true
                },
                {
                    name: "published",
                    type: "bool",
                    required: false
                },
                {
                    name: "tags",
                    type: "text",
                    required: false
                }
            ]
        };

        const createCollectionResponse = await fetch(`${POCKETBASE_URL}/api/collections`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(postsCollection)
        });

        if (createCollectionResponse.ok) {
            console.log('✅ posts 集合创建成功');
        } else {
            const error = await createCollectionResponse.text();
            console.log('⚠️  posts 集合可能已存在或创建失败:', error);
        }

        // 2. 创建一些示例记录
        console.log('2. 创建示例记录...');
        const samplePosts = [
            {
                title: "欢迎使用 PocketBase",
                content: "这是一个使用 PocketBase 创建的示例文章。PocketBase 是一个开源的实时后端和数据库。",
                author: "系统管理员",
                published: true,
                tags: "pocketbase,backend,database"
            },
            {
                title: "MCP 服务器集成",
                content: "我们成功集成了 Model Context Protocol (MCP) 服务器，可以通过 AI 助手来操作数据库。",
                author: "开发者",
                published: true,
                tags: "mcp,ai,integration"
            },
            {
                title: "测试文章",
                content: "这是一篇测试文章，用于验证我们的系统功能。",
                author: "测试用户",
                published: false,
                tags: "test,demo"
            }
        ];

        for (const post of samplePosts) {
            const createRecordResponse = await fetch(`${POCKETBASE_URL}/api/collections/posts/records`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(post)
            });

            if (createRecordResponse.ok) {
                const result = await createRecordResponse.json();
                console.log(`✅ 记录创建成功: ${post.title} (ID: ${result.id})`);
            } else {
                const error = await createRecordResponse.text();
                console.log(`❌ 记录创建失败: ${post.title}`, error);
            }
        }

        // 3. 创建 users 集合
        console.log('3. 创建 users 集合...');
        const usersCollection = {
            name: "users",
            type: "auth",
            schema: [
                {
                    name: "name",
                    type: "text",
                    required: true
                },
                {
                    name: "avatar",
                    type: "file",
                    required: false
                },
                {
                    name: "bio",
                    type: "text",
                    required: false
                }
            ]
        };

        const createUsersResponse = await fetch(`${POCKETBASE_URL}/api/collections`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(usersCollection)
        });

        if (createUsersResponse.ok) {
            console.log('✅ users 集合创建成功');
        } else {
            const error = await createUsersResponse.text();
            console.log('⚠️  users 集合可能已存在或创建失败:', error);
        }

        console.log('🎉 演示数据设置完成！');
        console.log('📊 现在你可以：');
        console.log('   - 访问 http://localhost:8090/_/ 查看 PocketBase 管理界面');
        console.log('   - 访问 http://localhost:8093/ 测试 MCP 功能');
        console.log('   - 使用测试页面上的按钮来测试各种功能');

    } catch (error) {
        console.error('❌ 设置演示数据时出错:', error.message);
    }
}

// 运行设置
setupDemo();
