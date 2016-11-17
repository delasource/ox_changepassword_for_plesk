package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DBServer        = "localhost"
	DBPort          = "3306"
	DBUser          = "oxmanager"
	DBPass          = "12345678"
	DBDatabaseOx    = "oxdatabase_5"
	OXAdmin         = "oxadmin"
	OXAdminPassword = "admin_password"
)

func main() {

	f, err := os.OpenFile("/var/log/open-xchange/pw.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening log file: " + err.Error())
		os.Exit(1)
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("\n")

	if len(os.Args) != 11 {
		log.Println("Syntax: ./changepassword - contextid - - - userid - - - newpassword")
		fmt.Println("Syntax: ./changepassword - contextid - - - userid - - - newpassword")
		os.Exit(1)
		return
	}

	newpassword := os.Args[10]
	userid := os.Args[6]
	contextid := os.Args[2]

	if len(newpassword) < 3 || len(userid) == 0 || len(contextid) == 0 {
		log.Println("Syntax: ./changepassword - contextid - - - userid - - - newpassword")
		fmt.Println("Syntax: ./changepassword - contextid - - - userid - - - newpassword")
		os.Exit(1)
		return
	}

	c0 := exec.Command("whoami")
	out0, err0 := c0.Output()
	if err0 != nil {
		log.Printf("FATAL0: %s\n", err0)
		os.Exit(1)
		return
	}
	alloutput0 := strings.TrimSpace(string(out0[:]))
	log.Println("Executor: " + alloutput0)

	log.Println("Changing password to " + newpassword + " for user " + userid + " at context " + contextid)

	var DbOx *sqlx.DB
	DbOx, err = sqlx.Open("mysql", DBUser+":"+DBPass+"@tcp("+DBServer+":"+DBPort+")/"+DBDatabaseOx)
	if err != nil {
		log.Printf("Error opening ox database: %v\n", err)
		os.Exit(1)
		return
	}
	err = DbOx.Ping()
	if err != nil {
		log.Printf("Error on opening ox database connection: %s\n", err.Error())
		os.Exit(1)
		return
	}
	DbOx = DbOx.Unsafe()
	log.Println("Step 1 ok.")

	var usermail string
	err = DbOx.Get(&usermail, "SELECT mail FROM user WHERE cid=? AND id=?", contextid, userid)
	if err != nil {
		log.Printf("Error selecting user: %s\n", err.Error())
		os.Exit(4)
		return
	}

	log.Printf("Calling sudo %s %s %s %s %s\n", "/usr/local/psa/bin/mail", "-u", usermail, "-passwd", newpassword)
	c1 := exec.Command("sudo", "/usr/local/psa/bin/mail", "-u", usermail, "-passwd", newpassword)
	out1, err1 := c1.Output()
	if err1 != nil {
		log.Printf("FATAL1: %s\n", err1)
		log.Println(strings.TrimSpace(string(out1[:])))
		os.Exit(3)
		return
	}
	alloutput1 := strings.TrimSpace(string(out1[:]))
	log.Println(alloutput1)
	log.Println("Step 2 ok.")

	log.Printf("Calling %s %s %s %s %s %s %s %s %s %s %s\n", "/opt/open-xchange/sbin/changeuser", "-c", contextid, "-A", OXAdmin, "-P", OXAdminPassword, "-i", userid, "-p", newpassword)
	c2 := exec.Command("/opt/open-xchange/sbin/changeuser", "-c", contextid, "-A", OXAdmin, "-P", OXAdminPassword, "-i", userid, "-p", newpassword)
	out2, err2 := c2.Output()
	if err2 != nil {
		log.Printf("FATAL2: %s\n", err2)
		log.Println(strings.TrimSpace(string(out2[:])))
		os.Exit(5)
		return
	}
	alloutput2 := strings.TrimSpace(string(out2[:]))
	log.Println(alloutput2)

	log.Println("Step 3 ok.")
	log.Println("done.")

	os.Exit(0)

}
