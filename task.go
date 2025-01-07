package main

import (
	"log/slog"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/tikaflow/app-go"
	"gorm.io/gorm"
	"memo/util"
)

var sqlSuccess = &app.Pong{
	Code: 0,
	Msg:  "OK",
}

func switchWindow() {
	win := app.GetUI()
	if win.IsOpen() {
		_ = win.Close()
	} else {
		app.OpenWindow(1200, 800)
	}
}

/*
 * @event: upsert
 */
func upsertNote(ws *websocket.Conn, event string, data any) {
	note := mapNote(data.(map[string]any))
	slog.Info("upsertNote", "event", event, "data", note)
	res := db.Save(note)

	if res.Error != nil {
		slog.Error("upsertNote", "error", res.Error)
		return
	}
	queryNotes(ws, "query", nil)
}

/*
 * @event: deletes
 */
func deleteNotes(ws *websocket.Conn, event string, data any) {
	ids := data.(map[string]any)["ids"].([]any)
	slog.Info("deleteNotes", "event", event, "ids", ids)

	var res *gorm.DB
	if len(ids) == 1 && event == "delete" {
		res = db.Delete(&Note{}, ids[0])
	} else {
		res = db.Delete(&Note{}, ids)
	}

	if res.Error != nil {
		slog.Error("deleteNotes", "error", res.Error)
		return
	}
	queryNotes(ws, "query", nil)
}

/*
 * @event: up-pin
 */
func upPin(ws *websocket.Conn, event string, data any) {
	ids := data.(map[string]any)["ids"].([]any)
	pin := data.(map[string]any)["pin"].(bool)
	slog.Info("upPin", "event", event, "ids", ids)

	res := db.Model(&Note{}).Where("id in (?)", ids).Update("pin", pin)
	if res.Error != nil {
		slog.Error("upPin", "error", res.Error)
		return
	}
	queryNotes(ws, "query", nil)
}

/*
 * @event: up-cate
 */
func upCategory(ws *websocket.Conn, event string, data any) {
	from := data.(map[string]any)["from"].(string)
	to := data.(map[string]any)["to"].(string)
	slog.Info("upCategory", "event", event, "from", from, "to", to)

	res := db.Model(&Note{}).Where("category = ?", from).Update("category", to)
	if res.Error != nil {
		slog.Error("upCategory", "error", res.Error)
		return
	}
	queryNotes(ws, "query", nil)
}

/*
 * @event: clear
 */
func clearNotes(ws *websocket.Conn, event string, data any) {
	slog.Info("clearNotes")
	res := db.Where("1=1").Delete(&Note{})

	if res.Error != nil {
		slog.Error("clearNotes", "error", res.Error)
		return
	}
	queryNotes(ws, "query", nil)
}

/*
 * @event: query
 */
func queryNotes(ws *websocket.Conn, event string, data any) {
	slog.Info("queryNotes")
	notes, err := getNotes()
	if err != nil {
		slog.Error("queryNotes", "error", err)
		return
	}

	sqlSuccess.Event = "re-" + event
	sqlSuccess.Data = util.ConvertToAnySlice(notes)
	_ = ws.WriteJSON(sqlSuccess)

	sqlSuccess.Event = "re-tags"
	tags := getAttrs(notes, func(nt Note) []string { return strings.Split(nt.Tags, "|") })
	sqlSuccess.Data = util.ConvertToAnySlice(tags)
	_ = ws.WriteJSON(sqlSuccess)

	sqlSuccess.Event = "re-cates"
	cates := getAttrs(notes, func(nt Note) []string { return []string{nt.Category} })
	sqlSuccess.Data = util.ConvertToAnySlice(cates)
	_ = ws.WriteJSON(sqlSuccess)
}

func task() {
	_ = app.When("Alt+F9", switchWindow)

	_ = app.On("upsert", upsertNote)
	_ = app.On("touch", upsertNote)
	_ = app.On("query", queryNotes)
	_ = app.On("up-pin", upPin)
	_ = app.On("up-cate", upCategory)
	_ = app.On("deletes", deleteNotes)
	_ = app.On("clear", clearNotes)
	_ = app.On("quit", func(ws *websocket.Conn, event string, data any) { app.Quit() })

	<-app.Done()
}
