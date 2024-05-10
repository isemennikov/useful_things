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

        if _, err := os.Stat(filepath.Join(dir, fileName)); err == nil {
            fileName = strings.TrimSuffix(fileName, ext) + "_duplicate_" + ext
        }

        if originalFileName != fileName {
            err := os.Rename(filepath.Join(dir, originalFileName), filepath.Join(dir, fileName))
            if err != nil {
                errorContent += fmt.Sprintf("Error renaming %s to %s: %s\n", originalFileName, fileName, err)
                continue
            }
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