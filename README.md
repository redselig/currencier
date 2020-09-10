[![Build Status](https://travis-ci.com/redselig/currencier.svg?branch=master)](https://travis-ci.com/redselig/currencier)

[![Go Report Card](https://goreportcard.com/badge/github.com/redselig/currencier)](https://goreportcard.com/report/github.com/redselig/currencier)

Currencier - валютный торговец =)

Необходимо реализовать сервис со следующим функционалом на языке Golang  >= 1.13.

В базе данных Mysql/PosgreSql должна быть таблица currency c колонками:
id — первичный ключ
name — название валюты
rate — курс валюты к рублю
insert_dt – время обновления валюты

Должна быть команда для обновления данных в таблице currency по расписанию. 
Данные по курсам валют можно взять отсюда: http://www.cbr.ru/scripts/XML_daily.asp
Реализовать 2 REST API метода:
GET /currencies — должен возвращать список курсов валют с возможность пагинации
GET /currency/ — должен возвращать курс валюты для переданного id

Добавил ручку GET /lazycurrencies для более оптимальной пагинации(ленивой) без offset в БД

Порт: 4444

/currency/R01589

/currencies?limit=10&offset=5

/lazycurrencies?limit=10&lastid=R01589