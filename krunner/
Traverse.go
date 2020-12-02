package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yanyiwu/gojieba"
)

var x *gojieba.Jieba

var exlist = []string{"cmake_install.cmake", "litmain.sh", "*.qrc", ".git", ".vscode", ".config", ".vim", ".local", ".oh-my-zsh", ".oh-my-bash", ".oh-my-fish", ".uic", "po", "*.map", "*.pc", "*.gbff", "qrc_*.cpp", ".npm", ".npmrc", ".yarnrc", ".m2", ".mc", ".mozilla", ".android", ".aria2", ".cache", ".dbus", ".docker", ".gnome", ".kde", ".ssh", ".kingsoft", ".steam", ".gitbook", ".dbvis", ".directory", ".ipython", ".ptpython", ".pip", ".nvm", ".pki", ".yarn", ".node-gyp", "*.tmp", ".vsce", "*.temp", ".z", ".zcache", ".viminfo", ".vimrc", ".bashrc", ".zshrc", ".Xauthority", ".albertignore", ".anydesk", "*.jsc", "*.vhd", ".obj", "node_packages", "*.orig", "vbox.log", ".helmignore", "*.loT", "*.rcore", ".helm", ".rancher", ".minikube", ".minio", ".ipfs", ".lib64", ".java", ".jupyter", ".hg", "*.vhdx", "*.m4", "*.vdi", ".pch", "*.fastq", "*.lo", "lost+found", "*.sql.gz", "*.nvram", "*.init", "CMakeTmpQmake", "ui_*.h", "*.img", ".xsession-errors*", "lzo", "*.a", "*.omf", ".svn", "CMakeTmp", "confdefs.h", "CVS", "*.fq", "*.gb", "*.aux", "moc_*.cpp", ".yarn-cache", "*.elc", "*.csproj", "*.faa", ".bzr", "*.vm*", "*.swap", "*.o", "*.db", "*.gmo", "*.qcow2", "conftest", "confstat", "config.status", "node_modules", "*.fna", "*.la", "__pycache__", "*.po,.histfile.*", "*.rej", ".yarn", "libtool", "*~", "*.moc", "*.pyo", "*.so", "CTestTestfile.cmake", ".moc", "_darcs", "*.fasta", "*.gcode", "autom4te", "*.class", "*.qmlc", "*.vbox*", "*.pyc", "core-dumps", "*.vmdk", "Makefile.am", "nbproject", "CMakeCache.txt", "CMakeFiles", "*.ini", "*.part", "*.sql"}

func isMatch(name string) bool {
	for _, ex := range exlist {
		if ok, _ := filepath.Match(ex, name); ok {
			return true
		}
	}
	return false
}

func jieba(isdir bool, name, path string) {
	words := x.Cut(name, true)
	fmt.Println(path, ":", words)
}

func getFilelist(path string, expaths []string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		for _, expath := range expaths {
			if path == expath {
				fmt.Println("EX", expath)
				return filepath.SkipDir
			}
		}

		pname := filepath.Base(path)

		if isMatch(pname) {
			if f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 或者不作分词, 直接以bytes为key做map, 搜索时使用bytes.Contains
		jieba(f.IsDir(), pname, path)

		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func main() {
	flag.Parse()
	x = gojieba.NewJieba()
	defer x.Free()
	root := flag.Arg(0)
	expaths := flag.Args()[1:]
	getFilelist(root, expaths)
}
