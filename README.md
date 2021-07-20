go版本：1.16.6

web框架 ：iris

ORM框架：xrom

数据库：Mysql

验证登陆为用AES加密账号密码字符串所得hex作为token验证

user表结构
|  username   | password  |
|  :----  | :----  |
|  jim  | 1234  |

user_information表结构
|  id   | username  |  age  | gender  | height |
|  :----  | :----  |:----|:----|
|  1  | jim | 12 | male| 165|


