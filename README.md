# README #

This script adds the ability to change a users password in the appsuite backend - when using Open-Xchange with Plesk.

### Installation ###

* **Clone** this repo anywhere you want (typically in your go/src path)
* Create a new mysql user for **local** connections only and give him **SELECT** permissions on your openxchange database. In the example config it is called **oxmanager**.
* Specify these username, password and the open-xchange tablename in the file `changepassword.go` you just cloned.
* Add execution permissions to the build script with `chmod +x build.sh`
* And run the build script with root privileges: `sudo ./build.sh` (**Note** that this requires a valid GO installation in the root user's environment)
* Install [this](http://oxpedia.org/wiki/index.php?title=ChangePasswordExternal) Open-Xchange Plugin and follow the installation process
* In `/opt/open-xchange/etc/change_pwd_script.properties` set the following value at the end:
```
com.openexchange.passwordchange.script.shellscript=/opt/open-xchange/scripts/changepassword
```
* Allow the user *open-xchange* to run plesk's mail tool as root user by editing `/etc/sudoers` and add the following line:
```
# User alias specification
open-xchange ALL=(root) NOPASSWD: /usr/local/psa/bin/mail
```
* After restarting the open-xchange service, you and your customers will see a "Change Password"-Button in the settings. If you encounter any errors, review the log file at */var/log/open-xchange/pw.log*

Tested on debian 8.6
