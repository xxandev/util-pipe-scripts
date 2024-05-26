function PrintParameters {
    param(
        [string[]]$params
    )
    Write-Host "COUNT PARAMS: $($params.Count)"
    Write-Host "PARAMS:"
    foreach ($param in $params) {
        Write-Host "$param"
    }
}

Write-Host "POWERSHELL"
PrintParameters -params $args