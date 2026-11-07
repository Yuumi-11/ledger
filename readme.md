# Ledger 个人记账后端

## 项目介绍

使用 Go、Gin、GORM 和 MySQL 开发的个人记账后端。
支持新增、查询列表、修改和删除收入、支出记录，通过 HTTP 接口返回 JSON。

## 技术栈

- Go：编写后端程序
- Gin：接收请求、分配路由、返回 JSON
- GORM：操作数据库
- MySQL：保存账目
- Postman：测试接口

## 项目结构

```text
ledger/
├── main.go                  程序入口
├── database/mysql.go        数据库连接
├── models/transaction.go    账目模型
├── handlers/transaction.go  增删改查处理函数
├── routes/routes.go         路由登记
├── SQL/schema.sql           数据库初始化脚本
├── go.mod                   模块与依赖
├── go.sum                   依赖校验信息
└── readme.md                项目说明
```

## 数据库设计

数据库名：ledger
表名：transactions

| 字段 | 含义 |
|---|---|
| id | 自动生成的账目编号，主键 |
| kind | income 为收入，expense 为支出 |
| amount_cents | 金额，以分为单位，必须大于 0 |
| category | 分类，不能为空或全为空白 |
| occurred_on | 收支发生日期 |
| note | 备注 |
| updated_at | 最近修改时间 |

例如，1250 分表示 12.50 元。

## 运行方法

1. 安装 Go 和 MySQL，启动 MySQL。
   Go 版本要求以 go.mod 中的 go 指令为准。

2. 在项目根目录打开终端，登录 MySQL：

   ```sh
   mysql -u root -p
   ```

   若 macOS 安装后 mysql 不在 PATH 中，可以使用：

   ```sh
   /usr/local/mysql/bin/mysql -u root -p
   ```

3. 在 mysql> 提示符下执行：

   ```sql
   SOURCE SQL/schema.sql;
   exit;
   ```

4. 回到普通终端，下载依赖并运行：

   ```sh
   go mod download
   go run .
   ```

5. 按提示输入 MySQL root 密码，输入时不显示字符。

当前连接配置：
- 数据库地址：127.0.0.1:3306
- 数据库名：ledger
- 用户名：root
- 密码：启动时输入，不写入源码

服务地址：http://127.0.0.1:8080

测试期间保持程序运行，Control + C 停止服务。

## 接口说明

| 方式 | 路径 | 功能 |
|---|---|---|
| GET | /ping | 检查服务是否启动 |
| GET | /transactions | 查询全部账目 |
| POST | /transactions | 新增账目 |
| PUT | /transactions/:id | 修改指定账目 |
| DELETE | /transactions/:id | 删除指定账目 |

:id 替换为实际编号，例如 /transactions/5。

POST 和 PUT 使用 Body → raw → JSON，
请求头 Content-Type 为 application/json。

请求示例：

```json
{
  "kind": "expense",
  "amount_cents": 1250,
  "category": "餐饮",
  "occurred_on": "2026-11-07",
  "note": "接口测试"
}
```

PUT 需要提交修改后的完整账目内容，不是只提交金额。
备注可以填写空字符串。

GET 和 DELETE 不需要请求 Body。

响应状态：
- 200：查询、修改或删除成功
- 201：新增成功
- 400：请求输入不合法
- 404：指定账目不存在
- 500：服务端处理失败

查询和新增返回 message 与 data；
修改和删除成功返回 message 与 id。

当前返回的账目字段使用 Go 结构体字段名，
例如 ID、KIND、AmountCents；
请求字段使用 kind、amount_cents 等 JSON 名称。

## Postman 测试步骤

### 正常流程

1. POST 新增账目，预期 201，记录返回的 ID。
2. GET 查询列表，确认包含刚才新增的 ID 和内容。
3. PUT 修改该 ID 的金额为 1500，预期 200。
4. GET 查询，确认同一笔账金额变为 1500。
5. DELETE 删除该 ID，预期 200。
6. GET 查询，确认该 ID 消失，其他账目仍在。
7. 再次 DELETE 同一 ID，预期 404。

### 异常输入

每次只改变一个字段，其余字段保持合法：

| 测试 | 预期 |
|---|---|
| POST 金额为 -100 | 400 |
| POST 类型为 abc | 400 |
| POST 分类全为空格 | 400 |
| POST 日期为 2026-02-30 | 400 |
| DELETE 编号为 abc | 400 |
| DELETE 编号为 0 | 400 |
| PUT 已删除的编号，Body 合法 | 404 |

错误新增请求执行后，再通过 GET 确认没有保存错误账目。