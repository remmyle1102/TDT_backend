Get-Service -Name *SQL* | Select Status, Name, DisplayName | ConvertTo-Json -Depth 2 
