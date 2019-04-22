package main

import (
  "database/sql"
  "fmt"
  _ "github.com/lib/pq"
  "strconv"
  "encoding/json"
  "net/http"
  "bytes"
  "os"
  "time"
  "math/rand"
)

//lib/pq is the driver for postgres
/*The last import, _ "github.com/lib/pq", might look funny at first, but the
short version is that we are importing the package so that it can register its
drivers with the database/sql package, and we use the _ identifier to tell Go
that we still want this included even though we will never directly reference
the package in our code.*/

type DatabaseInfo struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}


///////////////////////////////////////////////////
//Database Forwarding

type DB_STATE struct {
	Name   string `json:"Name"`
	State  string `json:"State"`
}

//used by the leader to change the state of a db in the raft log
func DB_State_Change(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Println("inside DB_State_Change")
  var db_state_obj DB_STATE
  _ = json.NewDecoder(req.Body).Decode(&db_state_obj)
  s.Set(db_state_obj.Name, db_state_obj.State)
  json.NewEncoder(w).Encode(db_state_obj)
}

//used by non-leader to notify leader of db state change
func notify_leader(db_name string, val string) () {
  fmt.Println("need to notify leader that " + db_name + " is down")
  //local server
  //leaderIP := "localhost:8888"

  //cloud server
  leaderIP := s.GetLeaderAddress() + ":8888"

  url := leaderIP + "/db_state_change"
  db_state_obj := DB_STATE{Name: db_name, State: val}
  jsonStr, _ := json.Marshal(db_state_obj)
  req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  client := &http.Client{}
  resp, _ := client.Do(req)
  defer resp.Body.Close()
}


//used by any node for database forwarding
func db_recover(db *sql.DB, down_time string, db_name string) (){
  m := s.GetAll()
  s.Locklog()
  defer s.Unlocklog()
  fmt.Println(db_name, "down time =", down_time)
  for key, query := range m {
    fmt.Println(key, query)
    if key != "AWS_DOWN" && key != "GCP_DOWN" && compare_times(key, down_time) > 0 {
      //if entry was committed to log after db went down, execute query
      fmt.Println("EXECUTE on ", db_name)
      _, _ = db.Query(query)
    }
  }
}

func compare_times(log_time string, down_time string) (int) {
  x, _ := strconv.ParseInt(log_time, 10, 64)
  y, _ := strconv.ParseInt(down_time, 10, 64)
  if x >= y {
    return 1
  }
  return 0
}


func dbForwarding(db *sql.DB, db_name string) (error){

    //local server
    //isLeader := true

    //cloud server
    isLeader := s.IsLeader()

	if len(s.GetAll()) == 0 {
		s.Set("AWS_DOWN", "0")
		s.Set("GCP_DOWN", "0")
	}

	err := db.Ping()
	if err != nil {
		aws_status, _ := s.Get(db_name)
		if aws_status == "0" {
		        cur_time := strconv.FormatInt(time.Now().Unix(), 10)
			if isLeader != true {
				notify_leader(db_name, cur_time)
			} else{
				s.Set(db_name, cur_time)
			}
		}
	} else {
		down_time, _ := s.Get(db_name)
		if down_time != "0" {
			if isLeader != true {
				notify_leader("db_name", "0")
			} else {
				s.Set(db_name, "0")
			}
				db_recover(db, down_time, db_name)
		}
        //db_recover(db, down_time) //for debugging
	}
    return err
}
///////////////////////////////////////////////////

func gcpLogin() (*sql.DB, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOSTGCP"),
		port,
		os.Getenv("USERGCP"),
		os.Getenv("PASSWORDGCP"),
		os.Getenv("NAMEGCP"))
	db, err := sql.Open("postgres", psqlInfo)

//	if err != nil {
//		return db, err
//	}

	/*err = db.Ping()
	if err != nil {
		return db, err
	}*/
	return db, err
}

func awsLogin() (*sql.DB, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfoAWS := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	  os.Getenv("HOSTAWS"),
		port,
		os.Getenv("USERAWS"),
		os.Getenv("PASSWORDAWS"),
		os.Getenv("NAMEAWS"))
	dbAWS, errAWS := sql.Open("postgres", psqlInfoAWS)

//	if errAWS != nil {
//		return dbAWS, errAWS
//	}

	/*errAWS = dbAWS.Ping()
	if errAWS != nil {
		return dbAWS, errAWS
	}*/

	return dbAWS, errAWS
}


func dbLoginread() (*sql.DB, error) {
   rand.Seed(time.Now().UTC().UnixNano())
   rand := rand.Intn(2)
   if (rand == 0) {
	db, _ := awsLogin()
	err := dbForwarding(db, "AWS_DOWN")
	if err == nil {
		fmt.Println("reading from AWS")
	} else {
		fmt.Println("error reading from AWS")
	}
	return db, err
  } else {
	//uncomment for gcp
	db, _ := gcpLogin()
	err := dbForwarding(db, "GCP_DOWN")
	if err == nil {
		fmt.Println("reading from GCP")
	} else {
		fmt.Println("error reading from GCP")
	}
	return db, err
  }
}

func dbLogin() ([]*sql.DB) {

  var working_dbs []*sql.DB

  dbAWS, err := awsLogin()
//  err = dbForwarding(dbAWS, "AWS_DOWN")
  if err == nil {
    working_dbs = append(working_dbs, dbAWS)
  }

  //uncomment for gcp
  dbGCP, _ := gcpLogin()
  //err = dbForwarding(dbGCP, "GCP_DOWN")
  if err == nil {
    working_dbs = append(working_dbs, dbGCP)
  }

	return working_dbs
}
