package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/svg"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type DocInfo struct {
	DocHash string   `json:"dochash"` // 文档哈希
	PageNum int      `json:"pagenum"` // 文档页数(转换为pdf后)
	PageImg []string `json:"pageimg"` // 文档页数据
}

func openOrCreate(filename string) *os.File {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return nil
		}
		return file
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil
	}

	return file
}

func ReplaceStr(str string) string {
	replacer := strings.NewReplacer("\\", "", " ", "", "\n", "", "\r", "", "\t", "", "/", "", ":", "", "*", "", "?", "", "|", "", "<", "", ">", "", "@", "")
	return replacer.Replace(str)
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

func GetMD5(filepath string) string {
	var filechunk uint64 = 8192 // we settle for 8KB

	file, err := os.Open(filepath)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(float64(filechunk), float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// 将 PDF 转成 SVG
func PDF2SVG(pdfFile, svgFile string, pageNO int) (err error) {

	//Usage: pdf2svg <in file.pdf> <out file.svg> [<page no>]
	pdfFile, err = filepath.Abs(pdfFile)
	if err != nil {
		return
	}

	svgFile, err = filepath.Abs(svgFile)
	if err != nil {
		return
	}

	args := []string{pdfFile, svgFile, strconv.Itoa(pageNO)}
	cmd := exec.Command("pdf2svg", args...)
	/*
			if strings.HasPrefix("pdf2svg", "sudo") {
				args = append([]string{strings.TrimPrefix("pdf2svg", "sudo")}, args...)
				cmd = exec.Command("sudo", args...)
		    }
	*/
	time.AfterFunc(20*time.Second, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	return
}

// 使用GhostScript将PDF页转换为PNG图片(单页)
func PDF2PNG(pdfname, pngname string, start, end, imgSize int) (err error) {
	// gs -sDEVICE=pngalpha -r256 -dNOPAUSE -dBATCH -dFirstPage=5 -dLastPage=5 -sOutputFile=output.png  input.pdf
	// gs -sDEVICE=pngalpha -r512 -dNOPAUSE -dBATCH -dFirstPage=1 -dLastPage=10 -sOutputFile=output-%d.png 算法图解.pdf
	// r似乎表示尺寸上限, 如r144表示不超过144KiB
	pdfname, err = filepath.Abs(pdfname)
	if err != nil {
		return
	}

	// -sstdout=/dev/null
	args := []string{"-sDEVICE=pngalpha", "-r" + strconv.Itoa(imgSize), "-dNOPAUSE", "-dBATCH", "-dFirstPage=" + strconv.Itoa(start), "-dLastPage=" + strconv.Itoa(end), "-sOutputFile=" + pngname, pdfname}
	cmd := exec.Command("gs", args...)
	/*
			if strings.HasPrefix("pdf2svg", "sudo") {
				args = append([]string{strings.TrimPrefix("pdf2svg", "sudo")}, args...)
				cmd = exec.Command("sudo", args...)
		    }
	*/
	time.AfterFunc(10*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()

	if err != nil {
		return
	}

	if start != end {
		seq := strconv.Itoa(start - 1)
		prefix := strings.TrimSuffix(pngname, "-"+seq+"+%d.png")
		for i := 1; i <= end-start+1; i++ {
			oldpngfile := prefix + "-" + seq + "+" + strconv.Itoa(i) + ".png"
			newpngfile := prefix + "-" + strconv.Itoa(start-1+i) + ".png"
			os.Rename(oldpngfile, newpngfile)
		}
	}
	return
}

//office文档转pdf，返回转化后的文档路径和错误
func OfficeToPDF(input string) (err error) {
	input, err = filepath.Abs(input)
	if err != nil {
		return
	}

	args := []string{"--headless", "--invisible", "--convert-to", "pdf", input, "--outdir", filepath.Dir(input)}
	cmd := exec.Command("soffice", args...)

	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	return
}

// pandoc 需要texlive转换pdf:
// sudo dnf install texlive texlive-xetex texlive-collection-langchinese
// pandoc --pdf-engine=xelatex --highlight-style zenburn -V geometry:"top=2cm, bottom=1.5cm, left=2cm, right=2cm" -V CJKoptions=BoldFont="SimHei" -V CJKmainfont="SimSun" -V CJKmonofont="SimSun" -s s.mobi -o s.pdf
// pandoc --pdf-engine=xelatex  -s -V mainfont="SimSun" -V monofont="SimSun"  c.epub -o c.pdf
// -N 给section加编号; --toc 加目录

// gs -sDEVICE=pdfwrite -dNOPAUSE -dBATCH -dFirstPage=5 -dLastPage=5 -sOutputFile=output.pdf input.pdf

//非office文档转pdf，返回转化后的文档路径和错误
func ConvertByPandoc(file string) (err error) {
	file, err = filepath.Abs(file)
	if err != nil {
		return
	}

	target := strings.TrimSuffix(file, path.Ext(file)) + ".pdf"
	args := []string{"--pdf-engine=xelatex", "-V geometry:'top=2cm, bottom=1.5cm, left=2cm, right=2cm'", "-V CJKoptions=BoldFont='SimHei'", "-V CJKmainfont='SimSun'", "-V CJKmonofont='SimSun'", "-o", target, file}
	cmd := exec.Command("pandoc", args...)

	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})

	err = cmd.Run()
	return
}

//非office文档(.txt,.mobi,.epub)转pdf文档
// calibre支持格式:EPUB、MOBI、AZW3、DOCX、FB2、HTMLZ、LIT、LRF、PDB、PDF、PMIZ、RB、RTF、SNB、TCR、TXT、TXTZ、ZIP
// 转换效果好
func FileToPDF(input string) (output string, err error) {
	//calibre := beego.AppConfig.DefaultString("calibre", "ebook-convert")
	/*
			if len(ext) > 0 {
				e = ext[0]
		    }
	*/
	input, err = filepath.Abs(input)
	if err != nil {
		return
	}

	output = filepath.Dir(input) + "/" + strings.TrimSuffix(filepath.Base(input), filepath.Ext(input)) + ".pdf"
	output, _ = filepath.Abs(output)
	args := []string{
		input,
		output,
	}

	/*
		args = append(args,
			"--paper-size", "a4",
			"--pdf-default-font-size", "16",
			"--pdf-page-margin-bottom", "36",
			"--pdf-page-margin-left", "36",
			"--pdf-page-margin-right", "36",
			"--pdf-page-margin-top", "36",
		)
	*/

	cmd := exec.Command("ebook-convert", args...)
	/*
			if strings.HasPrefix(calibre, "sudo") {
				calibre = strings.TrimPrefix(calibre, "sudo")
				args = append([]string{calibre}, args...)
				cmd = exec.Command("sudo", args...)
		    }
	*/
	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	return
}

// 使用pdfinfo获取pdf页数
func PageNumber(input string) int {
	ext := strings.ToLower(path.Ext(input))
	if ext != ".pdf" {
		return 0
	}
	input, err := filepath.Abs(input)
	if err != nil {
		return 0
	}

	args := "pdfinfo " + input + " | grep Pages | cut -d ':' -f2"
	cmd := exec.Command("bash", "-c", args)

	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})

	stdout, err := cmd.Output()
	if err != nil {
		return 0
	}

	res := strings.Trim(strings.Replace(string(stdout), "\n", "", -1), " ")
	pagenum, err := strconv.Atoi(res)
	if err != nil {
		return 0
	}

	return pagenum

	/*
		var out bytes.Buffer
		cmd.Stderr = &out

		cmd.Run()

		res := strings.Trim(out.String(), " ")
		pagenum, err := strconv.ParseInt(res, 10, 64)
		if err != nil {
			pagenum = 0
		}
	*/
}

func ImageToPNG(input string) (output string, err error) {
	input, err = filepath.Abs(input)
	if err != nil {
		return
	}
	ext := strings.ToLower(filepath.Ext(input))
	if ext == ".png" {
		output = input
		return
	}

	output = strings.TrimSuffix(input, ext) + ".png"

	// open files
	rf, err := os.Open(input)
	defer rf.Close()

	wf := openOrCreate(output)
	if wf == nil {
		err = errors.New("error create or open output file")
		return
	}
	defer wf.Close()

	// decode
	imageData, _, err := image.Decode(rf)
	if err != nil {
		return
	}

	// encode in new type
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(wf, imageData, nil)
	case ".webp":
		err = webp.Encode(wf, imageData, nil)
	case ".gif":
		err = gif.Encode(wf, imageData, nil)
	case ".bmp":
		err = bmp.Encode(wf, imageData)
	case ".tiff":
		err = tiff.Encode(wf, imageData, nil)
	}

	return
}

//将PDF、SVG文件转成jpg图片格式。注意：如果pdf只有一页，则文件后缀不会出现"-0.jpg"这种情况，否则会出现"-0.jpg,-1.jpg"等
func ConvertToJPEG(file string) (cover string, err error) {
	//convert := beego.AppConfig.DefaultString("imagick", "convert")
	file, err = filepath.Abs(file)
	if err != nil {
		return
	}

	cover = file + ".jpg"
	args := []string{"-density", "150", "-quality", "100", file, cover}
	/*
			if strings.HasPrefix(convert, "sudo") {
				args = append([]string{strings.TrimPrefix(convert, "sudo")}, args...)
				convert = "sudo"
		    }
	*/
	cmd := exec.Command("imagemagick convert", args...)
	time.AfterFunc(1*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	return
}

//获取PDF中指定页面的文本内容
//@param			file		PDF文件
//@param			from		起始页
//@param			to			截止页
func ExtractTextFromPDF(file string, from, to int) (content string) {
	file, _ = filepath.Abs(file)
	textfile := file + ".txt"
	args := []string{"-f", strconv.Itoa(from), "-l", strconv.Itoa(to), file, textfile}
	/*
			if strings.HasPrefix(pdftotext, "sudo") {
				args = append([]string{strings.TrimPrefix(pdftotext, "sudo")}, args...)
				pdftotext = "sudo"
		    }
	*/
	defer os.Remove(textfile)
	exec.Command("pdftotext", args...).Run()
	content = getTextFormTxtFile(textfile)
	if content == "" {
		os.Remove(textfile)
		textfile, _ = FileToPDF(file)
		content = getTextFormTxtFile(textfile)
	}
	return
}

func CompressSVG(input string, level ...int) (output string, err error) {
	// 经测试，level 值为4，压缩质量和清晰度都比较均衡
	input, err = filepath.Abs(input)
	if err != nil {
		return
	}

	output = strings.TrimSuffix(input, path.Ext(input)) + ".opti.png"
	lv := 4
	if len(level) > 0 {
		lv = level[0]
	}
	media := "image/svg+xml"
	min := &svg.Minifier{Decimals: lv}
	m := minify.New()
	m.AddFunc(media, min.Minify)
	file, err := os.Open(input)
	if err != nil {
		return
	}
	defer file.Close()

	w := &bytes.Buffer{}
	if err = m.Minify(media, w, file); err != nil {
		return
	}

	err = ioutil.WriteFile(output, w.Bytes(), os.ModePerm)
	return
}

func CompressPNG(input string) (output string, err error) {
	input, err = filepath.Abs(input)
	if err != nil {
		return
	}
	output = strings.TrimSuffix(input, path.Ext(input)) + ".opti.png"
	args := []string{"-e", ".opti.png", input}
	cmd := exec.Command("pngcrush", args...)

	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	return
}

func getTextFormTxtFile(textfile string) (content string) {
	b, err := ioutil.ReadFile(textfile)
	if err != nil {
		return
	}

	content = string(b)
	if len(content) > 0 {
		content = strings.Replace(content, "\t", " ", -1)
		content = strings.Replace(content, "\n", " ", -1)
		content = strings.Replace(content, "\r", " ", -1)
	}
	return
}

// 多线程转换
func Convert(filestr string) (dochash string, pagenum int, pngfiles []string) {
	defer elapse(time.Now())
	stdname := ReplaceStr(filestr)
	fullname, err := filepath.Abs(stdname)
	if err != nil {
		// 文件不存在或文件名不合法
		return
	}

	// dir := filepath.Dir(fullname)
	// file := filepath.Base(fullname)
	ext := filepath.Ext(fullname)

	dochash = GetMD5(fullname)

	isimg := false
	pdfname := fullname

	switch strings.ToLower(ext) {
	case ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".wps", ".rtf", ".pps", ".ppsx", ".dps", ".odp", ".pot", ".et", ".ods":
		// office文档转pdf
		// 需要为服务器安装相应字体, 否则可能导致转换后页数与原文件不一致
		// 常用字体: simsun.ttf(宋体SimSun), simhei.ttf(黑体SimHei)
		_ = OfficeToPDF(fullname)
		pdfname = strings.TrimSuffix(fullname, ext) + ".pdf"
	case ".epub", ".umd", ".chm", ".mobi", ".md", ".txt", ".azw3", ".fb2", ".htmlz", ".lit", ".lrf", ".pdb", ".pmiz", ".rb", ".snb", ".tcr", ".txtz":
		pdfname, _ = FileToPDF(fullname)
	case ".jpg", ".jpeg", ".bmp", ".gif", ".tiff", ".webp", ".png":
		isimg = true
	}

	if PathExists(pdfname) {
		defer func() {
			// 删除临时生成的pdf文件
			if !strings.HasSuffix(pdfname, fullname) {
				os.Remove(pdfname)
			}
		}()

		if isimg {
			pagenum = 1
			pngfile, err := ImageToPNG(pdfname)
			if err == nil {
				pngfiles = append(pngfiles, pngfile)
			}
			return
		}

		pagenum = PageNumber(pdfname)

		mc := 10 // 线程数

		var wg sync.WaitGroup
		// wg.Add(pagenum/10 + 1)
		for i := 0; i < pagenum/mc+1; i++ {
			start := i*mc + 1
			end := (i + 1) * mc

			if start > pagenum {
				break
			}
			if end > pagenum {
				end = pagenum
			}

			wg.Add(1)

			for i := start; i <= end; i++ {
				pngfile := strings.TrimSuffix(pdfname, ".pdf") + "-" + strconv.Itoa(i) + ".png"
				pngfiles = append(pngfiles, pngfile)
			}

			pngname := strings.TrimSuffix(pdfname, ".pdf") + "-" + strconv.Itoa(start-1) + "+%d.png"

			go func(pdfname, pngname string, start, end, size int) {
				PDF2PNG(pdfname, pngname, start, end, size)
				wg.Done()
			}(pdfname, pngname, start, end, 256)
		}

		wg.Wait()

	}

	return

}

// 单线程转换
func Convert2(filestr string) (dochash string, pagenum int, pngfiles []string) {
	defer elapse(time.Now())
	stdname := ReplaceStr(filestr)
	fullname, err := filepath.Abs(stdname)
	if err != nil {
		// 文件不存在或文件名不合法
		return
	}

	// dir := filepath.Dir(fullname)
	// file := filepath.Base(fullname)
	ext := filepath.Ext(fullname)

	dochash = GetMD5(fullname)

	isimg := false
	pdfname := fullname

	switch strings.ToLower(ext) {
	case ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".wps", ".rtf", ".pps", ".ppsx", ".dps", ".odp", ".pot", ".et", ".ods":
		// office文档转pdf
		// 需要为服务器安装相应字体, 否则可能导致转换后页数与原文件不一致
		// 常用字体: simsun.ttf(宋体SimSun), simhei.ttf(黑体SimHei)
		_ = OfficeToPDF(fullname)
		pdfname = strings.TrimSuffix(fullname, ext) + ".pdf"
	case ".epub", ".umd", ".chm", ".mobi", ".md", ".txt", ".azw3", ".fb2", ".htmlz", ".lit", ".lrf", ".pdb", ".pmiz", ".rb", ".snb", ".tcr", ".txtz":
		pdfname, _ = FileToPDF(fullname)
	case ".jpg", ".jpeg", ".bmp", ".gif", ".tiff", ".webp", ".png":
		isimg = true
	}

	if PathExists(pdfname) {
		defer func() {
			// 删除临时生成的pdf文件
			if !strings.HasSuffix(pdfname, fullname) {
				os.Remove(pdfname)
			}
		}()

		if isimg {
			pagenum = 1
			pngfile, err := ImageToPNG(pdfname)
			if err == nil {
				pngfiles = append(pngfiles, pngfile)
			}
			return
		}

		pagenum = PageNumber(pdfname)

		for i := 1; i <= pagenum; i++ {
			start := i
			end := i

			pngfile := strings.TrimSuffix(pdfname, ".pdf") + "-" + strconv.Itoa(i) + ".png"
			pngfiles = append(pngfiles, pngfile)

			PDF2PNG(pdfname, pngfile, start, end, 256)
		}
	}

	return

}

// 计算运行时间
func elapse(start time.Time) {
	// fmt.Println(time.Since(start))
}

func main() {
	/*
		app := "pdfinfo"
		//app := "buah"

		arg0 := "-e"
		arg1 := "Hello world"
		arg2 := "\n\tfrom"
		arg3 := "golang"
	*/
	//filestr := "alg.docx"
	//filestr := "arch.ppt"
	//filestr := "alg.pdf"

	args := os.Args
	if len(args) < 2 {
		return
	}

	dochash, pagenum, pageimg := Convert(args[1])
	docinfo := &DocInfo{
		DocHash: dochash,
		PageNum: pagenum,
		PageImg: pageimg,
	}

	docBytes, _ := json.Marshal(docinfo)

	println(string(docBytes))

	/* 直接保存文件页的base64编码字符串
	// 优点: 磁盘无文件残留, 调用者无需再进行文件读取和编码相关操作, 直接可以将数据入库
	// 缺点: 内存消耗过大(一页图片暂定256KB上限, 10000页将消耗2.56GB内存)
	var filestrs []string
	for _, pngfile := range pngfiles {
		if PathExists(pngfile) {
			if fileBytes, err := ioutil.ReadFile(pngfile); err == nil {
				filestr := base64.StdEncoding.EncodeToString(fileBytes)
				filestrs = append(filestrs, filestr)
			}
		}
	}
	*/

	return

	//app := "pdfinfo"
	cmdargs := "gs -sDEVICE=pngalpha -r512 -dNOPAUSE -dBATCH -dFirstPage=1 -dLastPage=10 -sOutputFile=output-%d.png ç®—æ³•å›¾è§£.pdf"

	//cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	cmd := exec.Command("bash", "-c", cmdargs)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		println(line)
	}
	/*
		stdout, err := cmd.Output()

		time.AfterFunc(5*time.Minute, func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		})

		if err != nil {
			println("error:", err.Error())
			println(string(stdout))
			return
		}
		s := strings.Replace(string(stdout), "\n", "", -1)
		println(s)
		res := strings.Trim(s, " ")
		pagenum, err := strconv.ParseInt(res, 10, 64)
		print(pagenum, err)
	*/
}
