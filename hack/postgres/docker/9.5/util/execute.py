import sys, os, subprocess, json, time
from datetime import datetime

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

def get_auth():
    Flag["username"] = "postgres"
    try:
        with open('/srv/postgres/secrets/.admin') as data_file:
            for line in data_file:
                s = line.rstrip().split("=",1)
                if s[0] == "POSTGRES_USERNAME":
                    Flag["username"] = s[1]
                elif s[0] == "POSTGRES_PASSWORD":
                    Flag["password"] = s[1]
    except:
        print "fail"
        exit(1)

def continuous_exec(process):
    code = 1
    start = datetime.utcnow()
    while True:
        code = subprocess.call(['./utils.sh', process, Flag["host"], Flag["username"], Flag["password"]])
        if code == 0:
            break
        now = datetime.utcnow()
        duration = (now - start).seconds
        if duration > 120:
            break
        time.sleep(30)

    if code != 0:
        print "fail"
        exit(1)

def main(argv):
    for flag in argv:
        if flag[:2]!= "--":
            continue
        v = flag.split("=", 1)
        Flag[v[0][2:]]=v[1]

    for flag in ["process", "cloud", "bucket", "snapshot", "host", "database"]:
        if flag not in Flag:
            print '--%s is required'%flag
            return

    if Flag["cloud"] == "gce":
        gce_cred_file()
    elif Flag["cloud"] == "aws":
        aws_cred_file()
    else:
        return

    if Flag["process"] == "backup":
        get_auth()

        continuous_exec("dump")

        code = subprocess.call(['./utils.sh', "push", Flag["cloud"], Flag["bucket"], Flag["database"], Flag["snapshot"]])
        if code != 0:
            print "fail"
            exit(1)

    if Flag["process"] == "restore":
        get_auth()

        code = subprocess.call(['./utils.sh', "pull", Flag["cloud"], Flag["bucket"], Flag["database"], Flag["snapshot"]])
        if code != 0:
            print "fail"
            exit(1)

        continuous_exec("restore")

    print "success"


if __name__ == "__main__":
    main(sys.argv[1:])
