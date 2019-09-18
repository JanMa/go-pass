package entry

import (
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

const (
	TestEmail                    = "test@example.com"
	TestPublicKeyFingerPrint     = "0484D2AB195D9FAD"
	TestPublicKeyFingerPrintLong = "EFF6DDA2A3D5837E977220190484D2AB195D9FAD"
	TestPublicSubKeyFingerPrint  = "F2089481896D82B9"
	TestPublicKey                = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQENBF1r9zMBCAC96l46ncxNYkeXz//le7qOMCh2reSFt7pgSd4pmjA0sR61uqIu
Z+h/Q6Ssk4lMeisgj5tV8EvU8OyKiuXCQZXmKJmswCtCqmq4a3pcqqcbwmmm7AqV
uB9be4NJlJy2YK5uSx37VAC/ohhsR1g8EhMnB9RiyVIkESTPbLZNj96SVApJYzJy
r8K5EeOKA4zFEoxfwwg8JvdbON3NGVbH5wpU5sks/vm3kw/gGL8MHNRvaTVUUJZg
T61X2HWkiexNDeKNgKct3KwRLnmcsIe6/eRQGratrYbqyyiRlfacaRW1E6btA4Rw
+4YkxvzhOKgqADrY3jL+GFwhujUC28ovJMpFABEBAAG0HlRlc3QgUGVyc29uIDx0
ZXN0QGV4YW1wbGUuY29tPokBVAQTAQgAPhYhBO/23aKj1YN+l3IgGQSE0qsZXZ+t
BQJda/czAhsDBQkDwmcABQsJCAcCBhUKCQgLAgQWAgMBAh4BAheAAAoJEASE0qsZ
XZ+th90H/RIoVT+qz/fHwXbtzpZYP6i8TqS5tQQ+3n7Xxotpx2+C3kfc/jTLNrlt
nIfOKQTbzwo4z+b1g3KhP/Cylw5LpNXggNHf9kKUL28idw+3mpLz5GnupIuvLSK+
Gh2vZWY6OuJOp2ajEOBT16rmIFcVgx9PFm6FDeBm48lCnnndmncuoTrNpagAwGRG
GR1TVrIjDitZi8g83d0VHkskSbhj1CbPa+bzlQSN9VJ0LTNlXGJxKPFgflGxghKX
Vn74v2+j7ArrMgHZQ2AyrW+3K6P98sJsuA7ZgcPNd4jtdty12W0elOwx3bZNWQWH
AjQtluhXoJS8UWBaL2ddI1IIYaDr66e5AQ0EXWv3MwEIALlCz1EzRbCHXm7r+F7J
4ZLw4q+xQYhWeX6lvxEOEPSy+hnM8WMGE2lcas26+gpg+LZTOVp7K41bzWVuUcJf
Kz/XJuJ8Ko7SBsgwLKLEt1Deq67t7YluJJDZKW9o5FT5D7FrQeBgk8oCgPO1jhsT
N57MV9K9jWrilQdxCKTaFGmK1dK6FZDEB8QFMkrImt4xykI0yUvWDjAHYd5llWcN
aAZXLob+jXcwXsMVgIYcUVP9mI7pSJQYdBr0S/cUzpUrKUMk3ricthb8kx63PFIv
Crn9D5Uak6xH2JZMVRfOIz7e5XSHoYazl3wiS0sVMgrENIBKLpMwerMoLxlCGk71
GjEAEQEAAYkBPAQYAQgAJhYhBO/23aKj1YN+l3IgGQSE0qsZXZ+tBQJda/czAhsM
BQkDwmcAAAoJEASE0qsZXZ+thqAIALzk52X+1rEcI3r3bDE9oRNzdS70NYXnDrhC
16QmNOOPYFmhDwok0aOfJOmJ6mlh73k9g7PMs1iiE8KnSnre3UiyZRr/gPdqrkpa
n8JgBtCfraNsIEMsol/pb0lLvR8Zu9SqTG7JPpO2IY77MBpaJgn48FiCAqOK5Kpb
jUTX4EvpEzimt0ywoR1cjwQEYtIvqTctzNRLDcObK17sOm9sDPqrHoz3nFtkg3OA
Y/nybO7JFvzT2usdSD4cXs/CgmMyS8rV1gijbZGu/O7wyHAS9puGkkpKyu2DG1oY
7wrjocQGevNRPe7UbDaFzQtwA2lkZsoYqkV4+9rC/pn0eBfZ01c=
=WvBQ
-----END PGP PUBLIC KEY BLOCK-----`
	TestPrivateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----

lQOYBF1r9zMBCAC96l46ncxNYkeXz//le7qOMCh2reSFt7pgSd4pmjA0sR61uqIu
Z+h/Q6Ssk4lMeisgj5tV8EvU8OyKiuXCQZXmKJmswCtCqmq4a3pcqqcbwmmm7AqV
uB9be4NJlJy2YK5uSx37VAC/ohhsR1g8EhMnB9RiyVIkESTPbLZNj96SVApJYzJy
r8K5EeOKA4zFEoxfwwg8JvdbON3NGVbH5wpU5sks/vm3kw/gGL8MHNRvaTVUUJZg
T61X2HWkiexNDeKNgKct3KwRLnmcsIe6/eRQGratrYbqyyiRlfacaRW1E6btA4Rw
+4YkxvzhOKgqADrY3jL+GFwhujUC28ovJMpFABEBAAEAB/4oIvCALdpAHn2n8W8Q
9kmyhC4BCvYpZm9uqx3XJ/15HVeyq3tcGFy0DL0wNcwGDCmyK3d2PZ8hJTuwefNd
zxOZfUohDsn1wGNmj1xgWjaP8NH4U5lXE/8Yi+1urwhBqQQkTh7Lf3DbDyxfMeLN
azp9/LMfn5GLCs5/QefzU20D+n6N8+VWqOLG0qpVgXm9te5+D5dzlIkmpXKn42kK
DNz5vdy+ueqthnzUBciqIeL+XOAYNoN79abyL7msUJp0NgakkvrtY+IHTz1Q155a
McQ3ZPg9fRCpjZ+y7mdSQEG7EWe1MW1MHqhRz4Dmt6OQ3sY2ZuDhltHs4XpBa7Zo
fUcfBADQ4uFiLLRASTR87hurxiym5zil/8FCyLqVuXKcCy9bUfmRh3ynjkVG22tu
hoVjSHu6hW59xhURZu0jOHBcwAxwhlJu3rrH0cPHVRXmT83SxMxfgXKI74+yQIa7
4ITeSjj5hT8hn3XzRuODo4JNd+w1PFX37G4aONBal+UKiOvNqwQA6MAdL66aEV8/
5023ctn+TK45H2N/nUe6NRailNHh6GPmsW4Q6RYIPG+LTVzfVGNEozOH5R5Jv6As
NZ2aCB6Rz24z99oLqG7SKZPIhTxGhvIdYFdWKzBWkB4160pw2siNPMSA54kPRhkM
BYYLjlsNtVrJP77uWVCftwsHKi6dd88D/1ganVUqLYQpLXVzl76pQaWqvnFAVDNh
0pHsMmTzL6E1SS+/x8f3yTJwJyJBEWwLjYAvqp7AtvM/BzJ/5cuhSuNCVnVmQiH6
sdJfqQCc6pa3vQquqB0fbDr625WCQDHxcjJqOyaFu8wsN7TArMHyq4fVNC15vyIA
X+qwfXiEbKPePyy0HlRlc3QgUGVyc29uIDx0ZXN0QGV4YW1wbGUuY29tPokBVAQT
AQgAPhYhBO/23aKj1YN+l3IgGQSE0qsZXZ+tBQJda/czAhsDBQkDwmcABQsJCAcC
BhUKCQgLAgQWAgMBAh4BAheAAAoJEASE0qsZXZ+th90H/RIoVT+qz/fHwXbtzpZY
P6i8TqS5tQQ+3n7Xxotpx2+C3kfc/jTLNrltnIfOKQTbzwo4z+b1g3KhP/Cylw5L
pNXggNHf9kKUL28idw+3mpLz5GnupIuvLSK+Gh2vZWY6OuJOp2ajEOBT16rmIFcV
gx9PFm6FDeBm48lCnnndmncuoTrNpagAwGRGGR1TVrIjDitZi8g83d0VHkskSbhj
1CbPa+bzlQSN9VJ0LTNlXGJxKPFgflGxghKXVn74v2+j7ArrMgHZQ2AyrW+3K6P9
8sJsuA7ZgcPNd4jtdty12W0elOwx3bZNWQWHAjQtluhXoJS8UWBaL2ddI1IIYaDr
66edA5gEXWv3MwEIALlCz1EzRbCHXm7r+F7J4ZLw4q+xQYhWeX6lvxEOEPSy+hnM
8WMGE2lcas26+gpg+LZTOVp7K41bzWVuUcJfKz/XJuJ8Ko7SBsgwLKLEt1Deq67t
7YluJJDZKW9o5FT5D7FrQeBgk8oCgPO1jhsTN57MV9K9jWrilQdxCKTaFGmK1dK6
FZDEB8QFMkrImt4xykI0yUvWDjAHYd5llWcNaAZXLob+jXcwXsMVgIYcUVP9mI7p
SJQYdBr0S/cUzpUrKUMk3ricthb8kx63PFIvCrn9D5Uak6xH2JZMVRfOIz7e5XSH
oYazl3wiS0sVMgrENIBKLpMwerMoLxlCGk71GjEAEQEAAQAH/RoeHlqg+xPsPaav
GyyH0waTca3ZtDarlETDyqg1nrPUsJgnEafHcUCjEAaRc9M3QRD5MiZ8o1LyLIZ+
c0XPA2qkYx2+agSI/P5Hdl97EqnyvmryrZB4p+yIxQPpGnmVRD5bs+WVT/iEelgB
UjekcaywO7hg0zlMmLx8Fb8h5ItNUvPmAzYNDO4D/qeZIPNQWAV97b1aNykFV2z0
s44dnQXSxlM5T6lJpiVunZU+XPBww0aBZRvk+3qZQmhIZr458PMjYlvnyVL6wj1W
76ITUVqd262oYtyZ6ZR7CrmK1otuWecW/RSn5w+TnzkMyyZkZox7Jn11VPWwIGzO
3ioPwg0EAMZUFEujBvb5ylnl/wRQQDOBPWO/CEtctsheoJJE9vAtSms/GDnJszX9
thndPX7l2mmBq5OXdCgXv9tHjlii2oomgGtDSU5fOYVwdJf4YKePUXavkOgCqY8A
qPixuvoZ6Qesege5sEXbZDTdPsMvId7zWZ0aWRZVDdae7qPGElhdBADvIfc54zPJ
lffNJg5PKOfcbgY1AUSkLlB7u4u51IIKR1/XPMyG7zzvKgAbNUSIVLlFeUwK6jMd
HyJ8PRjVc7n0f5lHI3HrRHeXhlCGPvHcwfFIBFOMswnPTHxnpxi4fZR8FARsvKUe
s47bsMXBGMGccAAQjuGJEO1aEOAmiX9b5QQAuxRb1uhyoUho0PlIknkvg6/mfjQ7
Ry6a2ct5kFhfN+LHhZAtVz3a6ljvaPVbwbqamJkqj2qTwf9EhXNPBQYhLet2Qvz1
OdJIn5us4TjJd1IT8ox5a8tLoohKzQNoXHpRm+VGy+8W4Duc0LJXnt/hONTzYF8O
AQ1uwmvZCVxmzBo67okBPAQYAQgAJhYhBO/23aKj1YN+l3IgGQSE0qsZXZ+tBQJd
a/czAhsMBQkDwmcAAAoJEASE0qsZXZ+thqAIALzk52X+1rEcI3r3bDE9oRNzdS70
NYXnDrhC16QmNOOPYFmhDwok0aOfJOmJ6mlh73k9g7PMs1iiE8KnSnre3UiyZRr/
gPdqrkpan8JgBtCfraNsIEMsol/pb0lLvR8Zu9SqTG7JPpO2IY77MBpaJgn48FiC
AqOK5KpbjUTX4EvpEzimt0ywoR1cjwQEYtIvqTctzNRLDcObK17sOm9sDPqrHoz3
nFtkg3OAY/nybO7JFvzT2usdSD4cXs/CgmMyS8rV1gijbZGu/O7wyHAS9puGkkpK
yu2DG1oY7wrjocQGevNRPe7UbDaFzQtwA2lkZsoYqkV4+9rC/pn0eBfZ01c=
=2m/S
-----END PGP PRIVATE KEY BLOCK-----
`
	// TestEncryptedMessage contains the plaintext "Passw0rd"
	TestEncryptedMessage = `-----BEGIN PGP MESSAGE-----

hQEMA/IIlIGJbYK5AQgArcfKZQLY41DE8/jTbQy/NSNgRmRdkO1uy3Ayf9cu/pGR
xC0ZP0yJwgZfNquQ33i20P+GP0ZDZ1CH6qlmTpn5q4KSXokEqKlZyELMaD/ruMzj
iG62/dHX/lHnAoO0vq4Gp6U08tIdUOZ6uxxJfdJRp5U/CJEY3I5bXvWrOa+7Ntqo
ZKnMy75z5ydMgzfTJOGsM4aF3JUg46aQoXLziaxvrC6Cp3QkC5eyvk3zBZMuNSng
qZ2QXrizPLc3Aq1aN17ITiA2XVW1MbQZMOsvNnyuOtz5mgzD2rZKxA3D07cFdbVF
ps72uhkHe6w0Hcl5UiWDpz+E1ZzQmmOqaBgfkYiTCNJEAbyLm+Oh+xMuCjhLj2h3
FKctaiBV3pc9iVC1Idawq6GSEOD7bsn06NaZs8oPuEO7cON2Wk7EpOU0CmK1m1Ia
+Ke6Eio=
=pyN7
-----END PGP MESSAGE-----`
)

func importGpgKey(key string, t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "go-pass-store")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatal("creating temp gpg-id:", err)
	}
	_, err = tmpFile.Write([]byte(key))
	if err != nil {
		t.Fatal("writing data to temp file", err)
	}
	gpgImport := exec.Command("gpg",
		"--import", tmpFile.Name())
	out, err := gpgImport.CombinedOutput()
	if err != nil {
		t.Log(string(out))
		t.Fatal("importing test GPG key:", err)
	}
}

func deleteGpgKey(secret bool, t *testing.T) {
	command := "--delete-keys"
	if secret {
		command = "--delete-secret-keys"
	}
	gpgDelete := exec.Command("gpg",
		"--batch", "--yes", command, TestPublicKeyFingerPrintLong)
	out, err := gpgDelete.CombinedOutput()
	if err != nil {
		t.Log(string(out))
		t.Fatal("deleting test GPG key:", err)
	}
}

func TestEntry_Decrypt(t *testing.T) {
	importGpgKey(TestPrivateKey, t)
	importGpgKey(TestPublicKey, t)
	defer func() {
		deleteGpgKey(true, t)
		deleteGpgKey(false, t)
	}()
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	err = ioutil.WriteFile(tmpPath+"/Test.gpg", []byte(TestEncryptedMessage), 0644)
	if err != nil {
		t.Fatal("writing to tmp file:", err)
	}

	type fields struct {
		Name  string
		Path  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Name:  tt.fields.Name,
				Path:  tt.fields.Path,
				value: tt.fields.value,
			}
			if err := e.Decrypt(); (err != nil) != tt.wantErr {
				t.Errorf("Entry.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntry_Encrypt(t *testing.T) {
	importGpgKey(TestPrivateKey, t)
	importGpgKey(TestPublicKey, t)
	defer func() {
		deleteGpgKey(true, t)
		deleteGpgKey(false, t)
	}()
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	err = ioutil.WriteFile(tmpPath+"/Test.gpg", []byte(TestEncryptedMessage), 0644)
	if err != nil {
		t.Fatal("writing to tmp file:", err)
	}
	testRunning = true

	type fields struct {
		Name  string
		Path  string
		value string
	}
	type args struct {
		recipients []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "Passw0rd",
			},
			args: args{
				recipients: []string{TestEmail},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Name:  tt.fields.Name,
				Path:  tt.fields.Path,
				value: tt.fields.value,
			}
			if err := e.Encrypt(tt.args.recipients); (err != nil) != tt.wantErr {
				t.Errorf("Entry.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntry_Show(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)

	type fields struct {
		Name  string
		Path  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "Passw0rd",
			},
			wantErr: false,
		},
		{
			name: "Simple - not decrypted",
			fields: fields{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Name:  tt.fields.Name,
				Path:  tt.fields.Path,
				value: tt.fields.value,
			}
			if err := e.Show(); (err != nil) != tt.wantErr {
				t.Errorf("Entry.Show() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	type args struct {
		name string
		path string
	}
	tests := []struct {
		name string
		args args
		want *Entry
	}{
		{
			name: "Simple",
			args: args{
				name: "Test",
				path: tmpPath + "/Test.gpg",
			},
			want: &Entry{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "",
			},
		},
		{
			name: "Simple - empty Entry",
			args: args{},
			want: &Entry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.name, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_Value(t *testing.T) {
	type fields struct {
		Name  string
		Path  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				Name:  "Test",
				value: "Passw0rd",
			},
			want:    "Passw0rd",
			wantErr: false,
		},
		{
			name: "Simple - no value",
			fields: fields{
				Name: "Test",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Name:  tt.fields.Name,
				Path:  tt.fields.Path,
				value: tt.fields.value,
			}
			got, err := e.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Entry.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Entry.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_Insert(t *testing.T) {
	type fields struct {
		Name  string
		Path  string
		value string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Simple",
			fields: fields{
				Name:  "Test",
				value: "Foo",
			},
			args: args{
				value: "Bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Name:  tt.fields.Name,
				Path:  tt.fields.Path,
				value: tt.fields.value,
			}
			e.Insert(tt.args.value)
			// This is probably unnecessary
			if e.value != tt.args.value {
				t.Errorf("Entry.value() = %v, want %v", e.value, tt.args.value)
			}
		})
	}
}

func TestEntry_Delete(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	err = ioutil.WriteFile(tmpPath+"/Test.gpg", []byte(TestEncryptedMessage), 0644)
	if err != nil {
		t.Fatal("writing to tmp file:", err)
	}
	type fields struct {
		Name  string
		Path  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "",
			},
			wantErr: false,
		},
		{
			name: "Simple - already deleted",
			fields: fields{
				Name:  "Test",
				Path:  tmpPath + "/Test.gpg",
				value: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Name:  tt.fields.Name,
				Path:  tt.fields.Path,
				value: tt.fields.value,
			}
			if err := e.Delete(); (err != nil) != tt.wantErr {
				t.Errorf("Entry.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
