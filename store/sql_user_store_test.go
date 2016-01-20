// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"strings"
	"testing"
	"time"
)

func TestUserStoreSave(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.Email = model.NewId()
	u1.Username = model.NewId()
	u1.TeamId = model.NewId()

	if err := (<-store.User().Save(utils.T, &u1)).Err; err != nil {
		t.Fatal("couldn't save user", err)
	}

	if err := (<-store.User().Save(utils.T, &u1)).Err; err == nil {
		t.Fatal("shouldn't be able to update user from save")
	}

	u1.Id = ""
	if err := (<-store.User().Save(utils.T, &u1)).Err; err == nil {
		t.Fatal("should be unique email")
	}

	u1.Email = ""
	if err := (<-store.User().Save(utils.T, &u1)).Err; err == nil {
		t.Fatal("should be unique username")
	}

	u1.Email = strings.Repeat("0123456789", 20)
	u1.Username = ""
	if err := (<-store.User().Save(utils.T, &u1)).Err; err == nil {
		t.Fatal("should be unique username")
	}

	for i := 0; i < 50; i++ {
		u1.Id = ""
		u1.Email = model.NewId()
		u1.Username = model.NewId()
		if err := (<-store.User().Save(utils.T, &u1)).Err; err != nil {
			t.Fatal("couldn't save item", err)
		}
	}

	u1.Id = ""
	u1.Email = model.NewId()
	u1.Username = model.NewId()
	if err := (<-store.User().Save(utils.T, &u1)).Err; err == nil {
		t.Fatal("should be the limit", err)
	}
}

func TestUserStoreUpdate(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	time.Sleep(100 * time.Millisecond)

	if err := (<-store.User().Update(utils.T, &u1, false)).Err; err != nil {
		t.Fatal(err)
	}

	u1.Id = "missing"
	if err := (<-store.User().Update(utils.T, &u1, false)).Err; err == nil {
		t.Fatal("Update should have failed because of missing key")
	}

	u1.Id = model.NewId()
	if err := (<-store.User().Update(utils.T, &u1, false)).Err; err == nil {
		t.Fatal("Update should have faile because id change")
	}
}

func TestUserStoreUpdateLastPingAt(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().UpdateLastPingAt(utils.T, u1.Id, 1234567890)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-store.User().Get(utils.T, u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).LastPingAt != 1234567890 {
			t.Fatal("LastPingAt not updated correctly")
		}
	}

}

func TestUserStoreUpdateLastActivityAt(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().UpdateLastActivityAt(utils.T, u1.Id, 1234567890)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-store.User().Get(utils.T, u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).LastActivityAt != 1234567890 {
			t.Fatal("LastActivityAt not updated correctly")
		}
	}

}

func TestUserStoreUpdateFailedPasswordAttempts(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().UpdateFailedPasswordAttempts(utils.T, u1.Id, 3)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-store.User().Get(utils.T, u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).FailedAttempts != 3 {
			t.Fatal("LastActivityAt not updated correctly")
		}
	}

}

func TestUserStoreUpdateUserAndSessionActivity(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	s1 := model.Session{}
	s1.UserId = u1.Id
	s1.TeamId = u1.TeamId
	Must(store.Session().Save(utils.T, &s1))

	if err := (<-store.User().UpdateUserAndSessionActivity(utils.T, u1.Id, s1.Id, 1234567890)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-store.User().Get(utils.T, u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).LastActivityAt != 1234567890 {
			t.Fatal("LastActivityAt not updated correctly for user")
		}
	}

	if r2 := <-store.Session().Get(utils.T, s1.Id); r2.Err != nil {
		t.Fatal(r2.Err)
	} else {
		if r2.Data.(*model.Session).LastActivityAt != 1234567890 {
			t.Fatal("LastActivityAt not updated correctly for session")
		}
	}

}

func TestUserStoreGet(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if r1 := <-store.User().Get(utils.T, u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).ToJson() != u1.ToJson() {
			t.Fatal("invalid returned user")
		}
	}

	if err := (<-store.User().Get(utils.T, "")).Err; err == nil {
		t.Fatal("Missing id should have failed")
	}
}

func TestUserCount(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if result := <-store.User().GetTotalUsersCount(utils.T); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		count := result.Data.(int64)
		if count <= 0 {
			t.Fatal()
		}
	}
}

func TestActiveUserCount(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if result := <-store.User().GetTotalActiveUsersCount(utils.T); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		count := result.Data.(int64)
		if count <= 0 {
			t.Fatal()
		}
	}
}

func TestUserStoreGetProfiles(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	u2 := model.User{}
	u2.TeamId = u1.TeamId
	u2.Email = model.NewId()
	Must(store.User().Save(utils.T, &u2))

	if r1 := <-store.User().GetProfiles(utils.T, u1.TeamId); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.(map[string]*model.User)
		if len(users) != 2 {
			t.Fatal("invalid returned users")
		}

		if users[u1.Id].Id != u1.Id {
			t.Fatal("invalid returned user")
		}
	}

	if r2 := <-store.User().GetProfiles(utils.T, "123"); r2.Err != nil {
		t.Fatal(r2.Err)
	} else {
		if len(r2.Data.(map[string]*model.User)) != 0 {
			t.Fatal("should have returned empty map")
		}
	}
}

func TestUserStoreGetSystemAdminProfiles(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	u2 := model.User{}
	u2.TeamId = u1.TeamId
	u2.Email = model.NewId()
	Must(store.User().Save(utils.T, &u2))

	if r1 := <-store.User().GetSystemAdminProfiles(utils.T); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.(map[string]*model.User)
		if len(users) <= 0 {
			t.Fatal("invalid returned system admin users")
		}
	}
}

func TestUserStoreGetByEmail(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().GetByEmail(utils.T, u1.TeamId, u1.Email)).Err; err != nil {
		t.Fatal(err)
	}

	if err := (<-store.User().GetByEmail(utils.T, "", "")).Err; err == nil {
		t.Fatal("Should have failed because of missing email")
	}
}

func TestUserStoreGetByAuthData(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	u1.AuthData = "123"
	u1.AuthService = "service"
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().GetByAuth(utils.T, u1.TeamId, u1.AuthData, u1.AuthService)).Err; err != nil {
		t.Fatal(err)
	}

	if err := (<-store.User().GetByAuth(utils.T, "", "", "")).Err; err == nil {
		t.Fatal("Should have failed because of missing auth data")
	}
}

func TestUserStoreGetByUsername(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	u1.Username = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().GetByUsername(utils.T, u1.TeamId, u1.Username)).Err; err != nil {
		t.Fatal(err)
	}

	if err := (<-store.User().GetByUsername(utils.T, "", "")).Err; err == nil {
		t.Fatal("Should have failed because of missing username")
	}
}

func TestUserStoreUpdatePassword(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	hashedPassword := model.HashPassword("newpwd")

	if err := (<-store.User().UpdatePassword(utils.T, u1.Id, hashedPassword)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-store.User().GetByEmail(utils.T, u1.TeamId, u1.Email); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		user := r1.Data.(*model.User)
		if user.Password != hashedPassword {
			t.Fatal("Password was not updated correctly")
		}
	}
}

func TestUserStoreDelete(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	if err := (<-store.User().PermanentDelete(utils.T, u1.Id)).Err; err != nil {
		t.Fatal(err)
	}
}

func TestUserStoreUpdateAuthData(t *testing.T) {
	Setup()

	u1 := model.User{}
	u1.TeamId = model.NewId()
	u1.Email = model.NewId()
	Must(store.User().Save(utils.T, &u1))

	service := "someservice"
	authData := "1"

	if err := (<-store.User().UpdateAuthData(utils.T, u1.Id, service, authData)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-store.User().GetByEmail(utils.T, u1.TeamId, u1.Email); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		user := r1.Data.(*model.User)
		if user.AuthService != service {
			t.Fatal("AuthService was not updated correctly")
		}
		if user.AuthData != authData {
			t.Fatal("AuthData was not updated correctly")
		}
		if user.Password != "" {
			t.Fatal("Password was not cleared properly")
		}
	}
}
