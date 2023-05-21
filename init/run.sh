#!/bin/bash

SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")
DATE_STR=`date '+%Y_%m_%d'`

if [ -f /etc/systemd/system/pipe-scripts.service ]; then echo "file exists"; sleep 5; exit; fi
if ! [ -d $SCRIPTPATH/pipe-scripts-repository/ ]; then mkdir $SCRIPTPATH/pipe-scripts-repository; fi

echo "creating a service"
echo "" > /etc/systemd/system/pipe-scripts.service
cat <<EOT >> /etc/systemd/system/pipe-scripts.service

[Unit]
Description=pipe-scripts
After=syslog.target 
After=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=$SCRIPTPATH/
ExecStart=$SCRIPTPATH/pipe-scripts $1
Restart=always
RestartSec=10
KillMode=process

[Install]
WantedBy=multi-user.target

EOT

#############################################################################################
### CREATE SCRIPT SERVICE STOP
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/stop-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/stop-pipe-scripts.sh
#!/bin/bash

systemctl stop pipe-scripts

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/stop-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE START
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/start-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/start-pipe-scripts.sh
#!/bin/bash

systemctl start pipe-scripts

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/start-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE RESTART
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/restart-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/restart-pipe-scripts.sh
#!/bin/bash

systemctl restart pipe-scripts

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/restart-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE STATUS
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/status-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/status-pipe-scripts.sh
#!/bin/bash

systemctl status pipe-scripts

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/status-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE JOURNAL
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/journal-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/journal-pipe-scripts.sh
#!/bin/bash

journalctl -u pipe-scripts.service -f

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/journal-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE ENABLE
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/enable-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/enable-pipe-scripts.sh
#!/bin/bash

systemctl enable pipe-scripts
systemctl start pipe-scripts

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/enable-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE DISABLE
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/disable-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/disable-pipe-scripts.sh
#!/bin/bash

systemctl stop pipe-scripts
systemctl disable pipe-scripts

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/disable-pipe-scripts.sh

#############################################################################################
### CREATE SCRIPT SERVICE REMOVE
#############################################################################################
echo "" > $SCRIPTPATH/pipe-scripts-repository/remove-pipe-scripts.sh
cat <<EOT >> $SCRIPTPATH/pipe-scripts-repository/remove-pipe-scripts.sh
#!/bin/bash

systemctl stop pipe-scripts
systemctl disable pipe-scripts
rm /etc/systemd/system/pipe-scripts.service
systemctl daemon-reload

EOT
chmod +x $SCRIPTPATH/pipe-scripts-repository/remove-pipe-scripts.sh

#############################################################################################
### START SERVICE
#############################################################################################
systemctl daemon-reload
systemctl enable pipe-scripts
systemctl start pipe-scripts
systemctl status pipe-scripts
