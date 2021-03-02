package print

import (
	"fmt"
	"heimatcli/src/heimat"

	"github.com/alexeyco/simpletable"
)

// Profile prints the profile and time balances of a given user
func Profile(u *heimat.User, b *heimat.Balances) {
	table := simpletable.New()
	emptyRow := []*simpletable.Cell{{}, {}}

	// cells := make([]*simpletable.Cell, 0)
	nameRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Name:"},
		{Align: simpletable.AlignLeft, Text: u.Name()},
	}
	table.Body.Cells = append(table.Body.Cells, nameRow)

	idRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "ID:"},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", u.ID)},
	}
	table.Body.Cells = append(table.Body.Cells, idRow)

	table.Body.Cells = append(table.Body.Cells, emptyRow)

	workingHoursRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Worked hours:"},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%.2f", b.WorkingHours)},
	}
	table.Body.Cells = append(table.Body.Cells, workingHoursRow)

	plannedHoursRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Planned hours:"},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%.2f", b.PlannedWorkingHours)},
	}
	table.Body.Cells = append(table.Body.Cells, plannedHoursRow)

	plusMinus := ""
	if b.BalanceWorkingHours < 0 {
		plusMinus = "-"
	} else {
		plusMinus = "+"
	}
	balanceRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Balance:"},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s%.2f", plusMinus, b.BalanceWorkingHours)},
	}
	table.Body.Cells = append(table.Body.Cells, balanceRow)

	holidaysRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Vacation days left:"},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%.2f", b.HolidayEntitlement)},
	}
	table.Body.Cells = append(table.Body.Cells, holidaysRow)

	holidaysTookRow := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Vacation days taken:"},
		{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%.2f", b.Holidays)},
	}
	table.Body.Cells = append(table.Body.Cells, holidaysTookRow)

	table.SetStyle(simpletable.StyleCompactLite)
	table.Println()
}

func printUser(u *heimat.User) {
	fmt.Printf("%#v\n", u)
}
