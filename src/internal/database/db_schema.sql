CREATE TABLE "Monitors" (
	"Monitor_id"	INTEGER NOT NULL UNIQUE,
	"UUID"	TEXT NOT NULL,
	"Url"	TEXT NOT NULL,
	"Type"	TEXT NOT NULL CHECK("Type" IN ('ICMP', 'TCP', 'HTTP')),
	"Interval_in_seconds"	INTEGER NOT NULL,
	"Timeout_in_seconds"	INTEGER NOT NULL,
	"Port"	INTEGER,
	PRIMARY KEY("Monitor_id" AUTOINCREMENT)
);

CREATE TABLE "Results" (
	"Check_id"	INTEGER NOT NULL,
	"Monitor_id"	INTEGER NOT NULL,
	"Is_up"	INTEGER NOT NULL,
	"Response_time_in_ms"	INTEGER NOT NULL,
	"Run_time"	INTEGER NOT NULL,
	FOREIGN KEY("Monitor_id") REFERENCES "Monitors"("Monitor_id"),
	PRIMARY KEY("Check_id" AUTOINCREMENT)
);

--add notifications table to organise which notifications are setup for each user
--add another table to organise which notification mechanisms are setup for which Monitor