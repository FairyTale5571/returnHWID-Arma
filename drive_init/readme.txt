Как скомпилить:
 - В основной папке запустить compile.bat
 - В drive_init папке запустить compile.bat

Как подключить гугл драйв
 - Проходишь сюда https://developers.google.com/drive/api/v3/quickstart/go 
 - В первом шаге тыкаешь на синюю кнопку, выбираешь Desktop App если стоит другое и скачиваешь credentials.json
 - Закидываешь credentials.json в папку drive_init, нажимаешь на init.bat
 - В терминале появляется ссылка на аутентификацию, переходишь по ней, аутентифицируешься
 - В конце тебе дают ключ и говорят вставить в своё приложение, вставляешь в терминал и нажимаешь Enter
 - Если не выбило никаких ошибок, то всё хорошо
 
Как использовать:
"returnHWID" callExtension ["init_creds",["{creds json}"]]
"returnHWID" callExtension ["init_token",["{creds json}"]]

"returnHWID" callExtension ["screenAndUpload",["username", "uid"]]


Важно это всё передавать как строки
