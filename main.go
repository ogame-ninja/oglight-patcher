package main

import (
	ep "github.com/ogame-ninja/extension-patcher"
)

func main() {
	const (
		webstoreURL    = "https://openuserjs.org/install/nullNaN/OGLight.user.js"
		oglight_sha256 = "978c7932426a23b2c69d7f4be32ea4e8e3abbb6a3ea84d7278381be6336f55c3"
	)

	files := []ep.FileAndProcessors{
		ep.NewFile("OGLight.user.js", processOGLight),
	}

	ep.MustNew(ep.Params{
		ExpectedSha256: oglight_sha256,
		WebstoreURL:    webstoreURL,
		Files:          files,
	}).Start()
}

var replN = ep.MustReplaceN

func processOGLight(by []byte) []byte {
	by = replN(by, `@name         OGLight`, `@name         OGLight Ninja`, 1)
	by = replN(by, "// @match        https://*.ogame.gameforge.com/game/*\r\n",
		`{old}// @match        *127.0.0.1*/bots/*/browser/html/*?page=*
// @match        *.ogame.ninja/bots/*/browser/html/*?page=*
`, 1)
	by = replN(by, `// ==/UserScript==`,
		`{old}

	const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(window.location.href)[1];
	const lang = /browser\/html\/s(\d+)-(\w+)/.exec(window.location.href)[2];
	const UNIVERSE = "s" + universeNum + "-" + lang;
	const PROTOCOL = window.location.protocol;
	const HOST = window.location.host;
`, 1)
	by = replN(by, "var cookieAccounts = document.cookie.match(/prsess\\_([0-9]+)=/g), cookieAccounts = cookieAccounts[cookieAccounts.length - 1].replace(/\\D/g, \"\");",
		`var cookieAccounts=document.querySelector('head meta[name="ogame-player-id"]').getAttribute('content').replace(/\D/g, '');`, 1)
	by = replN(by, `this.server.id = window.location.host.replace(/\D/g, "")`,
		`this.server.id=document.querySelector('head meta[name="ogame-universe"]').getAttribute('content').replace(/\D/g,"")`, 1)
	by = replN(by, `this.account.lang = /oglocale=([a-z]+);/.exec(document.cookie)[1]`,
		`this.account.lang=lang`, 1)
	by = replN(by, "cache = [ crypto.randomUUID(), 0 ]",
		`cache = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
			var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
			return v.toString(16);
		})`, 1)
	by = replN(by, "url: `https://${window.location.host}/api/playerData.xml",
		"url:`${PROTOCOL}//${HOST}/api/s${universeNum}/${lang}/playerData.xml", 1)
	by = replN(by, "url: `https://${window.location.host}/api/serverData.xml`,",
		"url:`${PROTOCOL}//${HOST}/api/s${universeNum}/${lang}/serverData.xml`,", 1)
	by = replN(by, `${player.name} <a href="https://${window.location.host}/game/index.php?`,
		`${player.name} <a href="${window.location.protocol}//${window.location.host}${window.location.pathname}?`, 1)
	by = replN(by, "https://${window.location.host}/game/index.php?page=componentOnly&component=messagedetails&messageId=` + message.id",
		"${window.location.protocol}//${window.location.host}${window.location.pathname}?page=componentOnly&component=messagedetails&messageId=` + message.id", 1)
	by = replN(by, "https://${window.location.host}/game/index.php", "", 25)
	return by
}
