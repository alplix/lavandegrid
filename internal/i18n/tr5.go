package i18n

func init() {
	reg("tr", map[string]string{
		"set.lang": "Dil", "set.appearance": "Gorunum", "set.theme": "Tema rengi",
		"set.startLocal": "Paket BOINC istemcisini baslat", "set.localStarted": "BOINC istemcisi baslatildi. Baglanti kabul etmesi bir kac saniye surebilir.",
		"set.localFailed": "Baslatma basarisiz", "set.localTitle": "Yerel BOINC istemcisi",
		"about.title": "Hakkinda", "about.desc": "BOINC filonuz icin lavanta temali bir yonetici.",
		"about.built": "Go + Fyne ile yapildi. BOINC istemcilerine GUI RPC ile baglanir.",
		"about.license": "2026 Alperen Yavuz. MIT Lisansi.",
	})
}