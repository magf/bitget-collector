Bitget Trade Collector
Bitget Trade Collector — это приложение на языке Go для сбора данных о сделках с криптовалютной биржи Bitget через WebSocket и их хранения в базах данных SQLite. Каждая торговая пара (например, BTCUSDT, ETHUSDT) обрабатывается отдельным экземпляром systemd-сервиса, что обеспечивает независимость и отказоустойчивость. Проект поставляется в виде DEB-пакета для удобной установки на системах Debian/Ubuntu.
Возможности

Сбор данных о сделках (ID сделки, символ, цена, объем, сторона, временная метка) для указанных торговых пар.
Хранение данных в отдельных базах SQLite для каждой пары (например, trades_BTCUSDT.db) с использованием Write-Ahead Logging (WAL) для неблокирующего доступа.
Поддержка нескольких торговых пар через шаблонные systemd-сервисы (collector@<pair>.service).
Опциональный режим отладки для диагностики.
Автоматическая настройка пользователя и каталога при установке DEB-пакета.
Автоматический запуск сервиса для BTCUSDT после установки.

Требования

Go 1.18 или новее (для сборки из исходников).
Система Debian/Ubuntu для установки DEB-пакета.
Утилита sqlite3 для просмотра базы данных.
Утилита dpkg для сборки и установки DEB-пакета.
Доступ в интернет для подключения к WebSocket API Bitget.

Установка
Сборка из исходников

Клонируйте репозиторий:git clone https://github.com/magf/bitget-collector.git
cd bitget-collector


Инициализируйте модули Go и установите зависимости:go mod init bitget-collector
go get github.com/gorilla/websocket
go get github.com/mattn/go-sqlite3


Соберите бинарный файл:make build



Сборка DEB-пакета

Убедитесь, что установлен dpkg-deb:
sudo apt install dpkg


Соберите DEB-пакет:
make deb


Установите DEB-пакет:
sudo dpkg -i bitget-collector.deb

Процесс установки:

Создаёт системного пользователя bitget.
Создаёт каталог /var/lib/bitget-collector с соответствующими правами.
Устанавливает шаблонный systemd-сервис.
Автоматически активирует и запускает сервис collector@BTCUSDT.



Запуск сервиса

Запустите экземпляр сервиса для конкретной торговой пары (например, ETHUSDT):sudo systemctl start collector@ETHUSDT
sudo systemctl enable collector@ETHUSDT


Проверьте статус сервиса:sudo systemctl status collector@ETHUSDT


Просмотрите собранные данные:sqlite3 /var/lib/bitget-collector/trades_ETHUSDT.db "SELECT * FROM trades LIMIT 10;"



Использование

Запустите вручную для тестирования с отладочным выводом:./collector -pair=BTCUSDT -debug


Используйте make run для запуска с парой по умолчанию (BTCUSDT): make run


Данные сохраняются в /var/lib/bitget-collector/trades_<pair>.db.

Структура проекта
bitget-collector/
├── cmd/
│   └── collector/        # Точка входа приложения
├── pkg/                  # Плейсхолдер для будущих Go-пакетов
├── debian/               # Конфигурация DEB-пакета (control, systemd-сервис, скрипты установки)
├── scripts/              # Скрипты для сборки DEB-пакета
├── LICENSE               # Лицензия MIT
├── README.md             # Документация на английском
├── README-ru.md          # Документация на русском
├── Makefile              # Автоматизация сборки
├── go.mod                # Зависимости Go-модуля
├── go.sum                # Контрольные суммы зависимостей
└── .gitignore            # Правила исключения для Git

Вклад в проект
Приветствуются любые улучшения! Пожалуйста, создавайте issues или pull requests в репозитории. Убедитесь, что изменения протестированы и соответствуют стилю кода проекта.
Лицензия
Проект распространяется под лицензией MIT. Подробности в файле LICENSE.
Автор
Максим Гайдай maxim.gajdaj@gmail.com
