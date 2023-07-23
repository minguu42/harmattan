# このファイルは以下の想定で作成されている
# - mysqlコマンドを使用できること
# - スキーマを作成するSQLファイル/schema.sqlが存在すること
# - ダミーデータを投入するSQLファイル/import.sqlが存在すること

mysql -u root -e "CREATE DATABASE IF NOT EXISTS opepe_test"
mysql -u root opepe_test < /schema.sql
mysql -u root opepe_test < /import.sql
