package store

import (
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"gitlab.com/JanMa/go-pass/pkg/store/entry"
)

const (
	TestEmail                   = "test@example.com"
	TestPublicKeyFingerPrint    = "0484D2AB195D9FAD"
	TestPublicSubKeyFingerPrint = "F2089481896D82B9"
	TestPublicKey               = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
)

func TestNew(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	tmpPathInvalid, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	os.RemoveAll(tmpPathInvalid)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		{
			name: "New Store",
			args: args{
				path: tmpPath,
			},
			want: &Store{
				Path:       tmpPath,
				entries:    nil,
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			wantErr: false,
		},
		{
			name: "New Store - no directory",
			args: args{
				path: tmpPathInvalid,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Fill(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	err = ioutil.WriteFile(tmpPath+"/entry.gpg", []byte("Passw0rd"), 0644)
	if err != nil {
		t.Fatal("creating dummy entry:", err)
	}

	tmpPathInvalid, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	os.RemoveAll(tmpPathInvalid)

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Fill Store",
			fields: fields{
				Path:       tmpPath,
				entries:    make(map[string]*entry.Entry),
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			wantErr: false,
		},
		{
			name: "Fill Store - no directory",
			fields: fields{
				Path:       tmpPathInvalid,
				entries:    make(map[string]*entry.Entry),
				gpgID:      tmpPathInvalid + "/.gpg-id",
				recipients: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			if err := s.Fill(); (err != nil) != tt.wantErr {
				t.Errorf("Store.Fill() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPasswordStore(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	tmpPathInvalid, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	err = ioutil.WriteFile(tmpPath+"/store", []byte("a File"), 0644)
	if err != nil {
		t.Fatal("creating dummy entry:", err)
	}

	os.RemoveAll(tmpPathInvalid)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		env     string
	}{
		{
			name:    "Get Store",
			want:    tmpPath,
			wantErr: false,
			env:     tmpPath,
		},
		{
			name:    "Get Store - directory not existing",
			want:    "",
			wantErr: true,
			env:     tmpPathInvalid,
		},
		{
			name:    "Get Store - no directory",
			want:    "",
			wantErr: true,
			env:     tmpPath + "/store",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("PASSWORD_STORE_DIR", tt.env)
			if err != nil {
				t.Fatal("setting environment variable:", err)
			}
			got, err := GetPasswordStore()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPasswordStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPasswordStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_FindEntry(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)

	testEntry := entry.New("Test", tmpPath+"/Test.gpg")

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	type args struct {
		e string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entry.Entry
		wantErr bool
	}{
		{
			name: "Find Entry",
			fields: fields{
				Path: tmpPath,
				entries: map[string]*entry.Entry{
					"Test": testEntry,
				},
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: "Test",
			},
			want:    testEntry,
			wantErr: false,
		},
		{
			name: "Find Entry - not present",
			fields: fields{
				Path: tmpPath,
				entries: map[string]*entry.Entry{
					"Test": testEntry,
				},
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: "Bar",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			got, err := s.FindEntry(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.FindEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.FindEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_FindEntries(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)

	testEntry := entry.New("Test", tmpPath+"/Test.gpg")

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	type args struct {
		e string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*entry.Entry
		wantErr bool
	}{
		{
			name: "Find Entries",
			fields: fields{
				Path: tmpPath,
				entries: map[string]*entry.Entry{
					"Test": testEntry,
				},
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: "Test",
			},
			want: map[string]*entry.Entry{
				"Test": testEntry,
			},
			wantErr: false,
		},
		{
			name: "Find Entries - not present",
			fields: fields{
				Path: tmpPath,
				entries: map[string]*entry.Entry{
					"Test": testEntry,
				},
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: "Bar",
			},
			want:    make(map[string]*entry.Entry),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			got, err := s.FindEntries(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.FindEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.FindEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortEntries(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	aEntry := entry.New("A", tmpPath+"/A.gpg")
	bEntry := entry.New("B", tmpPath+"/B.gpg")
	cEntry := entry.New("C", tmpPath+"/C.gpg")

	type args struct {
		entries map[string]*entry.Entry
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Sort Entries",
			args: args{
				entries: map[string]*entry.Entry{
					"A": aEntry,
					"B": bEntry,
					"C": cEntry,
				},
			},
			want: []string{"A", "B", "C"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortEntries(tt.args.entries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_InsertEntry(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	testEntry := entry.New("Test", tmpPath+"/Test.gpg")
	testMap := map[string]*entry.Entry{
		"Test": testEntry,
	}

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	type args struct {
		e *entry.Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Insert Entry",
			fields: fields{
				Path:       tmpPath,
				entries:    make(map[string]*entry.Entry),
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: testEntry,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			s.InsertEntry(tt.args.e)
			if !reflect.DeepEqual(s.entries, testMap) {
				t.Errorf("InsertEntry() = %v, want %v", s.entries, testMap)
			}
		})
	}
}

func TestStore_DeleteEntry(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	testEntry := entry.New("Test", tmpPath+"/Test.gpg")
	testEntryTwo := entry.New("Test2", tmpPath+"/Test2.gpg")
	testMap := map[string]*entry.Entry{
		"Test": testEntry,
	}

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	type args struct {
		e *entry.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete Entry - existing entry",
			fields: fields{
				Path:       tmpPath,
				entries:    testMap,
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: testEntry,
			},
			wantErr: false,
		},
		{
			name: "Delete Entry - not existing entry",
			fields: fields{
				Path:       tmpPath,
				entries:    testMap,
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				e: testEntryTwo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			if err := s.DeleteEntry(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("Store.DeleteEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_FindGpgID(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	err = os.Mkdir(tmpPath+"/Test", 0755)
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	err = ioutil.WriteFile(tmpPath+"/Test/.gpg-id", []byte("test@example.com"), 0644)
	if err != nil {
		t.Fatal("creating dummy gpg-id:", err)
	}

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Find GPG ID",
			fields: fields{
				Path:       tmpPath,
				entries:    nil,
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				path: tmpPath + "/Bar/example.gpg",
			},
			want: tmpPath + "/.gpg-id",
		},
		{
			name: "Find GPG ID - Subdir",
			fields: fields{
				Path:       tmpPath,
				entries:    nil,
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			args: args{
				path: tmpPath + "/Test/example.gpg",
			},
			want: tmpPath + "/Test/.gpg-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			if got := s.FindGpgID(tt.args.path); got != tt.want {
				t.Errorf("Store.FindGpgID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseGpgID(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "go-pass-store")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatal("creating temp gpg-id:", err)
	}
	_, err = tmpFile.Write([]byte("test@example.com"))
	if err != nil {
		t.Fatal("writing data to temp file", err)
	}
	tmpFileInvalid, err := ioutil.TempFile("", "go-pass-store")
	os.Remove(tmpFileInvalid.Name())

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Parse GPG ID",
			args: args{
				path: tmpFile.Name(),
			},
			want:    []string{"test@example.com"},
			wantErr: false,
		},
		{
			name: "Parse GPG ID - no file",
			args: args{
				path: tmpFileInvalid.Name(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGpgID(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGpgID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGpgID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_ShowAll(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "go-pass-store")
	if err != nil {
		t.Fatal("creating tmp dir:", err)
	}
	defer os.RemoveAll(tmpPath)
	testEntry := entry.New("Test", tmpPath+"/Test.gpg")
	testMap := map[string]*entry.Entry{
		"Test": testEntry,
	}

	type fields struct {
		Path       string
		entries    map[string]*entry.Entry
		gpgID      string
		recipients []string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*entry.Entry
	}{
		{
			name: "Show All",
			fields: fields{
				Path:       tmpPath,
				entries:    testMap,
				gpgID:      tmpPath + "/.gpg-id",
				recipients: nil,
			},
			want: testMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				Path:       tt.fields.Path,
				entries:    tt.fields.entries,
				gpgID:      tt.fields.gpgID,
				recipients: tt.fields.recipients,
			}
			if got := s.ShowAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.ShowAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGpgKeys(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "go-pass-store")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatal("creating temp gpg-id:", err)
	}
	_, err = tmpFile.Write([]byte(TestPublicKey))
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

	type args struct {
		recipients []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Get GPG Keys",
			args: args{
				recipients: []string{TestEmail},
			},
			want:    []string{TestPublicSubKeyFingerPrint},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGpgKeys(tt.args.recipients)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGpgKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGpgKeys() = %v, want %v", got, tt.want)
			}
		})
	}

	gpgDelete := exec.Command("gpg",
		"--batch", "--yes", "--delete-keys", TestEmail)
	out, err = gpgDelete.CombinedOutput()
	if err != nil {
		t.Log(string(out))
		t.Fatal("deleting test GPG key:", err)
	}
}
