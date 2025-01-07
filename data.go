package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

type Note struct {
	ID        uint      `json:"id"`
	Category  string    `json:"category"`
	Pin       bool      `json:"pin"`
	Tags      string    `json:"tags"`
	Lang      string    `json:"lang"`
	Comment   string    `json:"comment"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"ctime"`
	UpdatedAt time.Time `json:"mtime"`
}

func (this *Note) String() string {
	return fmt.Sprintf("Note(id=%d, category=%s, tags=%s, lang=%s, comment=%s, content=%s, ctime=%s, mtime=%s)",
		this.ID, this.Category, this.Tags, this.Lang, this.Comment, this.Content, this.CreatedAt, this.UpdatedAt)
}

func initDB() {
	exe, _ := os.Executable()
	exeRoot := filepath.Dir(exe)
	dbFile := filepath.Join(exeRoot, "memo.db")

	if db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}); err != nil {
		panic("数据库打开失败")
	}

	_ = db.AutoMigrate(&Note{})
}

func getNotes() ([]Note, error) {
	var notes []Note
	return notes, db.Find(&notes).Error
}

func getAttrs(notes []Note, get func(Note) []string) []string {
	tags := make([]string, 0)
	mapping := make(map[string]bool)

	for _, note := range notes {
		for _, tag := range get(note) {
			if _, ok := mapping[tag]; !ok {
				mapping[tag] = true
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

func mapNote(jsn map[string]any) *Note {
	note := &Note{}

	if id, ok := jsn["id"].(float64); ok {
		note.ID = uint(id)
	}
	if category, ok := jsn["category"].(string); ok {
		note.Category = category
	}
	if pin, ok := jsn["pin"].(bool); ok {
		note.Pin = pin
	}
	if tags, ok := jsn["tags"].(string); ok {
		note.Tags = tags
	}
	if lang, ok := jsn["lang"].(string); ok {
		note.Lang = lang
	}
	if comment, ok := jsn["comment"].(string); ok {
		note.Comment = comment
	}
	if content, ok := jsn["content"].(string); ok {
		note.Content = content
	}
	if ctime, ok := jsn["ctime"].(string); ok {
		note.CreatedAt, _ = time.Parse(time.RFC3339, ctime)
	}
	if mtime, ok := jsn["mtime"].(string); ok {
		note.UpdatedAt, _ = time.Parse(time.RFC3339, mtime)
	}

	return note
}
