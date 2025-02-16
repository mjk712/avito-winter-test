[![CI](https://github.com/mjk712/avito-winter-test/actions/workflows/ci.yaml/badge.svg)](https://github.com/mjk712/avito-winter-test/actions/workflows/ci.yaml)

[![Coverage Status](https://coveralls.io/repos/github/mjk712/avito-winter-test/badge.svg?branch=main)](https://coveralls.io/github/mjk712/avito-winter-test?branch=main)

Магазин мерча

Реализован сервис позволяющий обмениваться монетками и приобретать товары за монеты.

Подробное описание API находится в папке docs

Инструкция по запуску:

1.Настройка окружения - сервис настраивается при помощи переменных окружения:

SERVER_ADDRESS

POSTGRES_USERNAME

POSTGRES_PASSWORD

POSTGRES_HOST

POSTGRES_PORT

POSTGRES_DATABASE

POSTGRES_CONN

ENV

JWT_SECRET

2. Запуск сервиса осуществляется командой docker-compose up 

3. Процент покрытия тестами и статус сборки отображены в бейджах
