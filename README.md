# Farlogin

Provides access to a machine that does not have a static ip. It allows you to use the terminal screen of the remote machine.

It uses embedded nats.io server for communication.

I got this idea from this page: https://dev.to/napicella/linux-terminals-tty-pty-and-shell-192e

This method can be used to access nodes, especially during the development of iot projects.

It provides ease of use in case of a large number of nodes (hundreds), and can simplify asset management. It eliminates the need to use a gateway modem for a remote node.

**What else can be done:**
In a real installation, https should be used and all communication should be via TLS.

*flnode* (the agent copied to the node) can be conditionally compiled so that it runs only on a specific machine (by CPU id etc.).

In addition to TLS, two-way encryption can be used as an additional security measure. Login history can be kept, who connected to which machine from where, which commands they ran, etc.

In case you want to play with:

* Install nginx and use it as reverse proxy. Proxy pass all request to your domain to *farlogin* app. There is a sample proxy config in main folder.
* Import the *farlogin* database to mysql on the same machine *farlogin* app is running. Db dump file is in main folder.
* copy *fladmin* to the computer of admin user.
* Make the necessary changes in farlogin.ini file. App will work on localhost initially. if you want to run it on a server which has a public domain name you should change the ini file accordingly and don't forget to change the nats url in the app.go files for *farlogin*, *fladmin* and *flnode*.
* After running *farlogin* web app, login and create a node by giving it a name.
* Run *flnode* in remote machine (if you run *farlogin* app in localhost, just open a terminal and run the *flnode* from command line for a simple test) as a deamon and give it the same node name via command line parameter.
* You will see the node online in *farlogin* web app after a couple of seconds.
* Create a session key.
* Use *fladmin* to connect remote machine terminal with using the session key and the node name.

If you are making simple test with an offline machine, you run *farlogin*, *fladmin* and *flnode* in localhost. this means you are accessing the same machine via nats.io messaging.

A session key can only be used once.

Both *fladmin* and *flnode* give usage info when started without any paramatter.

Useful links:
https://www.linusakesson.net/programming/tty/index.php
