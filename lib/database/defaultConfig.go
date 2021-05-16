package database

type DefaultMysqlDBCfg struct {
	
}

func (DefaultMysqlDBCfg)GetDBUserName() string{
	return ""
}
func (DefaultMysqlDBCfg)GetDBPassword() string{
	return ""
}
func (DefaultMysqlDBCfg)GetDBName() string{
	return ""
}
func (DefaultMysqlDBCfg)GetDBHost() string{
	return "127.0.0.1"
}
func (DefaultMysqlDBCfg)GetDBPort() uint{
	return 3306
}
