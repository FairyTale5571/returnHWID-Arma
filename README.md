# returnHWID-Arma
Return hardware ID to Arma
 
 Approved by battleye!
 
 Build: 
    
    go build -o returnHWID_x64.dll -buildmode=c-shared

 Usage in Arma:
 
    "midika" -HWID
    "hardidi" -HDD 0 ID
    "windidi" -Product Key
    "macsie" -mac address
    "companiesname" - computer name
    "guidreas" - unique UUID generated this dll in registry windows
    "VSC" - проверяет записан ли UUID в реестре, если нет, то создает, возвращает в игру просто текст написанный мною

    -Result: STRING Crypted in sha256
