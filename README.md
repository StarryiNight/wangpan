# wangpan

|       |           |          |
| ----- | :-------: | -------- |
| POST  |   /signup | 用户注册 |
| POST  |  /login   | 用户登陆 |
| GET   |  /logout  | 用户登出 |
| GET   |  /upload   | 下载 |
| GET   |  /download  | 上传|


download : url dir filename


upload   : 选择下载文件

中间件 JWTAuthMiddleware

支持断点续传和登录唯一
