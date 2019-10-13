package codeanalysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	CodeDir    string
	GopathDir  string
	VendorDir  string
	IgnoreDirs []string
}

type AnalysisResult interface {
	OutputToFile(logfile string)
}

func AnalysisCode(config Config) AnalysisResult {
	tool := &analysisTool{
		interfaceMetas:              []*interfaceMeta{},
		structMetas:                 []*structMeta{},
		packagePathPackageNameCache: map[string]string{},
	}
	tool.analysis(config)
	return tool
}

func HasPrefixInSomeElement(value string, src []string) bool {
	result := false
	for _, srcValue := range src {
		if strings.HasPrefix(value, srcValue) {
			result = true
			break
		}
	}
	return result
}

func sliceContains(src []string, value string) bool {
	isContain := false
	for _, srcValue := range src {
		if srcValue == value {
			isContain = true
			break
		}
	}
	return isContain
}

func sliceContainsSlice(s []string, s2 []string) bool {
	for _, str := range s2 {
		if !sliceContains(s, str) {
			return false
		}
	}
	return true
}

func mapContains(src map[string]string, key string) bool {
	if _, ok := src[key]; ok {
		return true
	}
	return false
}

func findGoPackageNameInDirPath(dirpath string) string {

	dir_list, e := ioutil.ReadDir(dirpath)

	if e != nil {
		log.Printf("读取目录%s文件列表失败,%s\n", dirpath, e)
		return ""
	}

	for _, fileInfo := range dir_list {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".go") {
			packageName := ParsePackageNameFromGoFile(path.Join(dirpath, fileInfo.Name()))
			if packageName != "" {
				return packageName
			}
		}
	}

	return ""
}

func ParsePackageNameFromGoFile(filepath string) string {

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)

	if err != nil {
		log.Printf("解析文件%s失败, %s\n", filepath, err)
		return ""
	}

	return file.Name.Name

}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type baseInfo struct {
	// go文件路径
	FilePath string
	// 包路径, 例如 git.oschina.net/jscode/list-interface
	PackagePath string
}

type interfaceMeta struct {
	baseInfo
	Name string
	// interface的方法签名列表,
	MethodSigns []string
}

type structMeta struct {
	baseInfo
	Name string
	// struct的方法签名列表
	MethodSigns []string
}

type importMeta struct {
	// 例如 main
	Alias string
	// 例如 git.oschina.net/jscode/list-interface
	Path string
}

type analysisTool struct {
	config Config

	// 当前解析的go文件, 例如/appdev/go-demo/src/git.oschina.net/jscode/list-interface/a.go
	currentFile string
	// 当前解析的go文件,所在包路径, 例如git.oschina.net/jscode/list-interface
	currentPackagePath string
	// 当前解析的go文件,引入的其他包
	currentFileImports []*importMeta

	// 所有的interface
	interfaceMetas []*interfaceMeta
	// 所有的struct
	structMetas []*structMeta
	// package path与package name的映射关系,例如git.oschina.net/jscode/list-interface 对应的pakcage name为 main
	packagePathPackageNameCache map[string]string
}

func (this *analysisTool) analysis(config Config) {

	this.config = config

	if this.config.CodeDir == "" || !PathExists(this.config.CodeDir) {
		log.Printf("找不到代码目录%s\n", this.config.CodeDir)
		return
	}

	if this.config.GopathDir == "" || !PathExists(this.config.GopathDir) {
		log.Printf("找不到GOPATH目录%s\n", this.config.GopathDir)
		return
	}

	for _, lib := range stdlibs {
		this.mapPackagePath_PackageName(lib, path.Base(lib))
	}

	dir_walk_once := func(path string, info os.FileInfo, err error) error {
		// 过滤掉测试代码
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "test.go") {
			if config.IgnoreDirs != nil && HasPrefixInSomeElement(path, config.IgnoreDirs) {
				// ignore
			} else {
				log.Println("解析 " + path)
				this.visitTypeInFile(path)
			}
		}

		return nil
	}

	filepath.Walk(config.CodeDir, dir_walk_once)

	dir_walk_twice := func(path string, info os.FileInfo, err error) error {
		// 过滤掉测试代码
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "test.go") {
			if config.IgnoreDirs != nil && HasPrefixInSomeElement(path, config.IgnoreDirs) {
				// ignore
			} else {
				log.Println("解析 " + path)
				this.visitFuncInFile(path)
			}
		}

		return nil
	}

	filepath.Walk(config.CodeDir, dir_walk_twice)

}

func (this *analysisTool) initFile(path string) {
	log.Println("path=", path)

	this.currentFile = path
	this.currentPackagePath = this.filepathToPackagePath(path)

	if this.currentPackagePath == "" {
		log.Printf("packagePath为空,currentFile=%s\n", this.currentFile)
	}

}

func (this *analysisTool) mapPackagePath_PackageName(packagePath string, packageName string) {
	if packagePath == "" || packageName == "" {
		log.Printf("mapPackagePath_PackageName, packageName=%s, packagePath=%s\n, current_file=%s",
			packageName, packagePath, this.currentFile)
		return
	}

	if mapContains(this.packagePathPackageNameCache, packagePath) {
		return
	}

	log.Printf("mapPackagePath_PackageName, packageName=%s, packagePath=%s\n", packageName, packagePath)
	this.packagePathPackageNameCache[packagePath] = packageName

}

func (this *analysisTool) visitTypeInFile(path string) {

	this.initFile(path)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)

	if err != nil {
		log.Fatalln(err)
		return
	}

	this.mapPackagePath_PackageName(this.currentPackagePath, file.Name.Name)

	for _, decl := range file.Decls {

		genDecl, ok := decl.(*ast.GenDecl)

		if ok {
			for _, spec := range genDecl.Specs {

				typeSpec, ok := spec.(*ast.TypeSpec)

				if ok {

					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if ok {
						this.visitInterface(typeSpec.Name.Name, interfaceType)
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if ok {
						this.visitStruct(typeSpec.Name.Name, structType)
					}

				}
			}
		}

	}

}

func (this *analysisTool) filepathToPackagePath(filepath string) string {

	filepath = path.Dir(filepath)

	if this.config.VendorDir != "" {
		if strings.HasPrefix(filepath, this.config.VendorDir) {
			packagePath := strings.TrimPrefix(filepath, this.config.VendorDir)
			packagePath = strings.TrimPrefix(packagePath, "/")
			return packagePath
		}
	}

	if this.config.GopathDir != "" {
		srcdir := path.Join(this.config.GopathDir, "src")
		if strings.HasPrefix(filepath, srcdir) {
			packagePath := strings.TrimPrefix(filepath, srcdir)
			packagePath = strings.TrimPrefix(packagePath, "/")
			return packagePath
		}
	}

	log.Printf("无法确认包路径名, filepath=%s\n", filepath)

	return ""

}

func (this *analysisTool) visitFuncInFile(path string) {

	this.initFile(path)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)

	if err != nil {
		log.Fatal(err)
		return
	}

	this.currentFileImports = []*importMeta{}

	if file.Imports != nil {
		for _, import1 := range file.Imports {

			alias := ""
			packagePath := strings.TrimSuffix(strings.TrimPrefix(import1.Path.Value, "\""), "\"")

			if import1.Name != nil {
				alias = import1.Name.Name
			} else {
				aliasCache, ok := this.packagePathPackageNameCache[packagePath]
				log.Printf("findAliasInCache,packagePath=%s,alias=%s,ok=%t\n", packagePath, aliasCache, ok)
				if ok {
					alias = aliasCache
				} else {
					alias = this.findAliasByPackagePath(packagePath)
				}
			}

			log.Printf("current_file=%s packagePath=%s, alias=%s\n", this.currentFile, packagePath, alias)

			this.currentFileImports = append(this.currentFileImports, &importMeta{
				Alias: alias,
				Path:  packagePath,
			})
		}
	}

	for _, decl := range file.Decls {

		genDecl, ok := decl.(*ast.GenDecl)

		if ok {
			for _, spec := range genDecl.Specs {

				typeSpec, ok := spec.(*ast.TypeSpec)

				if ok {

					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if ok {
						this.visitInterfaceFunctions(typeSpec.Name.Name, interfaceType)
					}

				}
			}
		}

	}

	for _, decl := range file.Decls {

		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok {
			this.visitFunc(funcDecl)
		}

	}

}

func (this *analysisTool) visitStruct(name string, structType *ast.StructType) {

	strutMeta1 := &structMeta{
		baseInfo: baseInfo{
			FilePath:    this.currentFile,
			PackagePath: this.currentPackagePath,
		},
		Name:        name,
		MethodSigns: []string{},
	}

	this.structMetas = append(this.structMetas, strutMeta1)

}

func (this *analysisTool) structBodyToString(structType *ast.StructType) string {

	result := "{\n"

	for _, field := range structType.Fields.List {
		result += "  " + this.fieldToString(field) + "\n"
	}

	result += "}"

	return result

}

func (this *analysisTool) visitInterface(name string, interfaceType *ast.InterfaceType) {

	interfaceInfo1 := &interfaceMeta{
		baseInfo: baseInfo{
			FilePath:    this.currentFile,
			PackagePath: this.currentPackagePath,
		},
		Name: name,
	}

	this.interfaceMetas = append(this.interfaceMetas, interfaceInfo1)

}

func (this *analysisTool) funcParamsResultsToString(funcType *ast.FuncType) string {

	funcString := "("

	if funcType.Params != nil {
		for index, field := range funcType.Params.List {
			if index != 0 {
				funcString += ","
			}

			funcString += this.fieldToString(field)
		}
	}

	funcString += ")"

	if funcType.Results != nil {

		if len(funcType.Results.List) >= 2 {
			funcString += "("
		}

		for index, field := range funcType.Results.List {
			if index != 0 {
				funcString += ","
			}

			funcString += this.fieldToString(field)
		}

		if len(funcType.Results.List) >= 2 {
			funcString += ")"
		}
	}

	return funcString

}

func (this *analysisTool) findStruct(packagePath string, structName string) *structMeta {

	for _, structMeta1 := range this.structMetas {
		if structMeta1.Name == structName && structMeta1.PackagePath == packagePath {
			return structMeta1
		}
	}

	return nil
}

func (this *analysisTool) findInterfaceMeta(packagePath string, interfaceName string) *interfaceMeta {

	for _, interfaceMeta := range this.interfaceMetas {
		if interfaceMeta.Name == interfaceName && interfaceMeta.PackagePath == packagePath {
			return interfaceMeta
		}
	}

	return nil
}

func (this *analysisTool) visitFunc(funcDecl *ast.FuncDecl) {

	this.debugFunc(funcDecl)

	packageAlias, structName := this.findStructTypeOfFunc(funcDecl)

	if structName != "" {

		packagePath := ""
		if packageAlias == "" {
			packagePath = this.currentPackagePath
		}

		structMeta := this.findStruct(packagePath, structName)
		if structMeta != nil {
			methodSign := this.createMethodSign(funcDecl.Name.Name, funcDecl.Type)
			structMeta.MethodSigns = append(structMeta.MethodSigns, methodSign)
		}
	}

}

func (this *analysisTool) visitInterfaceFunctions(name string, interfaceType *ast.InterfaceType) {

	methods := []string{}

	for _, field := range interfaceType.Methods.List {

		funcType, ok := field.Type.(*ast.FuncType)

		if ok {
			methods = append(methods, this.createMethodSign(field.Names[0].Name, funcType))
		}
	}

	interfaceMeta := this.findInterfaceMeta(this.currentPackagePath, name)
	interfaceMeta.MethodSigns = methods

}

func (this *analysisTool) findStructTypeOfFunc(funcDecl *ast.FuncDecl) (packageAlias string, structName string) {

	if funcDecl.Recv != nil {

		for _, field := range funcDecl.Recv.List {

			t := field.Type

			ident, ok := t.(*ast.Ident)
			if ok {
				packageAlias = ""
				structName = ident.Name
			}

			starExpr, ok := t.(*ast.StarExpr)
			if ok {
				ident, ok := starExpr.X.(*ast.Ident)
				if ok {
					packageAlias = ""
					structName = ident.Name
				}

			}
		}
	}

	return
}

func (this *analysisTool) debugFunc(funcDecl *ast.FuncDecl) {

	log.Println("func name=", funcDecl.Name)

	if funcDecl.Recv != nil {
		for _, field := range funcDecl.Recv.List {
			log.Println("func recv, name=", field.Names, " type=", field.Type)
		}
	}

	if funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			log.Println("func param, name=", field.Names, " type=", field.Type)
		}
	}

	if funcDecl.Type.Results != nil {
		for _, field := range funcDecl.Type.Results.List {
			log.Println("func result, type=", field.Type)
		}
	}

}

func (this *analysisTool) IdentsToString(names []*ast.Ident) string {
	r := ""
	for index, name := range names {
		if index != 0 {
			r += ","
		}
		r += name.Name
	}

	return r
}

// 创建方法签名
func (this *analysisTool) createMethodSign(methodName string, funcType *ast.FuncType) string {

	methodSign := methodName + "("

	if funcType.Params != nil {
		for index, field := range funcType.Params.List {
			if index != 0 {
				methodSign += ","
			}
			methodSign += this.fieldToStringInMethodSign(field)
		}
	}

	methodSign += ")"

	if funcType.Results != nil {

		if len(funcType.Results.List) >= 2 {
			methodSign += "("
		}

		for index, field := range funcType.Results.List {
			if index != 0 {
				methodSign += ","
			}
			methodSign += this.fieldToStringInMethodSign(field)
		}

		if len(funcType.Results.List) >= 2 {
			methodSign += ")"
		}
	}

	return methodSign
}

func (this *analysisTool) fieldToStringInMethodSign(f *ast.Field) string {

	argCount := len(f.Names)

	if argCount == 0 {
		argCount = 1
	}

	sign := ""

	for i := 0; i < argCount; i++ {
		if i != 0 {
			sign += ","
		}
		sign += this.typeToString(f.Type)
	}

	return sign
}

func (this *analysisTool) fieldToString(f *ast.Field) string {

	r := ""

	if len(f.Names) > 0 {

		for index, name := range f.Names {
			if index != 0 {
				r += ","
			}

			r += name.Name
		}

		r += " "

	}

	r += this.typeToString(f.Type)

	return r

}

func (this *analysisTool) typeToString(t ast.Expr) string {

	ident, ok := t.(*ast.Ident)
	if ok {
		return this.addPackagePathWhenStruct(ident.Name)
	}

	starExpr, ok := t.(*ast.StarExpr)
	if ok {
		return "*" + this.typeToString(starExpr.X)
	}

	arrayType, ok := t.(*ast.ArrayType)
	if ok {
		return "[]" + this.typeToString(arrayType.Elt)
	}

	mapType, ok := t.(*ast.MapType)
	if ok {
		return "map[" + this.typeToString(mapType.Key) + "]" + this.typeToString(mapType.Value)
	}

	chanType, ok := t.(*ast.ChanType)
	if ok {
		return "chan " + this.typeToString(chanType.Value)
	}

	funcType, ok := t.(*ast.FuncType)
	if ok {
		return "func" + this.funcParamsResultsToString(funcType)
	}

	interfaceType, ok := t.(*ast.InterfaceType)
	if ok {
		return "interface " + strings.Replace(this.interfaceBodyToString(interfaceType), "\n", " ", -1)
	}

	selectorExpr, ok := t.(*ast.SelectorExpr)
	if ok {
		return this.findPackagePathByAlias(this.selectorExprToString(selectorExpr.X), selectorExpr.Sel.Name) + "." + selectorExpr.Sel.Name
	}

	structType, ok := t.(*ast.StructType)
	if ok {
		return "struct " + strings.Replace(this.structBodyToString(structType), "\n", " ", -1)
	}

	ellipsis, ok := t.(*ast.Ellipsis)
	if ok {
		return "... " + this.typeToString(ellipsis.Elt)
	}

	parenExpr, ok := t.(*ast.ParenExpr)
	if ok {
		return " (" + this.typeToString(parenExpr.X) + ")"
	}

	log.Println("typeToString ", reflect.TypeOf(t), " file=", this.currentFile, " expr=", this.content(t))

	return ""
}

func (this *analysisTool) selectorExprToString(t ast.Expr) string {

	ident, ok := t.(*ast.Ident)
	if ok {
		return ident.Name
	}

	log.Println("selectorExprToString ", reflect.TypeOf(t), " file=", this.currentFile, " expr=", this.content(t))

	return ""
}

func (this *analysisTool) addPackagePathWhenStruct(fieldType string) string {

	searchPackages := []string{this.currentPackagePath}

	for _, import1 := range this.currentFileImports {
		if import1.Alias == "." {
			searchPackages = append(searchPackages, import1.Path)
		}
	}

	for _, meta := range this.structMetas {
		if sliceContains(searchPackages, meta.PackagePath) && meta.Name == fieldType {
			return meta.PackagePath + "." + fieldType
		}
	}

	for _, meta := range this.interfaceMetas {
		if sliceContains(searchPackages, meta.PackagePath) && meta.Name == fieldType {
			return meta.PackagePath + "." + fieldType
		}
	}

	return fieldType
}

func (this *analysisTool) findAliasByPackagePath(packagePath string) string {
	result := ""

	if this.config.VendorDir != "" {
		absPath := path.Join(this.config.VendorDir, packagePath)
		if PathExists(absPath) {
			result = findGoPackageNameInDirPath(absPath)
		}
	}

	if this.config.GopathDir != "" {
		absPath := path.Join(this.config.GopathDir, "src", packagePath)
		if PathExists(absPath) {
			result = findGoPackageNameInDirPath(absPath)
		}
	}

	log.Println("packagepath=%s, alias=%s\n", packagePath, result)

	return result
}

func (this *analysisTool) findPackagePathByAlias(alias string, structName string) string {

	for _, importMeta := range this.currentFileImports {
		if importMeta.Path == alias {
			log.Printf("findPackagePathByAlias, alias=%s, packagePath=%s\n", alias, alias)
			return alias
		}
	}

	matchedImportMetas := []*importMeta{}

	for _, importMeta := range this.currentFileImports {
		if importMeta.Alias == alias {
			matchedImportMetas = append(matchedImportMetas, importMeta)
		}
	}

	if len(matchedImportMetas) == 1 {
		log.Println("findPackagePathByAlias, alias=%s, packagePath=%s\n", alias, matchedImportMetas[0].Path)
		return matchedImportMetas[0].Path
	}

	if len(matchedImportMetas) > 1 {

		for _, matchedImportMeta := range matchedImportMetas {

			for _, structMeta := range this.structMetas {
				if structMeta.Name == structName && structMeta.PackagePath == matchedImportMeta.Path {
					log.Println("findPackagePathByAlias, alias=%s, packagePath=%s\n", alias, matchedImportMeta.Path)
					return matchedImportMeta.Path
				}
			}
		}

	}

	log.Printf("找不到包的全路径，包名为%s，在%s文件, matchedImportMetas=%d", alias, this.currentFile, len(matchedImportMetas))

	return alias
}

func (this *analysisTool) interfaceBodyToString(interfaceType *ast.InterfaceType) string {

	result := " {\n"

	for _, field := range interfaceType.Methods.List {

		funcType, ok := field.Type.(*ast.FuncType)

		if ok {
			result += "  " + this.IdentsToString(field.Names) + this.funcParamsResultsToString(funcType) + "\n"
		}

	}

	result += "}"

	return result

}

func (this *analysisTool) content(t ast.Expr) string {
	bytes, err := ioutil.ReadFile(this.currentFile)
	if err != nil {
		log.Println("读取文件", this.currentFile, "失败", err)
		return ""
	}

	return string(bytes[t.Pos()-1 : t.End()-1])
}

/**
 * 查找interface有哪些实现的Struct
 */
func (this *analysisTool) findInterfaceImpls(interfaceMeta1 *interfaceMeta) []*structMeta {
	metas := []*structMeta{}

	for _, structMeta1 := range this.structMetas {
		if sliceContainsSlice(structMeta1.MethodSigns, interfaceMeta1.MethodSigns) {
			metas = append(metas, structMeta1)
		}
	}

	return metas
}

func (this *analysisTool) OutputToFile(logfile string) {

	file, err := os.Create(logfile)
	if err != nil {
		log.Printf("打开文件%s失败\n", logfile)
	}

	for _, interfaceMeta := range this.interfaceMetas {

		string := fmt.Sprintf("interface %s 在文件%s中\n", interfaceMeta.Name, interfaceMeta.FilePath)
		log.Println(string)
		file.WriteString(string)

		structMetas := this.findInterfaceImpls(interfaceMeta)

		if len(structMetas) == 0 {

			string = "没有找到实现该接口的struct\n"
			log.Println(string)
			file.WriteString(string)

		} else {

			string = fmt.Sprintf("有%d个struct实现了接口\n", len(structMetas))
			log.Println(string)
			file.WriteString(string)

			for _, structMeta := range structMetas {
				string = fmt.Sprintf("struct %s 在文件%s中\n", structMeta.Name, structMeta.FilePath)
				log.Println(string)
				file.WriteString(string)
			}

		}

		log.Printf("\n")
		file.WriteString("\n")
	}

	file.Close()

	log.Printf("解析结果已保存到%s\n", logfile)

}
