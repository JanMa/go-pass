package cmd

var (
	Clip        bool
	Copy        int
	Echo        bool
	ForceCp     bool
	ForceGen    bool
	ForceInsert bool
	ForceMv     bool
	ForceRm     bool
	GenQRCode   bool
	InPlace     bool
	MultiLine   bool
	NoSymbols   bool
	OtpClip     bool
	QRCode      int
	Subdir      string
	Version     string
	Zsh         bool
)

func init() {
	// show flags
	showCmd.Flags().IntVarP(&Copy, "clip", "c", 0, "Copy password to clipboard")
	showCmd.Flags().Lookup("clip").NoOptDefVal = "1"
	showCmd.Flags().IntVarP(&QRCode, "qrcode", "q", 0, "Display output as QR code")
	showCmd.Flags().Lookup("qrcode").NoOptDefVal = "1"

	// rm flags
	rmCmd.Flags().BoolVarP(&ForceRm, "force", "f", false, "Forcefully remove password(s).")

	// mv flags
	mvCmd.Flags().BoolVarP(&ForceMv, "force", "f", false, "Forcefully copy password or directory.")

	// insert flags
	insertCmd.Flags().BoolVarP(&Echo, "echo", "e", false, "Echo password back to console")
	insertCmd.Flags().BoolVarP(&MultiLine, "multiline", "m", false, "Multiline input")
	insertCmd.Flags().BoolVarP(&ForceInsert, "force", "f", false, "Overwrite existing password without prompt")

	// init flags
	initCmd.Flags().StringVarP(&Subdir, "path", "p", "", "A specific gpg-id or set of gpg-ids is assigned for that specific sub folder of the password store")

	// generate flags
	generateCmd.Flags().BoolVarP(&NoSymbols, "no-symbols", "n", false, "Generate password with no symbols.")
	generateCmd.Flags().BoolVarP(&Clip, "clip", "c", false, "Put generated password on the clipboard.")
	generateCmd.Flags().BoolVarP(&InPlace, "in-place", "i", false, "Replace only the first line of an existing file with a new password.")
	generateCmd.Flags().BoolVarP(&ForceGen, "force", "f", false, "Forcefully overwrite existing password.")
	generateCmd.Flags().BoolVarP(&GenQRCode, "qrcode", "q", false, "Display output as QR code.")

	//cp flags
	cpCmd.Flags().BoolVarP(&ForceCp, "force", "f", false, "Forcefully copy password or directory.")

	//completion flags
	completionCmd.Flags().BoolVarP(&Zsh, "zsh", "z", false, "Generate Zsh completion.")

	//otp flags
	otpCmd.Flags().BoolVarP(&OtpClip, "clip", "c", false, "Put OTP code on the clipboard.")
}
