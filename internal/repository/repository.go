package repository

import (
	"database/sql"

	_ "modernc.org/sqlite"
	"github.com/FabianRolfMatthiasNoll/yasl/internal/model"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository() (*Repository, error) {
	db, err := sql.Open("sqlite", "file:yasl.db?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}

	repo := &Repository{DB: db}
	if err := repo.init(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		list_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		category TEXT,
		checked BOOLEAN NOT NULL DEFAULT 0,
		FOREIGN KEY(list_id) REFERENCES lists(id) ON DELETE CASCADE
	);`

	_, err := r.DB.Exec(schema)
	return err
}

func (r *Repository) Close() error {
	return r.DB.Close()
}

func (r *Repository) CreateList(name string) (int64, error) {
	result, err := r.DB.Exec("INSERT INTO lists (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) DeleteList(id int64) error {
	_, err := r.DB.Exec("DELETE FROM lists WHERE id = ?", id)
	return err
}

func (r *Repository) GetLists() ([]model.List, error) {
	rows, err := r.DB.Query("SELECT id, name FROM lists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []model.List
	for rows.Next() {
		var list model.List
		if err := rows.Scan(&list.ID, &list.Name); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (r *Repository) CreateItem(listID int64, name, category string) (int64, error) {
	result, err := r.DB.Exec("INSERT INTO items (list_id, name, category) VALUES (?, ?, ?)", listID, name, category)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) DeleteItem(id int64) error {
	_, err := r.DB.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

func (r *Repository) GetItems(listID int64) ([]model.Item, error) {
	rows, err := r.DB.Query("SELECT id, list_id, name, category, checked FROM items WHERE list_id = ?", listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.ListID, &item.Name, &item.Category, &item.Checked); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) UpdateItemChecked(id int64, checked bool) error {
	_, err := r.DB.Exec("UPDATE items SET checked = ? WHERE id = ?", checked, id)
	return err
}

func (r *Repository) UpdateItemCategory(id int64, category string) error {
	_, err := r.DB.Exec("UPDATE items SET category = ? WHERE id = ?", category, id)
	return err
}

func (r *Repository) UpdateItemName(id int64, name string) error {
	_, err := r.DB.Exec("UPDATE items SET name = ? WHERE id = ?", name, id)
	return err
}

func (r *Repository) ClearCheckedItems(listID int64) error {
	_, err := r.DB.Exec("DELETE FROM items WHERE list_id = ? AND checked = 1", listID)
	return err
}

func (r *Repository) MoveItemToList(itemID, newListID int64) error {
	_, err := r.DB.Exec("UPDATE items SET list_id = ? WHERE id = ?", newListID, itemID)
	return err
}

