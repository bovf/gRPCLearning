#!/bin/bash

# Check if Docker is installed
docker --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "Docker is not installed. Please install Docker first."
  exit 1
fi

# Spin up an OpenLDAP server using Docker with modified security settings
docker run --name my-openldap -p 389:389 -p 636:636 \
  --env LDAP_ORGANISATION="My Company" \
  --env LDAP_DOMAIN="mycompany.com" \
  --env LDAP_ADMIN_PASSWORD="adminpassword" \
  --env LDAP_TLS=false \
  --env LDAP_READONLY_USER=true \
  --env LDAP_READONLY_USER_USERNAME=readonly \
  --env LDAP_READONLY_USER_PASSWORD=readonlypassword \
  --env LDAP_EXTRA_SCHEMAS=cosine,inetorgperson \
  --env LDAP_CONFIG_PASSWORD=configpassword \
  --detach osixia/openldap:1.5.0

if [ $? -eq 0 ]; then
  echo "OpenLDAP server is running on ports 389 (LDAP) and 636 (LDAPS)."
  echo "A readonly user 'readonly' has been created with password 'readonlypassword'."
  echo "TLS has been disabled for testing purposes."
else
  echo "Failed to start OpenLDAP server."
  exit 1
fi

