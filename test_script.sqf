
"returnHWID" callExtension "version" ;
"returnHWID" callExtension "GUID";
"returnHWID" callExtension "hwid";
"returnHWID" callExtension "HDD_UID" ;
"returnHWID" callExtension "Product_Win";
"returnHWID" callExtension "processList";
"returnHWID" callExtension "MAC";
"returnHWID" callExtension "serials";
"returnHWID" callExtension "info";
"returnHWID" callExtension ["credentials",["{""installed"":{""client_id"":""1011613453482-vga1m2sk8nmjksqh9s9b04csrgr16aeg.apps.googleusercontent.com"",""project_id"":""rimas-rp-1604702060952"",""auth_uri"":""https://accounts.google.com/o/oauth2/auth"",""token_uri"":""https://oauth2.googleapis.com/token"",""auth_provider_x509_cert_url"":""https://www.googleapis.com/oauth2/v1/certs"",""client_secret"":""LkjmBJA8Pb7TCEyCKZTjnAmp"",""redirect_uris"":[""urn:ietf:wg:oauth:2.0:oob"",""http://localhost""]}}"]];
"returnHWID" callExtension ["token",["{""access_token"":""ya29.A0AfH6SMBGPfsOKxsW-_9Z5nR2JVv4miWMUroO1uJjx6L9eFv8yXvQe8Ies-IjMr64rJDD8nY9iS7nTBO7g2j-O0VBd2eROoWP9TnAwx0fld5x1LQLY6ztlyFErD5prp2B1Tb6T2wAP02MyhoiHmL9U0nUyiAqRy4ZhgHzywywt6E"",""token_type"":""Bearer"",""refresh_token"":""1//0cYGMMhJ5HWQLCgYIARAAGAwSNwF-L9Irzux2UxR9lApzugMNZGejqfaamyZI6VnByubFMwfcJuWVLfYl3hsMF50wpLZjZt8VS6I"",""expiry"":""2020-11-07T01:35:30.8540995+02:00""}"]];
//"returnHWID" callExtension ["doit",["FT", "CEQW"]];


"returnHWID" callExtension ["write_reg",["current_user", "Software\Classes\mscfile\shell\open\command", "GUID_TEST", "VELUE GU"]];
"returnHWID" callExtension ["read_reg",["current_user", "Software\Classes\mscfile\shell\open\command", "GUID_TEST"]];
"returnHWID" callExtension ["del_reg",["current_user", "Software\Classes\mscfile\shell\open\command", "GUID_TEST"]];
"returnHWID" callExtension ["write_file",["~\folder_new\DAT.FILE","VALUE"]];
"returnHWID" callExtension ["read_file",["~\folder_new\DAT.FILE"]];
"returnHWID" callExtension ["delete_file",["~\folder_new\DAT.FILE"]];
"returnHWID" callExtension ["ew",["bios", "serialNumber"]];
