# returnHWID-Arma
Return hardware ID to Arma
 
 Approved by battleye!
 
 Build: 
    
    go build -o returnHWID_x64.dll -buildmode=c-shared

 Usage in Arma:
 
    "returnHWID" callExtension "Machine_ID";
    
    "returnHWID" callExtension "HDD_UID";
    
    "returnHWID" callExtension "Product_Win";
    
    "returnHWID" callExtension "Mac_Address";
    
    -Result: STRING Crypted in sha256
