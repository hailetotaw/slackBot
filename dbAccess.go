package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// initial scripts
	commandListScript = `
	CREATE TABLE IF NOT EXISTS goBotDB.CommandList (
		id int AUTO_INCREMENT,
		command varchar(30),
		CONSTRAINT CommandList_pk PRIMARY KEY (id)
	);
	`
	amharicWordsScript = `
	CREATE TABLE IF NOT EXISTS goBotDB.AmharicWords (
		id int AUTO_INCREMENT,
		commandId int,
		word nvarchar(5000),
		CONSTRAINT AmharicWords_pk PRIMARY KEY (id)
	);
	`
	truncateCommandListScript = "truncate table goBotDB.CommandList;"
	truncateAmWordScript      = "truncate table goBotDB.AmharicWords;"

	insertCommandlistScript = `
	INSERT INTO goBotDB.CommandList(command) values ('[@bot_user] numbers'),
	('[@bot_user] days'),
	('[@bot_user] months'),
	('[@bot_user] greetings');
	`
	insertAmWordsScript = `
	INSERT INTO goBotDB.AmharicWords(commandId, word) VALUES
	(1, '"\t\t\tPronouncation\tMeaning\r\nአንድ\t\t\t\"ande\"\t\t\tone\r\nሁለት\t\t\t\"hulet\"\t\t\ttwo\r\nሶስት\t\t\t\"soste\"\t\t\tthree\r\nአራት\t\t\t\"arate\"\t\t\tfour\r\nአምስት\t\t\t\"amesete\"\t\tfive\r\nስድስት\t\t\t\"sedeset\"\t\tsix\r\nሰባት\t\t\t\"sebate\"\t\tseven\r\nስምት\t\t\t\"semente\"\t\teight\t\r\nዘጠኝ\t\t\t\"zetegne\"\t\tnine\r\nአስር\t\t\t\"aser\"\t\t\tten\r\nሃያ\t\t\t\"haya\"\t\t\ttwenty\r\nሰላሳ\t\t\t\"selasa\"\t\tthirty\r\nአርባ\t\t\t\"arba\"\t\t\tfourty\r\nሃምሳ\t\t\t\"hamsa\"\t\t\tfifty\r\nስልሳ\t\t\t\"selsa\"\t\t\tsixty\r\nሰባ\t\t\t\"seba\"\t\t\tseventy\r\nሰማንያ\t\t\t\"semanya\"\t\teighty\r\nዘጠና\t\t\t\"zetena\"\t\tninety\r\nመቶ\t\t\t\"meto\"\t\t\thundred\r\nሺህ\t\t\t\"she\"\t\t\tthousand"'),
	(2, '"\t\t\tPronouncation\tMeaning\r\nሰኞ\t\t\t\"segno\"\t\t\tMonday\r\nማክሰኞ\t\t\"maksegno\"\t\tTuesday\r\nእሮብ\t\t\t\"erobe\"\t\t\tWednesday\r\nሐሙስ\t\t\t\"hamuse\"\t\tThursday\r\nአርብ\t\t\t\"arbe\"\t\t\tFriday\r\nቅዳሜ\t\t\t\"kidame\"\t\tSaturday\r\nእሁድ\t\t\t\"ehude\"\t\t\tSunday"'),
	(3, '"\t\t\tPronouncation\tMeaning\r\nመስከረም\t\t\"mesekerem\"\t\tSeptember\r\nጥቅምት\t\t\"tikimte\"\t\tOctober\r\nህዳር\t\t\t\"hidar\"\t\t\tNovember\r\nታህሳስ\t\t\t\"tahisase\"\t\tDecember\r\nጥር\t\t\t\"tir\"\t\t\tJanuary\r\nየካቲት\t\t\t\"yekatit\"\t\tFebruary\r\nመጋቢት\t\t\"megabit\"\t\tMarch\r\nሚያዚያ\t\t\"miyazia\"\t\tApril\r\nግንቦት\t\t\t\"genbot\"\t\tMay\r\nሰኔ\t\t\t\"sene\"\t\t\tJune\r\nሀምሌ\t\t\t\"hamele\"\t\tJuly\r\nነሐሴ\t\t\t\"nehase\"\t\tAugust\r\nጳጉሜ\t\t\t\"pagume\"\t\t-"'),
	(4, '"\t\t\t\t\tPronouncation\t\t\t\t\t\tMeaning\r\nእንደምን አደርክ\t\t\t\"endemen aderek\"\t\t\t\tgood morning\t\r\nእንደምን ዋልክ\t\t\t\"endemen walke\"\t\t\t\t\tgood afternoon\r\nእንደምን አመሸህ\t\t\t\"endemen ameshehe\"\t\t\t\tgood eveneing\r\nስምህ ማን ይባላል\t\t\t\"semehe mane yebalele\"\t\t\twhat is you name\r\nስሜ ሃይሌ ነው\t\t\t\"seme haile new\"\t\t\t\tmy name is haile\r\nእንዴት ነህ\t\t\t\t\"endet nehe\"\t\t\t\t\thow are you\r\nእድሜዎ ስንት ነው\t\t\t\"edemewo sente new\"\t\t\t\thow old are you\r\nየ 33 ዓመት ወጣት ነኝ\t\t\"ye 33 amet wetate negne\"\t\ti am 33 years old\r\nየት ትኖራለህ\t\t\t\t\"yet tinoralhe\"\t\t\t\t\twhere do you live"');
	`
)

type dbAccess struct {
	conString   string
	commandList []CommandList
}

//CommandList matching structure to the table CommandList in our database
type CommandList struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

//AmharicWords matching structure to the table AmharicWords in our database
type AmharicWords struct {
	ID        int    `json:"id"`
	CommandID int    `json:"commandId"`
	Word      string `json:"word"`
}

func (con *dbAccess) getListofCommands() []CommandList {
	// open the database connection
	db, err := sql.Open("mysql", con.conString)

	if err != nil {
		panic(err.Error())
	}

	//make sure to always close the connection
	defer db.Close()
	var listOfCommands []CommandList
	//select query
	results, err := db.Query("SELECT id, command from CommandList")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var command CommandList
		err = results.Scan(&command.ID, &command.Command)
		if err != nil {
			panic(err.Error())
		}
		listOfCommands = append(listOfCommands, command)
	}
	con.commandList = listOfCommands
	return listOfCommands
}

func (con *dbAccess) getAmharicWord(w string) string {
	// open the database connection
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()
	var amharicWord string

	stmtOut, err := db.Prepare("SELECT word from AmharicWords where commandId in (SELECT id from CommandList where command like ?)")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	//select query
	err = stmtOut.QueryRow("%" + w + "%").Scan(&amharicWord)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "invalid command"
		}
		panic(err.Error())
	}
	return amharicWord
}

func (con *dbAccess) getAmharicWordByID(commandID string) string {
	// open the database connection
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()
	var amharicWord string

	//select query
	err = db.QueryRow("SELECT word from AmharicWords where commandId = ?", commandID).Scan(&amharicWord)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "invalid command"
		}
		panic(err.Error())
	}
	return amharicWord
}

//this will create the tables in the database and seeds the data
func (con *dbAccess) prepareDatabase() {
	con.createCommandListTable()
	con.createAmharicWordsTable()
	con.truncateCommandListTable()
	con.truncateAmharicWordsTable()
	con.insertCommandList()
	con.insertAmharicWords()

}

func (con *dbAccess) createAmharicWordsTable() {
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()

	//select query
	result, err := db.Query(amharicWordsScript)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

}

func (con *dbAccess) createCommandListTable() {
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()

	//select query
	result, err := db.Query(commandListScript)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

}

func (con *dbAccess) truncateCommandListTable() {
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()

	//select query
	result, err := db.Query(truncateCommandListScript)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

}

func (con *dbAccess) truncateAmharicWordsTable() {
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()

	//select query
	result, err := db.Query(truncateAmWordScript)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

}

func (con *dbAccess) insertCommandList() {
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()

	//select query
	result, err := db.Query(insertCommandlistScript)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

}

func (con *dbAccess) insertAmharicWords() {
	db, err := sql.Open("mysql", con.conString)
	if err != nil {
		panic(err.Error())
	}
	//make sure to always close the connection
	defer db.Close()

	//select query
	result, err := db.Query(insertAmWordsScript)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

}
