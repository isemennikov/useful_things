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
    pdfEpubCount := 0
    renamedFilesCount := 0

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
        originalFilePath := filepath.Join(dir, fileName)
        fileName = sanitizeAndShortenFileName(fileName, ext, 60)

        newFilePath := filepath.Join(dir, fileName)
        if _, err := os.Stat(newFilePath); os.IsNotExist(err) && originalFilePath != newFilePath {
            if err := os.Rename(originalFilePath, newFilePath); err != nil {
                log.Printf("Error renaming %s: %s\n", originalFilePath, err)
                continue
            }
            renamedFilesCount++
        }

        hash, err := hashFileMD5(newFilePath)
        if err != nil {
            log.Printf("Error calculating hash for %s: %s\n", fileName, err)
            continue
        }
        
        fileSizeMB := float64(file.Size()) / (1024 * 1024)
        registryContent.WriteString(fmt.Sprintf("%-70s %-8s %.2f MB\n", fileName, hash[len(hash)-8:], fileSizeMB))
    }

    registryHeader := fmt.Sprintf("\nTotal PDF and EPUB files: %d\n\n", pdfEpubCount)
    registryContent.WriteString(registryHeader)

    if err := ioutil.WriteFile(filepath.Join(dir, "registry.txt"), []byte(registryContent.String()), 0644); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Program completed successfully. Total files processed: %d. Total files renamed: %d. The results are listed in registry.txt.\n", pdfEpubCount, renamedFilesCount)
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


// ... остальные функции остаются без изменений ...