# BT wifi auto login software

There are three login types for BT Wifi.\
BT Business Broadband has tested by me.\
BT Broadband and BT Wifi login types are propably working.\
I am not sure because I don't have accounts to test them.

You need to put your details in **app.txt** and run the sofware.\
You can access **app.txt** file from here https://github.com/s4l1h/bt-wifi-auto-login/blob/main/app.txt \
Download automatic prebuilt software from https://github.com/s4l1h/bt-wifi-auto-login/releases \
Or you can build your self. It's written in GO. https://go.dev/doc/tutorial/compile-install


#### We will check if we are connected every app_login_check_timer minutes
app_login_check_timer=1

#### Login Details will be posting this URL. It depends on your account type.
app_login_url=https://www.btwifi.com:8443/ante?partnerNetwork=btb

#### Usualy you dont need to change this. Checking rule. Basically we are cheking the keyword in the app_login_check_url 
app_login_check_keyword=now logged on to BT\
app_login_check_url=https://www.btwifi.com:8443/\

#### Anything start with header_ will be sending via header.
header_Referer=https://www.btwifi.com:8443/home\
header_Origin=https://www.btwifi.com:8443\

#### Anything start with post_ will be posting to app_login_url address.
post_username=youremailaddress@hotmail.co.uk\
post_password=yourpassword\
post_xhtmlLogon=https://www.btwifi.com:8443/ante ***-> it depends on your account type***
