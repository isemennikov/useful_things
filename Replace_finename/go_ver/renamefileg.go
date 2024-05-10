package main

import (
    "fmt"
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

    var registryContent string
    var errorContent string
    pdfEpubCount := 0

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
        fileName = strings.Replace(fileName, " ", "_", -1)
        fileName = strings.Map(sanitizeRune, fileName)
        if strings.HasPrefix(fileName, "dokumen.pub_") {
            fileName = strings.TrimPrefix(fileName, "dokumen.pub_")
        }

        // Обрезаем имя файла до 40 символов, если необходимо
        if len(fileName) > 60 {
            fileName = fileName[:60-len(ext)] + ext
        }

        newFilePath := filepath.Join(dir, fileName)
        if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
            fileName = strings.TrimSuffix(fileName, ext) + "_duplicate" + ext
            newFilePath = filepath.Join(dir, fileName)
        }

        if originalFileName != fileName {
            err := os.Rename(filepath.Join(dir, originalFileName), newFilePath)
            if err != nil {
                errorContent += fmt.Sprintf("Error renaming %s to %s: %s\n", originalFileName, fileName, err)
                continue
            }
            fmt.Printf("Renamed %s to %s\n", originalFileName, fileName)
        } else {
            fmt.Printf("No need to rename %s\n", originalFileName)
        }

        fileSizeMB := float64(file.Size()) / (1024 * 1024)
        registryContent += fmt.Sprintf("%s - %.2f MB\n", fileName, fileSizeMB)
    }

    registryContent = fmt.Sprintf("Total PDF and EPUB files: %d\n\n%s", pdfEpubCount, registryContent)
    if errorContent != "" {
        registryContent += "\nErrors:\n" + errorContent
    }

    err = ioutil.WriteFile(filepath.Join(dir, "registry.txt"), []byte(registryContent), 0644)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(registryContent)
    if errorContent != "" {
        fmt.Printf("Completed with errors. Please check the error log above.\n")
    } else {
        fmt.Printf("Program completed successfully. All files have been processed.\n")
    }
}

func sanitizeRune(r rune) rune {
    if r == ' ' {
        return '_'
    }
    if r == '.' || r == '-' || r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
        return r
    }
    return '-'
}