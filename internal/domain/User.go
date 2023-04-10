package domain

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

type UserSessions struct {
	SessionData *webauthn.SessionData `json:"-"`
	SessionCred *webauthn.Credential  `json:"-"`
	DisplayName string
	Jwt         string
	Expiration  uint64 `json:"-"`
}

func (session *UserSessions) DeleteAfter(sessions map[string]*UserSessions) {

	if session.Expiration > 0 {
		time.Sleep(time.Second)
		session.Expiration -= 1
		session.DeleteAfter(sessions)
		return
	}

	log.Printf("user delete")

	user := UserModel{}

	user.Username = session.DisplayName

	userModel := user.Get()

	if user.Password == "" && user.Incredentials == "" {
		userModel.Delete()
	}

	delete(sessions, session.DisplayName)
}

type UserModel struct {
	Id            uint      `gorm:"primarykey;autoIncrement;not null"`
	Icon          string    `gorm:"type:varchar(255);"`
	Username      string    `gorm:"type:varchar(255);not null"`
	Email         string    `gorm:"type:varchar(255);"`
	Password      string    `gorm:"type:varchar(255);"`
	Permission    uint64    `gorm:"type:bigint;default:4607"`
	Incredentials string    `gorm:"column:credentials type:text"`
	ValideAccount bool      `gorm:"type:bool; default false"`
	Disable       bool      `gorm:"type:bool; default false"`
	Subscribtion  []Channel `gorm:"many2many:channel_subscription;  onUpdate:CASCADE; onDelete:CASCADE"`
	Role          []Role    `gorm:"many2many:user_roles;  onUpdate:CASCADE; onDelete:CASCADE"`
	webauthn.User `gorm:"-" json:"-"`
	Credentials   []webauthn.Credential `gorm:"-"`
}

func (user *UserModel) TableName() string {
	return "users"
}

func (user *UserModel) SaveCredentials() error {
	// @todo asure that credentials are transform to string
	var publicKeys []string
	for _, v := range user.Credentials {
		b, _ := json.Marshal(v)

		publicKeys = append(publicKeys, string(b))
	}
	user.Incredentials = strings.Join(publicKeys, ";")
	tx := Db.Save(&user)

	return tx.Error
}

func (user *UserModel) ParseCredentials() {
	for _, v := range strings.Split(user.Incredentials, ";") {
		cred := new(webauthn.Credential)
		err := json.Unmarshal([]byte(v), cred)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		user.Credentials = append(user.Credentials, *cred)
	}
}

func (user *UserModel) Create() error {
	tx := Db.Create(user)

	return tx.Error
}

/*
if user fin return true else false
*/
func (user *UserModel) Find() bool {
	tx := Db.Where("username = ?", user.Username).Find(user)
	return tx.RowsAffected != 0
}
func (user *UserModel) Get() *UserModel {

	tx := Db.Where("username = ?", user.Username).Find(user)
	if tx.RowsAffected == 0 {
		return nil
	}
	return user

}
func (user *UserModel) Delete() {
	Db.Delete(user)
}

func (user *UserModel) Update() {
	Db.Save(&user)
}

/* func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
} */

// WebAuthnID returns the user's ID
func (u UserModel) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.Id))
	return buf
}

// WebAuthnName returns the user's username
func (u UserModel) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName returns the user's display name
func (u UserModel) WebAuthnDisplayName() string {
	return u.Username
}

// WebAuthnIcon is not (yet) implemented
func (u UserModel) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *UserModel) AddCredential(cred webauthn.Credential) {
	u.Credentials = append(u.Credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u UserModel) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u UserModel) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
