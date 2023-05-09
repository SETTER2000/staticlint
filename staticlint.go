package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "staticlint", // имя анализатора
	// Текст с описанием работы анализатора. Этот текст будет отображаться по команде help,
	// поэтому его нужно сделать многострочным и описать в нём все флаги анализатора.
	Doc:      `анализатор, запрещает использовать прямой вызов os.Exit в функции main пакета main`,
	Run:      run, // функция, которая отвечает за анализ исходного кода
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type Pass struct {
	Fset         *token.FileSet // информация о позиции токенов
	Files        []*ast.File    // AST для каждого файла
	OtherFiles   []string       // имена файлов не на Go в пакете
	IgnoredFiles []string       // имена игнорируемых исходных файлов в пакете
	Pkg          *types.Package // информация о типах пакета
	TypesInfo    *types.Info    // информация о типах в AST
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder(nil, func(n ast.Node) {

	})
	return nil, nil
}

// Чтобы отследить игнорируемые ошибки, нужно знать, возвращает ли функция ошибку.
// Для этого определим переменную errorType, которая равна интерфейсному типу error,
// а также функцию isErrorType, которая определяет, соответствует ли тип ошибке.
// Воспользуемся методами пакета стандартной библиотеки go/types
// — этот пакет предоставляет инструментарий для работы c типами языка Go.
var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}
