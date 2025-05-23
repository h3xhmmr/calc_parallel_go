# Веб-сервис для решения простейших арифметических выражений

## Принцип работы
Веб-сервис состоит из двух микросервисов:
* Оркестратор, принимающий от пользователя выражение, делящий его на подзадачи вида "a+b" с помощью AST и RPN(будь они прокляты), отправляющий их агенту, после чего он собирает ответы в новое выражение, пока не придет к окончательному ответу
* Агент, запускающий горутины для вычисления подзадач, отправляемых оркестратором

Агент и оркестратор коннектятся через gRPC.


Для получения ответа на выражение пользователь периодически отправляет на сервер запрос с помощью метода GET, чтобы проверить состояние выражения/получить список всех выражений


## Вводимые данные


Сервер принимает запрос по url-ом `localhost/api/v1/calculate` и телом  `{
    "expression": "выражение, которое ввёл пользователь"
}`


Сервер поддерживает следующие арифметические операции с числами (именно числами, т.е. имеющими более 1 цифры в записи):  
* умножение
* деление
* сложение
* вычитание


А также учитывает приоритетные знаки (в том числе и скобки)


## Ошибки


Сервер будет выдавать ошибки в следующих случаях:  
* Введенное выражение имеет незакрытые скобки или несколько операторов, стоящих подряд ("invalid expression")
* В результате вычислений будет происходить деление на ноль ("division by zero")
* Введенное выражение имеет символы кроме символов операторов `()+-.*` и цифр `1234567890` ("invalid expression")
* Была введена пустая строка ("invalid body")
* При попытке отправить запрос не POST методом ("wrong method")


## Запуск проекта

Для начала установите Go на свой компьютер [тык](https://go.dev/doc/install), VisualStudio [тык](https://code.visualstudio.com/) и Git [тык](https://git-scm.com/downloads)

Далее, зайдя в VStudio необходимо клонировать репозиторий `git clone https://github.com/h3xhmmr/calc_parallel_go.git`, а затем запустить сервер (по умолчанию работает на порте :8080)
Запуск сервера:
* Сначала запускаем оркестратор `go run cmd/orch/main.go`
* Затем запускаем агента в новом терминале `go run cmd/agent/main.go`


## Примеры использования

### Регистрация и авторизация


#### Регистрация:


Запрос:


`curl --location 'localhost/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{ "login": , 
"password": 
}'`


Ответ:


200+OK (в случае успеха)



#### Авторизация:


`curl --location 'localhost/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{ "login": , 
"password": 
}'`



Ответ:

200+OK и JWT-токен, например, 

{
  "alg": "HS512",
  "typ": "JWT"
}


{
  "sub": "12345",
  "name": "John Gold",
  "admin": true
}


eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NSIsIm5hbWUiOiJKb2huIEdvbGQiLCJhZG1pbiI6dHJ1ZX0K.
LIHjWCBORSWMEibq-tnT8ue_deUqZx1K0XxCOXZRrBI



### Отправление выражения
Запрос:


`curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2"
}'`


Ответ:


`{
    "id": "1"
}`


### Получение списка выражений 
Запрос:


`curl --location 'localhost/api/v1/expressions'`


Ответ:


`{
    "expressions": [
        {
            "id": "1",
            "status": "in process",
            "result": "2+3*6"
        },
        {
            "id": "2",
            "status": "done",
            "result": "5.0"
        }
    ]
}`


### Получение выражения по id 
Запрос:


`curl --location 'localhost/api/v1/expressions/1'`


Ответ:


`{
    "expression":
        {
            "id": "1",
            "status": "in process",
            "result": "2-7"
        }
}`


Для удобства проверки советую использовать Postman [тык](https://www.postman.com/downloads/)
Для этого в строке адреса выберите метод Post и введите адрес `http://localhost:8080/api/v1/calculate`, после чего в поле body введите свой запрос в формате `{"expression": "2+2*2"}`


Для получения списков выражений и конкретных выражений выберите метод GET и введите адрес `localhost/api/v1/expressions` или `localhost/api/v1/expressions/:id` соответственно


Примеры использования через Postman лежат в папке example_postman


















P.S. не факт, что успею доделать проект, поэтому пусть хоть наброски будут на гитхабе лежать