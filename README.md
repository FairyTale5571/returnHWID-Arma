# returnHWID-Arma
Return hardware ID to Arma
 
 Approved by battleye!
 
 Build: 
    
    go build -o returnHWID_x64.dll -buildmode=c-shared

 Usage in Arma:
 
    "upc" - серийный номер процессора
    "fontsHash" - хэш установленных шрифтов
    "midika" - HWID
    "macsie" - mac address
    "companiesname" - имя пк
    "guidreas" - Сгенерированый UUID командой VSC
    "VSC" - проверяет записан ли UUID в реестре, если нет, то создает, возвращает в игру просто текст написанный мною