import sys, os, subprocess, json, shutil
from elasticsearch import Elasticsearch

Flag = {}
aws_cred_dir = "/root/.aws"

#############################################
# this will be written in /root/.boto for gce
gce_cred_data = '''[Credentials]
gs_service_key_file = /var/credentials/gce

[Boto]
https_validate_certificates = True

[GSUtil]
content_language = en
default_api_version = 2
default_project_id = %s
'''
#############################################

#############################################
# this will be written in /root/.aws/credentials for aws
aws_cred_data = '''[default]
aws_access_key_id = %s
aws_secret_access_key = %s
'''
#############################################


#######################################
# this func write .boto file for gsutil
def gce_cred_file():
    print "Setting GCE credential.."
    with open('/var/credentials/gce') as data_file:
        json_data = json.load(data_file)

    f = open('/root/.boto', 'w+')
    f.write(gce_cred_data%json_data['project_id'])
    f.close()
    return True
#######################################

#######################################
# this func write .aws/credentials file for aws
def aws_cred_file():
    print "Setting AWS credential.."
    with open('/var/credentials/aws/keyid') as data_file:
        key = data_file.read().rstrip()
    with open('/var/credentials/aws/secret') as data_file:
        secret = data_file.read().rstrip()

    if not os.path.exists(aws_cred_dir):
        os.makedirs(aws_cred_dir)
    f = open(aws_cred_dir+'/credentials', 'w+')
    f.write(aws_cred_data%(key,secret))
    f.close()
#######################################

def backup_process():
    print "Backup process starting..."

    es = Elasticsearch(hosts=[{'host': Flag["host"], 'port': 9200}], timeout=120)
    indices=es.indices.get_alias()

    print "Total indices: " + str(len(indices))
    path = '/var/dump-backup/'+Flag["snapshot"]
    shutil.rmtree(path, ignore_errors=True)

    if not os.path.exists(path):
        os.makedirs(path)

    for index in indices:
        code = subprocess.call(['./utils.sh', "backup", Flag["host"], Flag["snapshot"], index])
        if code != 0:
            print "Fail to take backup for index: "+index
            exit(1)

    filep = open(path+"/indices.txt", "wb")
    for index in indices:
        print>>filep, index
    filep.close()

    code = subprocess.call(['./utils.sh', "push", Flag["cloud"], Flag["bucket"], Flag["database"], Flag["snapshot"]])
    if code != 0:
        print "Fail to push backup files to cloud..."
        exit(1)

def restore_process():
    print "Restore process starting..."

    code = subprocess.call(['./utils.sh', "pull", Flag["cloud"], Flag["bucket"], Flag["database"], Flag["snapshot"]])
    if code != 0:
        print "Fail to pull backup files from cloud..."
        exit(1)

    path = '/var/dump-restore/'+Flag["snapshot"]
    fileP = open(path+"/indices.txt", "r")

    for index in fileP.readlines():
        index = index.rstrip("\n")
        code = subprocess.call(['./utils.sh', "restore", Flag["host"], Flag["snapshot"], index])
        if code != 0:
            print "Fail to restore index: "+index
            exit(1)


def main(argv):
    for flag in argv:
        if flag[:2] != "--":
            continue
        v = flag.split("=", 1)
        Flag[v[0][2:]]=v[1]

    for flag in ["process", "host", "cloud", "bucket", "snapshot", "database"]:
        if flag not in Flag:
            print '--%s is required'%flag
            return

    if Flag["cloud"] == "gce":
        gce_cred_file()
    elif Flag["cloud"] == "aws":
        aws_cred_file()
    else:
        print "Invalid Cloud"
        exit(1)
        return

    if Flag["process"] == "backup":
        backup_process()
    elif Flag["process"] == "restore":
        restore_process()

    exit(0)


if __name__ == "__main__":
    main(sys.argv[1:])
