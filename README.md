# Задание

### В БД:
- Развернуть локально PostgreSQL.
- Создать свою базу данных. 
- Настроить своего пользователя. 
- Создать таблицы для хранения полученных данных.

### В сервисе:
- Подключение и подписка на канал в nats-streaming. 
- Полученные данные писать в PostgreSQL. 
- Так же полученные данные сохранить in-memory в сервисе (Кеш). 
- В случае падения сервиса восстанавливать Кеш из Postgres. 
- Поднять http сервер и выдавать данные по id из кеша. 
- Сделать простейший интерфейс отображения полученных данных, для их запроса по id.

### Дополнительная информация:
- Модель в файле model.json.
- Данные статичны, исходя из этого подумайте насчет модели хранения в Кеше и в pg. 
- В канал могут закинуть что угодно, подумайте как избежать проблем из-за этого.
- Чтобы проверить работает ли подписка онлайн, сделайте себе отдельны скрипт, для публикации данных в канал.
- Подумайте как не терять данные в случае ошибок или проблем с сервисом.
- Nats-streaming разверните локально (не путать с Nats).
   
### Бонус задание
- Покройте сервис автотестами. Будет плюсик вам в карму. 😊
- Устройте вашему сервису стресс тест, выясните на что он способен - воспользуйтесь утилитами WRK, Vegeta. Попробуйте оптимизировать код.

---

## Запуск

Чтобы выполненить миграции нужна утилита [migrate](https://github.com/golang-migrate/migrate).

Чтобы поднять PostgreSQL и nats-streaming потребуется [docker](https://www.docker.com).


Для запуска выполнить следующие команды:

```
make docker_up

make migration_up

make run_subscriber

make run_publisher
```

В качестве графического интерфейса открыть файл "L0/ui/index.html" в любом браузере.

Результаты стресс теста:

![Alt text](stress_test/Result.png)

![Alt text](stress_test/Plot.png)
