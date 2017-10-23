# Weburg. Бот Telegram

Вебург у вас в телеграме. Weburg Times: [@weburg_times](http://telegram.me/weburg_times). Новинки видео: [@weburg_movies](http://telegram.me/weburg_movies). Новинки музыки: [@weburg_music](http://telegram.me/weburg_music). Новинки сериалов: [@weburg_series](http://telegram.me/weburg_series).

## Запуск

```
$ make
$ ./bin/weburg-telegram-bot start -w
```

## Параметры запуска

```
--token, -t         none          Telegram API токен [$WEBURG_BOT_TOKEN]
--watch, -w         false         Запустить в режиме демона [$WEBURG_BOT_WATCH]
--debug, -d         false         Режим отладки [$WEBURG_BOT_DEBUG]
--no-color, -nc     false         Отключение цветов в логах [$WEBURG_BOT_NO_COLOR]
--listen-address    0.0.0.0:9109  Адрес для веб-интерфейса и телеметрии [$WEBURG_BOT_LISTEN_ADDR]
--metrics-path      /metrics      Путь, по которому будут доступны метрики [$WEBURG_BOT_METRICS_PATH]
--database-path     ./database.db Путь к файлу базы данных [$WEBURG_BOT_DATABASE_PATH]
--config-file       ./config.yaml Путь к файлу конфигурации [$WEBURG_BOT_CONFIG_FILE]
--disable-messenger false         Отключить отправку сообщений в Telegram [$WEBURG_BOT_DISABLE_MESSENGER]
```

## Файл конфигурации

```
---
token: 123456:1234567890
watch: true
listen_addr: 127.0.0.1:9109
metrics_path: /metrics
database_path: ./database.db
disable_messenger: false
agents:
  - name: movies
    endpoint:
      type: rss
      url: http://rss.weburg.net/movies/all.rss
    interval: 1m
    channel:
      type: channel
      username: weburg_movies
    cache_size: 3
    print_categories: true
    skip_categories:
      - 18
  - name: series
    endpoint:
      type: clever_title_series
      url: http://weburg.net/series/all/?clever_title=1&template=0&last=0&sorts=date_update
    interval: 1m
    channel:
      type: channel
      username: weburg_series
    cache_size: 2
    print_categories: true
    print_description: true
```

Все значения, за исключением `agents`, могут быть переопределены переменными окружения и параметрами запуска.

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
