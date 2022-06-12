# MEDODS-interview
JWT authentication service

## Про маршруты

### /login

Принимает как параметр правильный uuid\

Возвращает access и refresh токены и сообщение об ошибке в случае неудачи

**Request:**<br>
`POST /login `<br>
`{ "uuid" : <user-uuid>}`

**Response:**<br>
`200`<br>
`{  "access-token" : <access-token>,
"refresh-token" : <refresh-token> }`

![img.png](readme/img.png)

### /token

Принимает как параметры пару access, refresh токенов\

Возвращает access и refresh токены и сообщение об ошибке в случае неудачи

**Request:**<br>
`POST /token `<br>
`{  "access-token" : <access-token>,
    "refresh-token" : <refresh-token> }`

**Response:**<br>
`200`<br>
`{  "access-token" : <access-token>,
"refresh-token" : <refresh-token> }`

![img_1.png](readme/img_1.png)

## Про URI для подключения к MongoDB
В .env находится переменная MONGO_URI для подключения к своему кластеру.
