package bot

var ConstantMessages map[string]string = map[string]string{
	"/start": `Ciao, benvenuto in News dei Comuni!

Comandi disponibili:

	1. /iscrivi : Visualizza la lista di tutti i comuni presenti.
	2. /iscrizioni : Gestisci le tue iscrizioni.
	3. /informativa : Leggi l'informativa dati.

Per iniziare, scrivi /iscrivi e iscriviti alle notizie dei comuni che preferisci!.
`,
	"/informativa": `Le notizie pubblicate su questo bot Telegram sono fornite a scopo puramente informativo. Non possiamo garantire l'accuratezza o completezza delle informazioni condivise. Gli utenti sono incoraggiati a verificare le fonti e consultare altre fonti per decisioni importanti.

Il servizio fornito è offerto a titolo gratuito. Ciò significa che l'utilizzo del servizio non comporterà alcun costo o obbligo finanziario per l'utente. L'utente comprende e accetta che l'accesso e l'utilizzo del servizio sono forniti senza alcuna garanzia o responsabilità da parte del fornitore del servizio.

Il servizio fornito si basa su meccanismi automatici. Sebbene il servizio si impegna a fornire informazioni accurate e aggiornate, potrebbero verificarsi errori o inesattezze.

Il fornitore del servizio non sarà responsabile per eventuali danni o perdite derivanti dall'uso del servizio, inclusi, ma non limitati a, danni diretti, indiretti, accidentali, speciali o consequenziali. Inoltre, il fornitore del servizio si riserva il diritto di interrompere o modificare il servizio in qualsiasi momento senza preavviso.

L'utente è invitato a prendere atto di questo disclaimer prima di utilizzare il servizio. In caso di disaccordo con qualsiasi parte di questo disclaimer, si consiglia di non utilizzare il servizio fornito. L'uso continuato del servizio costituirà l'accettazione dei termini e delle condizioni esposte in questo disclaimer.`,
	"internal-error:no-towns": "Non sei iscritto a nessuna città. Scrivi /iscrivi per vedere la lista delle città alle quali puoi iscriverti.",
}
