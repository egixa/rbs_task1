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

	//Проверить запрос
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(time.Now().Format("01-02-2006 15:04:05"), "Неудачный ответ:", resp.Status)
	}
	defer resp.Body.Close()

	// Прочитать содержимое ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(time.Now().Format("01-02-2006 15:04:05"), "Ошибка чтения ответа:", err)
	}
	return body, nil
}

// Запись тела сайта в новый файл по указанному пути
func writeBody(content []byte, dstFolder *string, url string) {
	domen := strings.Split(url, "/")

	txtfile, err := os.Create(fmt.Sprintf("%s/%s.txt", *dstFolder, domen[len(domen)-2]))
	if err != nil {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Ошибка при создании файла:", err)
		return
	}

	_, err = txtfile.Write(content)
	if err != nil {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Ошибка записи в файл:", err)
		return
	}
	defer txtfile.Close()

	fmt.Println(time.Now().Format("01-02-2006 15:04:05"), fmt.Sprintf("Запись страницы %s завершена", domen[len(domen)-2]))
	return
}

func main() {

	srcUrl := flag.String("src", "", "Путь до текстового файла, содержащий ссылки")
	dstFolder := flag.String("dst", "", "Путь до папки для создания нового текстового файла")

	flag.Parse()
	if *srcUrl == "" || *dstFolder == "" {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Отсутствуют данные о местоположении файла и директории.")
		fmt.Println("Ожидаемые данные:")
		flag.PrintDefaults()

		if *dstFolder == "" {
			os.Mkdir(*dstFolder, 777)
		}
		return
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
	fmt.Println("цикл завершен")
}
