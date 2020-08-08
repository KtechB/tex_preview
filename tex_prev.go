package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	beginfile string
	files     []string
)

const dirname = "preview_dir"

// プログラム引数を取得
func init() {
	flag.StringVar(&beginfile, "beginfile", "begin.tex", "filename of header")
	flag.Parse()
	files = flag.Args()
}

func logfatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println(files)
	fmt.Println(beginfile)

	if len(files) == 0 {
		println("give filename as args")
		os.Exit(0)
	}
	path, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Dir(path) // 実行ファイルのディレクトリを取得

	bytes, err := ioutil.ReadFile(filepath.Join(path, beginfile))
	if err != nil {
		panic(err)
	}
	output := string(bytes)
	outFilename := ""
	for _, file := range files {
		// 行番号、ファイル名を表示する
		outFilename += file

		bytes, err := ioutil.ReadFile(filepath.Join(path, file))
		if err != nil {
			panic(err)
		}
		source := string(bytes)
		output += "\n" + source
	}

	output += "\n" + "\\end{document}"

	// ファイルを書き込み用にオープン (mode=0666)
	outFilename = "./PREVIEW.tex"
	file, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// テキストを書き込む
	_, err = file.WriteString(output)
	if err != nil {
		panic(err)
	}
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}

	texCommand := []string{"-l", "-ot", "-kanji=utf8", "-synctex=1", outFilename}

	outlog, err := exec.Command("ptex2pdf", texCommand...).Output()
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(outlog)
	logfatal(err)

}
