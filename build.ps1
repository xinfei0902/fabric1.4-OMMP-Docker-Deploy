if ( [String]::IsNullOrEmpty($Env:GOPATH) ) {

    $Env:GOPATH = $PSScriptRoot

}


$Folder = $args[0]

if ( [String]::IsNullOrEmpty($Folder) ) {
    Write-Output 'empty target'
    exit
}

$githash = git rev-parse HEAD
$build = date -u '+%Y-%m-%d_%H:%M:%S'

$commitDate = git log --pretty=format:"%h" -1
$headName = git rev-parse --abbrev-ref HEAD
$gitTagName = git describe --abbrev=0 --tags
$gitBranchName = git symbolic-ref --short -q HEAD

Write-Output $githash $build $commitDate $headName $gitTagName

$BaseFlag = " -X static.BuildDate=$build -X static.BuildVersion=$githash -X static.BuildName=$gitBranchName "


$Flag = $BaseFlag



if ( ! [String]::IsNullOrEmpty($gitTagName) ) {
    $Flag = "$Flag -X static.Version=$gitTagName"
}

$env:GOOS = 'windows'
go build                  -ldflags "$Flag"  $Folder

if ( !$? ) {
    exit
}


$env:GOOS = 'linux'
go build                  -ldflags "$Flag" -o "${Folder}.bin"  $Folder


