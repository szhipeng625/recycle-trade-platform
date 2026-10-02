# 废旧物资回收交易平台

## 技术栈
- **前端**: Vue 3 + Vite + Element Plus + ECharts + Axios
- **后端**: Go 1.24 + Gin + GORM + JWT
- **数据库**: MySQL 8.0

## 项目结构
```
废旧物资回收交易平台-v2/
├── server/              # Go 后端
│   ├── main.go          # 入口
│   ├── config/          # 配置
│   ├── models/          # 数据模型
│   ├── controllers/     # 控制器
│   ├── middleware/      # 中间件(JWT鉴权)
│   ├── routes/          # 路由
│   ├── utils/           # 工具函数
│   └── server.exe       # 编译后的可执行文件
├── client/              # Vue 前端
│   ├── src/
│   │   ├── views/       # 页面组件
│   │   │   ├── admin/   # 管理员页面
│   │   │   └── user/    # 用户页面
│   │   ├── layouts/     # 布局组件
│   │   ├── router/      # 路由
│   │   └── api/         # API封装
│   └── package.json
└── 启动.bat             # 一键启动脚本
```

## 功能模块

### 管理员端
- **数据概览**: 今日交易额/收货量/交易笔数统计、价格走势折线图、库存概览
- **每日定价**: 设置各类货物每日回收单价
- **交易开单**: 选择卖货人、货物、重量(支持吨/公斤/市斤)，自动计算金额，自动入库
- **货物管理**: 增删改查货物类别
- **卖货人管理**: 管理卖货人信息
- **仓库库存**: 库存总览、手动入库/出库、出入库流水记录

### 用户端
- **今日行情**: 查看今日各类货物回收价格、价格走势折线图
- **我要卖货**: 卖货估价计算器(支持吨/公斤/市斤换算)
- **我的记录**: 查看交易记录

## 单位换算
- 1 吨 = 1000 公斤
- 1 市斤 = 0.5 公斤
- 所有库存和统计统一以公斤为基准单位

## 快速启动

### 方式一：一键启动（推荐）
双击 `启动.bat`

### 方式二：手动启动
```bash
# 1. 确保 MySQL 已启动
# 2. 启动后端
cd server
server.exe

# 3. 启动前端（新开终端）
cd client
npx vite --host
```

## 访问地址
- 前端: http://localhost:5173
- 后端: http://localhost:9000
- 管理员账号: `admin` / `admin123`
- 用户可自行注册

## 数据库配置
默认配置（可在 server/config/config.go 修改）：
- 地址: 127.0.0.1:3306
- 数据库: recycle_trade
- 用户: root
- 密码: 12345678

首次启动会自动创建数据表并初始化默认数据（管理员、5种货物、今日定价）。

## API 接口
所有接口前缀 `/api`，需登录的接口在 Header 中携带 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /auth/login | 登录 | 公开 |
| POST | /auth/register | 注册 | 公开 |
| GET | /auth/profile | 当前用户信息 | 登录 |
| GET | /goods | 货物列表 | 登录 |
| POST | /goods | 新增货物 | 管理员 |
| GET | /pricing/today | 今日定价 | 登录 |
| GET | /pricing/trend | 价格走势 | 登录 |
| POST | /pricing | 设置定价 | 管理员 |
| GET | /transactions | 交易记录 | 登录 |
| POST | /transactions | 开单结算 | 管理员 |
| GET | /transactions/stats | 今日统计 | 登录 |
| GET | /sellers | 卖货人列表 | 登录 |
| POST | /sellers | 新增卖货人 | 管理员 |
| GET | /inventory | 库存列表 | 登录 |
| POST | /inventory/in | 入库 | 管理员 |
| POST | /inventory/out | 出库 | 管理员 |
| GET | /inventory/records | 出入库流水 | 登录 |
