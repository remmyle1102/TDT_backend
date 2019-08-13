Get-NetFirewallRule -All | Where-Object {$_.DisplayName -match ".*sql.*" } | Select-Object DisplayName,
@{Name='Protocol';Expression={($PSItem | Get-NetFirewallPortFilter).Protocol}},
@{Name='LocalPort';Expression={($PSItem | Get-NetFirewallPortFilter).LocalPort}},
@{Name='RemotePort';Expression={($PSItem | Get-NetFirewallPortFilter).RemotePort}},
@{Name='RemoteAddress';Expression={($PSItem | Get-NetFirewallAddressFilter).RemoteAddress}},
Enabled, Profile, Direction, Action | ConvertTo-Json -Depth 2 
