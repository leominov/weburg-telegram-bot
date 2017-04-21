# Weburg. Бот Telegram

Вебург у вас в телеграме. Weburg Times: [@weburg_times](http://telegram.me/weburg_times). Новинки видео: [@weburg_movies](http://telegram.me/weburg_movies). Новинки музыки: [@weburg_music](http://telegram.me/weburg_music). Новинки сериалов: [@weburg_series](http://telegram.me/weburg_series).

## Запуск

```
$ make
$ ./bin/weburg-telegram-bot start
```

## Параметры

```
--token, -t      Telegram API токен [$WEBURG_BOT_TOKEN]
--rss-watch, -r  Запустить в режиме демона [$WEBURG_BOT_RSS_WATCH]
--debug, -d      Режим отладки [$WEBURG_BOT_DEBUG]
--no-color, --nc Отключение цветов в логах [$WEBURG_BOT_NO_COLOR]
--listen-address Адрес для веб-интерфейса и телеметрии [$WEBURG_BOT_LISTEN_ADDR]
--metrics-path   Путь, по которому будут доступны метрики [$WEBURG_BOT_METRICS_PATH]
```

## Метрики

```
$ curl http://127.0.0.1:9109/metrics
```

* `pulls_total_count` – число запросов RSS-лент;
* `pulls_fail_count` – число ошибок, при запросах RSS-лент;
* `messages_total_count` – число отправленных сообщений в Telegram;
* `messages_fail_count` – число ошибок, при отправке сообщений в Telegram.
