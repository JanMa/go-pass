# go-pass
go-pass is a [pass](https://www.passwordstore.org/) clone written in Go. 

## Usage

```
Usage:
  go-pass [subfolder | command] [flags]
  go-pass [command]

Available Commands:
  completion  Generates completion scripts
  cp          Copies old-path to new-path, optionally forcefully, selectively reencrypting.
  edit        Insert a new password or edit an existing password using $EDITOR.
  find        List passwords that match pass-names
  generate    Generate a new password of pass-length (or 25 if unspecified) with optionally no symbols.
  git         If the password store is a git repository, execute a git command specified by git-command-args.
  grep        Search for password files containing search-string when decrypted.
  help        Help about any command
  init        Initialize new password storage and use gpg-id for encryption.
  insert      Insert new password.
  ls          List passwords.
  mv          Renames or moves old-path to new-path, optionally forcefully, selectively reencrypting.
  otp         Generate OTP code
  rm          Remove existing password or directory, optionally forcefully.
  show        Show existing password and optionally put it on the clipboard.
  version     Show version information

Flags:
  -h, --help   help for go-pass

Use "go-pass [command] --help" for more information about a command.
```

## Building

For a convenient usage, the repository includes a Makefile. 
```
git clone https://gitlab.com/JanMa/go-pass.git
cd go-pass
make
```

## Installing
After checking out the repository, you can install `go-pass` with
```
make install
```
