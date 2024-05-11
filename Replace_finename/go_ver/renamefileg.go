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
    totalFilesCount := 0

    // Добавление заголовков столбцов в registryContent
    registryContent.WriteString(fmt.Sprintf("%-70s\t%-8s\t%10s\n", "Filename", "Hash", "Size(MB)"))
    registryContent.WriteString(strings.Repeat("-", 70) + "\t" + strings.Repeat("-", 8) + "\t" + strings.Repeat("-", 10) + "\n")

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        fileName := file.Name()
        ext := filepath.Ext(fileName)
        if ext != ".pdf" && ext != ".epub" {
            continue
        }

        totalFilesCount++
        if strings.Contains(fileName, "_duplicate") {
            pdfEpubCount++
            continue // Пропускаем файлы с _duplicate в имени
        }

        originalFilePath := filepath.Join(dir, fileName)
        fileName = sanitizeAndShortenFileName(fileName, ext, 60) // Удаление префикса и сокращение имени файла

        newFilePath := originalFilePath
        if needRename(fileName) {
            fileName, newFilePath, err = renameFile(dir, fileName, ext)
            if err != nil {
                errorContent.WriteString(fmt.Sprintf("Error renaming %s: %s\n", fileName, err))
                continue
            }
            fmt.Printf("Renamed to %s\n", fileName)
        }

        hash, err := hashFileMD5(newFilePath)
        if err != nil {
            errorContent.WriteString(fmt.Sprintf("Error calculating hash for %s: %s\n", fileName, err))
            continue
        }

        fileSizeMB := float64(file.Size()) / (1024 * 1024)
        registryContent.WriteString(fmt.Sprintf("%-70s\t%-8s\t%.2f MB\n", fileName, hash[len(hash)-8:], fileSizeMB))
        pdfEpubCount++
    }

    // Вывод общего количества файлов в начало registryContent
    totalFilesContent := fmt.Sprintf("Total PDF and EPUB files: %d\n\n", totalFilesCount)
    registryContent.WriteString(totalFilesContent)

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

func needRename(fileName string) bool {
    return strings.Contains(fileName, " ") || strings.HasPrefix(fileName, "dokumen.pub_")
}

func renameFile(dir, fileName, ext string) (newFileName, newFilePath string, err error) {
    newFileName = strings.Replace(fileName, " ", "_", -1)
    newFileName = strings.Map(sanitizeRune, newFileName)
    if strings.HasPrefix(newFileName, "dokumen.pub_") {
        newFileName = strings.TrimPrefix(newFileName, "dokumen.pub_")
    }
    newFileName = shortenFileName(newFileName, 60, ext)

    newFilePath = filepath.Join(dir, newFileName)
    if newFileName != fileName {
        err = os.Rename(filepath.Join(dir, fileName), newFilePath)
    }
    return
}

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

func sanitizeAndShortenFileName(fileName, ext string, maxLength int) string {
    fileName = strings.Replace(fileName, " ", "_", -1)
    fileName = strings.Map(sanitizeRune, fileName)
    if strings.HasPrefix(fileName, "dokumen.pub_") {
        fileName = strings.TrimPrefix(fileName, "dokumen.pub_")
    }
    return shortenFileName(fileName, maxLength, ext)
}

func shortenFileName(fileName string, maxLength int, ext string) string {
    if len(fileName) > maxLength {
        fileName = fileName[:maxLength-len(ext)] + ext
    }
    return fileName
}

func sanitizeRune(r rune) rune {
    if r == '.' || r == '-' || r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
        return r
    }
    return '-'
}