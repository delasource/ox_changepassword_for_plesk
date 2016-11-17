#!/bin/bash
mkdir /opt/open-xchange/scripts/
rm /opt/open-xchange/scripts/changepassword
go build changepassword.go
mv changepassword /opt/open-xchange/scripts/
chown -hR open-xchange:open-xchange /opt/open-xchange/scripts/
