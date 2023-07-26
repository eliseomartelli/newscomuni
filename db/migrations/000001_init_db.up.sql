-- Create tables and populate the db.

CREATE TABLE IF NOT EXISTS towns (
	id VARCHAR(100) PRIMARY KEY,
	name VARCHAR(100),
	feed_url TEXT,
	last_updated int DEFAULT (strftime('%s', 'now') * 1000)
);

CREATE TABLE IF NOT EXISTS subscriptions (
	chat_id INTEGER,
	town    VARCHAR(100),
	FOREIGN KEY(town)
		REFERENCES towns(id)
		ON DELETE CASCADE
	PRIMARY KEY (chat_id, town)
);

INSERT INTO towns(id, name, feed_url)
	VALUES
		("frisa", "Frisa", "https://www.comune.frisa.ch.it/po/elenco_news_rss.php"),
		("lanciano", "Lanciano", "https://www.lanciano.eu/c069046/po/elenco_news_rss.php?tags=1&area=H"),
		("mozzagrogna", "Mozzagrogna", "https://comunemozzagrogna.it/web/feed/"),
		("ortona", "Ortona", "https://morss.it/:items=||*[class=u-nbfc]/https://www.comuneortona.ch.it/"),
		("sandamarmar", "Santa Maria Imbaro", "https://www.comune.santamariaimbaro.ch.it/po/elenco_news_rss.php"),
		("sanvitochietino", "San Vito Chietino", "https://www.comunesanvitochietino.it/feed/"),
		("sasi", "🛠️ UTENZE - SASI", "https://sasispa.it/category/aggiornamenti/feed/"),
		("torino", "Torino", "http://www.comune.torino.it/cgi-bin/torss/rssfeed.cgi?id=1"),
		("torinodisangro", "Torino di Sangro", "https://www.comune.torinodisangro.ch.it")
;
