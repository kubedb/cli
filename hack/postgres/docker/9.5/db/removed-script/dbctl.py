# -*- coding: utf-8 -*-
import sys, psycopg2

Flag = {}


def get_password():
    try:
        with open('/srv/postgres/secrets/.admin') as data_file:
            for line in data_file:
                s = line.rstrip().split("=",1)
                if s[0] == "POSTGRES_PASSWORD":
                    Flag["password"] = s[1]
    except:
        sys.exit()


def db_connection():
    try:
        get_password()
        conn = psycopg2.connect(host="127.0.0.1", port=5432, database="postgres", user="postgres", password=Flag["password"] )
        conn.set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)
        return conn
    except:
        sys.exit()


def streaming_list():
     # Connect database with postgres user
    try:
        conn = db_connection()
        cur = conn.cursor()
        cur.execute("select * from pg_stat_replication;")
        rows = cur.fetchall()
        for row in rows:
            if row[9] == "streaming":
                print "node: ", row[3]
        print "success"
        return
    except:
        print "fail"
        sys.exit()


def main(argv):
    for flag in argv:
        if flag[:2]!= "--":
            continue
        v = flag.split("=", 1)
        Flag[v[0][2:]]=v[1]

    for flag in ["process"]:
        if flag not in Flag:
            print '--%s is required'%flag
            print "fail"
            sys.exit()

    if Flag["process"] == "streaming":
        streaming_list()
    else:
        print "fail"


if __name__ == "__main__":
    main(sys.argv[1:])
