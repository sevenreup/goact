package types

import (
	"fmt"
	"go/ast"
	"strings"
)

func GenerateTs(structType *ast.TypeSpec) string {
	var tsFields []string
	if str, ok := structType.Type.(*ast.StructType); ok {
		for _, field := range str.Fields.List {
			fieldName := field.Names[0].Name
			goType := field.Type

			tsType := convertType(goType)
			tsFields = append(tsFields, fmt.Sprintf("%s: %s;", fieldName, tsType))
		}
		return fmt.Sprintf("export type %s = {%s};", structType.Name.Name, strings.Join(tsFields, " "))
	}
	return ""
}

func convertType(goType ast.Expr) string {
	switch t := goType.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			return "number"
		case "float32", "float64":
			return "number"
		case "string":
			return "string"
		case "bool":
			return "boolean"
		default:
			return "any"
		}
	case *ast.StarExpr:
		return convertType(t.X)
	case *ast.ArrayType:
		elemType := convertType(t.Elt)
		return fmt.Sprintf("Array<%s>", elemType)
	case *ast.MapType:
		keyType := convertType(t.Key)
		valueType := convertType(t.Value)
		return fmt.Sprintf("{[key: %s]: %s}", keyType, valueType)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel)
	default:
		return "any"
	}
}
