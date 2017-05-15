# Weburg. Бот Telegram

Вебург у вас в телеграме. Weburg Times: [@weburg_times](http://telegram.me/weburg_times). Новинки видео: [@weburg_movies](http://telegram.me/weburg_movies). Новинки музыки: [@weburg_music](http://telegram.me/weburg_music). Новинки сериалов: [@weburg_series](http://telegram.me/weburg_series).

## Запуск

```
$ make
$ ./bin/weburg-telegram-bot start -w
```

## Параметры

```
--token, -t      Telegram API токен [$WEBURG_BOT_TOKEN]
--watch, -w      Запустить в режиме демона [$WEBURG_BOT_WATCH]
--debug, -d      Режим отладки [$WEBURG_BOT_DEBUG]
--no-color, --nc Отключение цветов в логах [$WEBURG_BOT_NO_COLOR]
--listen-address Адрес для веб-интерфейса и телеметрии [$WEBURG_BOT_LISTEN_ADDR]
--metrics-path   Путь, по которому будут доступны метрики [$WEBURG_BOT_METRICS_PATH]
--database-path  Путь к файлу базы данных [$WEBURG_BOT_DATABASE_PATH]
```

## Метрики

```
$ curl http://127.0.0.1:9109/metrics
```

* `pulls_total_count` – число запросов RSS-лент;
* `pulls_feed_total_count` – число запросов по лентам;
* `pulls_fail_count` – число ошибок, при запросах RSS-лент;
* `pulls_feed_fail_count` – число ошибок, при запросах по лентам;
* `messages_total_count` – число отправленных сообщений в Telegram;
* `messages_feed_total_count` – число отправленных сообщений по лентам;
* `messages_fail_count` – число ошибок, при отправке сообщений в Telegram;
* `messages_feed_fail_count` – число ошибок, при отправке сообщений по лентам.
