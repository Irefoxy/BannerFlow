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
В корне проекта находится `Makefile`. Для начала работы выполнить `make + <command>`
Commands:
- unit_test
- test
- fuzze
- deploy
- stop

## Пояснения
# Api
Первоначально было сгенерировано с помощью `oapi-codegen`. После обнаружения изменений в файле спецификации и для решения дополнительных заданий было принято решение модифицировать api, 
изменения в спецификацию не вносились из-за нехватки времени.

# Endpoints
- В получении всех баннеров имеется возможность получить все банеры, соответствующие конкретному тэгу, соответствующие конкретной фиче или конкретным тэгам и фичам. Все с учетом лимита и офсета. 
- В создании нового баннера предполгается, что нельзя передать банер с отсутствующими полями. В том числе, поскольку в бд отслеживается уникальность банера.

# Сервис
Для решения доп задания по удалению банеров по фиче или тэгам было реализован простейший механизм отложенных действий на основе каналов. Условием взятия запроса в работы выставлено количество текущих
активных запросов. Объективно, решение должно приниматься на основе метрик нагрузки системы, или дожен быть добвлен планировщик. 

# Кэш 
Размер локального кэша устанавливается через конфигурационный файл. Для тестового задания выбран обычный клиент, а не кластер или кольцо. В случае масштабирования слоистая архитектура позволяет переключится
на нужную конфигурацию. На текщий момент в кэш кладется тэг, фича и сам баннер. Первоначально предполагалось класть uuid, банер и множество всех пар тэгов, фич, uuid. Однако, не имея статистики по тэгам, было решено выбрать вариант проще.

# База даееых
- Для увеличения производительности используется `pgpool`. 
- При решении задания предполагалось, что количество запросов на получение баннеров сильно превышает действия админов. Соответственно, для быстрого получения баннеров и
проверки уникальности, было решено реализовать текущую архитекуру бд. Большая часть дейсвтий по модификации и удалению баннров реализованы с помощью триггеров. На больших данных модификация баннеров дорогая в угоду получения баннеров.
- Было указано, что баннеры могут временно отключаться. Преполагается, что они отключаются временно и таких баннеров не много, из-за этого неактивные баннеры вынесены в отдельную таблицу.

# Аутентификация и авторизация
Хотелось бы отдельным контейнером поднять sso сервер, однако на текущий момент реализован генератор токенов

## Условия
1. Используйте этот [API](https://github.com/avito-tech/backend-trainee-assignment-2024/blob/main/api.yaml)
2. Тегов и фичей небольшое количество (до 1000), RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%
3. Для авторизации доступов должны использоваться 2 вида токенов: пользовательский и админский.  Получение баннера может происходить с помощью пользовательского или админского токена, а все остальные действия могут выполняться только с помощью админского токена.
4. Реализуйте интеграционный или E2E-тест на сценарий получения баннера.
5. Если при получении баннера передан флаг use_last_revision, необходимо отдавать самую актуальную информацию.  В ином случае допускается передача информации, которая была актуальна 5 минут назад.
6. Баннеры могут быть временно выключены. Если баннер выключен, то обычные пользователи не должны его получать, при этом админы должны иметь к нему доступ.

## Дополнительные задания:
Эти задания не являются обязательными, но выполнение всех или части из них даст вам преимущество перед другими кандидатами.
1. Адаптировать систему для значительного увеличения количества тегов и фичей, при котором допускается увеличение времени исполнения по редко запрашиваемым тегам и фичам
2. Провести нагрузочное тестирование полученного решения и приложить результаты тестирования к решению
3. Иногда получается так, что необходимо вернуться к одной из трех предыдущих версий баннера в связи с найденной ошибкой в логике, тексте и т.д.  Измените API таким образом, чтобы можно было просмотреть существующие версии баннера и выбрать подходящую версию
4. Добавить метод удаления баннеров по фиче или тегу, время ответа которого не должно превышать 100 мс, независимо от количества баннеров.  В связи с небольшим временем ответа метода, рекомендуется ознакомиться с механизмом выполнения отложенных действий
5. Реализовать интеграционное или E2E-тестирование для остальных сценариев
6. Описать конфигурацию линтера