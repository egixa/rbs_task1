package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	srcUrl := flag.String("src", "", "a src")
	dstFolder := flag.String("dst", "", "a dst")

	flag.Parse()

	// Открыть файл с URL-адресами
	file, err := os.Open(*srcUrl)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()
	// Читать строки из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		// Отправить GET-запрос
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Ошибка запроса:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Неудачный ответ:", resp.Status)
			continue
		}

		// Прочитать содержимое ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка чтения ответа:", err)
			continue
		}

		// Записать содержимое ответа в файл
		file, err = os.Create(fmt.Sprintf("%s/%s.txt", *dstFolder, strings.Split(resp.TLS.ServerName, ".")))
		if err != nil {
			fmt.Println("Ошибка создания файла:", err)
			continue
		}

		_, err = file.Write(body)
		if err != nil {
			fmt.Println("Ошибка записи в файл:", err)
			continue
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка сканирования:", err)
		}

		fmt.Println("Обработка завершена.")
	}
}
