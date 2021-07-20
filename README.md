go版本：1.16.6

web框架 ：iris

ORM框架：xrom

数据库：Mysql

验证登陆为用AES加密账号密码字符串所得hex作为token验证

登陆信息：
```json
{
    "code": 200,
    "msg": "Login success",
    "data": {
        "time": "2021-07-20 16:45:35",
        "token": "c8ecdee5f7668ffcee24486341c779016c6939a394f166c7953602cd191b73c64899309d376838c3dbb654b6636a84c6"
    }
}
```
个人信息：

```json
{
    "code": 200,
    "msg": "获取信息成功",
    "data": {
        "id": "1",
        "username": "jim",
        "age": 12,
        "gender": "male",
        "height": 165
    }
}

```



user表结构
|  username   | password  |
|  :----  | :----  |
|  jim  | 1234  |



user_information表结构

| id | username | age | gender | height |
| :-----| :----- | :----- | :----- | :----- |
| 1 | jim | 12 | male | 165 |


