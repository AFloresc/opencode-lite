package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// applyPatchTool
// applyPatchFuzzyTool
// applyPatchAutoTool
// applyPatchStructuredTool

// apply_patch ULTRA: soporte para múltiples hunks, validación de contexto y detección de conflictos
func applyPatchTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch", nil, "falta argumento obligatorio: path"}
	}

	patchRaw, ok := args["patch"]
	if !ok {
		return ToolResult{"apply_patch", nil, "falta argumento obligatorio: patch"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch", nil, "el argumento 'path' debe ser string"}
	}

	patch, ok := patchRaw.(string)
	if !ok {
		return ToolResult{"apply_patch", nil, "el argumento 'patch' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo original
	originalBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return ToolResult{"apply_patch", nil, fmt.Sprintf("error leyendo archivo original: %v", err)}
	}

	original := strings.Split(string(originalBytes), "\n")
	result := make([]string, 0, len(original))

	lines := strings.Split(patch, "\n")
	origIdx := 0 // índice en original

	// Estado de conflicto
	conflicts := []string{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Ignorar encabezados de archivo
		if strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") {
			continue
		}

		// Hunk: @@ -a,b +c,d @@
		if strings.HasPrefix(line, "@@") {
			// Ejemplo: @@ -1,3 +1,4 @@
			parts := strings.Split(line, " ")
			if len(parts) < 3 {
				conflicts = append(conflicts, fmt.Sprintf("hunk mal formado: %s", line))
				continue
			}

			// Rango original: "-a,b"
			oldRange := strings.TrimPrefix(parts[1], "-")
			oldStartStr := strings.Split(oldRange, ",")[0]
			oldStart, err := strconv.Atoi(oldStartStr)
			if err != nil {
				conflicts = append(conflicts, fmt.Sprintf("no se pudo parsear offset original en hunk: %s", line))
				continue
			}
			oldStart-- // 1-based → 0-based

			// Copiar líneas previas sin cambios
			for origIdx < oldStart && origIdx < len(original) {
				result = append(result, original[origIdx])
				origIdx++
			}

			// Procesar hunk
			i++
			for i < len(lines) {
				hline := lines[i]

				// Nuevo hunk → retrocede uno para que el for externo lo procese
				if strings.HasPrefix(hline, "@@") {
					i--
					break
				}

				// Encabezados de archivo dentro del diff
				if strings.HasPrefix(hline, "--- ") || strings.HasPrefix(hline, "+++ ") {
					// dejamos que el siguiente ciclo los ignore
					break
				}

				// Línea eliminada: debe coincidir con original
				if strings.HasPrefix(hline, "-") {
					expected := hline[1:]
					if origIdx >= len(original) || original[origIdx] != expected {
						conflicts = append(conflicts, fmt.Sprintf(
							"conflicto: se esperaba eliminar '%s' pero en el archivo hay '%s'",
							expected,
							func() string {
								if origIdx < len(original) {
									return original[origIdx]
								}
								return "<EOF>"
							}(),
						))
					} else {
						// Coincide → se elimina avanzando el índice
						origIdx++
					}
					i++
					continue
				}

				// Línea añadida
				if strings.HasPrefix(hline, "+") {
					result = append(result, hline[1:])
					i++
					continue
				}

				// Línea de contexto (sin prefijo): debe coincidir con original
				if !strings.HasPrefix(hline, "+") && !strings.HasPrefix(hline, "-") {
					if origIdx >= len(original) || original[origIdx] != hline {
						conflicts = append(conflicts, fmt.Sprintf(
							"conflicto en contexto: se esperaba '%s' pero en el archivo hay '%s'",
							hline,
							func() string {
								if origIdx < len(original) {
									return original[origIdx]
								}
								return "<EOF>"
							}(),
						))
					} else {
						result = append(result, original[origIdx])
						origIdx++
					}
					i++
					continue
				}
			}

			continue
		}
	}

	// Copiar el resto del archivo original
	for origIdx < len(original) {
		result = append(result, original[origIdx])
		origIdx++
	}

	// Si hubo conflictos, no escribimos nada
	if len(conflicts) > 0 {
		return ToolResult{
			ToolName: "apply_patch",
			Result:   nil,
			Error:    fmt.Sprintf("conflictos al aplicar el parche:\n%s", strings.Join(conflicts, "\n")),
		}
	}

	final := strings.Join(result, "\n")
	err = os.WriteFile(fullPath, []byte(final), 0644)
	if err != nil {
		return ToolResult{"apply_patch", nil, fmt.Sprintf("error escribiendo archivo modificado: %v", err)}
	}

	return ToolResult{"apply_patch", fmt.Sprintf("parche aplicado correctamente a %s", fullPath), ""}
}

// apply_patch_fuzzy: aplica un parche aunque el contexto no coincida exactamente.
// Permite duplicar cambios y modificar líneas parcialmente.
func applyPatchFuzzyTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "falta argumento obligatorio: path"}
	}

	patchRaw, ok := args["patch"]
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "falta argumento obligatorio: patch"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "el argumento 'path' debe ser string"}
	}

	patch, ok := patchRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "el argumento 'patch' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo original
	originalBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return ToolResult{"apply_patch_fuzzy", nil, fmt.Sprintf("error leyendo archivo original: %v", err)}
	}

	original := strings.Split(string(originalBytes), "\n")
	result := make([]string, 0, len(original))

	lines := strings.Split(patch, "\n")

	// Extraer líneas - y +
	var toRemove []string
	var toAdd []string

	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			toRemove = append(toRemove, line[1:])
		}
		if strings.HasPrefix(line, "+") {
			toAdd = append(toAdd, line[1:])
		}
	}

	// Aplicación fuzzy
	for _, line := range original {
		removed := false
		for _, r := range toRemove {
			if strings.Contains(line, r) {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, line)
		}
	}

	// Añadir líneas nuevas al final
	result = append(result, toAdd...)

	final := strings.Join(result, "\n")
	err = os.WriteFile(fullPath, []byte(final), 0644)
	if err != nil {
		return ToolResult{"apply_patch_fuzzy", nil, fmt.Sprintf("error escribiendo archivo modificado: %v", err)}
	}

	return ToolResult{"apply_patch_fuzzy", fmt.Sprintf("parche fuzzy aplicado correctamente a %s", fullPath), ""}
}

// apply_patch_auto: aplica un parche inteligente sin requerir contexto exacto
func applyPatchAutoTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "falta argumento obligatorio: path"}
	}

	patchRaw, ok := args["patch"]
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "falta argumento obligatorio: patch"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "el argumento 'path' debe ser string"}
	}

	patch, ok := patchRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "el argumento 'patch' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo original
	originalBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"apply_patch_auto", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"apply_patch_auto", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	original := string(originalBytes)

	// Dividir el parche en líneas
	lines := strings.Split(patch, "\n")

	var result strings.Builder
	result.WriteString(original)

	// Aplicación automática:
	// - Si la línea empieza con "+" → añadir al final
	// - Si empieza con "-" → eliminar todas las ocurrencias
	// - Si empieza con "~" → reemplazo inteligente: "~buscar => reemplazo"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "+"):
			// Añadir al final
			result.WriteString("\n" + strings.TrimPrefix(line, "+"))

		case strings.HasPrefix(line, "-"):
			// Eliminar todas las ocurrencias
			target := strings.TrimPrefix(line, "-")
			resultStr := result.String()
			resultStr = strings.ReplaceAll(resultStr, target, "")
			result.Reset()
			result.WriteString(resultStr)

		case strings.HasPrefix(line, "~"):
			// Reemplazo inteligente "~buscar => reemplazo"
			body := strings.TrimPrefix(line, "~")
			parts := strings.SplitN(body, "=>", 2)
			if len(parts) == 2 {
				search := strings.TrimSpace(parts[0])
				replace := strings.TrimSpace(parts[1])
				resultStr := result.String()
				resultStr = strings.ReplaceAll(resultStr, search, replace)
				result.Reset()
				result.WriteString(resultStr)
			}
		}
	}

	// Guardar archivo modificado
	err = os.WriteFile(fullPath, []byte(result.String()), 0644)
	if err != nil {
		return ToolResult{"apply_patch_auto", nil, fmt.Sprintf("error escribiendo archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "apply_patch_auto",
		Result:   fmt.Sprintf("parche inteligente aplicado correctamente a '%s'", path),
	}
}

// apply_patch_structured: modificaciones semánticas de alto nivel
func applyPatchStructuredTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "falta argumento obligatorio: path"}
	}

	opRaw, ok := args["op"]
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "falta argumento obligatorio: op"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "el argumento 'path' debe ser string"}
	}

	op, ok := opRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "el argumento 'op' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	content := string(contentBytes)
	updated := content

	switch op {

	// ---------------------------------------------------------
	// Insertar import
	// ---------------------------------------------------------
	case "insert_import":
		importRaw := args["import"]
		if importRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "falta argumento 'import'"}
		}
		importLine, _ := importRaw.(string)

		re := regexp.MustCompile(`(?m)^import\s*\(`)
		if re.MatchString(updated) {
			updated = re.ReplaceAllString(updated, "import (\n    "+importLine)
		} else {
			updated = "import (\n    " + importLine + "\n)\n\n" + updated
		}

	// ---------------------------------------------------------
	// Insertar antes de una función
	// ---------------------------------------------------------
	case "insert_before_func":
		nameRaw := args["name"]
		codeRaw := args["code"]
		if nameRaw == nil || codeRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'name' y/o 'code'"}
		}
		name := nameRaw.(string)
		code := codeRaw.(string)

		re := regexp.MustCompile(`(?m)^func\s+` + regexp.QuoteMeta(name) + `\s*\(`)
		loc := re.FindStringIndex(updated)
		if loc == nil {
			return ToolResult{"apply_patch_structured", nil, "función no encontrada"}
		}

		updated = updated[:loc[0]] + code + "\n" + updated[loc[0]:]

	// ---------------------------------------------------------
	// Insertar después de una función
	// ---------------------------------------------------------
	case "insert_after_func":
		nameRaw := args["name"]
		codeRaw := args["code"]
		if nameRaw == nil || codeRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'name' y/o 'code'"}
		}
		name := nameRaw.(string)
		code := codeRaw.(string)

		re := regexp.MustCompile(`(?s)func\s+` + regexp.QuoteMeta(name) + `\s*\([^)]*\)\s*{.*?}`)
		match := re.FindStringIndex(updated)
		if match == nil {
			return ToolResult{"apply_patch_structured", nil, "función no encontrada"}
		}

		updated = updated[:match[1]] + "\n" + code + "\n" + updated[match[1]:]

	// ---------------------------------------------------------
	// Reemplazar función completa
	// ---------------------------------------------------------
	case "replace_func":
		nameRaw := args["name"]
		codeRaw := args["code"]
		if nameRaw == nil || codeRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'name' y/o 'code'"}
		}
		name := nameRaw.(string)
		code := codeRaw.(string)

		re := regexp.MustCompile(`(?s)func\s+` + regexp.QuoteMeta(name) + `\s*\([^)]*\)\s*{.*?}`)
		updated = re.ReplaceAllString(updated, code)

	// ---------------------------------------------------------
	// Eliminar función completa
	// ---------------------------------------------------------
	case "delete_func":
		nameRaw := args["name"]
		if nameRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "falta argumento 'name'"}
		}
		name := nameRaw.(string)

		re := regexp.MustCompile(`(?s)func\s+` + regexp.QuoteMeta(name) + `\s*\([^)]*\)\s*{.*?}`)
		updated = re.ReplaceAllString(updated, "")

	// ---------------------------------------------------------
	// Reemplazo por regex
	// ---------------------------------------------------------
	case "regex_replace":
		regexRaw := args["regex"]
		replaceRaw := args["replace"]
		if regexRaw == nil || replaceRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'regex' y/o 'replace'"}
		}
		regex := regexRaw.(string)
		replace := replaceRaw.(string)

		re, err := regexp.Compile(regex)
		if err != nil {
			return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("regex inválida: %v", err)}
		}

		updated = re.ReplaceAllString(updated, replace)

	default:
		return ToolResult{"apply_patch_structured", nil, "operación desconocida"}
	}

	// Guardar archivo
	err = os.WriteFile(fullPath, []byte(updated), 0644)
	if err != nil {
		return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("error escribiendo archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "apply_patch_structured",
		Result:   fmt.Sprintf("parche estructurado aplicado correctamente a '%s'", path),
	}
}
