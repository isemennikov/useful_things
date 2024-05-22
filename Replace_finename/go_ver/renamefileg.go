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
    // Check if the user provided a directory as an argument
    if len(os.Args) < 2 {
        log.Fatal("Usage: program <directory>")
    }
    dir := os.Args[1]

    // Read the directory contents
    files, err := ioutil.ReadDir(dir) //read define dir 
    if err != nil {
        log.Fatal(err)
    }

    // Initialize variables to keep track of the registry content and counts
    var registryContent strings.Builder
    pdfEpubCount := 0
    renamedFilesCount := 0

    // Iterate over the files in the directory
    for _, file := range files {
        // Skip directories
        if file.IsDir() {
            continue
        }

        fileName := file.Name()
        ext := filepath.Ext(fileName)
        // Only process .pdf and .epub files
        if ext != ".pdf" && ext != ".epub" {
            continue
        }

        pdfEpubCount++
        originalFilePath := filepath.Join(dir, fileName)
        // Sanitize and shorten the file name if necessary
        fileName = sanitizeAndShortenFileName(fileName, ext, 60)

        newFilePath := filepath.Join(dir, fileName)
        // Rename the file if the new file name does not already exist
        if _, err := os.Stat(newFilePath); os.IsNotExist(err) && originalFilePath != newFilePath {
            if err := os.Rename(originalFilePath, newFilePath); err != nil {
                log.Printf("Error renaming %s: %s\n", originalFilePath, err)
                continue
            }
            renamedFilesCount++
        }

         // Calculate the MD5 hash of the file
        hash, err := hashFileMD5(newFilePath)
        if err != nil {
            log.Printf("Error calculating hash for %s: %s\n", fileName, err)
            continue
        }
        // Calculate the file size in MB
        fileSizeMB := float64(file.Size()) / (1024 * 1024)
        // Write the file info to the registry content
        registryContent.WriteString(fmt.Sprintf("%-70s %-8s %.2f MB\n", fileName, hash[len(hash)-8:], fileSizeMB))
    }
    // Write the registry header
    registryHeader := fmt.Sprintf("\nTotal PDF and EPUB files: %d\n\n", pdfEpubCount)
    registryContent.WriteString(registryHeader)
    // Write the registry content to a file
    if err := ioutil.WriteFile(filepath.Join(dir, "registry.txt"), []byte(registryContent.String()), 0644); err != nil {
        log.Fatal(err)
    }
    // Print the completion message
    fmt.Printf("Program completed successfully. Total files processed: %d. Total files renamed: %d. The results are listed in registry.txt.\n", pdfEpubCount, renamedFilesCount)
}

// sanitizeAndShortenFileName sanitizes and shortens the file name if necessary
func sanitizeAndShortenFileName(fileName, ext string, maxLength int) string {
    // Replace spaces with underscores
    fileName = strings.Replace(fileName, " ", "_", -1)
    // Sanitize the file name by replacing invalid characters
    fileName = strings.Map(sanitizeRune, fileName)
    // Remove "dokumen.pub_" prefix if present
    if strings.HasPrefix(fileName, "dokumen.pub_") {
        fileName = strings.TrimPrefix(fileName, "dokumen.pub_")
    }
    // Shorten the file name if it exceeds the maximum length
    return shortenFileName(fileName, maxLength, ext)
}
// shortenFileName shortens the file name if it exceeds the maximum length
func shortenFileName(fileName string, maxLength int, ext string) string {
    if len(fileName) > maxLength {
        fileName = fileName[:maxLength-len(ext)] + ext
    }
    return fileName
}
// sanitizeRune replaces invalid characters in the file name with a hyphen
func sanitizeRune(r rune) rune {
    if r == '.' || r == '-' || r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
        return r
    }
    return '-'
}
// hashFileMD5 calculates the MD5 hash of a file
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
