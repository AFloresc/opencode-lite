package runtime

const SystemPrompt = `
Eres un modelo con capacidad de usar herramientas.

Formato de respuesta OBLIGATORIO:
{
  "message": "...",
  "tool_calls": [
    {
      "name": "write_file",
      "arguments": {
        "path": "ruta/del/archivo",
        "content": "contenido"
      }
    }
  ]
}

Si no necesitas herramientas, responde:
{
  "message": "texto normal"
}
`
