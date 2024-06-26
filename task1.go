package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Функция проверки валидности url и возврат тела сайта
func getContent(url string) ([]byte, error) {

	// Отправить GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Ошибка запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	//Проверить запрос
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(time.Now().Format("01-02-2006 15:04:05"), "Неудачный ответ:", resp.Status)
	}

	// Прочитать содержимое ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(time.Now().Format("01-02-2006 15:04:05"), "Ошибка чтения ответа:", err)
	}
	return body, nil
}

// Запись тела сайта в новый файл по указанному пути
func writeBody(content []byte, dstFolder *string, url string) error {
	domen := strings.Split(url, "/")

	txtfile, err := os.Create(fmt.Sprintf("%s/%s.txt", *dstFolder, domen[len(domen)-2]))
	if err != nil {
		return fmt.Errorf(time.Now().Format("01-02-2006 15:04:05"), "Ошибка при создании файла:", err)
	}

	_, err = txtfile.Write(content)
	if err != nil {
		return fmt.Errorf(time.Now().Format("01-02-2006 15:04:05"), "Ошибка записи в файл:", err)
	}
	defer txtfile.Close()

	fmt.Println(time.Now().Format("01-02-2006 15:04:05"), fmt.Sprintf("Запись страницы %s завершена", domen[len(domen)-2]))
	return nil
}

func main() {
	start := time.Now()

	srcUrl := flag.String("src", "", "Путь до текстового файла, содержащий ссылки")
	dstFolder := flag.String("dst", "", "Путь до директории для создания нового текстового файла")

	flag.Parse()
	if *srcUrl == "" || *dstFolder == "" {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Отсутствуют данные о местоположении файла и директории.")
		fmt.Println("Ожидаемые параметры вызова программы:")
		flag.PrintDefaults()
		return
	}
	_, err := os.Stat(*dstFolder)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(*dstFolder, 0777)
		} else {
			fmt.Println("Ошибка при обнаружении директории для создания нового текстового файла:", err)
		}
	}

	// Открыть файл с URL-адресами
	file, err := os.Open(*srcUrl)
	if err != nil {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Ошибка открытия файла:", err)
		return
	}

	defer file.Close()

	// Читать строки из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		content, err := getContent(url)

		if err != nil {
			continue
		}

		writeBody(content, dstFolder, url)
	}
	duration := time.Since(start)

	fmt.Println("Программа завершена. Время выполнения:", duration)
}
