package tools

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ------------------------------------------------------------
// PATH HELPERS
// ------------------------------------------------------------

// safeJoinWorkspace asegura que el path siempre esté dentro de workspace/
func safeJoinWorkspace(rel string) (string, error) {
	base := filepath.Clean("workspace")
	full := filepath.Join(base, rel)
	clean := filepath.Clean(full)

	if !strings.HasPrefix(clean, base) {
		return "", errors.New("ruta fuera de workspace")
	}
	return clean, nil
}

// fileExists verifica si un archivo existe
func fileExists(rel string) bool {
	full, err := safeJoinWorkspace(rel)
	if err != nil {
		return false
	}
	_, err = os.Stat(full)
	return err == nil
}

// ------------------------------------------------------------
// FILE READERS
// ------------------------------------------------------------

// readFile lee un archivo y devuelve su contenido como string
func readFile(rel string) (string, error) {
	full, err := safeJoinWorkspace(rel)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// readLines devuelve un slice de líneas
func readLines(rel string) ([]string, error) {
	content, err := readFile(rel)
	if err != nil {
		return nil, err
	}
	return strings.Split(content, "\n"), nil
}

// readJSON lee un archivo JSON y lo deserializa
func readJSON(rel string, out interface{}) error {
	full, err := safeJoinWorkspace(rel)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

// ------------------------------------------------------------
// FILE WRITERS
// ------------------------------------------------------------

// writeFile escribe contenido en un archivo
func writeFile(rel string, content string) error {
	full, err := safeJoinWorkspace(rel)
	if err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0644)
}

// appendFile añade contenido al final de un archivo
func appendFile(rel string, content string) error {
	full, err := safeJoinWorkspace(rel)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(full, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

// ------------------------------------------------------------
// DIRECTORY HELPERS
// ------------------------------------------------------------

// listFilesRecursive devuelve todos los archivos dentro de un directorio
func listFilesRecursive(rel string) ([]string, error) {
	root, err := safeJoinWorkspace(rel)
	if err != nil {
		return nil, err
	}

	var files []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, _ := filepath.Rel(root, path)
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

// listDirsRecursive devuelve todos los directorios dentro de un path
func listDirsRecursive(rel string) ([]string, error) {
	root, err := safeJoinWorkspace(rel)
	if err != nil {
		return nil, err
	}

	var dirs []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			relPath, _ := filepath.Rel(root, path)
			dirs = append(dirs, relPath)
		}
		return nil
	})

	return dirs, err
}
