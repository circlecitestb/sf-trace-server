#!/bin/bash
function install_elastic_apm_server {
    cudir=$PWD
    # Require that this runs as root.
    [ "$UID" -eq 0 ] || exec sudo "$0" "$@"

    wget -P /etc  https://artifacts.elastic.co/downloads/apm-server/apm-server-oss-7.8.1-linux-x86_64.tar.gz
    cd /etc && tar -xzf apm-server-oss-7.8.1-linux-x86_64.tar.gz
    mv /etc/apm-server-7.8.1-linux-x86_64 /etc/apm-server
    cd $cudir

    RELEASE_URL="https://api.github.com/repos/snappyflow/sf-trace-server/releases/latest"
    curl -sL $RELEASE_URL \
    | grep -w "browser_download_url" \
    | grep "apm-server-oss" \
    | cut -d":" -f 2,3 \
    | tr -d '"' \
    | xargs wget -q

    tar -zxvf apm-server-oss-*linux*.tar.gz

    cp -f apm-server-oss /etc/apm-server/apm-server
    rm -rf apm-server-oss*
}

function configure_apm_server {
    mv /etc/apm-server/apm-server.yml /etc/apm-server/apm-server.yml.copy
    cat > /etc/apm-server/apm-server.yml <<EOF
apm-server:
  host: 0.0.0.0:8200
  rum.enabled: true
  jaeger:
    http:
      enabled: true
      host: 0.0.0.0:14268
output.elasticsearch:
  hosts:
  - localhost:9200
  indices:
  - index: trace-%{[labels._tag_profileId]}_write
  password: $2
  username: $1
  worker: 10
  bulk_max_size: 5120
setup.template.enabled: false
queue.mem:
  events: 102400
  flush.min_events: 10000
EOF

}


function create_service_file {
    if [ ! -f /etc/systemd/system/apm-server.service ]; then
	cat > /etc/systemd/system/apm-server.service <<EOF
[Unit]
Description=Elastic APM Server
Documentation=https://www.elastic.co/solutions/apm
Wants=network-online.target
After=network-online.target
[Service]
ExecStart=/etc/apm-server/apm-server -c /etc/apm-server/apm-server.yml -path.home /etc/apm-server -path.config /etc/apm-server -path.data /var/lib/apm-server -path.logs /var/log/apm-server
Restart=always
[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable apm-server.service
systemctl start apm-server.service

    fi
}

echo "Install Trace server"
install_elastic_apm_server
configure_apm_server $1 $2
create_service_file
