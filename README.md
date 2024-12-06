# W365toFET

```
W365toFET path/to/sp001_w365.json

    -> path/to/sp001_w365.log
    -> path/to/sp001.fet
    -> path/to/sp001.map
```

Die Stundenplan-Daten werden von der JSON-Datei eingelesen.

Ausgegeben wird eine FET-Datei im selben Ordner. Auch eine Logdatei (mit Fehlermeldungen, usw.) und eine Zuordnungsdatei für die FET-Activities werden erstellt.

## Kompilieren

Um mehrere ausführbare Dateien zu unterstützen, befinden sich die `main.go`-Dateien in Unterordnern des Ordners `cmd`. Zum Kompilieren (im Hauptordner):

```
go build cmd/W365toFET
```

Die ausführbaren Dateien können auch in einem anderen (schon existierenden!) Ordner abgelegt werden:

```
go build -o bin cmd/W365toFET
```

## Aktueller Stand (06.12.2024)

Bis auf die „Constraint“-Elemente werden alle Elemente in `docs/stundenplanschnittstelle.md` in einigermaßen entsprechende FET-Strukturen übertragen.

Die weiteren Constraints werden jetzt anfänglich übersetzt.

In dieser Version haben Lehrer und Klassen den gleichen Ansatz für Mittagspausen: Eine der Mittagsstunden muss frei sein.

Es wurde bisher nur wenig getestet!

In dieser Version werden die Daten in eine etwas andere interne Struktur gebracht – W365-unabhängig – bevor sie übersetzt werden, siehe „base-package“. Auch unabhängig von der FET-Ausgabe werden Grundlagen für die Stundenplanung in package "ttbase" vorbereitet.

## Neu: Druckausgabe

Stundenpläne können jetzt als PDF ausgegeben werden, aktuell die Klassentabellen und die Lehrertabellen. Dafür muss Typst installiert sein. Das Programm W365toTypst erstellt JSON-Dateien, die als Eingabe zu Typst-Skripten dienen. Es kann etwa so kompiliert werden:

```
go build -o bin cmd/W365toTypst
```

Um die Eingabedateien für die Typst-Skripten zu erstellen:

```
W365toTypst path/to/sp001_w365.json
```

Die resultierenden JSON-Dateien werden im Ordner `path/to/typst_files/_data` abgelegt. Ein Fehlerbericht kann, wie bei W365toFET, in der Log-Datei gefunden werden.

Der Befehl, um die PDF-Ausgabe zu erstellen, sieht etwa so aus:

```
typst compile --root "path/to/typst_files" --input ifile="/_data/sp001_teachers.json" "path/to/typst_files/scripts/print_timetable.typ" "path/to/typst_files/_pdf/sp001_teachers.pdf"
```

Der Ordner für die PDF-Ausgabe muss schon existieren.

Bei Erfolg wären die Ergebnisse dann im Ordner `path/to/_pdf` zu finden. Fehlermeldungen kann man von `stderr` lesen.


typst compile --root "/home/user/Development/github/W365toFET/_testdata1/typst_files" --input ifile="/_data/testx01_teachers.json" "/home/user/Development/github/W365toFET/_testdata1/typst_files/scripts/print_timetable.typ" "/home/user/Development/github/W365toFET/_testdata1/typst_files/_pdf/testx01_teachers.pdf"