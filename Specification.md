## 数据库

## Table

### user_basic

| 字段     | 类型    | 大小 | 备注            |
| -------- | ------- | ---- | --------------- |
| id       | bigint  |      | 主键自增id      |
| username | varchar | 100  | 用户名称        |
| password | char    | 16   | MD5后的用户密码 |
| email    | varchar | 50   | 用户邮箱        |
| phone    | char    | 11   | 用户手机号      |
| gender   | tinyint | 1    | 0 男 1 女       |
| name     | varchar | 50   | 用户姓名        |
| identity | char    | 18   | 身份证号码      |
| is_admin | tinyint | 1    | 0 普通 1 管理员 |
|          |         |      |                 |

### order_basic

| 字段         | 类型     | 大小 | 备注         |
| ------------ | -------- | ---- | ------------ |
| id           | bigint   | 8    | id           |
| identity     | char     | 18   | 订单唯一标识 |
| created_time | datetime |      | 创建时间     |
| user_id      | bigint   |      | 该订单的用户 |
| medicine_id  | bigin    |      | 该订单的药品 |

### medicine_basic

| 字段                | 类型    | 大小 | 备注                                        |
| ------------------- | ------- | ---- | ------------------------------------------- |
| id                  | bigint  |      | id                                          |
| medicine_sn         | char    | 13   | 药品条形码                                  |
| medicine_name       | varchar | 50   | 药品名称                                    |
| medicine_valid_date | date    |      | 有效日期                                    |
| medicine_price      | decimal |      | 价格                                        |
| medicine_stock      | int     |      | 库存                                        |
| medicine_type       | tinyint |      | 商品分类 0 功能性食品 1 特医食品 2 天然食品 |
| medicine_material   | varchar | 600  | 产品配方                                    |
| company_name        | varchar | 50   | 企业名称                                    |

### order_medicine

| 字段            | 类型    | 大小 | 备注     |
| --------------- | ------- | ---- | -------- |
| id              | bigint  |      | id       |
| order_id        |         |      |          |
| medicine_id     |         |      |          |
| medicine_name   | varchar | 50   | 药品名称 |
| medicine_price  | decimal |      | 价格     |
| medicine_amount | int     |      | 药品数量 |

## 功能

1. 登录功能
2. 注册功能
3. 添加药品信息
4. 修改药品信息
5. 删除药品信息
6. 分页展示药品列表
7. 分页展示用户列表
8. 分页展示订单列表
9. 关键字查询订单
10. 关键字查询药品
11. 关键字查询用户

## 后端

登录功能

注册功能

添加药品信息

修改药品信息

删除药品信息

分页展示药品列表

分页展示用户列表

分页展示订单列表

关键字查询订单

关键字查询药品

关键字查询用户

### 返回参数

## 前端

### 请求url
