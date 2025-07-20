# helper-sender-bot
Приложение для отслеживания чатов дежурств и упрощения работы с МР

## Документация
- [описание какие боли может закрывать сервис](docs/for_what.md)
- [описание поддерживаемых сценариев](docs/use_cases.md)
- [open API](docs/openapi.yaml)
- [seq диаграммы дежурств](docs/duty.puml)
- [seq диаграммы gitlab](docs/gitlab.puml)

## Локальный запуск
- Установите зависимости необходимые для проекта
- Поднимите все [контейнеры](docker/docker-compose.yml)
- Для gitlab может потребоваться сброс пароля:
  - `docker exec -it gitlab  gitlab-rake "gitlab:password:reset[root]"` где root имя пользователя
- Задайте локальные переменные за основу взять [пример енвов](.example.env)
- Запустите [миграции](cmd/db-init/main.go). Для миграции енв переменная должна быть `MIGRATION_COMMAND=up`
- Запустите [приложение](cmd/app/main.go)
- Запустите [воркеры](cmd/workers/main.go)

