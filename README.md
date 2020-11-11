# returnHWID-Arma
Return hardware ID to Arma
 
 Approved by battleye!
 
 Build: 
    
    go build -o returnHWID_x64.dll -buildmode=c-shared

 Usage in Arma:
 
    "version" - возврат версии екстеншна, запись и проверка поля в реестре с GUID
    "GUID" - чтение GUID из реестра
    "hwid" - MachineGuid
    "HDD_UID" - идентификатор диска 0
    "Product_Win" - ключ Windows
    "processList" - список процессов
    "MAC" - мак адрес
    "serials" - Список мак адресов и еще пачки устройств, парсится parseSimpleArray
    "info" - просто ))
    
    ["credentials",[creds_json]] - см drive_init    
    ["token",[]] - см drive_init
    ["doit",[name player,uid player]] - делает скрин и выгружает в облако, ойойо
    ["write_reg",[раздел, путь, ключ, значение]] - Записывает по заданому пути в реестре - !только current_user!
    ["read_reg",[раздел, путь, ключ]] - Читает по заданому пути в реестре
    ["del_reg",[раздел, путь, ключ]] - Удаляет ключ по заданому пути в реестре - !только current_user!
    
    ["read_file",[путь к файлу с файлом]] - прочитает файл
    ["write_file",[путь к файлу с файлом, значение]] - создаст файл и запишет любое значение, папка создается автоматом, если отсутствует
    ["delete_file",[путь к файлу с файлом] - удалит файл
    
    ["ew",["arg1","arg2"]] - читает wmic (wmic arg1 get arg2)
    
    
