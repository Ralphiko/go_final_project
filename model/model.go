package model

const DatePattern string = "20060102"

type Task struct {
  Id      string json:"id"
  Date    string json:"date"
  Title   string json:"title"
  Comment string json:"comment"
  Repeat  string json:"repeat"
}

type TaskIdResponse struct {
  Id int json:"id"
}

type Tasks struct {
  Tasks []Task json:"tasks"
}
type ErrorResponse struct {
  Error string json:"error"
}

type Sign struct {
  Password string json:"password"
}

type AuthToken struct {
  Token string json:"token"
}