# Сервис баннеров
## Тестовое задание для стажеров Backend Avito-tech

## Описание задачи
Необходимо реализовать сервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега
пользователя, а также управлять баннерами и связанными с ними тегами и фичами.

## Стек
- Язык программировния Golang
- Фреймворк Gin для веб-сервера
- База данных PostgreSQL
- Кэш на основе Redis и пакета `go-redis/cache`
- Деплой с помощью Docker и Docker compose
- Миграции с импользованием docker образа `go-migrate`

## Запуск
### Запуск деплоя
- `cd deployment`
- `docker-compose us --build -d`
- просмотреть логи `docker-compose logs <service>`
- остановить `docker-compose down`
### Заупск тестов 
- Поскольку тесты e2e, они запускаются на отдельном наборе контейнеров, аналогично основному деплою.

## Проблемы
Основная проблема - не хватило времени покрыть все юнит тестами и добавить больше интеграционных тестов и e2e. Возможно большое количество багов.

# Api

### Добавлено

- GET /get_token/*admin - получение токена аутентификации/авторизацией. Админовский только при точном соответствии с *admin == "admin"  
Предполагается что sso сервис будет отдельным. Реализованы отдельные интерфейсы под аутентификацию, авторизацию и проверку на админа(предполагется, что будет онлайн, например через grpc)
- GET /versions/:id - получение предыдущих записей банера (максимум 3) по `id`
- PUT /versions/:id/activate - выбор версии для `id`, требуется параметр `version`
- DELETE /banners - удаление баннеров по фичи или id в соответствии с заданием

# Уточнения
Для масштабирования и более удобного тестирования (моки) сервис был разбит на слои с множеством интерфейсов.
- основная программа обернута в app
- httpserver в ginapp
- хэндлеры и инициализация маршрутов в getHandlers
- В хэндлеры прокидываются сервис и 3 интерфейса аутентификации\авторизации
- В сервис редис и постгрес

## Endpoints
- В получении всех баннеров имеется возможность получить все банеры, соответствующие конкретному тэгу, соответствующие конкретной фиче или конкретным тэгам и фичам. Все с учетом лимита и офсета. 
- В создании нового баннера предполгается, что нельзя передать банер с отсутствующими полями. В том числе, поскольку в бд отслеживается уникальность банера.

## Сервис
Для решения доп задания по удалению банеров по фиче или тэгам было реализован простейший механизм отложенных действий на основе каналов. Условием взятия запроса в работы выставлено количество текущих
активных запросов. Объективно, решение должно приниматься на основе метрик нагрузки системы, или дожен быть добвлен планировщик. 

## Кэш 
Размер локального кэша устанавливается через конфигурационный файл. Для тестового задания выбран обычный клиент, а не кластер или кольцо. В случае масштабирования слоистая архитектура позволяет переключится
на нужную конфигурацию. На текщий момент в кэш кладется тэг, фича и сам баннер. Первоначально предполагалось класть uuid, банер и множество всех пар тэгов, фич, uuid. Однако, не имея статистики по тэгам, было решено выбрать вариант проще.

## База даееых
- Для увеличения производительности используется `pgpool`. 
- При решении задания предполагалось, что количество запросов на получение баннеров сильно превышает действия админов. Соответственно, для быстрого получения баннеров и
проверки уникальности, было решено реализовать текущую архитекуру бд. Большая часть дейсвтий по модификации и удалению баннров реализованы с помощью триггеров. На больших данных модификация баннеров дорогая в угоду получения баннеров.
- Было указано, что баннеры могут временно отключаться. Преполагается, что они отключаются временно и таких баннеров не много, из-за этого неактивные баннеры вынесены в отдельную таблицу.