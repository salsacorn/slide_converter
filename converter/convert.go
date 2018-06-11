package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	//      "regexp"
	"image/jpeg"
	"log"
	"strings"

	"github.com/nfnt/resize"
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

func PDFtoPPM(dir string, filename string) {
	filename = changeExtName(filename, ".pdf")
	//re := regexp.MustCompile("(.*).pdf")
	//dst := re.ReplaceAllString(filename, "$1")
	src := dir + filename
	dst := dir + "slide"
	err := exec.Command("pdftoppm", src, dst).Run()
	if err != nil {
		fmt.Println("Can't convert pdf to ppm", src, dst)
		panic(err)
	}
}

func PPTtoPDF(dir string, filename string) {
	src := dir + filename
	dst := changeExtName(filename, ".pdf")
	fmt.Println("結果:", src, dst)
	err := exec.Command("unoconv", "-f", "pdf", "-o", dst, src).Run()
	if err != nil {
		fmt.Println("Can't convert " + filename + " to pdf")
		panic(err)
	}
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

func GenerateJson(dir string, key string) {
	var imagelist []string
	jpgfiles := dir + "*jpg"
	dst := dir + "list.json"
	files, _ := filepath.Glob(jpgfiles)
	for _, file := range files {
		imagefile := "\"" + key + "/" + path.Base(file) + "\""
		imagelist = append(imagelist, imagefile)
	}
	content := "[" + strings.Join(imagelist, ",") + "]"
	file, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func resizeImage(src string, dst string, width uint, height uint) {
	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(width, height, img, resize.Lanczos3)

	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
}

func JPGtoThumbnail(dir string) {
	jpgfiles := dir + "slide-[0-9]*jpg"
	files, _ := filepath.Glob(jpgfiles)
	dst := dir + "tumbnail.jpg"
	resizeImage(files[0], dst, 320, 240)
	for _, file := range files {
		dst = dir + path.Base(file) + "-small.jpg"
		resizeImage(file, dst, 320, 240)
	}
}

func PDFtoTranscript(filename string) {
	pdfFile := changeExtName(filename, ".pdf")
	dst := dir + "transcript.txt"
	r, err := pdf.Open(pdfFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.NumPage())
	/*
	   for i = 1; i <= range r.NumPage(); i++ {
	           if i == 1 {
	                   _ := exec.Command("pdftotext", pdfFile, i, "-l", i, "-",">", dst).Run()
	                  str.gsub!(/([\r|\n|\t| |　|\u{2028}]+)/, ' ')
	           } else {
	                   _ := exec.Command("pdftotext", pdfFile, i, "-l", i, "-",">>", dst).Run()
	           }
	   }
	*/
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
		PPTtoPDF(dir, filename)
	} else if filetype == "pdf" {
		RenametoPDF(filepath)
	}
	PDFtoPPM(dir, filename)
	PPTtoJPG(dir)
	GenerateJson(dir, filename)
	JPGtoThumbnail(dir)
	//PDFtoTranscript(filepath)
}
