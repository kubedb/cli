import sys, os, subprocess, json
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


from influxdb import InfluxDBClient
def dump(host):
    username = ""
    password = ""
    with open('/srv/influxdb/.admin') as data_file:
        for line in data_file:
            s = line.rstrip().split("=",1)
            if s[0] == "INFLUX_ADMIN_USER":
                username = s[1]
            elif s[0] == "INFLUX_ADMIN_PASSWORD":
                password = s[1]

    client = InfluxDBClient(host, 8086, username, password)
    db_list = client.get_list_database()
    for db in db_list:
        code = subprocess.call(['./utils.sh', "backup", Flag["host"]+":8088", db["name"], Flag["snapshot"]])
        if code != 0:
            exit(1)


def main(argv):
    for flag in argv:
        if flag[:2]!= "--":
            continue
        v = flag.split("=", 1)
        Flag[v[0][2:]]=v[1]

    for flag in ["cloud", "bucket", "snapshot", "host", "database"]:
        if flag not in Flag:
            print '--%s is required'%flag
            return

    if Flag["cloud"] == "gce":
        gce_cred_file()
    elif Flag["cloud"] == "aws":
        aws_cred_file()
    else:
        return

    dump(Flag["host"])
    subprocess.call(['./utils.sh', "push", Flag["cloud"], Flag["bucket"], Flag["database"], Flag["snapshot"]])

if __name__ == "__main__":
    main(sys.argv[1:])
