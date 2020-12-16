# Reverse Market Backend

В данной части курсового проекта необходимо было разработать серверную часть приложение на языке программирования `Go`, которое бы предоставляло REST API для взаимодействия с клиентской частью.

## Участники
Проект выполнил студент группы 3530904/80102: Кожин С.В.

# Этапы проекта 

## Поставнока проблемы 

Проект решает проблему обмена вещами между людьми. Периодически каждому человеку нужно продавать ненужные вещи или покупать какие-то вещи. Проект позволяет найти определенные товары у людей, которые хотят их продать, и наоборот, найти человека, который готов купить ваш товар.

## Требования и диаграммы

Требования и диаграммы соотвтетсвуют таковым для клиентской части приложения в дополнении с компонентной диаграмой для серверной части.

- Component diagram (API Application)

![image](https://user-images.githubusercontent.com/27823412/102302406-0c471380-3f6a-11eb-907b-45b2eadb5f8a.png)

# Кодирование и отладка

Серверная часть системы реализована на языке `Go`, данные хранятся с использованием СУБД `PostgreSQL`. Взаимодействие между серверной и клиентской частью происходит по протоколу HTTP, посредством данных закодированных  в формате JSON по стандарту REST. Cборка и развёртка происходят в `Docker` контейнерах.

API представлет собой полный набор эндпоинтов для управленями данными пользователей, объявлений, предложений, а также просмотра категорий и тегов.

Сервер размещен на сервисе `AWS (Amazon Web Services)` по адресу:   
http://ec2-3-133-94-51.us-east-2.compute.amazonaws.com

# Тестирование
Для юнит-тестирования были использованы стандарные библиотеки `Go` [`testing`](https://golang.org/pkg/testing/) и [`httptest`](https://golang.org/pkg/net/http/httptest/).

# Сборка

Для начала необходимо склонить проект из github:  

    git clone https://github.com/reverse-market/backend.git

Далее необходимо установить [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/).

После установки сервис можно запустить следующей командой в корневой папке проекта:  

    docker-compose up 
    
Для завершения работы:

    docker-compose down

После запуска сервис будет работать по адресу:
http://localhost:8080

Для запуска юнит-тестов необходимо выполнить следующую команду:

    docker-compose -f docker-compose.test.yml up --exit-code-from test-app
