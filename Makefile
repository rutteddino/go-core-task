TASKS := 1 2 3 4 5 6 7 8 9
TASK  ?= 1

.PHONY: help test cover vet fmt check run build tidy

help: ## Показать список команд
	@echo "Доступные команды:"
	@echo "  make test          - запустить все тесты"
	@echo "  make cover         - тесты с покрытием кода"
	@echo "  make vet           - статический анализ (go vet)"
	@echo "  make fmt           - проверка форматирования (gofmt)"
	@echo "  make check         - fmt + vet + test"
	@echo "  make run TASK=N    - запустить задачу N (1-9)"
	@echo "  make build         - собрать все пакеты"
	@echo "  make tidy          - go mod tidy"

test: ## Запустить все тесты
	go test ./...

cover: ## Тесты с покрытием
	go test -cover ./...

vet: ## Статический анализ
	go vet ./...

fmt: ## Проверка форматирования
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "Не отформатированы файлы:"; echo "$$out"; exit 1; \
	fi

check: fmt vet test ## Полная проверка: форматирование, vet, тесты

run: ## Запуск задачи: make run TASK=3
	go run ./$(TASK)

build: ## Сборка всех пакетов
	go build ./...

tidy: ## Привести зависимости в порядок
	go mod tidy
