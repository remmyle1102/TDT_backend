package db

import (
	"TDT_backend/models"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

// AuditTaskDBInstance execute audit DB Instance query
func AuditTaskDBInstance(dbConn *sql.DB) (*models.AuditDBInstance, error) {

	// Check SQL Server Audit level
	// Field: AuditLevel
	auditLevelS := make([]*models.AuditLevel, 0)
	query := `DECLARE @AuditLevel int
	EXEC master.dbo.xp_instance_regread N'HKEY_LOCAL_MACHINE', 
   	N'Software\Microsoft\MSSQLServer\MSSQLServer', 
   	N'AuditLevel', @AuditLevel OUTPUT
	SELECT CASE WHEN @AuditLevel = 0 THEN 'None'
   	WHEN @AuditLevel = 1 THEN 'Successful logins only'
   	WHEN @AuditLevel = 2 THEN 'Failed logins only'
   	WHEN @AuditLevel = 3 THEN 'Both failed and successful logins' 
   	END AS [AuditLevel] `
	rows, err := dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		auditLevel := new(models.AuditLevel)
		err := rows.Scan(&auditLevel.AuditLevel)
		if err != nil {
			logrus.Error(err)
		}
		auditLevelS = append(auditLevelS, auditLevel)
	}
	// Check that Builtin\Administrators group removed from sysadmins role
	// Field: SrvRole, LoginName
	builtinCheckS := make([]*models.BuiltinCheck, 0)
	query = `SELECT r.name  as SrvRole, u.name  as LoginName  
	FROM sys.server_role_members m JOIN
  	sys.server_principals r ON m.role_principal_id = r.principal_id  JOIN
  	sys.server_principals u ON m.member_principal_id = u.principal_id 
	WHERE u.name = 'BUILTIN\Administrators'`
	rows, err = dbConn.Query(query)
	// Return no rows means OK
	if err != nil && err != sql.ErrNoRows {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		builtinCheck := new(models.BuiltinCheck)
		err := rows.Scan(&builtinCheck.ServerRole, &builtinCheck.LoginName)
		if err != nil {
			logrus.Error(err)
		}
		builtinCheckS = append(builtinCheckS, builtinCheck)
	}
	if len(builtinCheckS) == 0 {
		builtinCheck := models.BuiltinCheck{LoginName: " ", ServerRole: " "}
		builtinCheckS = append(builtinCheckS, &builtinCheck)
	}

	// Find members of the "Local Administrators" group on SQL Server
	// get results from the extended procedure below only if the BUILTIN\Administrators group exists as login on SQL Server.
	// field: account name, type, privilege, mapped login name, permission path
	localMemAdminS := make([]*models.LocalMemAdmin, 0)
	query = `EXEC master.sys.xp_logininfo 'BUILTIN\Administrators','members'`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		localMemAdmin := new(models.LocalMemAdmin)
		err := rows.Scan(&localMemAdmin.AccountName, &localMemAdmin.Type, &localMemAdmin.Privilege, &localMemAdmin.MappedLoginName, &localMemAdmin.PermissionPath)
		if err != nil {
			logrus.Error(err)
		}
		localMemAdminS = append(localMemAdminS, localMemAdmin)
	}
	if len(localMemAdminS) == 0 {
		localMemAdmin := models.LocalMemAdmin{" ", " ", " ", " ", " "}
		localMemAdminS = append(localMemAdminS, &localMemAdmin)
	}

	// -- Find Sysadmins server role's members (and other server level roles)
	// -- show all logins and what server level roles each login has been assigned.
	// -- field: ServerRole, MemberName, MemberSID
	sysAdMemS := make([]*models.SysAdMem, 0)
	query = `EXEC master.sys.sp_helpsrvrolemember`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		sysAdMem := new(models.SysAdMem)
		err := rows.Scan(&sysAdMem.ServerRole, &sysAdMem.MemberName, &sysAdMem.MemberSID)
		if err != nil {
			logrus.Error(err)
		}
		sysAdMemS = append(sysAdMemS, sysAdMem)
	}

	// provide all members of the db_owner database role in all databases in an instance
	// field: database_name, role, member
	ownerMemS := make([]*models.OwnerMem, 0)
	query = `exec sp_msForEachDb ' use [?] 
	select db_name() as [database_name], r.[name] as [role], p.[name] as [member] from  
    sys.database_role_members m 
	join 
    sys.database_principals r on m.role_principal_id = r.principal_id 
	join 
    sys.database_principals p on m.member_principal_id = p.principal_id 
	where 
    r.name = ''db_owner'''`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		ownerMem := new(models.OwnerMem)
		err := rows.Scan(&ownerMem.DatabaseName, &ownerMem.Role, &ownerMem.Member)
		if err != nil {
			logrus.Error(err)
		}
		ownerMemS = append(ownerMemS, ownerMem)
	}

	// -- all users & logins mapping
	// -- field: LoginName, DBName, UserName, AliasName
	loginMappingS := make([]*models.LoginMapping, 0)
	query = `exec sp_msloginmappings`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		loginMapping := new(models.LoginMapping)
		err := rows.Scan(&loginMapping.LoginName, &loginMapping.DBName, &loginMapping.UserName, &loginMapping.AliasName)
		if err != nil {
			logrus.Error(err)
		}
		loginMappingS = append(loginMappingS, loginMapping)
	}

	// 	-- check whether the password policy is turn on or off
	// -- field: name, is_policy_checked, is_expiration_checked
	policyUserCheckS := make([]*models.PolicyUserCheck, 0)
	query = `SELECT name, is_policy_checked, is_expiration_checked FROM sys.sql_logins 
 	WHERE  is_policy_checked=0 OR is_expiration_checked = 0`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		policyUserCheck := new(models.PolicyUserCheck)
		err := rows.Scan(&policyUserCheck.Name, &policyUserCheck.IsPolicyChecked, &policyUserCheck.IsExpirationChecked)
		if err != nil {
			logrus.Error(err)
		}
		policyUserCheckS = append(policyUserCheckS, policyUserCheck)
	}

	// 	-- Check that Production and Test databases are segregated (on different SQL Servers)
	// -- This will look for the value of "Test" or "Dev" in all your database names.
	// -- field: name
	checkTestServiceS := make([]*models.CheckTestService, 0)
	query = `SELECT name FROM master.sys.databases 
 	WHERE name LIKE '%Test%' OR name LIKE '%Dev%'`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		checkTestService := new(models.CheckTestService)
		err := rows.Scan(&checkTestService.Name)
		if err != nil {
			logrus.Error(err)
		}
		checkTestServiceS = append(checkTestServiceS, checkTestService)
	}
	if len(checkTestServiceS) == 0 {
		checkTestService := models.CheckTestService{" "}
		checkTestServiceS = append(checkTestServiceS, &checkTestService)
	}

	// -- check whether the sa password exists and if it does if the password policy is turned on for this login
	// -- field: name, Renamed, is_policy_checked, is_expiration_checked, is_disable
	policySACheckS := make([]*models.PolicySACheck, 0)
	query = `SELECT l.name, CASE WHEN l.name = 'sa' THEN 'NO' ELSE 'YES' END as Renamed,
  	s.is_policy_checked, s.is_expiration_checked, l.is_disabled
	FROM sys.server_principals AS l
 	LEFT OUTER JOIN sys.sql_logins AS s ON s.principal_id = l.principal_id
	WHERE l.sid = 0x01`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		policySACheck := new(models.PolicySACheck)
		err := rows.Scan(&policySACheck.Name, &policySACheck.ReNamed, &policySACheck.IsPolicyChecked, &policySACheck.IsExpirationChecked, &policySACheck.IsDisabled)
		if err != nil {
			logrus.Error(err)
		}
		policySACheckS = append(policySACheckS, policySACheck)
	}

	// 	-- This will check different server configuration settings such as: allow updates, cross db ownership chaining, clr enabled, SQL Mail XPs, Database Mail XPs, xp_cmdshell and Ad Hoc Distributed Queries
	// -- field: name, value_in_use, description
	configCheckS := make([]*models.ConfigCheck, 0)
	query = `SELECT name, value_in_use, description FROM sys.configurations
	WHERE configuration_id IN (102, 117,400, 1547, 1562, 16385, 16386, 16388, 16390, 16391, 16393)`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		configCheck := new(models.ConfigCheck)
		err := rows.Scan(&configCheck.Name, &configCheck.ValueInUse, &configCheck.Description)
		if err != nil {
			fmt.Println("hello")
			logrus.Error(err)
		}
		configCheckS = append(configCheckS, configCheck)
	}

	// 	-- list what permission the guest user has
	// -- Guest user by default has CONNECT permissions to the master, msdb and tempdb databases. Any other permissions will be returned by this query as potential problem
	// -- field: DatabaseName, class_desc, permission_name, ObjectName, CheckStatus
	guestPerCheckS := make([]*models.GuestPerCheck, 0)
	query = `SET NOCOUNT ON
	CREATE TABLE #guest_perms 
 	( db SYSNAME, class_desc SYSNAME, 
  	permission_name SYSNAME, ObjectName SYSNAME NULL)
	EXEC master.sys.sp_MSforeachdb
	'INSERT INTO #guest_perms
 	SELECT ''?'' as DBName, p.class_desc, p.permission_name, 
   	OBJECT_NAME (major_id, DB_ID(''?'')) as ObjectName
 	FROM [?].sys.database_permissions p JOIN [?].sys.database_principals l
  	ON p.grantee_principal_id= l.principal_id 
 	WHERE l.name = ''guest'' AND p.[state] = ''G'''
	SELECT db AS DatabaseName, class_desc, permission_name, 
 	CASE WHEN class_desc = 'DATABASE' THEN db ELSE ObjectName END as ObjectName, 
 	CASE WHEN DB_ID(db) IN (1, 2, 4) AND permission_name = 'CONNECT' THEN 'Default' 
  	ELSE 'Potential Problem!' END as CheckStatus
	FROM #guest_perms
	DROP TABLE #guest_perms`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		guestPerCheck := new(models.GuestPerCheck)
		err := rows.Scan(&guestPerCheck.DatabaseName, &guestPerCheck.ClassDesc, &guestPerCheck.PermissionName, &guestPerCheck.ObjectName, &guestPerCheck.CheckStatus)
		if err != nil {
			logrus.Error(err)
		}
		guestPerCheckS = append(guestPerCheckS, guestPerCheck)
	}

	// 	-- SQL Server Authentication mode
	// -- If this returns 0 the server uses both Windows and SQL Server security.
	// If the value is 1 it is only setup for Windows Authentication.
	serverAuthenticationCheckS := make([]*models.ServerAuthenticationCheck, 0)
	query = `SELECT SERVERPROPERTY ('IsIntegratedSecurityOnly') as IsIntegratedSecurityOnly`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		serverAuthenticationCheck := new(models.ServerAuthenticationCheck)
		err := rows.Scan(&serverAuthenticationCheck.IsIntegratedSecurityOnly)
		if err != nil {
			logrus.Error(err)
		}
		serverAuthenticationCheckS = append(serverAuthenticationCheckS, serverAuthenticationCheck)
	}

	// 	-- SQL Server version
	// -- field: Name, Character_Value
	sqlInfoCheckS := make([]*models.SQLInfoCheck, 0)
	query = `EXEC master.sys.xp_msver`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		sqlInfoCheck := new(models.SQLInfoCheck)
		err := rows.Scan(&sqlInfoCheck.Index, &sqlInfoCheck.Name, &sqlInfoCheck.InternalValue, &sqlInfoCheck.CharacterValue)
		if err != nil {
			logrus.Error(err)
		}
		sqlInfoCheckS = append(sqlInfoCheckS, sqlInfoCheck)
	}

	// -- list of the users
	dbUserS := make([]*models.DBUser, 0)
	query = `EXEC sys.sp_helpuser`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		dbUser := new(models.DBUser)
		err := rows.Scan(&dbUser.UserName, &dbUser.RoleName, &dbUser.LoginName, &dbUser.DefDBName, &dbUser.DefSchemaName, &dbUser.UserID, &dbUser.SID)
		if err != nil {
			logrus.Error(err)
		}
		dbUserS = append(dbUserS, dbUser)
	}

	// -- database permissions
	dbPermS := make([]*models.DBPerm, 0)
	query = `EXEC sys.sp_helprotect`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		dbPerm := new(models.DBPerm)
		err := rows.Scan(&dbPerm.Owner, &dbPerm.Object, &dbPerm.Grantee, &dbPerm.Grantor, &dbPerm.ProtectType, &dbPerm.Action, &dbPerm.Column)
		if err != nil {
			logrus.Error(err)
		}
		dbPermS = append(dbPermS, dbPerm)
	}

	// -- roles membership
	roleMembershipS := make([]*models.RoleMembership, 0)
	query = `EXEC sys.sp_helprolemember`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		roleMembership := new(models.RoleMembership)
		err := rows.Scan(&roleMembership.DBRole, &roleMembership.MemberName, &roleMembership.MemberSID)
		if err != nil {
			logrus.Error(err)
		}
		roleMembershipS = append(roleMembershipS, roleMembership)
	}

	// 	--gives the details about existing connections like when the connection was established and what protocol is being used by that particular connection.
	// -- field: connection_id, connect_time, net_transport, net_packet_size, client_net_address
	connectionInfoS := make([]*models.ConnectionInfo, 0)
	query = `SELECT connection_id, connect_time, net_transport, net_packet_size, client_net_address
	FROM sys.dm_exec_connections`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		connectionInfo := new(models.ConnectionInfo)
		err := rows.Scan(&connectionInfo.ConnectionID, &connectionInfo.ConnectTime, &connectionInfo.NetTransport, &connectionInfo.NetPacketSize, &connectionInfo.ClientNetAddress)
		if err != nil {
			logrus.Error(err)
		}
		connectionInfoS = append(connectionInfoS, connectionInfo)
	}

	// 	-- SQL Server Services Startup mode
	// -- field: servicename, statup_type_desc, status_desc, service_account, is_clustered, cluster_nodename
	sqlServiceStartupS := make([]*models.SQLServiceStartup, 0)
	query = `SELECT servicename, startup_type_desc, status_desc, service_account, is_clustered, cluster_nodename FROM sys.dm_server_services`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		sqlServiceStartup := new(models.SQLServiceStartup)
		err := rows.Scan(&sqlServiceStartup.Servicename, &sqlServiceStartup.StatupTypeDesc, &sqlServiceStartup.StatusDesc, &sqlServiceStartup.ServiceAccount, &sqlServiceStartup.IsClustered, &sqlServiceStartup.ClusterNodename)
		if err != nil {
			logrus.Error(err)
		}
		sqlServiceStartupS = append(sqlServiceStartupS, sqlServiceStartup)
	}

	// -- linked server logins
	linkedSvrLoginS := make([]*models.LinkedSvrLogin, 0)
	query = `EXEC master.sys.sp_helplinkedsrvlogin `
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		linkedSvrLogin := new(models.LinkedSvrLogin)
		err := rows.Scan(&linkedSvrLogin.LinkedServer, &linkedSvrLogin.LocalLogin, &linkedSvrLogin.IsSelfMapping, &linkedSvrLogin.RemoteLogin)
		if err != nil {
			logrus.Error(err)
		}
		linkedSvrLoginS = append(linkedSvrLoginS, linkedSvrLogin)
	}

	// 	-- Find orphaned users in all of the databases (no logins exist for the database users)
	// -- field: db, username, type_desc
	orphanedUserS := make([]*models.OrphanedUser, 0)
	query = `SET NOCOUNT ON
	CREATE TABLE #orph_users (db SYSNAME, username SYSNAME, 
    type_desc VARCHAR(30),type VARCHAR(30))
	EXEC master.sys.sp_msforeachdb  
	'INSERT INTO #orph_users
 	SELECT ''?'', u.name , u.type_desc, u.type
 	FROM  [?].sys.database_principals u 
  	LEFT JOIN  [?].sys.server_principals l ON u.sid = l.sid 
 	WHERE l.sid IS NULL 
  	AND u.type NOT IN (''A'', ''R'', ''C'') -- not a db./app. role or certificate
  	AND u.principal_id > 4 -- not dbo, guest or INFORMATION_SCHEMA
  	AND u.name NOT LIKE ''%DataCollector%'' 
  	AND u.name NOT LIKE ''mdw%'' -- not internal users in msdb or MDW databases'
 	SELECT * FROM #orph_users
 	DROP TABLE #orph_users`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		orphanedUser := new(models.OrphanedUser)
		err := rows.Scan(&orphanedUser.DB, &orphanedUser.Username, &orphanedUser.TypeDesc, &orphanedUser.Type)
		if err != nil {
			logrus.Error(err)
		}
		orphanedUserS = append(orphanedUserS, orphanedUser)
	}

	noPermLoginS := make([]*models.NoPermLogin, 0)
	query = `SET NOCOUNT ON
	CREATE TABLE #all_users (db VARCHAR(70), sid VARBINARY(85), stat VARCHAR(50))
	EXEC master.sys.sp_msforeachdb
	'INSERT INTO #all_users  
	SELECT ''?'', CONVERT(varbinary(85), sid) , 
	CASE WHEN  r.role_principal_id IS NULL AND p.major_id IS NULL 
	THEN ''no_db_permissions''  ELSE ''db_user'' END
	FROM [?].sys.database_principals u LEFT JOIN [?].sys.database_permissions p 
	ON u.principal_id = p.grantee_principal_id  
	AND p.permission_name <> ''CONNECT''
	LEFT JOIN [?].sys.database_role_members r 
	ON u.principal_id = r.member_principal_id
	WHERE u.SID IS NOT NULL AND u.type_desc <> ''DATABASE_ROLE'''
	IF EXISTS 
	(SELECT l.name FROM sys.server_principals l LEFT JOIN sys.server_permissions p 
	ON l.principal_id = p.grantee_principal_id  
	AND p.permission_name <> 'CONNECT SQL'
	LEFT JOIN sys.server_role_members r 
	ON l.principal_id = r.member_principal_id
	LEFT JOIN #all_users u 
	ON l.sid= u.sid
	WHERE r.role_principal_id IS NULL  AND l.type_desc <> 'SERVER_ROLE' 
	AND p.major_id IS NULL
	)
	BEGIN
	SELECT DISTINCT l.name LoginName, l.type_desc, l.is_disabled, 
	ISNULL(u.stat + ', but is user in ' + u.db  +' DB', 'no_db_users') db_perms, 
	CASE WHEN p.major_id IS NULL AND r.role_principal_id IS NULL  
	THEN 'no_srv_permissions' ELSE 'na' END srv_perms 
	FROM sys.server_principals l LEFT JOIN sys.server_permissions p 
	ON l.principal_id = p.grantee_principal_id  
	AND p.permission_name <> 'CONNECT SQL'
	LEFT JOIN sys.server_role_members r 
	ON l.principal_id = r.member_principal_id
	LEFT JOIN #all_users u 
	ON l.sid= u.sid
	WHERE  l.type_desc <> 'SERVER_ROLE' 
	AND ((u.db  IS NULL  AND p.major_id IS NULL 
		AND r.role_principal_id IS NULL )
	OR (u.stat = 'no_db_permissions' AND p.major_id IS NULL 
		AND r.role_principal_id IS NULL)) 
	ORDER BY 1, 4
	END
	DROP TABLE #all_users`
	rows, err = dbConn.Query(query)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		noPermLogin := new(models.NoPermLogin)
		err := rows.Scan(&noPermLogin.LoginName, &noPermLogin.TypeDesc, &noPermLogin.IsDisabled, &noPermLogin.DBPerm, &noPermLogin.SvrPerms)
		if err != nil {
			logrus.Error(err)
		}
		noPermLoginS = append(noPermLoginS, noPermLogin)
	}
	if len(noPermLoginS) == 0 {
		noPermLogin := models.NoPermLogin{" ", " ", false, " ", " "}
		noPermLoginS = append(noPermLoginS, &noPermLogin)
	}

	result := &models.AuditDBInstance{auditLevelS, builtinCheckS, localMemAdminS, sysAdMemS, ownerMemS, loginMappingS, policyUserCheckS, checkTestServiceS, policySACheckS, configCheckS, guestPerCheckS, serverAuthenticationCheckS, sqlInfoCheckS, dbUserS, dbPermS, roleMembershipS, connectionInfoS, sqlServiceStartupS, linkedSvrLoginS, orphanedUserS, noPermLoginS}

	return result, err
}
