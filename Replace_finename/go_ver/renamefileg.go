package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: program <directory>")
    }
    dir := os.Args[1]

    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    var registryContent strings.Builder
    var errorContent strings.Builder
    pdfEpubCount := 0
    hashColumnWidth := 8
    sizeColumnWidth := 10

    // Заголовки столбцов для registryContent
    registryContent.WriteString(fmt.Sprintf("%-40s %-8s %10s\n", "Filename", "Hash", "Size(MB)"))
    registryContent.WriteString(strings.Repeat("-", 40) + " " + strings.Repeat("-", hashColumnWidth) + " " + strings.Repeat("-", sizeColumnWidth) + "\n")

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        fileName := file.Name()
        ext := filepath.Ext(fileName)
        if ext != ".pdf" && ext != ".epub" {
            continue
        }

        pdfEpubCount++
        originalFileName := fileName
        // Очистка и сокращение имени файла
        fileName = sanitizeAndShortenFileName(fileName, ext)

        newFilePath := filepath.Join(dir, fileName)
        // Проверка на дублирование и переименование при необходимости
        if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
            fileName = strings.TrimSuffix(fileName, ext) + "_duplicate" + ext
            newFilePath = filepath.Join(dir, fileName)
        }

        // Переименование файла
        if originalFileName != fileName {
            err := os.Rename(filepath.Join(dir, originalFileName), newFilePath)
            if err != nil {
                errorContent.WriteString(fmt.Sprintf("Error renaming %s to %s: %s\n", originalFileName, fileName, err))
                continue
            }
        }

        // Вычисление хеша файла
        hash, err := hashFileMD5(newFilePath)
        if err != nil {
            errorContent.WriteString(fmt.Sprintf("Error calculating hash for %s: %s\n", fileName, err))
            continue
        }

        // Вычисление размера файла
        fileSizeMB := float64(file.Size()) / (1024 * 1024)
        // Форматирование вывода в столбики
        registryContent.WriteString(fmt.Sprintf("%-40s %-8s %*.2f MB\n", trimFileName(fileName, 40), hash[len(hash)-hashColumnWidth:], sizeColumnWidth, fileSizeMB))
    }

    fmt.Printf("Processed %d PDF and EPUB files.\n", pdfEpubCount)
    if errorContent.Len() > 0 {
        fmt.Printf("Completed with errors. Please check the error log in registry.txt.\n")
    } else {
        fmt.Printf("Program completed successfully. All files have been processed and listed in registry.txt.\n")
    }

    // Запись в файл registry.txt
    err = ioutil.WriteFile(filepath.Join(dir, "registry.txt"), []byte(registryContent.String()), 0644)
    if err != nil {
        log.Fatal(err)
    }
}

// Функции для очистки имени файла и его сокращения
func sanitizeAndShortenFileName(fileName, ext string) string {
    fileName = strings.Replace(fileName, " ", "_", -1)
    fileName = strings.Map(sanitizeRune, fileName)
    if strings.HasPrefix(fileName, "dokumen.pub_") {
        fileName = strings.TrimPrefix(fileName, "dokumen.pub_")
    }
    // Обрезаем имя файла до 40 символов, если необходимо
    if len(fileName) > 60 {
        fileName = fileName[:60-len(ext)] + ext
    }
    return fileName
}

// Функция для обрезания имени файла до заданной длины
func trimFileName(fileName string, maxLength int) string {
    if len(fileName) <= maxLength {
        return fileName
    }
    return fileName[:maxLength]
}

// Функция для получения MD5-хеша файла
func hashFileMD5(filePath string) (string, error) {
    var returnMD5String string
    file, err := os.Open(filePath)
    if err != nil {
        return returnMD5String, err
    }
    defer file.Close()

    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return returnMD5String, err
    }

    hashInBytes := hash.Sum(nil)[:16]
    returnMD5String = hex.EncodeToString(hashInBytes)
    return returnMD5String, nil
}

// Функция для фильтрации разрешенных символов в имени файла
func sanitizeRune(r rune) rune {
    if r == ' ' {
        return '_'
    }
    if r == '.' || r == '-' || r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
        return r
    }
    return '-'
}