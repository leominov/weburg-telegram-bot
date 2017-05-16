#!/bin/bash

if ! hash "php" 2>/dev/null; then
    echo "PHP built-in web server required to start."
    exit 1
fi

php -S 0.0.0.0:8081 -t ./rss/
