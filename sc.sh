make build

for i in $(seq 1 24); do
  if [ $((i % 6)) -eq 1 ]; then
    echo ./bin/kubectl-dba-linux-amd64 data insert postgres ha-postgres -n demo -r 10000
    echo $(./bin/kubectl-dba-linux-amd64 data insert postgres ha-postgres -n demo -r 10000)
  fi
  if [ $((i % 6)) -eq 2 ]; then
    echo ./bin/kubectl-dba-linux-amd64 data verify postgres ha-postgres -n demo -r 10000
    echo $(./bin/kubectl-dba-linux-amd64 data verify postgres ha-postgres -n demo -r 10000)
  fi
  if [ $((i % 6)) -eq 3 ]; then
    echo ./bin/kubectl-dba-linux-amd64 data drop postgres ha-postgres -n demo
    echo $(./bin/kubectl-dba-linux-amd64 data drop postgres ha-postgres -n demo)
  fi
  if [ $((i % 6)) -eq 4 ]; then
    echo ./bin/kubectl-dba-linux-amd64 data insert postgres pg -n demo -r 10000
    echo $(./bin/kubectl-dba-linux-amd64 data insert postgres pg -n demo -r 10000)
  fi
  if [ $((i % 6)) -eq 5 ]; then
    echo ./bin/kubectl-dba-linux-amd64 data verify postgres pg -n demo -r 10000
    echo $(./bin/kubectl-dba-linux-amd64 data verify postgres pg -n demo -r 10000)
  fi
  if [ $((i % 6)) -eq 0 ]; then
    echo ./bin/kubectl-dba-linux-amd64 data drop postgres pg -n demo
    echo $(./bin/kubectl-dba-linux-amd64 data drop postgres pg -n demo)
    echo make build
    make build
  fi
  echo
done
