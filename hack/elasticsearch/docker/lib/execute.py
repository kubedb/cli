from elasticsearch import Elasticsearch
import sys, time

Flag = {}
es = None

####
#
# ElasticSearch Part Start
#
####

###### Get Repository MAP for AWS ######
def get_aws_repo():
    json_map = {
        "type": "s3",
        "settings": {
            "bucket": Flag["bucket"],
            "region": Flag["region"],
            "base_path": Flag["database"],
            "access_key": Flag["keyid"],
            "secret_key": Flag["secret"],
            "max_retries": 5,
            "compress": True,
            "server_side_encryption": True,
        },
    }
    return json_map
########################################

###### Setup snapshot repository ######
def setup_repositor(json_map):
    print "Setting up snapshot..."
    global es
    es = Elasticsearch(hosts=[{'host': "127.0.0.1", 'port': 9200}], timeout=120)
    try:
        res = es.snapshot.create_repository(repository="es_backup", body=json_map, request_timeout=120, ignore=[400, 404, 500, 503])
        if "acknowledged" in res:
            print '''Setting complete for %s''' % (Flag["process"])
        else:
            print res["error"]["reason"]
            print "Fail to setup..."
            print "fail"
            exit(1)
    except:
        print "Unknown Error"
        print "Fail to setup..."
        print "fail"
        exit(1)
#######################################

###### Take Backup ######
def backup():
    print "Starting backup process..."
    global es
    try:
        res = es.snapshot.create("es_backup", Flag["snapshot"], request_timeout=120, ignore=[400, 404, 500, 503])
        if "accepted" in res:
            print '''Backup Started...'''
        elif "snapshot" in res:
            print '''Backup Started...'''
        else:
            print res["error"]["reason"]
            print "Fail to backup..."
            print "fail"
            exit(1)
    except:
        print "Unknown Error"
        print "Fail to backup..."
        print "fail"
        exit(1)

#########################

###### Restore Snapshot ######
def restore():
    print "Starting restore process..."
    global es
    try:
        res = es.indices.close(index="_all", request_timeout=120, ignore=[400, 404, 500, 503])
        if "acknowledged" in res:
            print "Indices closed..."
        else:
            print res["error"]["reason"]

        res = es.snapshot.restore("es_backup", Flag["snapshot"], request_timeout=120, ignore=[400, 404, 500, 503])
        if "accepted" in res:
            print "Restore Complete..."
        elif "snapshot" in res:
            print "Restore Complete..."
        else:
            print res["error"]["reason"]
            print "Fail to restore..."
            print "fail"
            exit(1)
    except:
        print "Unknown Error"
        print "Fail to restore..."
        print "fail"
        exit(1)
##############################

###### Check Restore Success ######
def check_cluster_health():
    global es
    try:
        res = es.cluster.health()
        if res["status"] == "red":
            return False
        else:
            return True
    except:
        return False

    return False
###################################

###### Check Backup Success ######
def check_backup_success():
    global es
    try:
        res = es.snapshot.status(repository="es_backup", snapshot=Flag["snapshot"])
        for snap in res["snapshots"]:
            if snap["state"] == "SUCCESS":
                return True
    except:
        return False

    return False
###################################

def check_startup_success():
    global es
    es = Elasticsearch(hosts=[{'host': "127.0.0.1", 'port': 9200}], timeout=120)
    while True:
        ok = check_cluster_health()
        if ok:
            print "success"
            return True
        time.sleep(30)
    print 'fail'
    exit(1)


####
#
# ElasticSearch Part End
# Controller Part Start
#
####

###### To control Process ######
def process_controller():
    if Flag["process"] == "backup":
        backup()
        while True:
            ok = check_backup_success()
            if ok:
                print '''Backup Completed...'''
                break
            time.sleep(30)

    if Flag["process"] == "restore":
        restore()
################################

###### Check Success ######
def check_success():
    while True:
        ok = check_cluster_health()
        if ok:
            break
        time.sleep(30)
###################################

def snap_process():
    for flag in ["bucket", "database", "snapshot", "region", "secret", "keyid"]:
        if flag not in Flag:
            print '--%s is required' % flag
            print 'fail'
            exit(1)

    setup_repositor(get_aws_repo())
    process_controller()
    check_success()
    print 'success'

####
#
# Controller Part End
#
####




def main(argv):
    print "Start processing..."
    for flag in argv:
        if flag[:2] != "--":
            continue
        v = flag.split("=", 1)
        Flag[v[0][2:]] = v[1]

    if "process" not in Flag:
        print '--process is required'
        print 'fail'
        exit(1)

    if Flag["process"] in ["backup", "restore"]:
        snap_process()
    elif Flag["process"] == "startup":
        check_startup_success()

    exit(0)

if __name__ == "__main__":
    main(sys.argv[1:])
