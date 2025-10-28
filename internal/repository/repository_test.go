package repository

import (
	"database/sql"
	"testing"

	"github.com/FabianRolfMatthiasNoll/yasl/internal/model"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *Repository {
	// Use an in-memory SQLite database for testing
	repo, err := setupTestRepository()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	return repo
}

func setupTestRepository() (*Repository, error) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	repo := &Repository{DB: db}
	if err := repo.init(); err != nil {
		return nil, err
	}

	return repo, nil
}

func TestCreateList(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	testCases := []struct {
		name     string
		listName string
		wantErr  bool
	}{
		{
			name:     "Valid list creation",
			listName: "Shopping List",
			wantErr:  false,
		},
		{
			name:     "Empty list name",
			listName: "",
			wantErr:  false, // SQLite allows empty strings
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := repo.CreateList(tc.listName)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateList() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && id <= 0 {
				t.Errorf("CreateList() got invalid id = %v", id)
			}
		})
	}
}

func TestGetLists(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create some test lists
	listNames := []string{"List 1", "List 2", "List 3"}
	for _, name := range listNames {
		_, err := repo.CreateList(name)
		if err != nil {
			t.Fatalf("Failed to create test list: %v", err)
		}
	}

	lists, err := repo.GetLists()
	if err != nil {
		t.Fatalf("GetLists() error = %v", err)
	}

	if len(lists) != len(listNames) {
		t.Errorf("GetLists() got %d lists, want %d", len(lists), len(listNames))
	}

	for i, list := range lists {
		if list.Name != listNames[i] {
			t.Errorf("GetLists() got list name %v, want %v", list.Name, listNames[i])
		}
	}
}

func TestCreateAndGetItems(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create a test list
	listID, err := repo.CreateList("Test List")
	if err != nil {
		t.Fatalf("Failed to create test list: %v", err)
	}

	// Test creating items
	testItems := []struct {
		name     string
		category string
	}{
		{"Item 1", "Category 1"},
		{"Item 2", "Category 2"},
		{"Item 3", "Category 1"},
	}

	for _, item := range testItems {
		id, err := repo.CreateItem(listID, item.name, item.category)
		if err != nil {
			t.Errorf("CreateItem() error = %v", err)
		}
		if id <= 0 {
			t.Errorf("CreateItem() got invalid id = %v", id)
		}
	}

	// Test getting items
	items, err := repo.GetItems(listID)
	if err != nil {
		t.Fatalf("GetItems() error = %v", err)
	}

	if len(items) != len(testItems) {
		t.Errorf("GetItems() got %d items, want %d", len(items), len(testItems))
	}

	for i, item := range items {
		if item.Name != testItems[i].name {
			t.Errorf("GetItems() got item name %v, want %v", item.Name, testItems[i].name)
		}
		if item.Category != testItems[i].category {
			t.Errorf("GetItems() got item category %v, want %v", item.Category, testItems[i].category)
		}
	}
}

func TestUpdateItem(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create a test list and item
	listID, err := repo.CreateList("Test List")
	if err != nil {
		t.Fatalf("Failed to create test list: %v", err)
	}

	itemID, err := repo.CreateItem(listID, "Test Item", "Test Category")
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Test updating item checked status
	if err := repo.UpdateItemChecked(itemID, true); err != nil {
		t.Errorf("UpdateItemChecked() error = %v", err)
	}

	// Test updating item category
	if err := repo.UpdateItemCategory(itemID, "New Category"); err != nil {
		t.Errorf("UpdateItemCategory() error = %v", err)
	}

	// Test updating item name
	if err := repo.UpdateItemName(itemID, "New Name"); err != nil {
		t.Errorf("UpdateItemName() error = %v", err)
	}

	// Verify updates
	var items []model.Item
	items, err = repo.GetItems(listID)
	if err != nil {
		t.Fatalf("GetItems() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	item := items[0]
	if !item.Checked {
		t.Error("Item checked status was not updated")
	}
	if item.Category != "New Category" {
		t.Errorf("Item category was not updated, got %v", item.Category)
	}
	if item.Name != "New Name" {
		t.Errorf("Item name was not updated, got %v", item.Name)
	}
}

func TestDeleteItem(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create a test list and item
	listID, err := repo.CreateList("Test List")
	if err != nil {
		t.Fatalf("Failed to create test list: %v", err)
	}

	itemID, err := repo.CreateItem(listID, "Test Item", "Test Category")
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Test deleting the item
	if err := repo.DeleteItem(itemID); err != nil {
		t.Errorf("DeleteItem() error = %v", err)
	}

	// Verify item was deleted
	items, err := repo.GetItems(listID)
	if err != nil {
		t.Fatalf("GetItems() error = %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items after deletion, got %d", len(items))
	}
}

func TestClearCheckedItems(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create a test list
	listID, err := repo.CreateList("Test List")
	if err != nil {
		t.Fatalf("Failed to create test list: %v", err)
	}

	// Create some items and mark some as checked
	items := []struct {
		name    string
		checked bool
	}{
		{"Item 1", true},
		{"Item 2", false},
		{"Item 3", true},
	}

	for _, item := range items {
		id, err := repo.CreateItem(listID, item.name, "Category")
		if err != nil {
			t.Fatalf("Failed to create test item: %v", err)
		}
		if item.checked {
			if err := repo.UpdateItemChecked(id, true); err != nil {
				t.Fatalf("Failed to update item checked status: %v", err)
			}
		}
	}

	// Test clearing checked items
	if err := repo.ClearCheckedItems(listID); err != nil {
		t.Errorf("ClearCheckedItems() error = %v", err)
	}

	// Verify only unchecked items remain
	remainingItems, err := repo.GetItems(listID)
	if err != nil {
		t.Fatalf("GetItems() error = %v", err)
	}

	if len(remainingItems) != 1 {
		t.Errorf("Expected 1 unchecked item to remain, got %d", len(remainingItems))
	}

	if remainingItems[0].Name != "Item 2" {
		t.Errorf("Expected remaining item to be 'Item 2', got %v", remainingItems[0].Name)
	}
}

func TestMoveItemToList(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create two test lists
	list1ID, err := repo.CreateList("List 1")
	if err != nil {
		t.Fatalf("Failed to create first test list: %v", err)
	}

	list2ID, err := repo.CreateList("List 2")
	if err != nil {
		t.Fatalf("Failed to create second test list: %v", err)
	}

	// Create an item in the first list
	itemID, err := repo.CreateItem(list1ID, "Test Item", "Test Category")
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Move the item to the second list
	if err := repo.MoveItemToList(itemID, list2ID); err != nil {
		t.Errorf("MoveItemToList() error = %v", err)
	}

	// Verify item was moved
	items1, err := repo.GetItems(list1ID)
	if err != nil {
		t.Fatalf("GetItems() error for list 1: %v", err)
	}
	if len(items1) != 0 {
		t.Errorf("Expected 0 items in list 1, got %d", len(items1))
	}

	items2, err := repo.GetItems(list2ID)
	if err != nil {
		t.Fatalf("GetItems() error for list 2: %v", err)
	}
	if len(items2) != 1 {
		t.Errorf("Expected 1 item in list 2, got %d", len(items2))
	}
}
