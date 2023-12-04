FROM eclipse-temurin:17-alpine

ENV GFR=/opt/app
ENV GFZ=/glassfish7/glassfish

COPY glassfish.zip /tmp/glassfish.zip

RUN unzip -d "$GFR" /tmp/glassfish.zip && rm /tmp/glassfish.zip

# Allow Derby to start as daemon (used by some Java EE app, such as Pet Store)
RUN echo "grant { permission java.net.SocketPermission \"localhost:1527\", \"listen\"; };" >> $JAVA_HOME/lib/security/java.policy

# Secure GF installation with a password and authorize network access
ADD password_1.txt /tmp/password_1.txt
ADD password_2.txt /tmp/password_2.txt

RUN "$GFR/$GFZ/bin/asadmin" --user admin --passwordfile /tmp/password_1.txt change-admin-password --domain_name domain1; \
    "$GFR/$GFZ/bin/asadmin" start-domain domain1; \
    "$GFR/$GFZ/bin/asadmin" --user admin --passwordfile /tmp/password_2.txt enable-secure-admin; \
    "$GFR/$GFZ/bin/asadmin" stop-domain domain1;

# RUN rm /tmp/password_?.txt

# Copy the WAR file to the autodeploy directory
COPY ./target/firepit.war "$GFR/$GFZ/domains/domain1/autodeploy/"

RUN mkdir -p "$GFR/$GFZ/domains/domain1/config/apps" && \
    mkdir -p "$GFR/$GFZ/domains/domain1/autodeploy"

VOLUME "$GFR/$GFZ/domains/domain1/autodeploy"
VOLUME "$GFR/$GFZ/domains/domain1/config/apps"

WORKDIR "$GFR/$GFZ"

# Start server
CMD ["java", "-jar", "./lib/client/appserver-cli.jar", "start-domain", "--verbose"]

EXPOSE 4848 8080
