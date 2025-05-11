package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Размер порции данных при считывании файла (4KB)
const chunkSize = 4096

func hashFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	buf := make([]byte, chunkSize)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}

		_, err = hash.Write(buf[:n])
		if err != nil {
			return "", err
		}
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString, nil
}

func walkDir(dir string, fileHashes map[string]string, mux *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash, err := hashFileMD5(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error hashing file %s: %v\n", path, err)
				return err
			}

			mux.Lock()
			if existingPath, found := fileHashes[hash]; found {
				msg := fmt.Sprintf("Duplicate file found: %s and %s have the same hash %s", existingPath, path, hash)
				fmt.Fprintln(os.Stdout, msg)
			} else {
				fileHashes[hash] = path
			}
			mux.Unlock()
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", dir, err)
	}
}


func printHelp() {
	fmt.Println(`Usage: filehash <dir1> <dir2> ...
This program recursively compares hash128 of files in the specified directories
and reports any duplicate files found.

Examples:
  docker run --rm -v /path/to/dir1:/dir1 -v /path/to/dir2:/dir2 filehash /dir1 /dir2`)
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		printHelp()
		return
	}

	dirs := flag.Args()

	fileHashes := make(map[string]string)
	var mux sync.Mutex
	var wg sync.WaitGroup

	for _, dir := range dirs {
		wg.Add(1)
		go walkDir(dir, fileHashes, &mux, &wg)
	}

	wg.Wait()
}
