$root = "v3";
$target_os = "windows", "linux";
$target_archs = "amd64", "386";

$env:GOPRIVATE="*";
$env:GONOSUMDB="*";

$old = pwd;

cd $root;

New-Item -Name build -ItemType directory -Force -Path $old;

Foreach ($os IN $target_os) {
    If ($os -eq "windows") {
        $extension = ".exe";
    } Else {
        $extension = "";
    }
    Foreach ($arch IN $target_archs) {
        $env:GOOS = $os;
        $env:GOARCH = $arch;
        go1.18beta2.exe build -v -o "$old/build/fullproxy-$os-$arch$extension" -ldflags="-s -w" -trimpath -buildvcs=false -mod vendor "cmd/fullproxy/main.go";
        go1.18beta2.exe build -v -o "$old/build/fullproxy-users-$os-$arch$extension" -ldflags="-s -w" -trimpath -buildvcs=false -mod vendor "cmd/fullproxy-users/main.go";
    }
}

cd $old;