This might help you with quick 'n dirty automated deployments.
Or at least give you ideas.

You could proxy it through your existing webserver.
Mind the auth though!

If you also need to restart a server, you might want to run this service as root.

```
sudo journalctl -u site-deployer
```