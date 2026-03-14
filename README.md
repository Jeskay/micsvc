## Description
Composed build creates a web service that allows operations on User database. All operations are logged and stored in elasticsearch. 
Successful and failed operations are transferred through kafka broker to websokcet server. Websocket server accepts connections, that can will receive all new operation events in a form of text messages.
## ENV
#### Individual service config
```
HOST="yourserverhost"
PORT="yourserverport"
TOKEN_EXPIRE=10 -auth token expires after (in hours)
SECRET_KEY="yoursecretkey" - auth token secret key
CONNECTION_TIMEOUT=5 - timeout (in seconds) to wait for clients response
KAFKA_ADDRESS="yourkafkaaddress" - address of kafka cluster
EVENT_TOPIC="yourtopicname" - name of kafka topic to store events
```
#### Docker compose config
```
KAFKA_CLUSTER_ID="yourclusterid" - id of your kafka cluster
NOTIFY_PORT=8080 - port of the notification service
WEB_PORT=8081 - port of the user web service
TOKEN_EXPIRE=6 - same as the Individual service config
AUTH_SECRET="yoursecretkey" - auth token secret key
NOTIFY_TOPIC="yourtopicname" - name of kafka topic to store events
GRAFANA_ADMIN_PASSWORD="yourpassword" - password for grafana admin panel
```
