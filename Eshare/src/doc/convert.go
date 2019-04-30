package doc

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/svg"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

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
	log.Println("PDF 转 SVG :", cmd.Args)
	err = cmd.Run()
	return
}

// 使用GhostScript将PDF页转换为PNG图片
func PDF2PNG(pdfFile, pngFile string, pageNO, imgSize int) (err error) {
	// gs -sDEVICE=pngalpha -r256 -dNOPAUSE -dBATCH -dFirstPage=5 -dLastPage=5 -sOutputFile=output.png  input.pdf
	// r似乎表示尺寸上限, 如r144表示不超过144KiB
	pdfFile, err = filepath.Abs(pdfFile)
	if err != nil {
		return
	}

	pngFile, err = filepath.Abs(pngFile)
	if err != nil {
		return
	}

	pageNum := strconv.Itoa(pageNO)
	args := []string{"-sDEVICE=pngalpha", "-r" + strconv.Itoa(imgSize), "-dNOPAUSE", "-dBATCH", "-dFirstPage=" + pageNum, "-dLastPage=" + pageNum, "-sOutputFile=" + pngFile, pdfFile}
	cmd := exec.Command("gs", args...)
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
	log.Println("PDF 转 PNG :", cmd.Args)
	err = cmd.Run()
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
	log.Println("office 文档转 PDF:%v", cmd.Args)
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

	log.Println("非office 文档转 PDF:%v", cmd.Args)
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
	log.Println("calibre文档转换：%v %v", "ebook-convert", strings.Join(args, " "))
	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	return
}

// 使用pdfinfo获取pdf页数
func PageNumber(input string) int64 {
	if strings.ToLower(path.Ext(input)) != ".pdf" {
		return -1
	}
	input, err := filepath.Abs(input)
	if err != nil {
		return -1
	}

	cmdstr := "pdfinfo " + input + " | grep Pages | cut -d ':' -f2"
	cmd := exec.Command(cmdstr)

	log.Println("calibre文档转换：%v", cmdstr)
	time.AfterFunc(5*time.Minute, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})

	var out bytes.Buffer
	cmd.Stderr = &out

	cmd.Run()

	res := strings.Trim(out.String(), " ")
	pagenum, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		pagenum = -1
	}

	return pagenum
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
	log.Println("转化封面图片：", cmd.Args)
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
	err := exec.Command("pdftotext", args...).Run()
	if err != nil {
		log.Println(err.Error())
	}
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
		log.Println("========================", err)
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
	log.Println("PNG 图片压缩:%v", cmd.Args)
	err = cmd.Run()
	return
}

func getTextFormTxtFile(textfile string) (content string) {
	b, err := ioutil.ReadFile(textfile)
	if err != nil {
		log.Println(err.Error())
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
