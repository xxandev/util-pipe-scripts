#!/bin/bash

SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")
DATE_STR=`date '+%Y_%m_%d'`

if [ -f /etc/systemd/system/pipe.service ]; then echo "file service exists"; sleep 5; exit; fi
if ! [ -d $SCRIPTPATH/pipe-rep/ ]; then mkdir $SCRIPTPATH/pipe-rep; fi

echo "creating a service"
echo "" > /etc/systemd/system/pipe.service
cat <<EOT >> /etc/systemd/system/pipe.service
[Unit]
Description=pipe
After=syslog.target 
After=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=$SCRIPTPATH/
ExecStart=$SCRIPTPATH/pipe $1
Restart=always
RestartSec=10
KillMode=process

[Install]
WantedBy=multi-user.target

EOT

#############################################################################################
### CREATE SCRIPT SERVICE STOP
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-stop.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-stop.sh
#!/bin/bash

systemctl stop pipe

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-stop.sh

#############################################################################################
### CREATE SCRIPT SERVICE START
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-start.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-start.sh
#!/bin/bash

systemctl start pipe

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-start.sh

#############################################################################################
### CREATE SCRIPT SERVICE RESTART
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-restart.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-restart.sh
#!/bin/bash

systemctl restart pipe

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-restart.sh

#############################################################################################
### CREATE SCRIPT SERVICE STATUS
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-status.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-status.sh
#!/bin/bash

systemctl status pipe

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-status.sh

#############################################################################################
### CREATE SCRIPT SERVICE JOURNAL
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-journal.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-journal.sh
#!/bin/bash

journalctl -u pipe.service -f

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-journal.sh

#############################################################################################
### CREATE SCRIPT SERVICE ENABLE
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-enable.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-enable.sh
#!/bin/bash

systemctl enable pipe
systemctl start pipe

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-enable.sh

#############################################################################################
### CREATE SCRIPT SERVICE DISABLE
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-disable.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-disable.sh
#!/bin/bash

systemctl stop pipe
systemctl disable pipe

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-disable.sh

#############################################################################################
### CREATE SCRIPT SERVICE REMOVE
#############################################################################################
echo "" > $SCRIPTPATH/pipe-rep/pipe-remove.sh
cat <<EOT >> $SCRIPTPATH/pipe-rep/pipe-remove.sh
#!/bin/bash

systemctl stop pipe
systemctl disable pipe
rm /etc/systemd/system/pipe.service
systemctl daemon-reload

EOT
chmod +x $SCRIPTPATH/pipe-rep/pipe-remove.sh

#############################################################################################
### START SERVICE
#############################################################################################
systemctl daemon-reload
systemctl enable pipe
systemctl start pipe
systemctl status pipe
