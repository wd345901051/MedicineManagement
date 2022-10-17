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
| identity | char    | 18   | 身份证号        |
| is_admin | tinyint | 1    | 0 普通 1 管理员 |

### order_basic

| 字段          | 类型     | 大小 | 备注         |
| ------------- | -------- | ---- | ------------ |
| id            | bigint   | 8    | id           |
| identity      | char     | 18   | 订单唯一标识 |
| created_time  | datetime |      | 创建时间     |
| user_id       | bigint   |      | 该订单的用户 |
| medicine_list | varchar  | 500  | 该订单的药品 |
| total_price   | decimal  | (,2) | 订单总金额   |

### medicine_basic

| 字段                   | 类型     | 大小 | 备注                                                               |
| ---------------------- | -------- | ---- | ------------------------------------------------------------------ |
| id                     | bigint   |      | id                                                                 |
| medicine_sn            | char     | 13   | 药品条形码                                                         |
| medicine_name          | varchar  | 100  | 药品名称                                                           |
| company_name           | varchar  | 50   | 供应商                                                             |
| medicine_valid_date    | datetime |      | 保质期                                                             |
| medicine_price         | decimal  | (,2) | 价格                                                               |
| medicine_stock         | int      |      | 库存                                                               |
| medicine_specification | varchar  | 50   | 规格                                                               |
| medicine_type          | tinyint  | 1    | 商品分类 0 全营养配方食品 1 特定全营养配方食品  2 非全营养配方食品 |
| medicine_apply         | varchar  | 50   | 适用人群                                                           |
| medicine_material      | varchar  | 600  | 产品配方                                                           |

### order_medicine

| 字段            | 类型   | 大小 | 备注     |
| --------------- | ------ | ---- | -------- |
| id              | bigint |      | id       |
| order_id        | bigint |      | 订单id   |
| medicine_id     | bigint |      | 药品id   |
| medicine_amount | int    |      | 药品数量 |

## 功能

1. 登录功能
2. 添加药品信息
3. 修改药品信息
4. 删除药品信息
5. 分页展示药品列表
6. 分页展示用户列表
7. 分页展示订单列表
8. 关键字查询订单
9. 关键字查询药品
10. 关键字查询用户

## 后端

### 登录功能

##### 请求消息

```http
POST  /mgr/login  HTTP/1.1
Content-Type:   application/x-www-form-urlencoded
```

##### 请求参数

http请求消息体中，参数以x-www-form-urlencoded格式存储：

需要携带以下参数：

username：用户名

password：密码

##### 响应消息

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http响应消息体中，数据以json格式存储

如果登录成功，则返回如下信息：

```json
{
    "ret": 0
}
```

ret 为 0 表示登录成功

如果登录失败，则返回状态及失败的原因：

```json
{
    "ret": 1,  
    "msg":  "用户名或者密码错误"
}
```

ret 为1表示登录失败， msg字段描述登录失败的原因

### 药品

#### 添加药品信息

##### 请求消息

```http
POST  /mgr/medicines  HTTP/1.1
Content-Type:   application/json
```

##### 请求参数

http 请求消息体携带添加药品的信息，消息体的格式是json，如下示例：

```json
{
    "action":"add_medicine",
    "data":{
        "medicine_name": "贝因美特殊医学用途婴儿无乳糖配方食品",
        "medicine_sn": "TY20180001",
        "company_name": "杭州贝因美母婴营养品有限公司",
        "medicine_valid_date":"xxxx-xx-xx",
        "medicine_price": 218.00,
        "medicine_stock": 99,
        "medicine_specification": "400g",
        "medicine_type":1,
        "medicine_material":"水解乳清蛋白粉"  
    }
}
```

其中

`action` 字段固定填写 `add_medicine` 表示添加一个药品

`data` 字段中存储了要添加的药品的信息

服务端接收到该请求后，应在系统中增加该药品。

##### 响应消息

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http 响应消息 body 中， 数据以json格式存储，

如果添加成功，返回如下

```json
{
    "ret": 0,
    "id" : 1
}
```

ret 为 0 表示成功，id 为 添加药品的id号。

如果添加失败，返回失败的原因，示例如下

```json
{
    "ret": 1,  
    "msg": "药品名已经存在"
}
```

ret 为1表示失败， msg字段描述添加失败的原因

#### 修改药品信息

##### 请求消息

```http
POST  /mgr/medicines  HTTP/1.1
Content-Type:   application/json
```

##### 请求参数

http 请求消息体携带添加药品的信息，消息体的格式是json，如下示例：

```json
{
    "action":"update_medicine",
    "id":1,
    "newdata":{
        "medicine_name": "贝因美特殊医学用途婴儿无乳糖配方食品",
        "medicine_sn": "TY20180001",
        "company_name": "杭州贝因美母婴营养品有限公司",
        "medicine_valid_date":"xxxx-xx-xx",
        "medicine_price": 218.00,
        "medicine_stock": 99,
        "medicine_specification": "400g",
        "medicine_type":1,
        "medicine_material":"水解乳清蛋白粉"  
    }
}
```

其中

`action` 字段固定填写 `update_medicine` 表示修改一个药品的信息

`id`表示要修改的药品 `id`

`newdata` 字段中存储了要修改的药品的信息

服务端接收到该请求后，应在系统中修改该药品的信息。

##### 响应消息

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http 响应消息 body 中， 数据以json格式存储，

如果修改成功，返回如下

```json
{
    "ret": 0
}
```

ret 为 0 表示成功，id 为 添加药品的id号。

如果添加失败，返回失败的原因，示例如下

```json
{
    "ret": 1,  
    "msg": "修改失败"
}
```

ret 为1表示失败， msg字段描述添加失败的原因

#### 删除药品信息

##### 请求消息

```http
DELETE  /api/mgr/medicines  HTTP/1.1
Content-Type:   application/json
```

##### 请求参数

http 请求消息 body 携带要删除药品的id

消息体的格式是json，如下示例：

```json
{
    "action":"del_medicine",
    "id": 6
}
```

其中

`action` 字段固定填写 `del_medicine` 表示删除一个药品

`id` 字段为要删除的药品的id号

服务端接受到该请求后，应该在系统中尝试删除该id对应的药品。

##### 响应消息

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http 响应消息 body 中， 数据以json格式存储，

如果删除成功，返回如下

```json
{
    "ret": 0
}
```

ret 为 0 表示成功。

如果删除失败，返回失败的原因，示例如下

```json
{
    "ret": 1,  
    "msg": "删除失败"
}
```

ret 不为 0 表示失败， msg字段描述失败的原因

#### 分页展示药品列表

##### 请求消息

```
GET  /medicines  HTTP/1.1
```

##### 请求参数

http 请求消息 url 中 需要携带如下参数，

- action

  必填项，填写值为 list_medicine
- pagesize

  必填项，用来分页，确定每页最大记录条数
- pagenum

  必填项，获取第几页的信息
- keywords

  可选项， 里面包含多个过滤关键字，关键字之间用空格分开

##### 响应消息

```
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http 响应消息体中， 数据以json格式存储，

如果获取信息成功，返回如下

```json
{
    "ret": 0,
    "retlist": [
        {
            "id":1,
            "medicine_name": "贝因美特殊医学用途婴儿无乳糖配方食品",
            "medicine_sn": "TY20180001",
            "company_name": "杭州贝因美母婴营养品有限公司",
            "medicine_valid_date":"xxxx-xx-xx",
            "medicine_price": 218.00,
            "medicine_stock": 99,
            "medicine_specification": "400g",
            "medicine_type":1,
            "medicine_material":"水解乳清蛋白粉" 
        }
     ],
     "total":1
        
}
```

ret 为0表示登录成功

retlist 里面包含了当前请求页的药品信息列表

total 表示药品的数量

每个药品信息以如下格式存储：

```json
{
    "id":1,
    "medicine_name": "贝因美特殊医学用途婴儿无乳糖配方食品",
    "medicine_sn": "TY20180001",
    "company_name": "杭州贝因美母婴营养品有限公司",
    "medicine_valid_date":"xxxx-xx-xx",
    "medicine_price": 218.00,
    "medicine_stock": 99,
    "medicine_specification": "400g",
    "medicine_type":1,
    "medicine_material":"水解乳清蛋白粉" 
} 
```

### 用户

#### 分页展示用户列表

##### 请求消息

```http
GET  /mgr/customers  HTTP/1.1
```

##### 请求参数

http 请求消息 url 中 需要携带如下参数，

- action

  必填项，填写值为 list_customer
- pagesize

  必填项，分页的 每页获取多少条记录
- pagenum

  必填项，获取第几页的信息
- keywords

  可选项， 里面包含的多个过滤关键字，关键字之间用 `空格` 分开

##### 响应消息

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http 响应消息 body 中， 数据以json格式存储，

如果获取信息成功，返回如下

```json
{
    "ret": 0,
    "retlist": [
        {
            "name":"张三",
            "gender":1,
            "email":"1236159221@163.com",
            "phone":"17391370206",
            "identity":"610721******5"
        },
        {
            "name":"李四",
            "gender":0,
            "email":"158624229@163.com",
            "phone":"18322657453",
            "identity":"821721******5"
        }
    ] , 
    "total": 2           
}
```

ret 为0表示登录成功

retlist 包含了当前请求页的客户信息列表。

total 表示系统中所有用户的数量

每个客户信息以如下格式存储:

```json
{
    "name":"李四",
    "gender":0,
    "email":"158624229@163.com",
    "phone":"18322657453",
    "identity":"821721******25"
}
```

### 订单

#### 分页展示订单列表

##### 请求消息

```http
GET  /mgr/orders  HTTP/1.1
```

##### 请求参数

http 请求消息 url 中 需要携带如下参数，

- action

  必填项，填写值为 list_customer
- pagesize

  必填项，分页的 每页获取多少条记录
- pagenum

  必填项，获取第几页的信息
- keywords

  可选项， 里面包含的多个过滤关键字，关键字之间用 `空格` 分开

##### 响应消息

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

##### 响应内容

http 响应消息 body 中， 数据以json格式存储，

如果获取信息成功，返回如下

```json
{
    "ret": 0,
    "retlist": [
        {
            "name":"张三",
            "created_time":"2022-10-16 15:04:05",
            "medicinelist":[
            {"id":1,"name":"贝因美特殊医学用途婴儿无乳糖配方食品","amount":1,"price":218.00},
        	{"id":2,"name":"易关舒","amount":1,"price":12.34}
            ]
        }
    ] , 
    "total": 2           
}
```

ret 为0表示登录成功

retlist 包含了当前请求页的订单信息列表。

total 表示系统中所有订单的数量

每个客户信息以如下格式存储:

```json
{
    "name":"张三",
    "created_time":"2022-10-16 15:04:05",
    "medicinelist":[
        {"id":1,"name":"贝因美特殊医学用途婴儿无乳糖配方食品","amount":1,"price":218.00},
        {"id":2,"name":"易关舒","amount":1,"price":12.34}
    ]
}
```

## 前端

### 请求url
