package crons

import (
	"fmt"
	"time"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/robfig/cron"
	"gopkg.in/guregu/null.v4"
)

func Setup() {
	c := cron.New()
	checkStates()
	c.AddFunc("@midnight", func() { checkStates() })
	c.Start()
}

func checkStates() {
	fmt.Println("Running")

	// Fetch all states
	states := make([]models.State, 0)
	if err := database.Gorm.Find(&states).Error; err != nil {
		fmt.Printf("ERROR - States not found | %s", err.Error())
		return
	}

	now := time.Now()

	// For each state check if it should be active
	// If not -> If before "from", set inactive. If after "to", set inactive and null "from" and "to"
	for i, state := range states {
		if state.From.Valid && state.To.Valid {

			if state.From.Time.Before(now) && state.To.Time.After(now) { // now is between start and finish
				states[i].Active = true
			} else if state.From.Time.After(now) { // now is before start
				states[i].Active = false
			} else { // now is after finish
				states[i].Active = false
				states[i].From = null.Time{}
				states[i].To = null.Time{}
			}

			if err := database.Gorm.Save(&states[i]).Error; err != nil {
				fmt.Printf("ERROR - Error while saving state | %s\n\n", err.Error())
				return
			}

			if states[i].Active {
				fmt.Printf("%s state is: Active\n\n", state.Name)
			} else {
				fmt.Printf("%s state is: Inactive\n\n", state.Name)
			}
		}
	}
}
