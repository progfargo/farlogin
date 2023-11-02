# farlogin

provides access to a machine that does not have a static ip. it allows you to use the terminal screen of the remote machine.

it uses embedded nats.io server for communication.

i got this idea from this page: https://dev.to/napicella/linux-terminals-tty-pty-and-shell-192e

this method can be used to access nodes, especially during the development of iot projects.

it provides ease of use in case of a large number of nodes (hundreds), and can simplify asset management because it works independently of ip addresses.

*what else can be done:*
in a real installation, https should be used and all communication should be via tls.

flnode (the agent copied to the node) can be conditionally compiled so that it runs only on a specific machine (by cpu id etc.).

in addition to tls, two-way encryption can be used as an additional security measure. login history can be kept, who connected to which machine from where, which commands they ran, etc.

if case you want to play with:

* install nginx and use it as reverse proxy. proxy pass all request to your domain to farlogin app. there is a sample proxy config in asset directory.
* import the farlogin database to mysql on the same machine farlogin app is running.
* copy fladmin to the computer of admin user.
* make the necessary changes in farlogin.ini file. app will work on localhost initially. if you want to run it on a server which has a public domain name you should change the ini file accordingly and dont forget to change the nats url in the app.go files for farlogin, fladmin and flnode.
* after running farlogin web app, login and create a node by giving it a name.
* run flnode in remote machine (if you run farlogin app in localhost, just open a terminal and run the flnode from command line for a simple test) as a deamon and give it the same node name via command line parameter.
* you will see the node online in farlogin web app after a couple of seconds.
* create a session key.
* use fladmin to connect remote machine terminal with using the session key and the node name.

if you are making simple test with an offline machine, you run farlogin, fladmin and flnode in localhost. this means you are accessing the same machine via nats.io messaging.

a session key can only be used once.

both fladmin and flnode give usage info when started without any paramatter.

useful links:
https://www.linusakesson.net/programming/tty/index.php
