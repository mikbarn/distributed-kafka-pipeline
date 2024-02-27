CTL_HOSTNAME_PREFIX="${CTL_HOSTNAME_PREFIX:=kafka-kraft}"
SERVICE_FQDN="${SERVICE_FQDN:=kafka-svc.kafka.svc.cluster.local}"
CTL_REPLICAS="${CTL_REPLICAS:=3}"
CFG_FILE="${CFG_FILE:=tst.properties}"
IS_CONTROLLER="${IS_CONTROLLER:=NO}"
FMT_DIR="${FMT_DIR:=}"
TIMEOUT="${TIMEOUT:=30}"
EXTERNAL_IP="${EXTERNAL_IP:=}"

get_server() {
    y=( $(echo $1 | awk '{n=split($0,a,"-"); for (i=1; i<=n; i++) print a[i]}') );
    count=${#y[@]};
    other="y";
    if [[ $# -gt 1 ]]; then
        if [[ $2 == "y" ]]; then
            idx=$((count-1));
            echo "${y[idx]}" ;
            other="n";
        fi
    fi
    if [[ $other == "y" ]]; then
        keep="${y[0]}";
        for (( i = 1; i < count - 1; i++ )); do
            keep=$keep"-${y[$i]}";
        done
        echo $keep;
    fi
}

servers=();
serv_ids=()

fetch_servers() {
    servers=();
    serv_ids=()

    ips=( $(host $SERVICE_FQDN | awk '{print $4}') );

    echo "Found "${#ips[@]}" IPs!";
    for x in "${ips[@]}"; do echo $x; done

    svc="${SERVICE_FQDN%%\.*}"

    for x in "${ips[@]}"; do
        serv=$(host $x | grep "\."$svc"\." | awk '{print $5}')
        serv=$(echo $serv | awk '{split($0,a,"."); print a[1]}')

        serv_prefix=$(get_server $serv)
        serv_id=$(get_server $serv y)

        if [[ $serv_prefix == $CTL_HOSTNAME_PREFIX ]]; then
            echo "Adding server "$serv" ...";
            servers+=( $CTL_HOSTNAME_PREFIX"-"$serv_id"."$SERVICE_FQDN );
            serv_ids+=($serv_id)
        else
            echo "$serv is not matching prefix $CTL_HOSTNAME_PREFIX!"
        fi
    done

}


ok="NO"
for ((i=0; i< $TIMEOUT; i++)); do
    fetch_servers
    num_servers=${#servers[@]}
    if [[ $CTL_REPLICAS -ne $num_servers ]]; then
        echo "CTL_REPLICAS variable not matching found DNS: $CTL_REPLICAS != $num_servers! Sleeping ...";
        sleep 1;
    else
        ok="YES";
        echo "Good to go!";
        break;
    fi
done

if [[ $ok != "YES" ]]; then
    echo "Timed out wating for CTL_REPLICAS variable to match DNS!";
    exit -1;
fi

self_addr=$(hostname)"."$SERVICE_FQDN

echo "Adding $num_servers servers!";

server_list="";

for (( i=0; i < $num_servers; i++ )); do
    node_id=$((${serv_ids[i]} + 1))
    server_list=$server_list","$node_id"@"${servers[i]}":9093";
done

server_list=${server_list:1};

echo "File to edit: $CFG_FILE";
echo "Updating quorum voters to: ";
echo $server_list;

sed -i "s/controller.quorum.voters=.*/controller.quorum.voters=$server_list/" $CFG_FILE;

node_offset=$(get_server $(hostname) y);

if [[ $IS_CONTROLLER == "YES" ]]; then
    node_id=$(($node_offset + 1))
    echo "Updating node.id to $node_id";
    sed -i "s/node.id=.*/node.id=$node_id/" $CFG_FILE;
    sed -i "s/listeners=.*/listeners=CONTROLLER:\/\/$self_addr:9093/" $CFG_FILE;
else
    node_id=$(($node_offset + 5001))
    echo "Updating node.id to $node_id";
    sed -i "s/node.id=.*/node.id=$node_id/" $CFG_FILE;
    sed -i "s/listeners=.*/listeners=PLAINTEXT:\/\/$self_addr:9092/" $CFG_FILE;
    #sed -i "s/advertised.listeners=.*/advertised.listeners=PLAINTEXT:\/\/$self_addr:9092/" $CFG_FILE;
    if [[ -n $EXTERNAL_IP  ]]; then
        if [[ $EXTERNAL_IP == "FILE" ]]; then
            ok="NO"
            for ((i=0; i< $TIMEOUT; i++)); do
                if [[ -f ./external_ip.txt ]]; then
                    ok="YES";
                    echo "Good to go!";
                    external_ip=$(cat external_ip.txt);
                    break;
                else
                    echo "External ip file not found! Sleeping ...";
                    sleep 1;
                fi
            done
            if [[ $ok != "YES" ]]; then
                echo "Timed out wating for External ip!";
                exit -1;
            fi
        else
            external_ip=$EXTERNAL_IP;
        fi
        echo "Updating advertised listeners IP to $external_ip";
        my_port=$(($node_offset+30000))
        sed -i "s/advertised.listeners=.*/advertised.listeners=PLAINTEXT:\/\/$external_ip:$my_port/" $CFG_FILE;
    fi
fi

if [[ -n $FMT_DIR ]]; then
    bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c $CFG_FILE
fi

cmd="bin/kafka-server-start.sh $CFG_FILE";
echo "Running command: $cmd";
$cmd

