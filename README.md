# 🧠 Agent Runtime — Filesystem & Patch Engine

Este proyecto implementa un **runtime de agente** en Go con un conjunto completo de herramientas para manipular un workspace local.  
El objetivo es proporcionar a un agente de IA un entorno seguro, modular y extensible para:

- leer y escribir archivos  
- aplicar parches inteligentes  
- realizar refactors estructurados  
- explorar el filesystem  
- comprimir y descomprimir  
- formatear código  
- buscar patrones complejos  
- gestionar proyectos completos  

El resultado es un **mini‑sistema operativo para agentes**, con capacidades equivalentes a un editor de código moderno.

---

## 📦 Estructura general

El runtime expone un **registry de herramientas**, cada una implementada como una función Go que recibe argumentos dinámicos (`map[string]interface{}`) y devuelve un `ToolResult`.

Todas las operaciones se realizan dentro del directorio:

```text
    workspace/
```


Esto garantiza aislamiento y seguridad.

---

# 🛠️ Herramientas disponibles

A continuación se listan todas las herramientas implementadas, agrupadas por categoría.

---

## 📁 Filesystem básico

### **read_file**
Lee el contenido de un archivo.

### **write_file**
Sobrescribe completamente un archivo.

### **create_file**
Crea un archivo nuevo (con contenido opcional).

### **delete_file**
Elimina un archivo.

### **rename_file**
Renombra un archivo o carpeta.

### **copy_file**
Copia un archivo.

### **move_file**
Mueve un archivo.

### **file_exists**
Comprueba si un archivo o directorio existe.

### **read_dir**
Lista el contenido de un directorio con metadatos.

### **stat_file**
Devuelve metadatos de un archivo o directorio.

### **touch_file**
Crea un archivo vacío o actualiza su timestamp.

---

## ✏️ Edición de archivos

### **append_file**
Añade contenido al final de un archivo.

### **truncate_file**
Vacía completamente un archivo sin eliminarlo.

---

## 🔍 Búsqueda y análisis

### **search_in_file**
Búsqueda exacta dentro de un archivo.

### **grep**
Búsqueda en múltiples archivos.

### **search_replace**
Buscar y reemplazar texto plano.

### **search_regex**
Búsqueda avanzada con expresiones regulares.

---

## 🧩 Parches y refactors

### **apply_patch**
Parche estricto estilo diff.

### **apply_patch_fuzzy**
Parche tolerante a cambios en el contexto.

### **apply_patch_auto**
Parche inteligente sin contexto exacto:
- `+` añadir  
- `-` eliminar  
- `~a => b` reemplazar  

### **apply_patch_structured**
Refactor semántico:
- insertar imports  
- insertar antes/después de funciones  
- reemplazar funciones  
- eliminar funciones  
- reemplazo por regex estructural  

---

## 🧰 Utilidades avanzadas

### **diff_files**
Comparación estilo unified diff entre dos archivos.

### **zip_dir**
Comprime un directorio en un `.zip`.

### **unzip**
Descomprime un archivo `.zip`.

### **format_code**
Autoformateo de código:
- Go (gofmt real)
- JSON
- YAML (limpieza básica)
- Genérico (normalización)

---

# 🖥️ Comandos sandbox (`run_command`)

El runtime incluye un **intérprete seguro** que ejecuta comandos internos sin tocar el sistema operativo real.

---

## 📏 Comandos básicos

| Comando | Descripción |
|--------|-------------|
| `count_lines <archivo>` | Cuenta líneas |
| `file_size <archivo>` | Tamaño en bytes |
| `validate_json <archivo>` | Valida JSON |
| `echo <texto>` | Devuelve el texto |
| `word_count <archivo>` | Cuenta palabras |
| `char_count <archivo>` | Cuenta caracteres |
| `sha256 <archivo>` | Hash SHA‑256 |
| `list_dir <directorio>` | Lista archivos |
| `head <archivo> <n>` | Primeras n líneas |
| `tail <archivo> <n>` | Últimas n líneas |
| `search <archivo> <texto>` | Búsqueda exacta |
| `now` | Fecha/hora actual |

---

## 🔧 Comandos de análisis de código

| Comando | Descripción |
|--------|-------------|
| `count_funcs <archivo.go>` | Cuenta funciones |
| `count_imports <archivo.go>` | Cuenta imports |
| `find_structs <archivo.go>` | Lista structs |
| `find_interfaces <archivo.go>` | Lista interfaces |

---

## 📦 Comandos de análisis de proyecto

| Comando | Descripción |
|--------|-------------|
| `project_stats` | Número de archivos y directorios |
| `largest_files` | Top 10 archivos más grandes |
| `file_tree` | Árbol completo del workspace |

---

## 🧠 Comandos inteligentes

| Comando | Descripción |
|--------|-------------|
| `detect_language <archivo>` | Detecta lenguaje por extensión |
| `summarize_file <archivo>` | Primeras 5 líneas |
| `extract_comments <archivo.go>` | Extrae comentarios `//` |

---

# 🧩 Registro de herramientas

Todas las herramientas se registran en:

```go
var toolRegistry = map[string]func(map[string]interface{}) ToolResult{
    // filesystem
    "read_file": readFileTool,
    "write_file": writeFileTool,
    "create_file": createFileTool,
    "delete_file": deleteFileTool,
    "rename_file": renameFileTool,
    "copy_file": copyFileTool,
    "move_file": moveFileTool,
    "file_exists": fileExistsTool,
    "read_dir": readDirTool,
    "stat_file": statFileTool,
    "touch_file": touchFileTool,

    // edición
    "append_file": appendFileTool,
    "truncate_file": truncateFileTool,

    // búsqueda
    "search_in_file": searchInFileTool,
    "grep": grepTool,
    "search_replace": searchReplaceTool,
    "search_regex": searchRegexTool,

    // parches
    "apply_patch": applyPatchTool,
    "apply_patch_fuzzy": applyPatchFuzzyTool,
    "apply_patch_auto": applyPatchAutoTool,
    "apply_patch_structured": applyPatchStructuredTool,

    // utilidades
    "diff_files": diffFilesTool,
    "zip_dir": zipDirTool,
    "unzip": unzipTool,
    "format_code": formatCodeTool,
}
```
# 🚀 Capacidades del agente
Con este runtime, un agente puede:

- navegar el filesystem como un IDE

- modificar código de forma segura

- aplicar refactors semánticos

- formatear código automáticamente

- analizar proyectos completos

- generar y restaurar backups

- realizar búsquedas avanzadas

- comparar versiones de archivos

Es un entorno de desarrollo completo, diseñado para agentes autónomos.

----