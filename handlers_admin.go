package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// handleAppStats gets the applications stats
func (s *server) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AppStats, err := s.database.GetAppStats()

		if err != nil {
			http.NotFound(w, r)
			return
		}

		ActiveBattleUserCount := 0
		for _, s := range h.arenas {
			ActiveBattleUserCount = ActiveBattleUserCount + len(s)
		}

		AppStats.ActiveBattleCount = len(h.arenas)
		AppStats.ActiveBattleUserCount = ActiveBattleUserCount

		s.respondWithJSON(w, http.StatusOK, AppStats)
	}
}

// handleGetRegisteredUsers gets a list of registered users
func (s *server) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Users := s.database.GetRegisteredUsers(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Users)
	}
}

// handleUserCreate registers a new authenticated user
func (s *server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		UserName, UserEmail, UserPassword, accountErr := ValidateUserAccount(
			keyVal["warriorName"].(string),
			strings.ToLower(keyVal["warriorEmail"].(string)),
			keyVal["warriorPassword1"].(string),
			keyVal["warriorPassword2"].(string),
		)

		if accountErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUser, VerifyID, err := s.database.CreateUserRegistered(UserName, UserEmail, UserPassword, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendWelcome(UserName, UserEmail, VerifyID)

		s.respondWithJSON(w, http.StatusOK, newUser)
	}
}

// handleAdminUserDelete attempts to delete a users account
func (s *server) handleAdminUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]

		User, UserErr := s.database.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateErr := s.database.DeleteUser(UserID)
		if updateErr != nil {
			log.Println("error attempting to delete user : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendDeleteConfirmation(User.UserName, User.UserEmail)

		return
	}
}

// handleUserPromote handles promoting a user to admin
func (s *server) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		err := s.database.PromoteUser(keyVal["warriorId"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDemote handles demoting a user to registered
func (s *server) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		err := s.database.DemoteUser(keyVal["warriorId"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleCleanBattles handles cleaning up old battles (ADMIN Manaually Triggered)
func (s *server) handleCleanBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_battles_days_old")

		err := s.database.CleanBattles(DaysOld)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleCleanGuests handles cleaning up old guests (ADMIN Manaually Triggered)
func (s *server) handleCleanGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_guests_days_old")

		err := s.database.CleanGuests(DaysOld)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleLowercaseUserEmails handles lowercasing any user emails that have any uppercase letters
func (s *server) handleLowercaseUserEmails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lowercasedUsers, err := s.database.LowercaseUserEmails()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("Lowercased", len(lowercasedUsers), "user emails")
		for _, u := range lowercasedUsers {
			s.email.SendEmailUpdate(u.UserName, u.UserEmail)
		}

		mergedUsers, err := s.database.MergeDuplicateAccounts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("Merged", len(mergedUsers), "user accounts")
		for _, u := range mergedUsers {
			s.email.SendMergedUpdate(u.UserName, u.UserEmail)
		}

		return
	}
}

// handleGetOrganizations gets a list of organizations
func (s *server) handleGetOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Organizations := s.database.OrganizationList(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Organizations)
	}
}

// handleGetTeams gets a list of teams
func (s *server) handleGetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.TeamList(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetAPIKeys gets a list of APIKeys
func (s *server) handleGetAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.GetAPIKeys(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}
