package tools

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// zipDirTool
// unzipTool

// zip_dir: comprime un directorio dentro del workspace en un archivo .zip
func zipDirTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"zip_dir", nil, "falta argumento obligatorio: path"}
	}

	outRaw, ok := args["output"]
	if !ok {
		return ToolResult{"zip_dir", nil, "falta argumento obligatorio: output"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"zip_dir", nil, "el argumento 'path' debe ser string"}
	}

	output, ok := outRaw.(string)
	if !ok {
		return ToolResult{"zip_dir", nil, "el argumento 'output' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)
	fullOutput := filepath.Join("workspace", output)

	// Verificar que el directorio existe
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"zip_dir", nil, fmt.Sprintf("el directorio '%s' no existe", path)}
		}
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error accediendo al directorio: %v", err)}
	}

	if !info.IsDir() {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("'%s' no es un directorio", path)}
	}

	// Crear directorios destino si no existen
	if err := os.MkdirAll(filepath.Dir(fullOutput), 0755); err != nil {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
	}

	// Crear archivo zip
	outFile, err := os.Create(fullOutput)
	if err != nil {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error creando archivo zip: %v", err)}
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Recorrer el directorio y añadir archivos
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(fullPath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error comprimiendo directorio: %v", err)}
	}

	return ToolResult{
		ToolName: "zip_dir",
		Result:   fmt.Sprintf("directorio '%s' comprimido correctamente en '%s'", path, output),
	}
}

// unzip: descomprime un archivo .zip dentro del workspace
func unzipTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"unzip", nil, "falta argumento obligatorio: path"}
	}

	destRaw, ok := args["dest"]
	if !ok {
		return ToolResult{"unzip", nil, "falta argumento obligatorio: dest"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"unzip", nil, "el argumento 'path' debe ser string"}
	}

	dest, ok := destRaw.(string)
	if !ok {
		return ToolResult{"unzip", nil, "el argumento 'dest' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)
	fullDest := filepath.Join("workspace", dest)

	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"unzip", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"unzip", nil, fmt.Sprintf("error accediendo al archivo: %v", err)}
	}

	// Abrir archivo zip
	r, err := zip.OpenReader(fullPath)
	if err != nil {
		return ToolResult{"unzip", nil, fmt.Sprintf("error abriendo zip: %v", err)}
	}
	defer r.Close()

	// Crear directorio destino
	if err := os.MkdirAll(fullDest, 0755); err != nil {
		return ToolResult{"unzip", nil, fmt.Sprintf("error creando directorio destino: %v", err)}
	}

	// Extraer archivos
	for _, f := range r.File {
		fpath := filepath.Join(fullDest, f.Name)

		// Evitar path traversal
		if !strings.HasPrefix(fpath, filepath.Clean(fullDest)+string(os.PathSeparator)) {
			return ToolResult{"unzip", nil, "zip contiene rutas inválidas (path traversal)"}
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}

		// Crear directorio si no existe
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return ToolResult{"unzip", nil, fmt.Sprintf("error creando directorio: %v", err)}
		}

		// Extraer archivo
		rc, err := f.Open()
		if err != nil {
			return ToolResult{"unzip", nil, fmt.Sprintf("error leyendo archivo del zip: %v", err)}
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			return ToolResult{"unzip", nil, fmt.Sprintf("error creando archivo: %v", err)}
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return ToolResult{"unzip", nil, fmt.Sprintf("error extrayendo archivo: %v", err)}
		}
	}

	return ToolResult{
		ToolName: "unzip",
		Result:   fmt.Sprintf("archivo '%s' descomprimido correctamente en '%s'", path, dest),
	}
}
