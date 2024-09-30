package utils

import (
	"errors"
	"gurusaranm0025/cbak/pkg/conf"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

func CopyDir(srcDir, destDir string) error {

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			err := CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err := CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateCacheDir(cacheDirName string) (string, error) {
	var err error
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDirPath := filepath.Join(homeDir, conf.CachePath, cacheDirName)
	pathInfo, err := os.Stat(cacheDirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(cacheDirPath, os.ModePerm)
		if err != nil {
			return "", err
		}
		return cacheDirPath, nil
	}

	if err != nil {
		return "", err
	}

	if !pathInfo.IsDir() {
		return "", errors.New("path is not a directory")
	}

	if err = os.RemoveAll(cacheDirPath); err != nil {
		return "", err
	}

	err = os.MkdirAll(cacheDirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return cacheDirPath, nil
}

func RunAsSudo() {
	cmd := exec.Command("sudo", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}
