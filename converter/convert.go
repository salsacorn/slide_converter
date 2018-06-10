package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
)

func _rename_to_pdf(filename string) error {
	cmd := "mv " + filename + " " + filename + ".pdf"
	c, err := shellwords.Parse(cmd)
	exec.Command(c[0]).Output()
	if err != nil {
		return err
	}
	return nil
}

func _pdf_to_ppm(filename string) error {
	cmd := "pdftoppm " + filename + " " + filename + ".ppm"
	c, err := shellwords.Parse(cmd)
	exec.Command(c[0]).Output()
	if err != nil {
		return err
	}
	return nil
}

func _ppt_to_pdf(filename string) error {
	cmd := "unoconv -f pdf -o " + filename + " " + filename + ".pdf"
	c, err := shellwords.Parse(cmd)
	exec.Command(c[0]).Output()
	if err != nil {
		return err
	}
	return nil
}

func filetype_check(filename string) (string, error) {
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

func pdf_to_ppm() {
	data_path := "data/download/"
	files, err := ioutil.ReadDir(data_path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filename := data_path + file.Name()
		// pptx -> pdf変換
		filetype, err := filetype_check(filename)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf(filetype)
		/*
			if err = os.Rename(filename, filename+".pptx"); err != nil {
				fmt.Println(err)
			}
		*/
	}
	out, _ := exec.Command("pwd").Output()
	fmt.Printf("結果: %s", out)
	out, _ = exec.Command("ls", "data/download/").Output()
	fmt.Printf("結果: %s", out)
}
