CREATE TABLE moz_perms (
  id INTEGER PRIMARY KEY,
  origin TEXT,
  TYPE TEXT,
  permission INTEGER,
  expireType INTEGER,
  expireTime INTEGER,
  modificationTime INTEGER
);


CREATE TABLE moz_hosts (
  id INTEGER PRIMARY KEY,
  HOST TEXT,
  TYPE TEXT,
  permission INTEGER,
  expireType INTEGER,
  expireTime INTEGER,
  modificationTime INTEGER,
  isInBrowserElement INTEGER
);
