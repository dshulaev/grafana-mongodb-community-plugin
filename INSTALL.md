# Установка плагина MongoDB для Grafana (с поддержкой алертов)

Плагин `meln5674-mongodb-community` версии **0.3.0**: источник данных MongoDB для Grafana с поддержкой Grafana Alerting.

## Проверенные версии

| Grafana | MongoDB | Результат |
|---|---|---|
| 13.0.2 (OSS) | 7.0 | работает, алерты работают |
| 12.3.7 (OSS) | 8.0.11 | работает, алерты работают |

- **Grafana 9.0 и новее.** На 8.x и ниже плагин не работает.
- **Self-hosted Grafana (OSS или Enterprise).** В Grafana Cloud неподписанный плагин установить нельзя.
- **Бэкенд собран под** Linux (amd64, arm64, arm), Windows (amd64), macOS (amd64, arm64).

## 1. Скачать плагин

Скачайте `meln5674-mongodb-community-0.3.0.zip` со страницы релиза:
https://github.com/dshulaev/grafana-mongodb-community-plugin/releases/tag/v0.3.0

Репозиторий приватный: скачивать нужно под аккаунтом с доступом. Из командной строки:

```bash
gh release download v0.3.0 -R dshulaev/grafana-mongodb-community-plugin
```

Контрольная сумма (SHA-256) лежит в релизе в файле `.sha256`. Проверить её можно так:
- Linux: `sha256sum meln5674-mongodb-community-0.3.0.zip`
- Windows: `certutil -hashfile meln5674-mongodb-community-0.3.0.zip SHA256`

## 2. Разрешить неподписанный плагин

Плагин не подписан. Без этого шага Grafana его не загрузит.

Добавьте в конфиг Grafana:

```ini
[plugins]
allow_loading_unsigned_plugins = meln5674-mongodb-community
```

Где лежит конфиг:
- **Linux (deb/rpm):** `/etc/grafana/grafana.ini`
- **Windows:** `<папка Grafana>\conf\custom.ini`. Создайте файл, если его нет. `defaults.ini` не редактируйте.
- **Docker:** переменная окружения `GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=meln5674-mongodb-community`

Если в этой строке уже перечислены другие плагины, добавьте наш через запятую.

## 3. Установить плагин

Выберите любой из способов.

### Способ А: через `grafana cli`

**Linux:**
```bash
sudo grafana cli --pluginUrl /путь/к/meln5674-mongodb-community-0.3.0.zip plugins install meln5674-mongodb-community
```

**Windows** (командная строка, из папки `bin` Grafana):
```
grafana.exe cli --pluginUrl C:\путь\к\meln5674-mongodb-community-0.3.0.zip plugins install meln5674-mongodb-community
```

Если плагины Grafana лежат не в стандартной папке, добавьте `--pluginsDir <папка плагинов>`.

### Способ Б: распаковать вручную

Распакуйте архив в папку плагинов Grafana. Должна получиться папка `meln5674-mongodb-community`, а в ней `plugin.json`.

| Установка | Папка плагинов |
|---|---|
| Linux (deb/rpm) | `/var/lib/grafana/plugins/` |
| Windows | `<папка Grafana>\data\plugins\` |
| Docker | `/var/lib/grafana/plugins/` |

**Linux:**
```bash
sudo unzip meln5674-mongodb-community-0.3.0.zip -d /var/lib/grafana/plugins/
sudo chmod +x /var/lib/grafana/plugins/meln5674-mongodb-community/gpx_*
sudo chown -R grafana:grafana /var/lib/grafana/plugins/meln5674-mongodb-community
```

### Способ В: Docker / docker-compose

Распакуйте архив рядом с `docker-compose.yaml` и подключите папку как volume:

```yaml
services:
  grafana:
    image: grafana/grafana:12.3.7
    ports:
      - 3000:3000
    environment:
      GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS: meln5674-mongodb-community
    volumes:
      - ./meln5674-mongodb-community:/var/lib/grafana/plugins/meln5674-mongodb-community
      - grafana-data:/var/lib/grafana
volumes:
  grafana-data:
```

## 4. Перезапустить Grafana

- **Linux:** `sudo systemctl restart grafana-server`
- **Windows (служба):** `services.msc` → Grafana → Перезапустить. Если Grafana запущена вручную, закройте её и запустите снова.
- **Docker:** `docker compose restart grafana`

Как понять, что плагин загрузился:
- в **Administration → Plugins and data → Plugins** появился `mongodb-community`, версия 0.3.0;
- в логе Grafana есть предупреждение `Permitting unsigned plugin`. Это нормально.

## 5. Подключить MongoDB

**Connections → Data sources → Add new data source → mongodb-community**

| Поле | Что указать |
|---|---|
| URL | `mongodb://host:27017`. Для кластера: `mongodb://host1:27017,host2:27017/?replicaSet=rs0`. Для Atlas: `mongodb+srv://cluster.xxxxx.mongodb.net` |
| Username / Password | логин и пароль пользователя MongoDB. Пустые, если авторизация не нужна |
| TLS Enabled и др. | включите, если MongoDB требует TLS. Для Atlas TLS включается через `mongodb+srv` автоматически |

Если пользователь создан не в базе `admin`, добавьте в URL `?authSource=<база>`.

Нажмите **Save & test**. Должно появиться **«MongoDB is Responding»**.

## 6. Настроить запрос

Пример для коллекции, где документы выглядят так:
`{ "metadata": { "sensorId": 1 }, "timestamp": ISODate(...), "value": 22 }`

| Поле | Значение |
|---|---|
| Database.Collection | `test` . `weather` |
| QueryType | `Timeseries` (для алертов обязательно) |
| Timestamp Field | `timestamp` |
| Label Fields | `sensorId`: для каждого значения строится отдельная линия и создаётся отдельный алерт |
| Value Fields | `value` : `int32` |
| Automatic Time-Bound | вкл.: берутся только данные за выбранный период |
| Automatic Time-Sort | вкл. |
| Aggregation | `[{"$project": {"timestamp": 1, "sensorId": "$metadata.sensorId", "value": 1, "_id": 0}}]` |

Важно:
- **Вложенные поля не поддерживаются.** Вытаскивайте их наверх через `$project`, как `sensorId` в примере выше.
- **Тип поля значения должен совпадать с типом в MongoDB.** Тип видно в Compass рядом со значением:

  | Тип в MongoDB (Compass) | Тип в плагине |
  |---|---|
  | Int32 | `int32` |
  | Int64 / Long | `int64` |
  | Double | `float64` |
  | String | `string` |

  Если поле может отсутствовать в части документов, поставьте перед типом звёздочку: `*float64`.
- Вместо ручного указания типов можно включить **Infer Schema**: плагин определит типы по первым документам.
- Для фильтра по времени в Aggregation можно использовать макросы `$__from` и `$__to` (миллисекунды). Они работают и в алертах:
  ```json
  [{"$match": {"timestamp": {"$gte": {"$date": $__from}, "$lte": {"$date": $__to}}}}]
  ```
  Переменные дашборда (`$myvar`) в алертах не работают. Это ограничение Grafana.

## 7. Настроить алерт

**Alerting → Alert rules → New alert rule**

1. **Name:** например `Sensor too hot`.
2. **Query:** датасорс MongoDB, запрос как в разделе 6, период **Last 10 minutes**.
3. **Expressions:** **Reduce → Last**, затем **Threshold → IS ABOVE 80**.
4. **Folder и Evaluation group:** выберите или создайте. Интервал проверки, например **1m**.
5. **Contact point:** куда слать уведомления (см. раздел 8).
6. **Save rule and exit.**

Проверка: в списке правил состояние сменится на **Firing**, если условие выполняется. Если данных за период нет, правило уйдёт в **No Data**.

## 8. Уведомления

**Alerting → Contact points → Add contact point**

### Telegram
1. Создайте бота у **@BotFather** командой `/newbot` и получите **токен**.
2. Напишите боту любое сообщение (или добавьте его в группу и напишите там).
3. Узнайте **chat ID**: откройте `https://api.telegram.org/bot<ТОКЕН>/getUpdates`, нужное значение лежит в поле `"chat":{"id": ...}`. У групп ID отрицательный.
4. В Grafana: Integration **Telegram**, укажите **BOT API Token** и **Chat ID**, нажмите **Test** и сохраните.

### Email
Сначала настройте SMTP в конфиге Grafana и перезапустите её:

```ini
[smtp]
enabled = true
host = smtp.yandex.ru:465
user = alerts@yandex.ru
password = <пароль приложения>
from_address = alerts@yandex.ru
from_name = Grafana
```

Адреса SMTP: Gmail `smtp.gmail.com:587`, Яндекс `smtp.yandex.ru:465`, Mail.ru `smtp.mail.ru:465`. Нужен **пароль приложения**, обычный пароль от почты не подойдёт. Затем в Grafana создайте contact point: Integration **Email**, адрес получателя.

## Решение проблем

| Симптом | Причина и решение |
|---|---|
| Плагина нет в списке, в логе `plugin is unsigned` / `not allowed` | Не выполнен шаг 2, или Grafana не перезапущена после изменения конфига |
| В логе `permission denied` / `exec format error` (Linux) | Нет права на запуск: `chmod +x .../meln5674-mongodb-community/gpx_*` |
| Save & test: `Ping failed` / `server selection timeout` | Grafana не видит MongoDB. Проверьте URL, firewall и `bindIp` в настройках MongoDB |
| `authentication failed` | Неверный логин или пароль, либо нужен `?authSource=<база>` в URL |
| `Type mismatch for field ...` | Тип в Value Fields не совпадает с типом в MongoDB (см. таблицу типов в разделе 6) |
| `Field ... was null or absent` | Поле есть не во всех документах: поставьте `*` перед типом |
| Алерт в **No Data** | За период проверки нет документов. Проверьте Timestamp Field и период |
| Датасорс нельзя выбрать в алерте | Установлена старая версия плагина без поддержки алертов. Нужна 0.3.0 |
