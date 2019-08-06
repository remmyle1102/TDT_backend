package models

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// CUSTOM NULL Handling structures

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 sql.NullInt64

// Scan implements the Scanner interface for NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}
	return nil
}

// NullBool is an alias for sql.NullBool data type
type NullBool sql.NullBool

// Scan implements the Scanner interface for NullBool
func (nb *NullBool) Scan(value interface{}) error {
	var b sql.NullBool
	if err := b.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nb = NullBool{b.Bool, false}
	} else {
		*nb = NullBool{b.Bool, true}
	}

	return nil
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 sql.NullFloat64

// Scan implements the Scanner interface for NullFloat64
func (nf *NullFloat64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nf = NullFloat64{f.Float64, false}
	} else {
		*nf = NullFloat64{f.Float64, true}
	}

	return nil
}

// NullString is an alias for sql.NullString data type
type NullString sql.NullString

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for NullBool
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nb.Bool)
	nb.Valid = (err == nil)
	return err
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

type AuditDBInstance struct {
	AuditLevel                []*AuditLevel
	BuiltinCheck              []*BuiltinCheck
	LocalMemAdmin             []*LocalMemAdmin
	SysAdMem                  []*SysAdMem
	OwnerMem                  []*OwnerMem
	LoginMapping              []*LoginMapping
	PolicyUserCheck           []*PolicyUserCheck
	CheckTestService          []*CheckTestService
	PolicySACheck             []*PolicySACheck
	ConfigCheck               []*ConfigCheck
	GuestPerCheck             []*GuestPerCheck
	ServerAuthenticationCheck []*ServerAuthenticationCheck
	SQLInfoCheck              []*SQLInfoCheck
	DBUser                    []*DBUser
	DBPerm                    []*DBPerm
	RoleMembership            []*RoleMembership
	ConnectionInfo            []*ConnectionInfo
	SQLServiceStartup         []*SQLServiceStartup
	LinkedSvrLogin            []*LinkedSvrLogin
	OrphanedUser              []*OrphanedUser
	NoPermLogin               []*NoPermLogin
}

type AuditLevel struct {
	AuditLevel string
}

type BuiltinCheck struct {
	ServerRole string
	LoginName  string
}
type LocalMemAdmin struct {
	AccountName     string
	Type            string
	Privilege       string
	MappedLoginName string
	PermissionPath  string
}

type SysAdMem struct {
	ServerRole string
	MemberName string
	MemberSID  string
}

type OwnerMem struct {
	DatabaseName string
	Role         string
	Member       string
}

type LoginMapping struct {
	LoginName NullString
	DBName    NullString
	UserName  NullString
	AliasName NullString
}

type PolicyUserCheck struct {
	Name                string
	IsPolicyChecked     bool
	IsExpirationChecked bool
}

type CheckTestService struct {
	Name string
}

type PolicySACheck struct {
	Name                string
	ReNamed             string
	IsPolicyChecked     bool
	IsExpirationChecked bool
	IsDisabled          bool
}

type ConfigCheck struct {
	Name        string
	ValueInUse  bool
	Description NullString
}

type GuestPerCheck struct {
	DatabaseName   string
	ClassDesc      string
	PermissionName string
	ObjectName     string
	CheckStatus    string
}

type ServerAuthenticationCheck struct {
	IsIntegratedSecurityOnly int
}

type SQLInfoCheck struct {
	Index          int
	Name           string
	InternalValue  NullInt64
	CharacterValue NullString
}

type DBUser struct {
	UserName      string
	RoleName      string
	LoginName     NullString
	DefDBName     NullString
	DefSchemaName NullString
	UserID        NullString
	SID           NullString
}

type DBPerm struct {
	Owner       string
	Object      string
	Grantee     string
	Grantor     string
	ProtectType string
	Action      string
	Column      string
}

type RoleMembership struct {
	DBRole     string
	MemberName string
	MemberSID  string
}

type ConnectionInfo struct {
	ConnectionID     string
	ConnectTime      string
	NetTransport     string
	NetPacketSize    int
	ClientNetAddress string
}

type SQLServiceStartup struct {
	Servicename     string
	StatupTypeDesc  string
	StatusDesc      string
	ServiceAccount  string
	IsClustered     string
	ClusterNodename NullString
}

type LinkedSvrLogin struct {
	LinkedServer  string
	LocalLogin    NullString
	IsSelfMapping bool
	RemoteLogin   NullBool
}

type OrphanedUser struct {
	DB       string
	Username string
	TypeDesc string
	Type     string
}

type NoPermLogin struct {
	LoginName  string
	TypeDesc   string
	IsDisabled bool
	DBPerm     string
	SvrPerms   string
}
