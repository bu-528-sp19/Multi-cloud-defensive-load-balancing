FROM golang
WORKDIR '/app'
RUN go get -u github.com/gorilla/mux
RUN go get github.com/lib/pq
COPY . .
EXPOSE 8888
CMD ["go", "run", "*.go"]


func dbLogin() *sql.DB {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("NAME"))
	db, err := sql.Open("postgres", psqlInfo)