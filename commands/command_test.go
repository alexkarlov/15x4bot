package commands_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexkarlov/15x4bot/commands"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"reflect"
	"testing"
)

func TestCommandAddRehearsalAdmin(t *testing.T) {
	mock, err := store.InitTest()
	rows := sqlmock.NewRows([]string{"id", "name", "address"}).AddRow(1, "name1", "address1").AddRow(2, "name2", "address2")
	mock.ExpectQuery("SELECT id, name, address FROM places").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	u := &store.User{
		Role: store.USER_ROLE_ADMIN,
	}
	cmdText := lang.MARKUP_BUTTON_CREATE_REHEARSAL
	typeCmd := "addRehearsal"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Errorf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	// test first step
	m, err := cmd.NextStep("")
	if err != nil {
		t.Errorf("unexpected error on first step: %s", err)
	}
	if m.Text != lang.ADD_REHEARSAL_WHEN {
		t.Errorf("wrong text of first cmd. expected %s actual %s", m.Text, lang.ADD_REHEARSAL_WHEN)
	}

	// test second step
	m, err = cmd.NextStep("")
	if err != nil {
		t.Errorf("unexpected error on second step: %s", err)
	}
	fmt.Printf("%#v", m.Buttons)
	if m.Text != lang.PLACES_CHOSE_PLACE {
		t.Errorf("wrong text of first cmd. expected %s actual %s", m.Text, lang.PLACES_CHOSE_PLACE)
	}
}

func TestCommandAddRehearsalGuest(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_GUEST,
	}
	cmdText := lang.MARKUP_BUTTON_CREATE_REHEARSAL
	typeCmd := "unknown"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Errorf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
}
