USE fabricexplorer;
TRUNCATE TABLE `blockchain_chaincode`;
TRUNCATE TABLE `blockchain_chaincode_category_rel`;
TRUNCATE TABLE `blockchain_chaincode_file`;
TRUNCATE TABLE `blockchain_chaincode_file_content_log`;
TRUNCATE TABLE `blockchain_chaincode_org_rel`;
TRUNCATE TABLE `blockchain_chaincode_status_log`;
TRUNCATE TABLE `blockchain_channel`;
TRUNCATE TABLE `blockchain_node_order`;
TRUNCATE TABLE `blockchain_node_peer`;
TRUNCATE TABLE `blockchain_org`;
UPDATE blockchain_config SET `value`=0 WHERE id=1;