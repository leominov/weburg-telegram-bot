# Weburg. Бот Telegram

Вебург у вас в телеграме. Weburg Times: [@weburg_times](http://telegram.me/weburg_times). Новинки видео: [@weburg_movies](http://telegram.me/weburg_movies). Новинки музыки: [@weburg_music](http://telegram.me/weburg_music). Новинки сериалов: [@weburg_series](http://telegram.me/weburg_series).

## Запуск

```
$ make
$ ./bin/weburg-telegram-bot
```

## Статус

```
$ curl http://127.0.0.1:9109/metrics
```

## Параметры

```
--token, -t      Your Telegram API token [$WEBURG_BOT_TOKEN]
--rss-watch, -r  Enable RSS watching [$WEBURG_BOT_RSS_WATCH]
--debug, -d      Enable debug mode [$WEBURG_BOT_DEBUG]
--no-color, --nc Don't show colors in logging [$WEBURG_BOT_NO_COLOR]
--listen-address Address to listen on for web interface and telemetry [$WEBURG_BOT_LISTEN_ADDR]
--metrics-path   Path under which to expose metrics [$WEBURG_BOT_METRICS_PATH]
```
