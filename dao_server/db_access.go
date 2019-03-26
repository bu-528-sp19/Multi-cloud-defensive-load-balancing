package main

import (
  "database/sql"
  "fmt"
  "os"
  "strconv"
  "math/rand"
  "time"
  _ "github.com/lib/pq"
  "bytes"
  "net/http"
  "encoding/json"
)

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
func db_recover(db *sql.DB, down_time string) (){
  m := s.GetAll()
  for key, query := range m {
    fmt.Println(key, query)
    if key != "AWS_DOWN" && compare_times(key, down_time) > 0 {
      //if entry was committed to log after db went down, execute query
      fmt.Println("execute", query)
      //db.Query(query)
    }
  }
}

func compare_times(down_time string, up_time string) (int) {
  x := compare_strings(down_time[0:4], up_time[0:4])
  if x != 0 {
    return x
  }
  x = compare_strings(down_time[5:7], up_time[5:7])
  if x != 0 {
    return x
  }
  x = compare_strings(down_time[8:10], up_time[8:10])
  if x != 0 {
    return x
  }
  x = compare_strings(down_time[11:13], up_time[11:13])
  if x != 0 {
    return x
  }
  x = compare_strings(down_time[14:16], up_time[14:16])
  if x != 0 {
    return x
  }
  x = compare_strings(down_time[17:19], up_time[17:19])
  return x
}

func compare_strings(s1 string, s2 string) (int){
  i,_ := strconv.Atoi(s1)
  j,_ := strconv.Atoi(s2)
  if i > j {
    return 1
  } else if j > i {
    return -1
  }
  return 0

}
///////////////////////////////////////////////////






type DatabaseInfo struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}

//deleted test select function

func dbLoginread() (*sql.DB) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("NAME"))
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}



//GCP

	port, _ = strconv.Atoi(os.Getenv("PORT"))
	psqlInfoAWS := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOSTAWS"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORDAWS"),
		os.Getenv("NAME"))
	dbAWS, errAWS := sql.Open("postgres", psqlInfoAWS)

	if errAWS != nil {
		panic(errAWS)
	}

	errAWS = dbAWS.Ping()
	if errAWS != nil {
		panic(errAWS)
	}


	if (errAWS == nil) && (err!=nil){
		return dbAWS
	}else if (errAWS!=nil)&&(err==nil){
		return db
	}else {
		rand.Seed(time.Now().UnixNano())
		dice :=rand.Intn(20)
		if (dice%2 == 0){
			return db
		}else{
			return dbAWS
		}

	}

}

func dbLogin() (*sql.DB, *sql.DB) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("NAME"))
	//db, err := sql.Open("postgres", psqlInfo)//err declared and not used
	db, _ := sql.Open("postgres", psqlInfo)
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	psqlInfoAWS := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOSTAWS"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORDAWS"),
		os.Getenv("NAME"))
	//dbAWS, errAWS := sql.Open("postgres", psqlInfoAWS)//errAWs declared and not used
	dbAWS, _ := sql.Open("postgres", psqlInfoAWS)

	return db,dbAWS
}
