$target_os = "windows", "linux";
$target_archs = "amd64", "386";

$env:GOPRIVATE="*";
$env:GONOSUMDB="*";

Remove-Item -LiteralPath "build" -Force -Recurse

Foreach ($os IN $target_os) {
    If ($os -eq "windows") {
        $extension = ".exe";
    } Else {
        $extension = "";
    }
    Foreach ($arch IN $target_archs) {
        $env:GOOS = $os;
        $env:GOARCH = $arch;
        go build -v -o "build/fullproxy-$os-$arch$extension" -ldflags="-s -w" -trimpath -buildvcs=false -mod vendor;
    }
}
