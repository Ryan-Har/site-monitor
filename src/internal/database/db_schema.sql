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
	PRIMARY KEY("Check_id" AUTOINCREMENT),
	FOREIGN KEY("Monitor_id") REFERENCES "Monitors"("Monitor_id") ON DELETE CASCADE
);

CREATE TABLE "Notifications" (
	"Notification_id"	INTEGER NOT NULL UNIQUE,
	"UUID"	TEXT NOT NULL,
	"Type"	TEXT NOT NULL CHECK("Type" in ('discord', 'slack', 'email')),
	"Additional_info"	TEXT NOT NULL,
	PRIMARY KEY("Notification_id")
);

CREATE TABLE "Monitor_Notifications" (
	"Monitor_id"	INTEGER NOT NULL,
	"Notification_id"	INTEGER NOT NULL,
	FOREIGN KEY("Notification_id") REFERENCES "Notifications"("Notification_id") ON DELETE CASCADE,
	FOREIGN KEY("Monitor_id") REFERENCES "Monitors"("Monitor_id")  ON DELETE CASCADE
);

CREATE TABLE "Indcidents" (
	"Incident_id"	INTEGER NOT NULL UNIQUE,
	"Start_time"	INTEGER NOT NULL,
	"End_time"	INTEGER,
	"Monitor_id"	INTEGER NOT NULL,
	PRIMARY KEY("Incident_id"),
	FOREIGN KEY("Monitor_id") REFERENCES "Monitors"("Monitor_id")
);
--add notifications table to organise which notifications are setup for each user
--add another table to organise which notification mechanisms are setup for which Monitor