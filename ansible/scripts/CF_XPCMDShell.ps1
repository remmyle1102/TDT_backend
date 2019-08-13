[Reflection.Assembly]::LoadWithPartialName("Microsoft.SqlServer.Smo")  | Out-Null;
 $srv = new-object Microsoft.SqlServer.Management.Smo.Server("MSSQLSERVER")
 if ($srv.Configuration.XPCmdShellEnabled -eq $TRUE)
 {
 # Write-Host "xp_cmdshell is enabled in instance" $srv.Name
 Write-Output "XP_CMDSHELL is ENABLED in this instance" $srv.Name >> C:\1.txt
 }
 else
 {
 # Write-Host "XP_CMDSHELL is DISABLED in this instance" $srv.Name
 Write-Output "XP_CMDSHELL is DISABLED in this instance" $srv.Name | ConvertTo-Json -Depth 2
 }
