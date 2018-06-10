package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"rsc.io/pdf"
)

func changeExtName(filename string, extention string) string {
	pos := strings.LastIndex(filename, ".")
	changedName := filename
	if pos == -1 {
		changedName = filename + extention
	} else {
		if filename[pos:] != extention {
			return ""
		} else {
			changedName = filename[:pos] + extention
		}
	}
	return changedName
}

func RenametoPDF(filename string) (string, error) {
	dst := changeExtName(filename, ".pdf")
	if dst == "" {
		return filename, nil
	} else {
		_, err := exec.Command("mv", filename, dst).CombinedOutput()
		if err != nil {
			return "", errors.New("Can't move " + filename + " to " + dst)
		}
	}
	return dst, nil
}

func PDFtoPPM(filename string) string {
	filename = changeExtName(filename, ".pdf")
	re := regexp.MustCompile("(.*).pdf")
	dst := re.ReplaceAllString(filename, "$1")
	out, err := exec.Command("pdftoppm", filename, dst).CombinedOutput()
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	fmt.Println("結果: %s", out)
	return dst
}

func PPTtoPDF(filename string) (string, error) {
	dst := changeExtName(filename, ".pdf")
	if dst == "" {
		return "", errors.New("file with extention [.pdf] can't convert")
	}
	out, err := exec.Command("unoconv", "-f", "pdf", "-o", dst, filename).CombinedOutput()
	if err != nil {
		return "", errors.New("Can't convert " + filename + " to pdf")
	}
	fmt.Println("結果: %s", out)
	return dst, nil
}

func PPTtoJPG(dir string) (error, []string) {
	ppmfiles := dir + "*ppm"
	err := exec.Command("mogrify", "-format", "jpg", ppmfiles).Run()
	if err == nil {
		files, _ := filepath.Glob(ppmfiles)
		fmt.Println(files)
		return err, files
	}
	return err, nil
}

func PDFtoTranscript(filename string) {
	pdfFile := changeExtName(filename, ".pdf")
	fmt.Println(pdfFile)
	r, err := pdf.Open(pdfFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.NumPage())
}

func filetypeCheck(filename string) (string, error) {
	out, _ := exec.Command("file", "--mime-type", "", filename).Output()
	filetype := strings.Split(string(out), ":")[2]
	pdf := "application/pdf"
	ppt := "application/vnd.ms-powerpoint"
	pptx := "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	if strings.Contains(filetype, pdf) {
		return "pdf", nil
	} else if strings.Contains(filetype, ppt) {
		return "ppt", nil
	} else if strings.Contains(filetype, pptx) {
		return "pptx", nil
	} else {
		return "", errors.New("[ERROR] invalide fileformat :[" + filetype + "]")
	}
}

func convert(dir string, filename string) {
	filepath := dir + filename
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("file doesn't exist")
	}
	filetype, err := filetypeCheck(filepath)
	if err != nil {
		fmt.Println(err)
	}
	if filetype == "ppt" || filetype == "pptx" {
		PPTtoPDF(filepath)
	} else if filetype == "pdf" {
		RenametoPDF(filepath)
	}
	PDFtoPPM(filepath)
	PPTtoJPG(dir)
	PDFtoTranscript(filepath)
}
