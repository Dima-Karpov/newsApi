

### migrate
    1. docker-compose exec db bash
    2. psql -h db -U news -d api -f /migrations/schema/up.sql (накатываем миграции)
    3. psql -U news -d api
    4. \dt (увидеть все таблицы)

    5. docker-compose run migrate (можно накатить миграции так но попросит пароль от db)

