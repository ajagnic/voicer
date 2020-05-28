# voicer

## Requirements
*   Go >= 1.14 (https://golang.org/)
*   Google Cloud Platform Service-Account Key
    *   Cloud Text-to-Speech API enabled

## Build
``` go build github.com/ajagnic/voicer ```

## Run
#### Windows
By environment variable:
```powershell
$env:GOOGLE_APPLICATION_CREDENTIALS="[PATH]"

./voicer.exe
```
By JSON file:
```powershell
./voicer.exe -key='[PATH]'
```

## Author
Adrian Agnic (https://github.com/ajagnic)
