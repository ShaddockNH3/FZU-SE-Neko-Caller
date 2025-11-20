# Neko Caller - 课堂随机点名系统

> 福州大学软件工程结对编程作业

一个支持多种点名模式、积分管理和趣味事件的智能课堂点名系统。

## ✨ 功能特性

- 🎲 **多种点名模式** - 加权随机、顺序、逆序、低分优先
- 📊 **积分管理** - 自动计算积分，支持自定义分数
- 🎉 **趣味事件** - 双倍积分、疯狂星期四、1024福报等
- 🔄 **点名转移权** - 每被点名2次获得1次转移机会
- 📥 **数据导入导出** - 支持Excel格式批量导入学生名单
- 📈 **排行榜统计** - 查看班级积分排名和数据分析

## 🚀 快速开始

### 后端启动

```bash
cd NekoCallerBackend
go mod tidy
go run main.go
```

服务运行在 `http://localhost:8888`

### 前端启动

```bash
cd NekoCallerFrountend
npm install
npm run dev
```

前端运行在 `http://localhost:5173`

## 🛠️ 技术栈

**后端**
- Go + CloudWeGo Hertz
- GORM + MySQL
- Thrift IDL

**前端**
- Vue 3 + Element Plus
- Vite + Vue Router
- Axios

## 📖 使用说明

详细功能说明和API文档请参考 [docs](/docs/)

## 📝 开发团队

- [ShaddockNH3](https://github.com/ShaddockNH3)
- [nieie](https://github.com/nieie)

## 📄 许可证

[MIT License](LICENSE)
