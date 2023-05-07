package main

// https://practicum.yandex.ru/learn/go-advanced/courses/1b5a1e55-b4dc-4dbe-ab70-207b0c2a3b59/sprints/89167/topics/8e2d6b94-5ee9-4c72-995f-8c51f31fd24c/lessons/30605f8d-fe8f-4a70-923b-c80a17d822a0/
import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/unitchecker"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

var excludedStyleChecks = map[string]struct{}{
	// Incorrect or missing package comment
	"ST1000": {},
	// The documentation of an exported function should start with the function's name
	"ST1020": {},
	// The documentation of an exported type should start with type's name
	"ST1021": {},
	// The documentation of an exported variable or constant should start with variable's name
	"ST1022": {},
}

var excludedStaticChecks = map[string]struct{}{
	//Storing non-pointer values in `sync.Pool` allocates memory
	//"SA6002": {},
	//// Field assignment that will never be observed
	//"SA4005": {},
}

var excludedQuickFixChecks = map[string]struct{}{}

func main() {
	analyzers := []*analysis.Analyzer{
		// report mismatches between assembly files and Go declarations
		asmdecl.Analyzer,
		// check for useless assignments
		assign.Analyzer,
		// check for common mistakes using the sync/atomic package
		atomic.Analyzer,
		// checks for non-64-bit-aligned arguments to sync/atomic functions
		atomicalign.Analyzer,
		// check for common mistakes involving boolean operators
		bools.Analyzer,
		buildssa.Analyzer,
		// check that +build tags are well-formed and correctly located
		buildtag.Analyzer,
		// detect some violations of the cgo pointer passing rules. Not working in Arcadia for now
		cgocall.Analyzer,
		// check for unkeyed composite literals
		composite.Analyzer,
		// check for locks erroneously passed by value
		copylock.Analyzer,
		ctrlflow.Analyzer,
		// check for the use of reflect.DeepEqual with error values
		deepequalerrors.Analyzer,
		directive.Analyzer,
		// check that the second argument to errors.As is a pointer to a type implementing error
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		// check for mistakes using HTTP response
		httpresponse.Analyzer,
		// check for impossible interface-to-interface type assertions
		ifaceassert.Analyzer,
		inspect.Analyzer,
		// check references to loop variables from within nested functions
		loopclosure.Analyzer,
		// check cancel func returned by context.WithCancel is called
		lostcancel.Analyzer,
		// check for useless comparisons between functions and nil
		nilfunc.Analyzer,
		// inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons
		nilness.Analyzer,
		pkgfact.Analyzer,
		// check consistency of Printf format strings and arguments
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		// check for possible unintended shadowing of variables EXPERIMENTAL
		shadow.Analyzer,
		// check for shifts that equal or exceed the width of the integer
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		// check signature of methods of well-known interfaces
		stdmethods.Analyzer,
		// check for string(int) conversions
		stringintconv.Analyzer,
		// check that struct field tags conform to reflect.StructTag.Get
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		// check for common mistaken usages of tests and examples
		tests.Analyzer,
		timeformat.Analyzer,
		// report passing non-pointer or non-interface values to unmarshal
		unmarshal.Analyzer,
		// check for unreachable code
		unreachable.Analyzer,
		// check for invalid conversions of uintptr to unsafe.Pointer
		unsafeptr.Analyzer,
		// check for unused results of calls to some functions
		unusedresult.Analyzer,
		// check for unused writes
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}

	// staticcheck
	// S
	for i, v := range simple.Analyzers {
		if i > 0 {
			break
		}
		analyzers = append(analyzers, v.Analyzer)
	}

	// SA
	for _, v := range staticcheck.Analyzers {
		if _, ok := excludedStaticChecks[v.Analyzer.Name]; ok {
			continue
		}
		analyzers = append(analyzers, v.Analyzer)
	}

	// ST
	for i, v := range stylecheck.Analyzers {
		if i > 0 {
			break
		}
		if _, ok := excludedStyleChecks[v.Analyzer.Name]; ok {
			continue
		}
		analyzers = append(analyzers, v.Analyzer)
	}

	// QF
	for i, v := range quickfix.Analyzers {
		if i > 0 {
			break
		}
		if _, ok := excludedQuickFixChecks[v.Analyzer.Name]; ok {
			continue
		}
		analyzers = append(analyzers, v.Analyzer)
	}

	unitchecker.Main(analyzers...)
}
