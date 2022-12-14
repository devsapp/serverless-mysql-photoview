package models_test

import (
	"testing"
	"time"

	"github.com/photoview/photoview/api/dataloader"
	"github.com/photoview/photoview/api/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestUserRegistrationAuthorization(t *testing.T) {
	//db := test_utils.DatabaseTest(t)

	t.Run("Register user", func(t *testing.T) {
		password := "1234"
		user, err := RegisterUser("admin", &password, true)
		if !assert.NoError(t, err) {
			return
		}

		assert.NotNil(t, user)
		assert.EqualValues(t, "admin", user.Username)
		assert.NotNil(t, user.Password)
		assert.NotEqualValues(t, "1234", user.Password) // should be hashed
		assert.True(t, user.Admin)
	})

	t.Run("Authorize user", func(t *testing.T) {
		user, err := AuthorizeUser("admin", "1234")
		if !assert.NoError(t, err) {
			return
		}

		assert.NotNil(t, user)
		assert.EqualValues(t, "admin", user.Username)
	})

	t.Run("Authorize invalid credentials", func(t *testing.T) {
		user, err := AuthorizeUser("invalid_username", "1234")
		assert.ErrorIs(t, err, ErrorInvalidUserCredentials)
		assert.Nil(t, user)

		user, err = AuthorizeUser("admin", "invalid_password")
		assert.ErrorIs(t, err, ErrorInvalidUserCredentials)
		assert.Nil(t, user)
	})
}

func TestAccessToken(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	pass := "<hashed_password>"
	user := User{
		Username: "user1",
		Password: &pass,
		Admin:    false,
	}

	if !assert.NoError(t, db.Save(&user).Error) {
		return
	}

	access_token, err := user.GenerateAccessToken()
	if !assert.NoError(t, err) {
		return
	}

	assert.NotNil(t, access_token)
	assert.Equal(t, user.ID, access_token.UserID)
	assert.NotEmpty(t, access_token.Value)
	assert.True(t, access_token.Expire.After(time.Now()))
}

func TestUserFillAlbums(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user := User{
		Username: "user",
	}

	if !assert.NoError(t, db.Save(&user).Error) {
		return
	}

	err := user.FillAlbums()
	assert.NoError(t, err)
	assert.Empty(t, user.Albums)

	albums := []Album{
		{
			Title: "album1",
			Path:  "/photos/album1",
		},
		{
			Title: "album2",
			Path:  "/photos/album2",
		},
	}

	if !assert.NoError(t, db.Model(&user).Association("Albums").Append(&albums)) {
		return
	}

	user.Albums = make([]Album, 0)

	err = user.FillAlbums()
	assert.NoError(t, err)
	assert.Len(t, user.Albums, 2)

}

func TestUserOwnsAlbum(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user := User{
		Username: "user",
	}

	if !assert.NoError(t, db.Save(&user).Error) {
		return
	}

	albums := []Album{
		{
			Title: "album1",
			Path:  "/photos/album1",
		},
		{
			Title: "album2",
			Path:  "/photos/album2",
		},
	}

	if !assert.NoError(t, db.Model(&user).Association("Albums").Append(&albums)) {
		return
	}

	sub_albums := []Album{
		{
			Title:         "subalbum1",
			Path:          "/photos/album2/subalbum1",
			ParentAlbumID: &albums[1].ID,
		},
		{
			Title:         "another_sub",
			Path:          "/photos/album2/another_sub",
			ParentAlbumID: &albums[1].ID,
		},
		{
			Title:         "subalbum2",
			Path:          "/photos/album1/subalbum2",
			ParentAlbumID: &albums[0].ID,
		},
	}

	if !assert.NoError(t, db.Model(&user).Association("Albums").Append(&sub_albums)) {
		return
	}

	for _, album := range albums {
		owns, err := user.OwnsAlbum( /*db, */ &album)
		assert.NoError(t, err)
		assert.True(t, owns)
	}

	for _, album := range sub_albums {
		owns, err := user.OwnsAlbum( /*db,*/ &album)
		assert.NoError(t, err)
		assert.True(t, owns)
	}

	separate_album := Album{
		Title: "separate_album",
		Path:  "/my_media/album123",
	}

	if !assert.NoError(t, db.Save(&separate_album).Error) {
		return
	}

	owns, err := user.OwnsAlbum( /*db, */ &separate_album)
	assert.NoError(t, err)
	assert.False(t, owns)
}

func TestUserFavoriteMedia(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, err := RegisterUser( /*db, */ "user1", nil, false)
	assert.NoError(t, err)

	rootAlbum := Album{
		Title: "root",
		Path:  "/photos",
	}

	assert.NoError(t, db.Save(&rootAlbum).Error)
	assert.NoError(t, db.Model(&user).Association("Albums").Append(&rootAlbum))

	media := Media{
		Title:   "pic1",
		Path:    "/photos/pic1",
		AlbumID: rootAlbum.ID,
	}

	assert.NoError(t, db.Save(&media).Error)

	// test that it starts out being false
	favourite, err := dataloader.NewUserFavoriteLoader(db).Load(&UserMediaData{
		UserID:  user.ID,
		MediaID: media.ID,
	})

	assert.NoError(t, err)
	assert.False(t, favourite)

	favMedia, err := user.FavoriteMedia( /*db, */ media.ID, true)
	assert.NoError(t, err)
	assert.NotNil(t, favMedia)

	// test that it is now true
	favourite, err = dataloader.NewUserFavoriteLoader(db).Load(&UserMediaData{
		UserID:  user.ID,
		MediaID: media.ID,
	})

	assert.NoError(t, err)
	assert.True(t, favourite)
}
