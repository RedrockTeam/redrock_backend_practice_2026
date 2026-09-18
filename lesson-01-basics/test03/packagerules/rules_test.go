package packagerules

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const packageProjectModule = "example.com/hello"

var normalizedPackageName = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

type fixtureSource struct {
	file        string
	dir         string
	importPath  string
	packageName string
	imports     []string
}

// package 规范题位于 testdata/package-project。
// testdata 不会被 go test 当成普通包编译，因此可以安全保存故意写错的包声明。
func TestPackageNamesAndDirectories(t *testing.T) {
	sources := loadFixtureSources(t, "package-project", packageProjectModule)
	namesByDir := make(map[string]string)

	for _, source := range sources {
		if previous, ok := namesByDir[source.dir]; ok && previous != source.packageName {
			t.Errorf(
				"%s 同时声明了 package %s 和 package %s；同一目录只能有一个包",
				displayFixturePath(source.file),
				previous,
				source.packageName,
			)
		} else {
			namesByDir[source.dir] = source.packageName
		}

		if source.importPath == packageProjectModule {
			if source.packageName != "main" {
				t.Errorf("%s 位于可执行程序根目录，应声明 package main", displayFixturePath(source.file))
			}
			continue
		}

		dirName := path.Base(source.importPath)
		if !normalizedPackageName.MatchString(dirName) {
			t.Errorf("目录 %q 应改成小写、简短且不含下划线的名称", dirName)
		}
		if !normalizedPackageName.MatchString(source.packageName) {
			t.Errorf("%s 中的 package %q 不符合小写包名规范", displayFixturePath(source.file), source.packageName)
		}
		if source.packageName != dirName {
			t.Errorf(
				"%s 声明 package %q，但目录名是 %q；本题要求二者保持一致",
				displayFixturePath(source.file),
				source.packageName,
				dirName,
			)
		}
		if source.packageName == "common" || source.packageName == "util" || source.packageName == "utils" {
			t.Errorf("%s 的包名 %q 过于宽泛，请根据实际职责重新命名", displayFixturePath(source.file), source.packageName)
		}
	}
}

func TestPackageImportBoundaries(t *testing.T) {
	projectSources := loadFixtureSources(t, "package-project", packageProjectModule)
	externalSources := loadFixtureSources(t, "external-module", "example.net/outside")
	projectPackages := make(map[string]bool)

	for _, source := range projectSources {
		projectPackages[source.importPath] = true
	}

	checkImports := func(source fixtureSource) {
		for _, imported := range source.imports {
			if imported != packageProjectModule && !strings.HasPrefix(imported, packageProjectModule+"/") {
				continue
			}

			if !projectPackages[imported] {
				t.Errorf(
					"%s 导入了不存在的本地路径 %q；重命名目录后也要更新 import",
					displayFixturePath(source.file),
					imported,
				)
				continue
			}

			if imported == packageProjectModule && source.importPath != packageProjectModule {
				t.Errorf("%s 尝试导入 %q，但 package main 不能被普通功能包导入", displayFixturePath(source.file), imported)
			}

			if !internalImportAllowed(source.importPath, imported) {
				t.Errorf(
					"%s 位于 module %q，不能越过边界导入 %q",
					displayFixturePath(source.file),
					moduleOf(source.importPath),
					imported,
				)
			}
		}
	}

	for _, source := range projectSources {
		checkImports(source)
	}
	for _, source := range externalSources {
		checkImports(source)
	}
}

func TestPackagesDoNotImportEachOtherInACycle(t *testing.T) {
	sources := loadFixtureSources(t, "package-project", packageProjectModule)
	graph := make(map[string]map[string]bool)

	for _, source := range sources {
		if graph[source.importPath] == nil {
			graph[source.importPath] = make(map[string]bool)
		}
	}
	for _, source := range sources {
		for _, imported := range source.imports {
			if graph[imported] != nil {
				graph[source.importPath][imported] = true
			}
		}
	}

	if cycle := findImportCycle(graph); len(cycle) > 0 {
		t.Fatalf("发现循环导入：%s；请重新设计单向依赖关系", strings.Join(cycle, " -> "))
	}
}

func loadFixtureSources(t *testing.T, fixtureName, modulePath string) []fixtureSource {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位 package 规范题")
	}
	root := filepath.Join(filepath.Dir(currentFile), "testdata", fixtureName)
	var sources []fixtureSource

	err := filepath.WalkDir(root, func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(filename) != ".go" {
			return nil
		}

		parsed, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		relDir, err := filepath.Rel(root, filepath.Dir(filename))
		if err != nil {
			return err
		}

		importPath := modulePath
		if relDir != "." {
			importPath += "/" + filepath.ToSlash(relDir)
		}

		source := fixtureSource{
			file:        filename,
			dir:         filepath.Dir(filename),
			importPath:  importPath,
			packageName: parsed.Name.Name,
		}
		for _, spec := range parsed.Imports {
			imported, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			source.imports = append(source.imports, imported)
		}
		sources = append(sources, source)
		return nil
	})
	if err != nil {
		t.Fatalf("读取 package 规范题失败：%v", err)
	}
	if len(sources) == 0 {
		t.Fatal("package 规范题中没有找到 Go 文件")
	}
	return sources
}

func internalImportAllowed(importer, imported string) bool {
	const marker = "/internal/"
	index := strings.Index(imported, marker)
	if index < 0 {
		return true
	}

	allowedRoot := imported[:index]
	return importer == allowedRoot || strings.HasPrefix(importer, allowedRoot+"/")
}

func moduleOf(importPath string) string {
	parts := strings.Split(importPath, "/")
	if len(parts) < 3 {
		return importPath
	}
	return strings.Join(parts[:3], "/")
}

func findImportCycle(graph map[string]map[string]bool) []string {
	const (
		unvisited = iota
		visiting
		visited
	)

	state := make(map[string]int)
	var stack []string
	var cycle []string

	var visit func(string) bool
	visit = func(current string) bool {
		state[current] = visiting
		stack = append(stack, current)

		var nextPackages []string
		for next := range graph[current] {
			nextPackages = append(nextPackages, next)
		}
		sort.Strings(nextPackages)

		for _, next := range nextPackages {
			switch state[next] {
			case unvisited:
				if visit(next) {
					return true
				}
			case visiting:
				start := 0
				for stack[start] != next {
					start++
				}
				cycle = append(cycle, stack[start:]...)
				cycle = append(cycle, next)
				return true
			}
		}

		stack = stack[:len(stack)-1]
		state[current] = visited
		return false
	}

	var packages []string
	for packagePath := range graph {
		packages = append(packages, packagePath)
	}
	sort.Strings(packages)

	for _, packagePath := range packages {
		if state[packagePath] == unvisited && visit(packagePath) {
			return cycle
		}
	}
	return nil
}

func displayFixturePath(filename string) string {
	index := strings.Index(filepath.ToSlash(filename), "/testdata/")
	if index < 0 {
		return filename
	}
	return filepath.ToSlash(filename)[index+1:]
}
