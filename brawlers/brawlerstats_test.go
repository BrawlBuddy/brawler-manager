package brawlers

import (
	"testing"
)

func TestGetAllBrawlers(t *testing.T) {
	result := GetAllBrawlers()
	t.Log(result)
}
