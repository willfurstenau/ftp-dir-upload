package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
)

func main() {
	host := os.Getenv("INPUT_HOST")
	port := os.Getenv("INPUT_PORT")
	user := os.Getenv("INPUT_USER")
	pass := os.Getenv("INPUT_PASSWORD")
	localDir := os.Getenv("INPUT_LOCALDIR")
	remoteDir := os.Getenv("INPUT_REMOTEDIR")

	addr := fmt.Sprintf("%s:%s", host, port)
	c, err := ftp.Dial(addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer c.Quit()

	if err := c.Login(user, pass); err != nil {
		log.Fatalf("login: %v", err)
	}

	// ensure remote root exists
	mkdirRecursive(c, remoteDir)

	root := filepath.Join(os.Getenv("GITHUB_WORKSPACE"), localDir)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Fatalf("local directory %q does not exist", root)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		remotePath := filepath.Join(remoteDir, filepath.ToSlash(rel))

		if info.IsDir() {
			return mkdirRecursive(c, remotePath)
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		return c.Stor(remotePath, file)
	})

	if err != nil {
		log.Fatalf("walk/upload: %v", err)
	}
	log.Println("FTP upload complete")
}

func mkdirRecursive(c *ftp.ServerConn, path string) error {
	parts := []string{}
	for _, p := range filepath.SplitList(path) {
		parts = append(parts, p)
		sub := filepath.Join(parts...)
		if err := c.MakeDir(sub); err != nil {
			// ignore "already exists" errors
			if !isFtpAlreadyExists(err) {
				return fmt.Errorf("mkdir %q: %w", sub, err)
			}
		}
	}
	return nil
}

func isFtpAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "550")
}
