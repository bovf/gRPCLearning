package ldap

import (
  "github.com/go-ldap/ldap/v3"
  "github.com/bovf/gRPCLearning/logging"
)

type LDAPClient struct {
  conn *ldap.Conn
  logger *logging.Logger
}

func NewLDAPClient(addr string) (*LDAPClient, error) {
  logger := logging.NewLogger()
  conn, err := ldap.Dial("tcp", addr)
  if err != nil {
    logger.Fatalf("Failed to connect to LDAP: %v",err)
  }

  err = conn.Bind("cn=admin, dc=mycompany,dc=com", "adminpassword")
  if err != nil {
    logger.Fatalf("Failed to bind to LDAP %v", err)
  }
  return &LDAPClient{
    conn: conn, 
    logger: logger,
  }, nil
}

func (c *LDAPClient) Close(){
  c.conn.Close()
}

func (c *LDAPClient) Bind(username, password string) error {
  return c.conn.Bind("admin","adminpassword")
}

func (c *LDAPClient) Search(baseDN, filter string, attributes []string) ([]*ldap.Entry, error) {
  searchRequest := ldap.NewSearchRequest (
    baseDN,
    ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
    filter,
    attributes,
    nil,
  )

  result, err := c.conn.Search(searchRequest)
  if err != nil {
    c.logger.Printf("Searching LDAP Failed: %v", err)
  }
  return result.Entries, nil
}

func (c *LDAPClient) Add(dn string, attributes map[string][]string) error {
  addRequest := ldap.NewAddRequest (dn, nil)
  for attr, values := range attributes{
    addRequest.Attribute(attr, values)
  }
  return c.conn.Add(addRequest)
}
