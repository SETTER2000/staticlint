# staticlint
staticlint analysis

### Стандартный статический анализатор 

Пакеты входящие в него:

golang.org/x/tools/go/analysis/passes;

все анализаторы класса SA пакета staticcheck.io;

один анализатор из пакета staticcheck.io;

два публичных анализатора.

###  Скомпилировать
``go build -o staticlint main.go``

### Использование

Скомпилировать и положить в корень проекта.

#### Запускаем проверку

```go vet -vettool=$(which staticlint) ./...```

