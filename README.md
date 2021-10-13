# returnHWID-Arma

Return hardware ID to Arma

NOT Approved by battleye!

Build:

    go build -o returnHWID_x64.dll -buildmode=c-shared

Usage in Arma:

    - DRM не требуется для:
    "goarch" - архитектура
    "getPlayerUID" - uid напрямую со стима
    "checkInfiBan" - Инфистар бан
    "isBan" - бан на сервере, обрабатывать parseSimpleArray [bool, string]
    "close" - закрыть игру
    "version" - возврат версии екстеншна, запись и проверка поля в реестре с GUID
    "checkDRM" - открыт ли екстеншн, ответ YES || NO
    "4_c" - чистит папку со скринами
    "4_c_a" - удаляет все файлы с названием ScreenX_ в temp
    "info" - просто ))
    ["backdoor_unlock",["pass"]] - на всякий
    ["unlockDRM",["key"]] - Открыть dll через сервер 
    ["addServer",["ips"]] - добавить сервер для проверки, если все зашитые умрут
    ["Message",["string"]] - откроет окно с сообщением
    ["Sentry",["strings"]] - отправит в сентри как ошибку
    ["NewHardware",["strings"]] - зарегистрирует нового игрока

    - DRM требуется для:
    "isAdmin" - имеются ли права админа
    "get_HWID" - MachineGuid
    "get_HDDUID" - идентификатор диска 0
    "get_Product" - ключ Windows
    "get_Process" - список процессов
    "get_MAC" - мак адрес
    "get_GUID" - чтение GUID из реестра
    "get_IP" - получим Ip адрес
    "get_GeoIP" - получим массив данных о его ip адресе (страна, провайдер и т.д.)
    "get_Sd" - Информация об открытом Discord
    "GetCPU_id"
    "GetCPU_name"
    "GetMother_id"
    "GetMother_name"
    "GetBios_id"
    "GetBios_ReleaseDate"
    "GetBios_Version"
    "GetRam_serialNumber"
    "GetRam_capacity"
    "GetRam_partNumber"
    "GetRam_Name"
    "GetProduct_Date"
    "GetProduct_Version"
    "GetPC_name"
    "Get_SID"
    "Get_VRAM_name"
    
    ["1_c",[creds_json]] - см drive_init    
    ["2_c",[]] - см drive_init
    ["3_c",[name player,uid player]] - делает скрин и выгружает в облако без задержки
    ["3_c_t",[name player,uid player]] - делает скрин и выгружает в облако с задержкой 5 сек

    ["1_r",[раздел, путь, ключ, значение]] - Записывает по заданому пути в реестре - !только current_user!
    ["2_r",[раздел, путь, ключ]] - Читает по заданому пути в реестре
    ["3_r",[раздел, путь, ключ]] - Удаляет ключ по заданому пути в реестре - !только current_user!
    
    ["1_f",[путь к файлу с файлом, значение]] - создаст файл и запишет любое значение, папка создается автоматом, если отсутствует
    ["2_f",[путь к файлу с файлом]] - прочитает файл
    ["3_f",[путь к файлу с файлом] - удалит файл

    ["loadSqf",["script"]] - вызовет с лк исполняемый код
    ["initInfistarVision",["public","private"]] - инициализирует Infistar Vision API
    ["insertHardware",_data] - вносит данные в бд о новом акке


    