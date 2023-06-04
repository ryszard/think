Message from the user: {{.Message}}


{{if .CommandRun}}

This is the result of runnign the previous command. This might not be relevant if the previous command was successful.

Standard Output (truncated to 1000 bytes):
```
{{ .Stdout }}
```

Standard Error (truncated to 1000 bytes):

```
{{ .Stderr }}
```
Exit Code: {{ .ExitCode }}

{{end}}

Please remember the format in which you should reply: one line for a concise explanation, one plain text line for the proposed shell code. The explanation line should be up to 160 characters long. The code line should contain only the code, no markdown blocks. SERIOUSLY, please DO NOT PUT THE CODE IN A MARKDOWN BLOCK! This format is very important, as I will be parsing your output.