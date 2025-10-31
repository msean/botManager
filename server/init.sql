-- MySQL dump 10.13  Distrib 8.0.43, for macos15 (arm64)
--
-- Host: test.c96oo22ccbq0.ap-east-1.rds.amazonaws.com    Database: bot_manager
-- ------------------------------------------------------
-- Server version	8.0.41

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '';

--
-- Table structure for table `casbin_rule`
--

DROP TABLE IF EXISTS `casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v0` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v1` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v2` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v3` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v4` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v5` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=1360 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES (1200,'p','1','/api/createApi','POST','','',''),(1199,'p','1','/api/deleteApi','POST','','',''),(1194,'p','1','/api/deleteApisByIds','DELETE','','',''),(1191,'p','1','/api/enterSyncApi','POST','','',''),(1196,'p','1','/api/getAllApis','POST','','',''),(1195,'p','1','/api/getApiById','POST','','',''),(1192,'p','1','/api/getApiGroups','GET','','',''),(1197,'p','1','/api/getApiList','POST','','',''),(1190,'p','1','/api/ignoreApi','POST','','',''),(1193,'p','1','/api/syncApi','GET','','',''),(1198,'p','1','/api/updateApi','POST','','',''),(1086,'p','1','/attachmentCategory/addCategory','POST','','',''),(1085,'p','1','/attachmentCategory/deleteCategory','POST','','',''),(1087,'p','1','/attachmentCategory/getCategoryList','GET','','',''),(1189,'p','1','/authority/copyAuthority','POST','','',''),(1188,'p','1','/authority/createAuthority','POST','','',''),(1187,'p','1','/authority/deleteAuthority','POST','','',''),(1185,'p','1','/authority/getAuthorityList','POST','','',''),(1184,'p','1','/authority/setDataAuthority','POST','','',''),(1186,'p','1','/authority/updateAuthority','PUT','','',''),(1110,'p','1','/authorityBtn/canRemoveAuthorityBtn','POST','','',''),(1111,'p','1','/authorityBtn/getAuthorityBtn','POST','','',''),(1112,'p','1','/authorityBtn/setAuthorityBtn','POST','','',''),(1137,'p','1','/autoCode/addFunc','POST','','',''),(1145,'p','1','/autoCode/createPackage','POST','','',''),(1153,'p','1','/autoCode/createTemp','POST','','',''),(1142,'p','1','/autoCode/delPackage','POST','','',''),(1138,'p','1','/autoCode/delSysHistory','POST','','',''),(1151,'p','1','/autoCode/getColumn','GET','','',''),(1155,'p','1','/autoCode/getDB','GET','','',''),(1141,'p','1','/autoCode/getMeta','POST','','',''),(1143,'p','1','/autoCode/getPackage','POST','','',''),(1139,'p','1','/autoCode/getSysHistory','POST','','',''),(1154,'p','1','/autoCode/getTables','GET','','',''),(1144,'p','1','/autoCode/getTemplates','GET','','',''),(1150,'p','1','/autoCode/installPlugin','POST','','',''),(1148,'p','1','/autoCode/mcp','POST','','',''),(1146,'p','1','/autoCode/mcpList','POST','','',''),(1147,'p','1','/autoCode/mcpTest','POST','','',''),(1152,'p','1','/autoCode/preview','POST','','',''),(1149,'p','1','/autoCode/pubPlug','POST','','',''),(1140,'p','1','/autoCode/rollback','POST','','',''),(1071,'p','1','/bot_mgr/createBot','POST','','',''),(1070,'p','1','/bot_mgr/deleteBot','DELETE','','',''),(1069,'p','1','/bot_mgr/deleteBotByIds','DELETE','','',''),(1067,'p','1','/bot_mgr/findBot','GET','','',''),(1066,'p','1','/bot_mgr/getBotList','GET','','',''),(1068,'p','1','/bot_mgr/updateBot','PUT','','',''),(1077,'p','1','/bot_msg_mgr/createBotMsgMgr','POST','','',''),(1076,'p','1','/bot_msg_mgr/deleteBotMsgMgr','DELETE','','',''),(1075,'p','1','/bot_msg_mgr/deleteBotMsgMgrByIds','DELETE','','',''),(1073,'p','1','/bot_msg_mgr/findBotMsgMgr','GET','','',''),(1072,'p','1','/bot_msg_mgr/getBotMsgMgrList','GET','','',''),(1074,'p','1','/bot_msg_mgr/updateBotMsgMgr','PUT','','',''),(1182,'p','1','/casbin/getPolicyPathByAuthorityId','POST','','',''),(1183,'p','1','/casbin/updateCasbin','POST','','',''),(1158,'p','1','/customer/customer','DELETE','','',''),(1157,'p','1','/customer/customer','GET','','',''),(1159,'p','1','/customer/customer','POST','','',''),(1160,'p','1','/customer/customer','PUT','','',''),(1156,'p','1','/customer/customerList','GET','','',''),(1114,'p','1','/email/emailTest','POST','','',''),(1113,'p','1','/email/sendEmail','POST','','',''),(1171,'p','1','/fileUploadAndDownload/breakpointContinue','POST','','',''),(1170,'p','1','/fileUploadAndDownload/breakpointContinueFinish','POST','','',''),(1167,'p','1','/fileUploadAndDownload/deleteFile','POST','','',''),(1166,'p','1','/fileUploadAndDownload/editFileName','POST','','',''),(1172,'p','1','/fileUploadAndDownload/findFile','GET','','',''),(1165,'p','1','/fileUploadAndDownload/getFileList','POST','','',''),(1164,'p','1','/fileUploadAndDownload/importURL','POST','','',''),(1169,'p','1','/fileUploadAndDownload/removeChunk','POST','','',''),(1168,'p','1','/fileUploadAndDownload/upload','POST','','',''),(1100,'p','1','/info/createInfo','POST','','',''),(1099,'p','1','/info/deleteInfo','DELETE','','',''),(1098,'p','1','/info/deleteInfoByIds','DELETE','','',''),(1096,'p','1','/info/findInfo','GET','','',''),(1095,'p','1','/info/getInfoList','GET','','',''),(1097,'p','1','/info/updateInfo','PUT','','',''),(1212,'p','1','/jwt/jsonInBlacklist','POST','','',''),(1181,'p','1','/menu/addBaseMenu','POST','','',''),(1173,'p','1','/menu/addMenuAuthority','POST','','',''),(1179,'p','1','/menu/deleteBaseMenu','POST','','',''),(1177,'p','1','/menu/getBaseMenuById','POST','','',''),(1175,'p','1','/menu/getBaseMenuTree','POST','','',''),(1180,'p','1','/menu/getMenu','POST','','',''),(1174,'p','1','/menu/getMenuAuthority','POST','','',''),(1176,'p','1','/menu/getMenuList','POST','','',''),(1178,'p','1','/menu/updateBaseMenu','POST','','',''),(1116,'p','1','/simpleUploader/checkFileMd5','GET','','',''),(1115,'p','1','/simpleUploader/mergeFileMd5','GET','','',''),(1117,'p','1','/simpleUploader/upload','POST','','',''),(1127,'p','1','/sysDictionary/createSysDictionary','POST','','',''),(1126,'p','1','/sysDictionary/deleteSysDictionary','DELETE','','',''),(1124,'p','1','/sysDictionary/findSysDictionary','GET','','',''),(1123,'p','1','/sysDictionary/getSysDictionaryList','GET','','',''),(1125,'p','1','/sysDictionary/updateSysDictionary','PUT','','',''),(1135,'p','1','/sysDictionaryDetail/createSysDictionaryDetail','POST','','',''),(1134,'p','1','/sysDictionaryDetail/deleteSysDictionaryDetail','DELETE','','',''),(1133,'p','1','/sysDictionaryDetail/findSysDictionaryDetail','GET','','',''),(1129,'p','1','/sysDictionaryDetail/getDictionaryDetailsByParent','GET','','',''),(1128,'p','1','/sysDictionaryDetail/getDictionaryPath','GET','','',''),(1131,'p','1','/sysDictionaryDetail/getDictionaryTreeList','GET','','',''),(1130,'p','1','/sysDictionaryDetail/getDictionaryTreeListByType','GET','','',''),(1132,'p','1','/sysDictionaryDetail/getSysDictionaryDetailList','GET','','',''),(1136,'p','1','/sysDictionaryDetail/updateSysDictionaryDetail','PUT','','',''),(1109,'p','1','/sysExportTemplate/createSysExportTemplate','POST','','',''),(1108,'p','1','/sysExportTemplate/deleteSysExportTemplate','DELETE','','',''),(1107,'p','1','/sysExportTemplate/deleteSysExportTemplateByIds','DELETE','','',''),(1103,'p','1','/sysExportTemplate/exportExcel','GET','','',''),(1102,'p','1','/sysExportTemplate/exportTemplate','GET','','',''),(1105,'p','1','/sysExportTemplate/findSysExportTemplate','GET','','',''),(1104,'p','1','/sysExportTemplate/getSysExportTemplateList','GET','','',''),(1101,'p','1','/sysExportTemplate/importExcel','POST','','',''),(1106,'p','1','/sysExportTemplate/updateSysExportTemplate','PUT','','',''),(1122,'p','1','/sysOperationRecord/createSysOperationRecord','POST','','',''),(1119,'p','1','/sysOperationRecord/deleteSysOperationRecord','DELETE','','',''),(1118,'p','1','/sysOperationRecord/deleteSysOperationRecordByIds','DELETE','','',''),(1121,'p','1','/sysOperationRecord/findSysOperationRecord','GET','','',''),(1120,'p','1','/sysOperationRecord/getSysOperationRecordList','GET','','',''),(1094,'p','1','/sysParams/createSysParams','POST','','',''),(1093,'p','1','/sysParams/deleteSysParams','DELETE','','',''),(1092,'p','1','/sysParams/deleteSysParamsByIds','DELETE','','',''),(1090,'p','1','/sysParams/findSysParams','GET','','',''),(1088,'p','1','/sysParams/getSysParam','GET','','',''),(1089,'p','1','/sysParams/getSysParamsList','GET','','',''),(1091,'p','1','/sysParams/updateSysParams','PUT','','',''),(1163,'p','1','/system/getServerInfo','POST','','',''),(1162,'p','1','/system/getSystemConfig','POST','','',''),(1161,'p','1','/system/setSystemConfig','POST','','',''),(1079,'p','1','/sysVersion/deleteSysVersion','DELETE','','',''),(1078,'p','1','/sysVersion/deleteSysVersionByIds','DELETE','','',''),(1082,'p','1','/sysVersion/downloadVersionJson','GET','','',''),(1081,'p','1','/sysVersion/exportVersion','POST','','',''),(1084,'p','1','/sysVersion/findSysVersion','GET','','',''),(1083,'p','1','/sysVersion/getSysVersionList','GET','','',''),(1080,'p','1','/sysVersion/importVersion','POST','','',''),(1210,'p','1','/user/admin_register','POST','','',''),(1204,'p','1','/user/changePassword','POST','','',''),(1211,'p','1','/user/deleteUser','DELETE','','',''),(1206,'p','1','/user/getUserInfo','GET','','',''),(1209,'p','1','/user/getUserList','POST','','',''),(1202,'p','1','/user/resetPassword','POST','','',''),(1207,'p','1','/user/setSelfInfo','PUT','','',''),(1201,'p','1','/user/setSelfSetting','PUT','','',''),(1205,'p','1','/user/setUserAuthorities','POST','','',''),(1203,'p','1','/user/setUserAuthority','POST','','',''),(1208,'p','1','/user/setUserInfo','PUT','','',''),(1347,'p','888','/api/createApi','POST','','',''),(1346,'p','888','/api/deleteApi','POST','','',''),(1341,'p','888','/api/deleteApisByIds','DELETE','','',''),(1338,'p','888','/api/enterSyncApi','POST','','',''),(1343,'p','888','/api/getAllApis','POST','','',''),(1342,'p','888','/api/getApiById','POST','','',''),(1339,'p','888','/api/getApiGroups','GET','','',''),(1344,'p','888','/api/getApiList','POST','','',''),(1337,'p','888','/api/ignoreApi','POST','','',''),(1340,'p','888','/api/syncApi','GET','','',''),(1345,'p','888','/api/updateApi','POST','','',''),(1233,'p','888','/attachmentCategory/addCategory','POST','','',''),(1232,'p','888','/attachmentCategory/deleteCategory','POST','','',''),(1234,'p','888','/attachmentCategory/getCategoryList','GET','','',''),(1336,'p','888','/authority/copyAuthority','POST','','',''),(1335,'p','888','/authority/createAuthority','POST','','',''),(1334,'p','888','/authority/deleteAuthority','POST','','',''),(1332,'p','888','/authority/getAuthorityList','POST','','',''),(1331,'p','888','/authority/setDataAuthority','POST','','',''),(1333,'p','888','/authority/updateAuthority','PUT','','',''),(1257,'p','888','/authorityBtn/canRemoveAuthorityBtn','POST','','',''),(1258,'p','888','/authorityBtn/getAuthorityBtn','POST','','',''),(1259,'p','888','/authorityBtn/setAuthorityBtn','POST','','',''),(1284,'p','888','/autoCode/addFunc','POST','','',''),(1292,'p','888','/autoCode/createPackage','POST','','',''),(1300,'p','888','/autoCode/createTemp','POST','','',''),(1289,'p','888','/autoCode/delPackage','POST','','',''),(1285,'p','888','/autoCode/delSysHistory','POST','','',''),(1298,'p','888','/autoCode/getColumn','GET','','',''),(1302,'p','888','/autoCode/getDB','GET','','',''),(1288,'p','888','/autoCode/getMeta','POST','','',''),(1290,'p','888','/autoCode/getPackage','POST','','',''),(1286,'p','888','/autoCode/getSysHistory','POST','','',''),(1301,'p','888','/autoCode/getTables','GET','','',''),(1291,'p','888','/autoCode/getTemplates','GET','','',''),(1297,'p','888','/autoCode/installPlugin','POST','','',''),(1295,'p','888','/autoCode/mcp','POST','','',''),(1293,'p','888','/autoCode/mcpList','POST','','',''),(1294,'p','888','/autoCode/mcpTest','POST','','',''),(1299,'p','888','/autoCode/preview','POST','','',''),(1296,'p','888','/autoCode/pubPlug','POST','','',''),(1287,'p','888','/autoCode/rollback','POST','','',''),(1218,'p','888','/bot_mgr/createBot','POST','','',''),(1217,'p','888','/bot_mgr/deleteBot','DELETE','','',''),(1216,'p','888','/bot_mgr/deleteBotByIds','DELETE','','',''),(1214,'p','888','/bot_mgr/findBot','GET','','',''),(1213,'p','888','/bot_mgr/getBotList','GET','','',''),(1215,'p','888','/bot_mgr/updateBot','PUT','','',''),(1224,'p','888','/bot_msg_mgr/createBotMsgMgr','POST','','',''),(1223,'p','888','/bot_msg_mgr/deleteBotMsgMgr','DELETE','','',''),(1222,'p','888','/bot_msg_mgr/deleteBotMsgMgrByIds','DELETE','','',''),(1220,'p','888','/bot_msg_mgr/findBotMsgMgr','GET','','',''),(1219,'p','888','/bot_msg_mgr/getBotMsgMgrList','GET','','',''),(1221,'p','888','/bot_msg_mgr/updateBotMsgMgr','PUT','','',''),(1329,'p','888','/casbin/getPolicyPathByAuthorityId','POST','','',''),(1330,'p','888','/casbin/updateCasbin','POST','','',''),(1305,'p','888','/customer/customer','DELETE','','',''),(1304,'p','888','/customer/customer','GET','','',''),(1306,'p','888','/customer/customer','POST','','',''),(1307,'p','888','/customer/customer','PUT','','',''),(1303,'p','888','/customer/customerList','GET','','',''),(1261,'p','888','/email/emailTest','POST','','',''),(1260,'p','888','/email/sendEmail','POST','','',''),(1318,'p','888','/fileUploadAndDownload/breakpointContinue','POST','','',''),(1317,'p','888','/fileUploadAndDownload/breakpointContinueFinish','POST','','',''),(1314,'p','888','/fileUploadAndDownload/deleteFile','POST','','',''),(1313,'p','888','/fileUploadAndDownload/editFileName','POST','','',''),(1319,'p','888','/fileUploadAndDownload/findFile','GET','','',''),(1312,'p','888','/fileUploadAndDownload/getFileList','POST','','',''),(1311,'p','888','/fileUploadAndDownload/importURL','POST','','',''),(1316,'p','888','/fileUploadAndDownload/removeChunk','POST','','',''),(1315,'p','888','/fileUploadAndDownload/upload','POST','','',''),(1247,'p','888','/info/createInfo','POST','','',''),(1246,'p','888','/info/deleteInfo','DELETE','','',''),(1245,'p','888','/info/deleteInfoByIds','DELETE','','',''),(1243,'p','888','/info/findInfo','GET','','',''),(1242,'p','888','/info/getInfoList','GET','','',''),(1244,'p','888','/info/updateInfo','PUT','','',''),(1359,'p','888','/jwt/jsonInBlacklist','POST','','',''),(1328,'p','888','/menu/addBaseMenu','POST','','',''),(1320,'p','888','/menu/addMenuAuthority','POST','','',''),(1326,'p','888','/menu/deleteBaseMenu','POST','','',''),(1324,'p','888','/menu/getBaseMenuById','POST','','',''),(1322,'p','888','/menu/getBaseMenuTree','POST','','',''),(1327,'p','888','/menu/getMenu','POST','','',''),(1321,'p','888','/menu/getMenuAuthority','POST','','',''),(1323,'p','888','/menu/getMenuList','POST','','',''),(1325,'p','888','/menu/updateBaseMenu','POST','','',''),(1263,'p','888','/simpleUploader/checkFileMd5','GET','','',''),(1262,'p','888','/simpleUploader/mergeFileMd5','GET','','',''),(1264,'p','888','/simpleUploader/upload','POST','','',''),(1274,'p','888','/sysDictionary/createSysDictionary','POST','','',''),(1273,'p','888','/sysDictionary/deleteSysDictionary','DELETE','','',''),(1271,'p','888','/sysDictionary/findSysDictionary','GET','','',''),(1270,'p','888','/sysDictionary/getSysDictionaryList','GET','','',''),(1272,'p','888','/sysDictionary/updateSysDictionary','PUT','','',''),(1282,'p','888','/sysDictionaryDetail/createSysDictionaryDetail','POST','','',''),(1281,'p','888','/sysDictionaryDetail/deleteSysDictionaryDetail','DELETE','','',''),(1280,'p','888','/sysDictionaryDetail/findSysDictionaryDetail','GET','','',''),(1276,'p','888','/sysDictionaryDetail/getDictionaryDetailsByParent','GET','','',''),(1275,'p','888','/sysDictionaryDetail/getDictionaryPath','GET','','',''),(1278,'p','888','/sysDictionaryDetail/getDictionaryTreeList','GET','','',''),(1277,'p','888','/sysDictionaryDetail/getDictionaryTreeListByType','GET','','',''),(1279,'p','888','/sysDictionaryDetail/getSysDictionaryDetailList','GET','','',''),(1283,'p','888','/sysDictionaryDetail/updateSysDictionaryDetail','PUT','','',''),(1256,'p','888','/sysExportTemplate/createSysExportTemplate','POST','','',''),(1255,'p','888','/sysExportTemplate/deleteSysExportTemplate','DELETE','','',''),(1254,'p','888','/sysExportTemplate/deleteSysExportTemplateByIds','DELETE','','',''),(1250,'p','888','/sysExportTemplate/exportExcel','GET','','',''),(1249,'p','888','/sysExportTemplate/exportTemplate','GET','','',''),(1252,'p','888','/sysExportTemplate/findSysExportTemplate','GET','','',''),(1251,'p','888','/sysExportTemplate/getSysExportTemplateList','GET','','',''),(1248,'p','888','/sysExportTemplate/importExcel','POST','','',''),(1253,'p','888','/sysExportTemplate/updateSysExportTemplate','PUT','','',''),(1269,'p','888','/sysOperationRecord/createSysOperationRecord','POST','','',''),(1266,'p','888','/sysOperationRecord/deleteSysOperationRecord','DELETE','','',''),(1265,'p','888','/sysOperationRecord/deleteSysOperationRecordByIds','DELETE','','',''),(1268,'p','888','/sysOperationRecord/findSysOperationRecord','GET','','',''),(1267,'p','888','/sysOperationRecord/getSysOperationRecordList','GET','','',''),(1241,'p','888','/sysParams/createSysParams','POST','','',''),(1240,'p','888','/sysParams/deleteSysParams','DELETE','','',''),(1239,'p','888','/sysParams/deleteSysParamsByIds','DELETE','','',''),(1237,'p','888','/sysParams/findSysParams','GET','','',''),(1235,'p','888','/sysParams/getSysParam','GET','','',''),(1236,'p','888','/sysParams/getSysParamsList','GET','','',''),(1238,'p','888','/sysParams/updateSysParams','PUT','','',''),(1310,'p','888','/system/getServerInfo','POST','','',''),(1309,'p','888','/system/getSystemConfig','POST','','',''),(1308,'p','888','/system/setSystemConfig','POST','','',''),(1226,'p','888','/sysVersion/deleteSysVersion','DELETE','','',''),(1225,'p','888','/sysVersion/deleteSysVersionByIds','DELETE','','',''),(1229,'p','888','/sysVersion/downloadVersionJson','GET','','',''),(1228,'p','888','/sysVersion/exportVersion','POST','','',''),(1231,'p','888','/sysVersion/findSysVersion','GET','','',''),(1230,'p','888','/sysVersion/getSysVersionList','GET','','',''),(1227,'p','888','/sysVersion/importVersion','POST','','',''),(1357,'p','888','/user/admin_register','POST','','',''),(1351,'p','888','/user/changePassword','POST','','',''),(1358,'p','888','/user/deleteUser','DELETE','','',''),(1353,'p','888','/user/getUserInfo','GET','','',''),(1356,'p','888','/user/getUserList','POST','','',''),(1349,'p','888','/user/resetPassword','POST','','',''),(1354,'p','888','/user/setSelfInfo','PUT','','',''),(1348,'p','888','/user/setSelfSetting','PUT','','',''),(1352,'p','888','/user/setUserAuthorities','POST','','',''),(1350,'p','888','/user/setUserAuthority','POST','','',''),(1355,'p','888','/user/setUserInfo','PUT','','',''),(139,'p','8881','/api/createApi','POST','','',''),(142,'p','8881','/api/deleteApi','POST','','',''),(144,'p','8881','/api/getAllApis','POST','','',''),(141,'p','8881','/api/getApiById','POST','','',''),(140,'p','8881','/api/getApiList','POST','','',''),(143,'p','8881','/api/updateApi','POST','','',''),(145,'p','8881','/authority/createAuthority','POST','','',''),(146,'p','8881','/authority/deleteAuthority','POST','','',''),(147,'p','8881','/authority/getAuthorityList','POST','','',''),(148,'p','8881','/authority/setDataAuthority','POST','','',''),(167,'p','8881','/casbin/getPolicyPathByAuthorityId','POST','','',''),(166,'p','8881','/casbin/updateCasbin','POST','','',''),(173,'p','8881','/customer/customer','DELETE','','',''),(174,'p','8881','/customer/customer','GET','','',''),(171,'p','8881','/customer/customer','POST','','',''),(172,'p','8881','/customer/customer','PUT','','',''),(175,'p','8881','/customer/customerList','GET','','',''),(163,'p','8881','/fileUploadAndDownload/deleteFile','POST','','',''),(164,'p','8881','/fileUploadAndDownload/editFileName','POST','','',''),(162,'p','8881','/fileUploadAndDownload/getFileList','POST','','',''),(165,'p','8881','/fileUploadAndDownload/importURL','POST','','',''),(161,'p','8881','/fileUploadAndDownload/upload','POST','','',''),(168,'p','8881','/jwt/jsonInBlacklist','POST','','',''),(151,'p','8881','/menu/addBaseMenu','POST','','',''),(153,'p','8881','/menu/addMenuAuthority','POST','','',''),(155,'p','8881','/menu/deleteBaseMenu','POST','','',''),(157,'p','8881','/menu/getBaseMenuById','POST','','',''),(152,'p','8881','/menu/getBaseMenuTree','POST','','',''),(149,'p','8881','/menu/getMenu','POST','','',''),(154,'p','8881','/menu/getMenuAuthority','POST','','',''),(150,'p','8881','/menu/getMenuList','POST','','',''),(156,'p','8881','/menu/updateBaseMenu','POST','','',''),(169,'p','8881','/system/getSystemConfig','POST','','',''),(170,'p','8881','/system/setSystemConfig','POST','','',''),(138,'p','8881','/user/admin_register','POST','','',''),(158,'p','8881','/user/changePassword','POST','','',''),(176,'p','8881','/user/getUserInfo','GET','','',''),(159,'p','8881','/user/getUserList','POST','','',''),(160,'p','8881','/user/setUserAuthority','POST','','','');
/*!40000 ALTER TABLE `casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `exa_attachment_category`
--

DROP TABLE IF EXISTS `exa_attachment_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `exa_attachment_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类名称',
  `pid` bigint DEFAULT '0' COMMENT '父节点ID',
  PRIMARY KEY (`id`),
  KEY `idx_exa_attachment_category_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exa_attachment_category`
--

LOCK TABLES `exa_attachment_category` WRITE;
/*!40000 ALTER TABLE `exa_attachment_category` DISABLE KEYS */;
/*!40000 ALTER TABLE `exa_attachment_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `exa_customers`
--

DROP TABLE IF EXISTS `exa_customers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `exa_customers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `customer_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户名',
  `customer_phone_data` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户手机号',
  `sys_user_id` bigint unsigned DEFAULT NULL COMMENT '管理ID',
  `sys_user_authority_id` bigint unsigned DEFAULT NULL COMMENT '管理角色ID',
  PRIMARY KEY (`id`),
  KEY `idx_exa_customers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exa_customers`
--

LOCK TABLES `exa_customers` WRITE;
/*!40000 ALTER TABLE `exa_customers` DISABLE KEYS */;
/*!40000 ALTER TABLE `exa_customers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `exa_file_chunks`
--

DROP TABLE IF EXISTS `exa_file_chunks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `exa_file_chunks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `exa_file_id` bigint unsigned DEFAULT NULL,
  `file_chunk_number` bigint DEFAULT NULL,
  `file_chunk_path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_exa_file_chunks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exa_file_chunks`
--

LOCK TABLES `exa_file_chunks` WRITE;
/*!40000 ALTER TABLE `exa_file_chunks` DISABLE KEYS */;
/*!40000 ALTER TABLE `exa_file_chunks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `exa_file_upload_and_downloads`
--

DROP TABLE IF EXISTS `exa_file_upload_and_downloads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `exa_file_upload_and_downloads` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `class_id` bigint DEFAULT '0' COMMENT '分类id',
  `url` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件地址',
  `tag` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件标签',
  `key` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '编号',
  PRIMARY KEY (`id`),
  KEY `idx_exa_file_upload_and_downloads_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exa_file_upload_and_downloads`
--

LOCK TABLES `exa_file_upload_and_downloads` WRITE;
/*!40000 ALTER TABLE `exa_file_upload_and_downloads` DISABLE KEYS */;
INSERT INTO `exa_file_upload_and_downloads` VALUES (1,'2025-10-30 15:17:00.448','2025-10-30 15:17:00.448',NULL,'10.png',0,'https://qmplusimg.henrongyi.top/gvalogo.png','png','158787308910.png'),(2,'2025-10-30 15:17:00.448','2025-10-30 15:17:00.448',NULL,'logo.png',0,'https://qmplusimg.henrongyi.top/1576554439myAvatar.png','png','1587973709logo.png');
/*!40000 ALTER TABLE `exa_file_upload_and_downloads` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `exa_files`
--

DROP TABLE IF EXISTS `exa_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `exa_files` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `file_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_md5` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chunk_total` bigint DEFAULT NULL,
  `is_finish` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_exa_files_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exa_files`
--

LOCK TABLES `exa_files` WRITE;
/*!40000 ALTER TABLE `exa_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `exa_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gva_announcements_info`
--

DROP TABLE IF EXISTS `gva_announcements_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gva_announcements_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '公告标题',
  `content` text COLLATE utf8mb4_unicode_ci COMMENT '公告内容',
  `user_id` bigint DEFAULT NULL COMMENT '发布者',
  `attachments` json DEFAULT NULL COMMENT '相关附件',
  PRIMARY KEY (`id`),
  KEY `idx_gva_announcements_info_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gva_announcements_info`
--

LOCK TABLES `gva_announcements_info` WRITE;
/*!40000 ALTER TABLE `gva_announcements_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `gva_announcements_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jwt_blacklists`
--

DROP TABLE IF EXISTS `jwt_blacklists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jwt_blacklists` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `jwt` text COLLATE utf8mb4_unicode_ci COMMENT 'jwt',
  PRIMARY KEY (`id`),
  KEY `idx_jwt_blacklists_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jwt_blacklists`
--

LOCK TABLES `jwt_blacklists` WRITE;
/*!40000 ALTER TABLE `jwt_blacklists` DISABLE KEYS */;
/*!40000 ALTER TABLE `jwt_blacklists` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_apis`
--

DROP TABLE IF EXISTS `sys_apis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_apis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api路径',
  `description` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api中文描述',
  `api_group` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api组',
  `method` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'POST' COMMENT '方法',
  PRIMARY KEY (`id`),
  KEY `idx_sys_apis_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=154 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_apis`
--

LOCK TABLES `sys_apis` WRITE;
/*!40000 ALTER TABLE `sys_apis` DISABLE KEYS */;
INSERT INTO `sys_apis` VALUES (1,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/jwt/jsonInBlacklist','jwt加入黑名单(退出，必选)','jwt','POST'),(2,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/deleteUser','删除用户','系统用户','DELETE'),(3,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/admin_register','用户注册','系统用户','POST'),(4,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/getUserList','获取用户列表','系统用户','POST'),(5,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/setUserInfo','设置用户信息','系统用户','PUT'),(6,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/setSelfInfo','设置自身信息(必选)','系统用户','PUT'),(7,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/getUserInfo','获取自身信息(必选)','系统用户','GET'),(8,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/setUserAuthorities','设置权限组','系统用户','POST'),(9,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/changePassword','修改密码（建议选择)','系统用户','POST'),(10,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/setUserAuthority','修改用户角色(必选)','系统用户','POST'),(11,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/resetPassword','重置用户密码','系统用户','POST'),(12,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/user/setSelfSetting','用户界面配置','系统用户','PUT'),(13,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/createApi','创建api','api','POST'),(14,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/deleteApi','删除Api','api','POST'),(15,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/updateApi','更新Api','api','POST'),(16,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/getApiList','获取api列表','api','POST'),(17,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/getAllApis','获取所有api','api','POST'),(18,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/getApiById','获取api详细信息','api','POST'),(19,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/deleteApisByIds','批量删除api','api','DELETE'),(20,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/syncApi','获取待同步API','api','GET'),(21,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/getApiGroups','获取路由组','api','GET'),(22,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/enterSyncApi','确认同步API','api','POST'),(23,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/api/ignoreApi','忽略API','api','POST'),(24,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authority/copyAuthority','拷贝角色','角色','POST'),(25,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authority/createAuthority','创建角色','角色','POST'),(26,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authority/deleteAuthority','删除角色','角色','POST'),(27,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authority/updateAuthority','更新角色信息','角色','PUT'),(28,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authority/getAuthorityList','获取角色列表','角色','POST'),(29,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authority/setDataAuthority','设置角色资源权限','角色','POST'),(30,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/casbin/updateCasbin','更改角色api权限','casbin','POST'),(31,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/casbin/getPolicyPathByAuthorityId','获取权限列表','casbin','POST'),(32,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/addBaseMenu','新增菜单','菜单','POST'),(33,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/getMenu','获取菜单树(必选)','菜单','POST'),(34,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/deleteBaseMenu','删除菜单','菜单','POST'),(35,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/updateBaseMenu','更新菜单','菜单','POST'),(36,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/getBaseMenuById','根据id获取菜单','菜单','POST'),(37,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/getMenuList','分页获取基础menu列表','菜单','POST'),(38,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/getBaseMenuTree','获取用户动态路由','菜单','POST'),(39,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/getMenuAuthority','获取指定角色menu','菜单','POST'),(40,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/menu/addMenuAuthority','增加menu和角色关联关系','菜单','POST'),(41,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/findFile','寻找目标文件（秒传）','分片上传','GET'),(42,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/breakpointContinue','断点续传','分片上传','POST'),(43,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/breakpointContinueFinish','断点续传完成','分片上传','POST'),(44,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/removeChunk','上传完成移除文件','分片上传','POST'),(45,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/upload','文件上传（建议选择）','文件上传与下载','POST'),(46,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/deleteFile','删除文件','文件上传与下载','POST'),(47,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/editFileName','文件名或者备注编辑','文件上传与下载','POST'),(48,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/getFileList','获取上传文件列表','文件上传与下载','POST'),(49,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/fileUploadAndDownload/importURL','导入URL','文件上传与下载','POST'),(50,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/system/getServerInfo','获取服务器信息','系统服务','POST'),(51,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/system/getSystemConfig','获取配置文件内容','系统服务','POST'),(52,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/system/setSystemConfig','设置配置文件内容','系统服务','POST'),(53,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/customer/customer','更新客户','客户','PUT'),(54,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/customer/customer','创建客户','客户','POST'),(55,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/customer/customer','删除客户','客户','DELETE'),(56,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/customer/customer','获取单一客户','客户','GET'),(57,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/customer/customerList','获取客户列表','客户','GET'),(58,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getDB','获取所有数据库','代码生成器','GET'),(59,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getTables','获取数据库表','代码生成器','GET'),(60,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/createTemp','自动化代码','代码生成器','POST'),(61,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/preview','预览自动化代码','代码生成器','POST'),(62,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getColumn','获取所选table的所有字段','代码生成器','GET'),(63,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/installPlugin','安装插件','代码生成器','POST'),(64,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/pubPlug','打包插件','代码生成器','POST'),(65,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/mcp','自动生成 MCP Tool 模板','代码生成器','POST'),(66,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/mcpTest','MCP Tool 测试','代码生成器','POST'),(67,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/mcpList','获取 MCP ToolList','代码生成器','POST'),(68,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/createPackage','配置模板','模板配置','POST'),(69,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getTemplates','获取模板文件','模板配置','GET'),(70,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getPackage','获取所有模板','模板配置','POST'),(71,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/delPackage','删除模板','模板配置','POST'),(72,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getMeta','获取meta信息','代码生成器历史','POST'),(73,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/rollback','回滚自动生成代码','代码生成器历史','POST'),(74,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/getSysHistory','查询回滚记录','代码生成器历史','POST'),(75,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/delSysHistory','删除回滚记录','代码生成器历史','POST'),(76,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/autoCode/addFunc','增加模板方法','代码生成器历史','POST'),(77,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/updateSysDictionaryDetail','更新字典内容','系统字典详情','PUT'),(78,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/createSysDictionaryDetail','新增字典内容','系统字典详情','POST'),(79,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/deleteSysDictionaryDetail','删除字典内容','系统字典详情','DELETE'),(80,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/findSysDictionaryDetail','根据ID获取字典内容','系统字典详情','GET'),(81,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/getSysDictionaryDetailList','获取字典内容列表','系统字典详情','GET'),(82,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/getDictionaryTreeList','获取字典数列表','系统字典详情','GET'),(83,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/getDictionaryTreeListByType','根据分类获取字典数列表','系统字典详情','GET'),(84,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/getDictionaryDetailsByParent','根据父级ID获取字典详情','系统字典详情','GET'),(85,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionaryDetail/getDictionaryPath','获取字典详情的完整路径','系统字典详情','GET'),(86,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionary/createSysDictionary','新增字典','系统字典','POST'),(87,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionary/deleteSysDictionary','删除字典','系统字典','DELETE'),(88,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionary/updateSysDictionary','更新字典','系统字典','PUT'),(89,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionary/findSysDictionary','根据ID获取字典（建议选择）','系统字典','GET'),(90,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysDictionary/getSysDictionaryList','获取字典列表','系统字典','GET'),(91,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysOperationRecord/createSysOperationRecord','新增操作记录','操作记录','POST'),(92,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysOperationRecord/findSysOperationRecord','根据ID获取操作记录','操作记录','GET'),(93,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysOperationRecord/getSysOperationRecordList','获取操作记录列表','操作记录','GET'),(94,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysOperationRecord/deleteSysOperationRecord','删除操作记录','操作记录','DELETE'),(95,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysOperationRecord/deleteSysOperationRecordByIds','批量删除操作历史','操作记录','DELETE'),(96,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/simpleUploader/upload','插件版分片上传','断点续传(插件版)','POST'),(97,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/simpleUploader/checkFileMd5','文件完整度验证','断点续传(插件版)','GET'),(98,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/simpleUploader/mergeFileMd5','上传完成合并文件','断点续传(插件版)','GET'),(99,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/email/emailTest','发送测试邮件','email','POST'),(100,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/email/sendEmail','发送邮件','email','POST'),(101,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authorityBtn/setAuthorityBtn','设置按钮权限','按钮权限','POST'),(102,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authorityBtn/getAuthorityBtn','获取已有按钮权限','按钮权限','POST'),(103,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/authorityBtn/canRemoveAuthorityBtn','删除按钮','按钮权限','POST'),(104,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/createSysExportTemplate','新增导出模板','导出模板','POST'),(105,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/deleteSysExportTemplate','删除导出模板','导出模板','DELETE'),(106,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/deleteSysExportTemplateByIds','批量删除导出模板','导出模板','DELETE'),(107,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/updateSysExportTemplate','更新导出模板','导出模板','PUT'),(108,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/findSysExportTemplate','根据ID获取导出模板','导出模板','GET'),(109,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/getSysExportTemplateList','获取导出模板列表','导出模板','GET'),(110,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/exportExcel','导出Excel','导出模板','GET'),(111,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/exportTemplate','下载模板','导出模板','GET'),(112,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysExportTemplate/importExcel','导入Excel','导出模板','POST'),(113,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/info/createInfo','新建公告','公告','POST'),(114,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/info/deleteInfo','删除公告','公告','DELETE'),(115,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/info/deleteInfoByIds','批量删除公告','公告','DELETE'),(116,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/info/updateInfo','更新公告','公告','PUT'),(117,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/info/findInfo','根据ID获取公告','公告','GET'),(118,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/info/getInfoList','获取公告列表','公告','GET'),(119,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/createSysParams','新建参数','参数管理','POST'),(120,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/deleteSysParams','删除参数','参数管理','DELETE'),(121,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/deleteSysParamsByIds','批量删除参数','参数管理','DELETE'),(122,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/updateSysParams','更新参数','参数管理','PUT'),(123,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/findSysParams','根据ID获取参数','参数管理','GET'),(124,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/getSysParamsList','获取参数列表','参数管理','GET'),(125,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysParams/getSysParam','获取参数列表','参数管理','GET'),(126,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/attachmentCategory/getCategoryList','分类列表','媒体库分类','GET'),(127,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/attachmentCategory/addCategory','添加/编辑分类','媒体库分类','POST'),(128,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/attachmentCategory/deleteCategory','删除分类','媒体库分类','POST'),(129,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/findSysVersion','获取单一版本','版本控制','GET'),(130,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/getSysVersionList','获取版本列表','版本控制','GET'),(131,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/downloadVersionJson','下载版本json','版本控制','GET'),(132,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/exportVersion','创建版本','版本控制','POST'),(133,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/importVersion','同步版本','版本控制','POST'),(134,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/deleteSysVersion','删除版本','版本控制','DELETE'),(135,'2025-10-30 15:16:44.434','2025-10-30 15:16:44.434',NULL,'/sysVersion/deleteSysVersionByIds','批量删除版本','版本控制','DELETE'),(136,'2025-10-30 18:26:09.788','2025-10-30 18:26:09.788','2025-10-30 20:51:36.745','/botList/createBot','新增表','机器人列表','POST'),(137,'2025-10-30 18:26:10.085','2025-10-30 18:26:10.085','2025-10-30 20:51:36.745','/botList/deleteBot','删除表','机器人列表','DELETE'),(138,'2025-10-30 18:26:10.384','2025-10-30 18:26:10.384','2025-10-30 20:51:19.891','/botList/deleteBotByIds','批量删除表','机器人列表','DELETE'),(139,'2025-10-30 18:26:10.781','2025-10-30 18:26:10.781','2025-10-30 20:51:19.891','/botList/updateBot','更新表','机器人列表','PUT'),(140,'2025-10-30 18:26:11.083','2025-10-30 18:26:11.083','2025-10-30 20:51:19.891','/botList/findBot','根据ID获取表','机器人列表','GET'),(141,'2025-10-30 18:26:11.381','2025-10-30 18:27:43.864','2025-10-30 20:51:19.891','/botList/getBotList','获取表列表','表','GET'),(142,'2025-10-30 20:42:27.394','2025-10-30 20:42:27.394',NULL,'/bot_msg_mgr/createBotMsgMgr','新增机器人消息管理','机器人消息管理','POST'),(143,'2025-10-30 20:42:27.694','2025-10-30 20:42:27.694',NULL,'/bot_msg_mgr/deleteBotMsgMgr','删除机器人消息管理','机器人消息管理','DELETE'),(144,'2025-10-30 20:42:27.979','2025-10-30 20:42:27.979',NULL,'/bot_msg_mgr/deleteBotMsgMgrByIds','批量删除机器人消息管理','机器人消息管理','DELETE'),(145,'2025-10-30 20:42:28.265','2025-10-30 20:42:28.265',NULL,'/bot_msg_mgr/updateBotMsgMgr','更新机器人消息管理','机器人消息管理','PUT'),(146,'2025-10-30 20:42:28.553','2025-10-30 20:42:28.553',NULL,'/bot_msg_mgr/findBotMsgMgr','根据ID获取机器人消息管理','机器人消息管理','GET'),(147,'2025-10-30 20:42:28.842','2025-10-30 20:42:28.842',NULL,'/bot_msg_mgr/getBotMsgMgrList','获取机器人消息管理列表','机器人消息管理','GET'),(148,'2025-10-30 21:27:46.321','2025-10-30 21:27:46.321',NULL,'/bot_mgr/createBot','新增机器人','机器人','POST'),(149,'2025-10-30 21:27:46.903','2025-10-30 21:27:46.903',NULL,'/bot_mgr/deleteBot','删除机器人','机器人','DELETE'),(150,'2025-10-30 21:27:47.229','2025-10-30 21:27:47.229',NULL,'/bot_mgr/deleteBotByIds','批量删除机器人','机器人','DELETE'),(151,'2025-10-30 21:27:47.611','2025-10-30 21:27:47.611',NULL,'/bot_mgr/updateBot','更新机器人','机器人','PUT'),(152,'2025-10-30 21:27:48.030','2025-10-30 21:27:48.030',NULL,'/bot_mgr/findBot','根据ID获取机器人','机器人','GET'),(153,'2025-10-30 21:27:48.410','2025-10-30 21:27:48.410',NULL,'/bot_mgr/getBotList','获取机器人列表','机器人','GET');
/*!40000 ALTER TABLE `sys_apis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_authorities`
--

DROP TABLE IF EXISTS `sys_authorities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_authorities` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `authority_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `authority_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色名',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父角色ID',
  `default_router` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'dashboard' COMMENT '默认菜单',
  PRIMARY KEY (`authority_id`),
  UNIQUE KEY `uni_sys_authorities_authority_id` (`authority_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9529 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_authorities`
--

LOCK TABLES `sys_authorities` WRITE;
/*!40000 ALTER TABLE `sys_authorities` DISABLE KEYS */;
INSERT INTO `sys_authorities` VALUES ('2025-10-30 21:17:20.228','2025-10-30 23:32:58.898',NULL,1,'管理员',0,'bot_mgr'),('2025-10-30 15:16:45.987','2025-10-30 21:46:43.814',NULL,888,'普通用户',0,'dashboard'),('2025-10-30 15:16:45.987','2025-10-30 15:16:59.246',NULL,8881,'普通用户子角色',888,'dashboard');
/*!40000 ALTER TABLE `sys_authorities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_authority_btns`
--

DROP TABLE IF EXISTS `sys_authority_btns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_authority_btns` (
  `authority_id` bigint unsigned DEFAULT NULL COMMENT '角色ID',
  `sys_menu_id` bigint unsigned DEFAULT NULL COMMENT '菜单ID',
  `sys_base_menu_btn_id` bigint unsigned DEFAULT NULL COMMENT '菜单按钮ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_authority_btns`
--

LOCK TABLES `sys_authority_btns` WRITE;
/*!40000 ALTER TABLE `sys_authority_btns` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_authority_btns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_authority_menus`
--

DROP TABLE IF EXISTS `sys_authority_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_authority_menus` (
  `sys_base_menu_id` bigint unsigned NOT NULL,
  `sys_authority_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_base_menu_id`,`sys_authority_authority_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_authority_menus`
--

LOCK TABLES `sys_authority_menus` WRITE;
/*!40000 ALTER TABLE `sys_authority_menus` DISABLE KEYS */;
INSERT INTO `sys_authority_menus` VALUES (1,1),(1,888),(1,8881),(3,1),(3,888),(3,8881),(4,8881),(4,9528),(5,8881),(6,1),(6,888),(6,8881),(9,8881),(10,1),(10,888),(11,1),(11,888),(12,1),(12,888),(13,1),(13,888),(14,1),(14,888),(15,1),(15,888),(16,1),(16,888),(17,8881),(18,8881),(19,8881),(20,1),(20,888),(20,8881),(21,1),(21,888),(21,8881),(22,1),(22,888),(22,8881),(23,1),(23,888),(23,8881),(24,1),(24,888),(24,8881),(25,1),(25,888),(25,8881),(26,1),(26,888),(26,8881),(27,1),(27,888),(27,8881),(28,1),(28,888),(28,8881),(29,1),(29,888),(29,8881),(30,1),(30,888),(30,8881),(37,1),(37,888),(38,1),(38,888);
/*!40000 ALTER TABLE `sys_authority_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_auto_code_histories`
--

DROP TABLE IF EXISTS `sys_auto_code_histories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_auto_code_histories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `table_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表名',
  `package` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模块名/插件名',
  `request` text COLLATE utf8mb4_unicode_ci COMMENT '前端传入的结构化信息',
  `struct_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '结构体名称',
  `abbreviation` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '结构体名称缩写',
  `business_db` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '业务库',
  `description` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Struct中文名称',
  `templates` text COLLATE utf8mb4_unicode_ci COMMENT '模板信息',
  `Injections` text COLLATE utf8mb4_unicode_ci COMMENT '注入路径',
  `flag` bigint DEFAULT NULL COMMENT '[0:创建,1:回滚]',
  `api_ids` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api表注册内容',
  `menu_id` bigint unsigned DEFAULT NULL COMMENT '菜单ID',
  `export_template_id` bigint unsigned DEFAULT NULL COMMENT '导出模板ID',
  `package_id` bigint unsigned DEFAULT NULL COMMENT '包ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_auto_code_histories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_auto_code_histories`
--

LOCK TABLES `sys_auto_code_histories` WRITE;
/*!40000 ALTER TABLE `sys_auto_code_histories` DISABLE KEYS */;
INSERT INTO `sys_auto_code_histories` VALUES (2,'2025-10-30 20:42:29.569','2025-10-30 20:42:29.569',NULL,'bot_msg_mgr','bot','{\"package\":\"bot\",\"tableName\":\"bot_msg_mgr\",\"businessDB\":\"\",\"structName\":\"BotMsgMgr\",\"packageName\":\"bot_msg_mgr\",\"description\":\"机器人消息管理\",\"abbreviation\":\"bot_msg_mgr\",\"humpPackageName\":\"bot_msg_mgr\",\"gvaModel\":true,\"autoMigrate\":true,\"autoCreateResource\":false,\"autoCreateApiToSql\":true,\"autoCreateMenuToSql\":true,\"autoCreateBtnAuth\":false,\"onlyTemplate\":false,\"isTree\":false,\"treeJson\":\"\",\"isAdd\":false,\"fields\":[{\"fieldName\":\"Ban_content\",\"fieldDesc\":\"禁用内容\",\"fieldType\":\"string\",\"fieldJson\":\"ban_content\",\"dataTypeLong\":\"1024\",\"comment\":\"ban_content\",\"columnName\":\"ban_content\",\"fieldSearchType\":\"LIKE\",\"fieldSearchHide\":false,\"dictType\":\"\",\"form\":true,\"table\":true,\"desc\":true,\"excel\":false,\"require\":false,\"defaultValue\":\"\",\"errorText\":\"\",\"clearable\":true,\"sort\":false,\"primaryKey\":false,\"dataSource\":{\"dbName\":\"\",\"table\":\"\",\"label\":\"\",\"value\":\"\",\"association\":1,\"hasDeletedAt\":false},\"checkDataSource\":false,\"fieldIndexType\":\"\"},{\"fieldName\":\"Bot_id\",\"fieldDesc\":\"机器人\",\"fieldType\":\"int\",\"fieldJson\":\"bot_id\",\"dataTypeLong\":\"\",\"comment\":\"机器人ID\",\"columnName\":\"bot_id\",\"fieldSearchType\":\"\",\"fieldSearchHide\":false,\"dictType\":\"\",\"form\":true,\"table\":true,\"desc\":true,\"excel\":false,\"require\":false,\"defaultValue\":\"\",\"errorText\":\"\",\"clearable\":true,\"sort\":false,\"primaryKey\":false,\"dataSource\":{\"dbName\":\"\",\"table\":\"\",\"label\":\"\",\"value\":\"\",\"association\":1,\"hasDeletedAt\":false},\"checkDataSource\":false,\"fieldIndexType\":\"\"}],\"generateWeb\":true,\"generateServer\":true,\"primaryField\":{\"fieldName\":\"ID\",\"fieldDesc\":\"ID\",\"fieldType\":\"uint\",\"fieldJson\":\"ID\",\"dataTypeLong\":\"20\",\"comment\":\"主键ID\",\"columnName\":\"id\",\"fieldSearchType\":\"\",\"fieldSearchHide\":false,\"dictType\":\"\",\"form\":false,\"table\":false,\"desc\":false,\"excel\":false,\"require\":false,\"defaultValue\":\"\",\"errorText\":\"\",\"clearable\":false,\"sort\":false,\"primaryKey\":false,\"dataSource\":null,\"checkDataSource\":false,\"fieldIndexType\":\"\"}}','BotMsgMgr','BotMsgMgr','','机器人消息管理','{\"resource/package/server/api/api.go.tpl\":\"api/v1/bot/bot_msg_mgr.go\",\"resource/package/server/model/model.go.tpl\":\"model/bot/bot_msg_mgr.go\",\"resource/package/server/model/request/request.go.tpl\":\"model/bot/request/bot_msg_mgr.go\",\"resource/package/server/router/router.go.tpl\":\"router/bot/bot_msg_mgr.go\",\"resource/package/server/service/service.go.tpl\":\"service/bot/bot_msg_mgr.go\",\"resource/package/web/api/api.js.tpl\":\"api/bot/bot_msg_mgr.js\",\"resource/package/web/view/form.vue.tpl\":\"view/bot/bot_msg_mgr/bot_msg_mgrForm.vue\",\"resource/package/web/view/table.vue.tpl\":\"view/bot/bot_msg_mgr/bot_msg_mgr.vue\"}','{\"PackageApiEnter\":\"{\\\"Type\\\":\\\"PackageApiEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/api/v1/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/api/v1/bot\\\\\\\"\\\",\\\"StructName\\\":\\\"BotApiGroup\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"api/v1/enter.go\\\",\\\"PackageStructName\\\":\\\"ApiGroup\\\"}\",\"PackageApiModuleEnter\":\"{\\\"Type\\\":\\\"PackageApiModuleEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/api/v1/bot/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/service\\\\\\\"\\\",\\\"RelativePath\\\":\\\"api/v1/bot/enter.go\\\",\\\"StructName\\\":\\\"BotMsgMgrApi\\\",\\\"AppName\\\":\\\"ServiceGroupApp\\\",\\\"GroupName\\\":\\\"BotServiceGroup\\\",\\\"ModuleName\\\":\\\"bot_msg_mgrService\\\",\\\"PackageName\\\":\\\"service\\\",\\\"ServiceName\\\":\\\"BotMsgMgrService\\\"}\",\"PackageInitializeGorm\":\"{\\\"Type\\\":\\\"PackageInitializeGorm\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/initialize/gorm_biz.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/model/bot\\\\\\\"\\\",\\\"Business\\\":\\\"\\\",\\\"StructName\\\":\\\"BotMsgMgr\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"initialize/gorm_biz.go\\\",\\\"IsNew\\\":true}\",\"PackageInitializeRouter\":\"{\\\"Type\\\":\\\"PackageInitializeRouter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/initialize/router_biz.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/router\\\\\\\"\\\",\\\"RelativePath\\\":\\\"initialize/router_biz.go\\\",\\\"AppName\\\":\\\"RouterGroupApp\\\",\\\"GroupName\\\":\\\"Bot\\\",\\\"ModuleName\\\":\\\"botRouter\\\",\\\"PackageName\\\":\\\"router\\\",\\\"FunctionName\\\":\\\"InitBotMsgMgrRouter\\\",\\\"RouterGroupName\\\":\\\"\\\",\\\"LeftRouterGroupName\\\":\\\"privateGroup\\\",\\\"RightRouterGroupName\\\":\\\"publicGroup\\\"}\",\"PackageRouterEnter\":\"{\\\"Type\\\":\\\"PackageRouterEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/router/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/router/bot\\\\\\\"\\\",\\\"StructName\\\":\\\"Bot\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"router/enter.go\\\",\\\"PackageStructName\\\":\\\"RouterGroup\\\"}\",\"PackageRouterModuleEnter\":\"{\\\"Type\\\":\\\"PackageRouterModuleEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/router/bot/enter.go\\\",\\\"ImportPath\\\":\\\"api \\\\\\\"github.com/msean/botmanager/server/api/v1\\\\\\\"\\\",\\\"RelativePath\\\":\\\"router/bot/enter.go\\\",\\\"StructName\\\":\\\"BotMsgMgrRouter\\\",\\\"AppName\\\":\\\"ApiGroupApp\\\",\\\"GroupName\\\":\\\"BotApiGroup\\\",\\\"ModuleName\\\":\\\"bot_msg_mgrApi\\\",\\\"PackageName\\\":\\\"api\\\",\\\"ServiceName\\\":\\\"BotMsgMgrApi\\\"}\",\"PackageServiceEnter\":\"{\\\"Type\\\":\\\"PackageServiceEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/service/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/service/bot\\\\\\\"\\\",\\\"StructName\\\":\\\"BotServiceGroup\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"service/enter.go\\\",\\\"PackageStructName\\\":\\\"ServiceGroup\\\"}\",\"PackageServiceModuleEnter\":\"{\\\"Type\\\":\\\"PackageServiceModuleEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/service/bot/enter.go\\\",\\\"ImportPath\\\":\\\"\\\",\\\"RelativePath\\\":\\\"service/bot/enter.go\\\",\\\"StructName\\\":\\\"BotMsgMgrService\\\",\\\"AppName\\\":\\\"\\\",\\\"GroupName\\\":\\\"\\\",\\\"ModuleName\\\":\\\"\\\",\\\"PackageName\\\":\\\"\\\",\\\"ServiceName\\\":\\\"\\\"}\"}',0,'[142,143,144,145,146,147]',37,0,0),(3,'2025-10-30 21:27:49.667','2025-10-30 21:27:49.667',NULL,'bot','bot','{\"package\":\"bot\",\"tableName\":\"bot\",\"businessDB\":\"\",\"structName\":\"Bot\",\"packageName\":\"bot\",\"description\":\"机器人\",\"abbreviation\":\"bot_mgr\",\"humpPackageName\":\"bot\",\"gvaModel\":true,\"autoMigrate\":true,\"autoCreateResource\":false,\"autoCreateApiToSql\":true,\"autoCreateMenuToSql\":true,\"autoCreateBtnAuth\":false,\"onlyTemplate\":false,\"isTree\":false,\"treeJson\":\"\",\"isAdd\":false,\"fields\":[{\"fieldName\":\"机器人名称\",\"fieldDesc\":\"机器人名称\",\"fieldType\":\"string\",\"fieldJson\":\"name\",\"dataTypeLong\":\"\",\"comment\":\"机器人名称\",\"columnName\":\"name\",\"fieldSearchType\":\"LIKE\",\"fieldSearchHide\":false,\"dictType\":\"\",\"form\":true,\"table\":true,\"desc\":true,\"excel\":false,\"require\":false,\"defaultValue\":\"\",\"errorText\":\"\",\"clearable\":true,\"sort\":false,\"primaryKey\":false,\"dataSource\":{\"dbName\":\"\",\"table\":\"\",\"label\":\"\",\"value\":\"\",\"association\":1,\"hasDeletedAt\":false},\"checkDataSource\":false,\"fieldIndexType\":\"\"},{\"fieldName\":\"机器人token\",\"fieldDesc\":\"机器人token\",\"fieldType\":\"string\",\"fieldJson\":\"token\",\"dataTypeLong\":\"256\",\"comment\":\"token\",\"columnName\":\"token\",\"fieldSearchType\":\"LIKE\",\"fieldSearchHide\":false,\"dictType\":\"\",\"form\":true,\"table\":true,\"desc\":true,\"excel\":false,\"require\":false,\"defaultValue\":\"\",\"errorText\":\"\",\"clearable\":true,\"sort\":false,\"primaryKey\":false,\"dataSource\":{\"dbName\":\"\",\"table\":\"\",\"label\":\"\",\"value\":\"\",\"association\":1,\"hasDeletedAt\":false},\"checkDataSource\":false,\"fieldIndexType\":\"\"}],\"generateWeb\":true,\"generateServer\":true,\"primaryField\":{\"fieldName\":\"ID\",\"fieldDesc\":\"ID\",\"fieldType\":\"uint\",\"fieldJson\":\"ID\",\"dataTypeLong\":\"20\",\"comment\":\"主键ID\",\"columnName\":\"id\",\"fieldSearchType\":\"\",\"fieldSearchHide\":false,\"dictType\":\"\",\"form\":false,\"table\":false,\"desc\":false,\"excel\":false,\"require\":false,\"defaultValue\":\"\",\"errorText\":\"\",\"clearable\":false,\"sort\":false,\"primaryKey\":false,\"dataSource\":null,\"checkDataSource\":false,\"fieldIndexType\":\"\"}}','Bot','Bot','','机器人','{\"resource/package/server/api/api.go.tpl\":\"api/v1/bot/bot.go\",\"resource/package/server/model/model.go.tpl\":\"model/bot/bot.go\",\"resource/package/server/model/request/request.go.tpl\":\"model/bot/request/bot.go\",\"resource/package/server/router/router.go.tpl\":\"router/bot/bot.go\",\"resource/package/server/service/service.go.tpl\":\"service/bot/bot.go\",\"resource/package/web/api/api.js.tpl\":\"api/bot/bot.js\",\"resource/package/web/view/form.vue.tpl\":\"view/bot/bot/botForm.vue\",\"resource/package/web/view/table.vue.tpl\":\"view/bot/bot/bot.vue\"}','{\"PackageApiEnter\":\"{\\\"Type\\\":\\\"PackageApiEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/api/v1/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/api/v1/bot\\\\\\\"\\\",\\\"StructName\\\":\\\"BotApiGroup\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"api/v1/enter.go\\\",\\\"PackageStructName\\\":\\\"ApiGroup\\\"}\",\"PackageApiModuleEnter\":\"{\\\"Type\\\":\\\"PackageApiModuleEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/api/v1/bot/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/service\\\\\\\"\\\",\\\"RelativePath\\\":\\\"api/v1/bot/enter.go\\\",\\\"StructName\\\":\\\"BotApi\\\",\\\"AppName\\\":\\\"ServiceGroupApp\\\",\\\"GroupName\\\":\\\"BotServiceGroup\\\",\\\"ModuleName\\\":\\\"bot_mgrService\\\",\\\"PackageName\\\":\\\"service\\\",\\\"ServiceName\\\":\\\"BotService\\\"}\",\"PackageInitializeGorm\":\"{\\\"Type\\\":\\\"PackageInitializeGorm\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/initialize/gorm_biz.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/model/bot\\\\\\\"\\\",\\\"Business\\\":\\\"\\\",\\\"StructName\\\":\\\"Bot\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"initialize/gorm_biz.go\\\",\\\"IsNew\\\":true}\",\"PackageInitializeRouter\":\"{\\\"Type\\\":\\\"PackageInitializeRouter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/initialize/router_biz.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/router\\\\\\\"\\\",\\\"RelativePath\\\":\\\"initialize/router_biz.go\\\",\\\"AppName\\\":\\\"RouterGroupApp\\\",\\\"GroupName\\\":\\\"Bot\\\",\\\"ModuleName\\\":\\\"botRouter\\\",\\\"PackageName\\\":\\\"router\\\",\\\"FunctionName\\\":\\\"InitBotRouter\\\",\\\"RouterGroupName\\\":\\\"\\\",\\\"LeftRouterGroupName\\\":\\\"privateGroup\\\",\\\"RightRouterGroupName\\\":\\\"publicGroup\\\"}\",\"PackageRouterEnter\":\"{\\\"Type\\\":\\\"PackageRouterEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/router/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/router/bot\\\\\\\"\\\",\\\"StructName\\\":\\\"Bot\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"router/enter.go\\\",\\\"PackageStructName\\\":\\\"RouterGroup\\\"}\",\"PackageRouterModuleEnter\":\"{\\\"Type\\\":\\\"PackageRouterModuleEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/router/bot/enter.go\\\",\\\"ImportPath\\\":\\\"api \\\\\\\"github.com/msean/botmanager/server/api/v1\\\\\\\"\\\",\\\"RelativePath\\\":\\\"router/bot/enter.go\\\",\\\"StructName\\\":\\\"BotRouter\\\",\\\"AppName\\\":\\\"ApiGroupApp\\\",\\\"GroupName\\\":\\\"BotApiGroup\\\",\\\"ModuleName\\\":\\\"bot_mgrApi\\\",\\\"PackageName\\\":\\\"api\\\",\\\"ServiceName\\\":\\\"BotApi\\\"}\",\"PackageServiceEnter\":\"{\\\"Type\\\":\\\"PackageServiceEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/service/enter.go\\\",\\\"ImportPath\\\":\\\"\\\\\\\"github.com/msean/botmanager/server/service/bot\\\\\\\"\\\",\\\"StructName\\\":\\\"BotServiceGroup\\\",\\\"PackageName\\\":\\\"bot\\\",\\\"RelativePath\\\":\\\"service/enter.go\\\",\\\"PackageStructName\\\":\\\"ServiceGroup\\\"}\",\"PackageServiceModuleEnter\":\"{\\\"Type\\\":\\\"PackageServiceModuleEnter\\\",\\\"Path\\\":\\\"/Users/msean/project/go/bot_manamger/server/service/bot/enter.go\\\",\\\"ImportPath\\\":\\\"\\\",\\\"RelativePath\\\":\\\"service/bot/enter.go\\\",\\\"StructName\\\":\\\"BotService\\\",\\\"AppName\\\":\\\"\\\",\\\"GroupName\\\":\\\"\\\",\\\"ModuleName\\\":\\\"\\\",\\\"PackageName\\\":\\\"\\\",\\\"ServiceName\\\":\\\"\\\"}\"}',0,'[148,149,150,151,152,153]',38,0,0);
/*!40000 ALTER TABLE `sys_auto_code_histories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_auto_code_packages`
--

DROP TABLE IF EXISTS `sys_auto_code_packages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_auto_code_packages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `desc` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '描述',
  `label` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '展示名',
  `template` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模版',
  `package_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '包名',
  `module` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sys_auto_code_packages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_auto_code_packages`
--

LOCK TABLES `sys_auto_code_packages` WRITE;
/*!40000 ALTER TABLE `sys_auto_code_packages` DISABLE KEYS */;
INSERT INTO `sys_auto_code_packages` VALUES (3,'2025-10-30 18:17:46.075','2025-10-30 18:17:46.075',NULL,'系统自动读取example包','example包','package','example','github.com/msean/botmanager/server'),(4,'2025-10-30 18:17:46.075','2025-10-30 18:17:46.075',NULL,'系统自动读取system包','system包','package','system','github.com/msean/botmanager/server'),(5,'2025-10-30 18:17:46.075','2025-10-30 18:17:46.075',NULL,'系统自动读取announcement插件，使用前请确认是否为v2版本插件','announcement插件','plugin','announcement','github.com/msean/botmanager/server'),(6,'2025-10-30 18:17:46.075','2025-10-30 18:17:46.075',NULL,'系统自动读取，但是缺少 initialize、plugin 结构，不建议自动化和mcp使用','email插件','plugin','email','github.com/msean/botmanager/server'),(7,'2025-10-30 18:17:46.075','2025-10-30 18:17:46.075',NULL,'系统自动读取，但是缺少 initialize、plugin、router、service、api、config 结构，不建议自动化和mcp使用','plugin-tool插件','plugin','plugin-tool','github.com/msean/botmanager/server'),(8,'2025-10-30 18:18:27.778','2025-10-30 18:18:27.778',NULL,'bot','bot包','package','bot','github.com/msean/botmanager/server');
/*!40000 ALTER TABLE `sys_auto_code_packages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menu_btns`
--

DROP TABLE IF EXISTS `sys_base_menu_btns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menu_btns` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '按钮关键key',
  `desc` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sys_base_menu_id` bigint unsigned DEFAULT NULL COMMENT '菜单ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menu_btns_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menu_btns`
--

LOCK TABLES `sys_base_menu_btns` WRITE;
/*!40000 ALTER TABLE `sys_base_menu_btns` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_base_menu_btns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menu_parameters`
--

DROP TABLE IF EXISTS `sys_base_menu_parameters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menu_parameters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `sys_base_menu_id` bigint unsigned DEFAULT NULL,
  `type` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地址栏携带参数为params还是query',
  `key` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地址栏携带参数的key',
  `value` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地址栏携带参数的值',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menu_parameters_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menu_parameters`
--

LOCK TABLES `sys_base_menu_parameters` WRITE;
/*!40000 ALTER TABLE `sys_base_menu_parameters` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_base_menu_parameters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_base_menus`
--

DROP TABLE IF EXISTS `sys_base_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_base_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `menu_level` bigint unsigned DEFAULT NULL,
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由path',
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由name',
  `hidden` tinyint(1) DEFAULT NULL COMMENT '是否在列表隐藏',
  `component` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '对应前端文件路径',
  `sort` bigint DEFAULT NULL COMMENT '排序标记',
  `active_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '附加属性',
  `keep_alive` tinyint(1) DEFAULT NULL COMMENT '附加属性',
  `default_menu` tinyint(1) DEFAULT NULL COMMENT '附加属性',
  `title` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '附加属性',
  `icon` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '附加属性',
  `close_tab` tinyint(1) DEFAULT NULL COMMENT '附加属性',
  `transition_type` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '附加属性',
  PRIMARY KEY (`id`),
  KEY `idx_sys_base_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_base_menus`
--

LOCK TABLES `sys_base_menus` WRITE;
/*!40000 ALTER TABLE `sys_base_menus` DISABLE KEYS */;
INSERT INTO `sys_base_menus` VALUES (1,'2025-10-30 15:16:53.592','2025-10-30 16:00:15.576',NULL,0,0,'dashboard','dashboard',1,'view/dashboard/index.vue',1,'',0,0,'仪表盘','odometer',0,''),(2,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592','2025-10-30 15:58:30.217',0,0,'about','about',1,'view/about/index.vue',9,'',0,0,'关于我们','info-filled',0,''),(3,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592',NULL,0,0,'admin','superAdmin',0,'view/superAdmin/index.vue',3,'',0,0,'超级管理员','user',0,''),(4,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592','2025-10-30 15:58:30.217',0,0,'person','person',1,'view/person/person.vue',4,'',0,0,'个人信息','message',0,''),(5,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592','2025-10-30 15:58:30.217',0,0,'example','example',1,'view/example/index.vue',7,'',0,0,'示例文件','management',0,''),(6,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592',NULL,0,0,'systemTools','systemTools',0,'view/systemTools/index.vue',5,'',0,0,'系统工具','tools',0,''),(7,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592','2025-10-30 15:56:09.234',0,0,'https://www.gin-vue-admin.com','https://www.gin-vue-admin.com',1,'/',0,'',0,0,'官方网站','customer-gva',0,''),(8,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592',NULL,0,0,'state','state',1,'view/system/state.vue',8,'',0,0,'服务器状态','cloudy',0,''),(9,'2025-10-30 15:16:53.592','2025-10-30 15:16:53.592','2025-10-30 15:58:30.217',0,0,'plugin','plugin',1,'view/routerHolder.vue',6,'',0,0,'插件系统','cherry',0,''),(10,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'authority','authority',0,'view/superAdmin/authority/authority.vue',1,'',0,0,'角色管理','avatar',0,''),(11,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'menu','menu',0,'view/superAdmin/menu/menu.vue',2,'',1,0,'菜单管理','tickets',0,''),(12,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'api','api',0,'view/superAdmin/api/api.vue',3,'',1,0,'api管理','platform',0,''),(13,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'user','user',0,'view/superAdmin/user/user.vue',4,'',0,0,'用户管理','coordinate',0,''),(14,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'dictionary','dictionary',0,'view/superAdmin/dictionary/sysDictionary.vue',5,'',0,0,'字典管理','notebook',0,''),(15,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'operation','operation',0,'view/superAdmin/operation/sysOperationRecord.vue',6,'',0,0,'操作历史','pie-chart',0,''),(16,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,3,'sysParams','sysParams',0,'view/superAdmin/params/sysParams.vue',7,'',0,0,'参数管理','compass',0,''),(17,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,5,'upload','upload',0,'view/example/upload/upload.vue',5,'',0,0,'媒体库（上传下载）','upload',0,''),(18,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,5,'breakpoint','breakpoint',0,'view/example/breakpoint/breakpoint.vue',6,'',0,0,'断点续传','upload-filled',0,''),(19,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,5,'customer','customer',0,'view/example/customer/customer.vue',7,'',0,0,'客户列表（资源示例）','avatar',0,''),(20,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'autoCode','autoCode',0,'view/systemTools/autoCode/index.vue',1,'',1,0,'代码生成器','cpu',0,''),(21,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'formCreate','formCreate',0,'view/systemTools/formCreate/index.vue',3,'',1,0,'表单生成器','magic-stick',0,''),(22,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'system','system',0,'view/systemTools/system/system.vue',4,'',0,0,'系统配置','operation',0,''),(23,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'autoCodeAdmin','autoCodeAdmin',0,'view/systemTools/autoCodeAdmin/index.vue',2,'',0,0,'自动化代码管理','magic-stick',0,''),(24,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'autoCodeEdit/:id','autoCodeEdit',1,'view/systemTools/autoCode/index.vue',0,'',0,0,'自动化代码-${id}','magic-stick',0,''),(25,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'autoPkg','autoPkg',0,'view/systemTools/autoPkg/autoPkg.vue',0,'',0,0,'模板配置','folder',0,''),(26,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'exportTemplate','exportTemplate',0,'view/systemTools/exportTemplate/exportTemplate.vue',5,'',0,0,'导出模板','reading',0,''),(27,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'picture','picture',0,'view/systemTools/autoCode/picture.vue',6,'',0,0,'AI页面绘制','picture-filled',0,''),(28,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'mcpTool','mcpTool',0,'view/systemTools/autoCode/mcp.vue',7,'',0,0,'Mcp Tools模板','magnet',0,''),(29,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'mcpTest','mcpTest',0,'view/systemTools/autoCode/mcpTest.vue',7,'',0,0,'Mcp Tools测试','partly-cloudy',0,''),(30,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,6,'sysVersion','sysVersion',0,'view/systemTools/version/version.vue',8,'',0,0,'版本管理','server',0,''),(31,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,9,'https://plugin.gin-vue-admin.com/','https://plugin.gin-vue-admin.com/',0,'https://plugin.gin-vue-admin.com/',0,'',0,0,'插件市场','shop',0,''),(32,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,9,'installPlugin','installPlugin',0,'view/systemTools/installPlugin/index.vue',1,'',0,0,'插件安装','box',0,''),(33,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,9,'pubPlug','pubPlug',0,'view/systemTools/pubPlug/pubPlug.vue',3,'',0,0,'打包插件','files',0,''),(34,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,9,'plugin-email','plugin-email',0,'plugin/email/view/index.vue',4,'',0,0,'邮件插件','message',0,''),(35,'2025-10-30 15:16:53.899','2025-10-30 15:16:53.899',NULL,1,9,'anInfo','anInfo',0,'plugin/announcement/view/info.vue',5,'',0,0,'公告管理[示例]','scaleToOriginal',0,''),(37,'2025-10-30 20:42:29.277','2025-10-30 20:42:29.277',NULL,0,0,'bot_msg_mgr','bot_msg_mgr',0,'view/bot/bot_msg_mgr/bot_msg_mgr.vue',0,'',0,0,'机器人消息管理','',0,''),(38,'2025-10-30 21:27:49.074','2025-10-30 21:27:49.074',NULL,0,0,'bot_mgr','bot_mgr',0,'view/bot/bot/bot.vue',0,'',0,0,'机器人','',0,'');
/*!40000 ALTER TABLE `sys_base_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_data_authority_id`
--

DROP TABLE IF EXISTS `sys_data_authority_id`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_data_authority_id` (
  `sys_authority_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `data_authority_id_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_authority_authority_id`,`data_authority_id_authority_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_data_authority_id`
--

LOCK TABLES `sys_data_authority_id` WRITE;
/*!40000 ALTER TABLE `sys_data_authority_id` DISABLE KEYS */;
INSERT INTO `sys_data_authority_id` VALUES (1,1),(888,888),(888,8881);
/*!40000 ALTER TABLE `sys_data_authority_id` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dictionaries`
--

DROP TABLE IF EXISTS `sys_dictionaries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dictionaries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '字典名（中）',
  `type` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '字典名（英）',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态',
  `desc` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '描述',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级字典ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_dictionaries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dictionaries`
--

LOCK TABLES `sys_dictionaries` WRITE;
/*!40000 ALTER TABLE `sys_dictionaries` DISABLE KEYS */;
INSERT INTO `sys_dictionaries` VALUES (1,'2025-10-30 15:16:48.280','2025-10-30 15:16:48.888',NULL,'性别','gender',1,'性别字典',NULL),(2,'2025-10-30 15:16:48.280','2025-10-30 15:16:49.645',NULL,'数据库int类型','int',1,'int类型对应的数据库类型',NULL),(3,'2025-10-30 15:16:48.280','2025-10-30 15:16:50.415',NULL,'数据库时间日期类型','time.Time',1,'数据库时间日期类型',NULL),(4,'2025-10-30 15:16:48.280','2025-10-30 15:16:51.168',NULL,'数据库浮点型','float64',1,'数据库浮点型',NULL),(5,'2025-10-30 15:16:48.280','2025-10-30 15:16:51.922',NULL,'数据库字符串','string',1,'数据库字符串',NULL),(6,'2025-10-30 15:16:48.280','2025-10-30 15:16:52.684',NULL,'数据库bool类型','bool',1,'数据库bool类型',NULL);
/*!40000 ALTER TABLE `sys_dictionaries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dictionary_details`
--

DROP TABLE IF EXISTS `sys_dictionary_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dictionary_details` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `label` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '展示值',
  `value` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '字典值',
  `extend` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '扩展值',
  `status` tinyint(1) DEFAULT NULL COMMENT '启用状态',
  `sort` bigint DEFAULT NULL COMMENT '排序标记',
  `sys_dictionary_id` bigint unsigned DEFAULT NULL COMMENT '关联标记',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级字典详情ID',
  `level` bigint DEFAULT NULL COMMENT '层级深度',
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '层级路径',
  PRIMARY KEY (`id`),
  KEY `idx_sys_dictionary_details_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dictionary_details`
--

LOCK TABLES `sys_dictionary_details` WRITE;
/*!40000 ALTER TABLE `sys_dictionary_details` DISABLE KEYS */;
INSERT INTO `sys_dictionary_details` VALUES (1,'2025-10-30 15:16:49.038','2025-10-30 15:16:49.038',NULL,'男','1','',1,1,1,NULL,0,''),(2,'2025-10-30 15:16:49.038','2025-10-30 15:16:49.038',NULL,'女','2','',1,2,1,NULL,0,''),(3,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'smallint','1','mysql',1,1,2,NULL,0,''),(4,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'mediumint','2','mysql',1,2,2,NULL,0,''),(5,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'int','3','mysql',1,3,2,NULL,0,''),(6,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'bigint','4','mysql',1,4,2,NULL,0,''),(7,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'int2','5','pgsql',1,5,2,NULL,0,''),(8,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'int4','6','pgsql',1,6,2,NULL,0,''),(9,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'int6','7','pgsql',1,7,2,NULL,0,''),(10,'2025-10-30 15:16:49.796','2025-10-30 15:16:49.796',NULL,'int8','8','pgsql',1,8,2,NULL,0,''),(11,'2025-10-30 15:16:50.563','2025-10-30 15:16:50.563',NULL,'date','0','mysql',1,0,3,NULL,0,''),(12,'2025-10-30 15:16:50.563','2025-10-30 15:16:50.563',NULL,'time','1','mysql',1,1,3,NULL,0,''),(13,'2025-10-30 15:16:50.563','2025-10-30 15:16:50.563',NULL,'year','2','mysql',1,2,3,NULL,0,''),(14,'2025-10-30 15:16:50.563','2025-10-30 15:16:50.563',NULL,'datetime','3','mysql',1,3,3,NULL,0,''),(15,'2025-10-30 15:16:50.563','2025-10-30 15:16:50.563',NULL,'timestamp','5','mysql',1,5,3,NULL,0,''),(16,'2025-10-30 15:16:50.563','2025-10-30 15:16:50.563',NULL,'timestamptz','6','pgsql',1,5,3,NULL,0,''),(17,'2025-10-30 15:16:51.317','2025-10-30 15:16:51.317',NULL,'float','0','mysql',1,0,4,NULL,0,''),(18,'2025-10-30 15:16:51.317','2025-10-30 15:16:51.317',NULL,'double','1','mysql',1,1,4,NULL,0,''),(19,'2025-10-30 15:16:51.317','2025-10-30 15:16:51.317',NULL,'decimal','2','mysql',1,2,4,NULL,0,''),(20,'2025-10-30 15:16:51.317','2025-10-30 15:16:51.317',NULL,'numeric','3','pgsql',1,3,4,NULL,0,''),(21,'2025-10-30 15:16:51.317','2025-10-30 15:16:51.317',NULL,'smallserial','4','pgsql',1,4,4,NULL,0,''),(22,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'char','0','mysql',1,0,5,NULL,0,''),(23,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'varchar','1','mysql',1,1,5,NULL,0,''),(24,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'tinyblob','2','mysql',1,2,5,NULL,0,''),(25,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'tinytext','3','mysql',1,3,5,NULL,0,''),(26,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'text','4','mysql',1,4,5,NULL,0,''),(27,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'blob','5','mysql',1,5,5,NULL,0,''),(28,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'mediumblob','6','mysql',1,6,5,NULL,0,''),(29,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'mediumtext','7','mysql',1,7,5,NULL,0,''),(30,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'longblob','8','mysql',1,8,5,NULL,0,''),(31,'2025-10-30 15:16:52.073','2025-10-30 15:16:52.073',NULL,'longtext','9','mysql',1,9,5,NULL,0,''),(32,'2025-10-30 15:16:52.834','2025-10-30 15:16:52.834',NULL,'tinyint','1','mysql',1,0,6,NULL,0,''),(33,'2025-10-30 15:16:52.834','2025-10-30 15:16:52.834',NULL,'bool','2','pgsql',1,0,6,NULL,0,'');
/*!40000 ALTER TABLE `sys_dictionary_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_export_template_condition`
--

DROP TABLE IF EXISTS `sys_export_template_condition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_export_template_condition` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模板标识',
  `from` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '条件取的key',
  `column` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '作为查询条件的字段',
  `operator` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作符',
  PRIMARY KEY (`id`),
  KEY `idx_sys_export_template_condition_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_export_template_condition`
--

LOCK TABLES `sys_export_template_condition` WRITE;
/*!40000 ALTER TABLE `sys_export_template_condition` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_export_template_condition` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_export_template_join`
--

DROP TABLE IF EXISTS `sys_export_template_join`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_export_template_join` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模板标识',
  `joins` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联',
  `table` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联表',
  `on` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联条件',
  PRIMARY KEY (`id`),
  KEY `idx_sys_export_template_join_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_export_template_join`
--

LOCK TABLES `sys_export_template_join` WRITE;
/*!40000 ALTER TABLE `sys_export_template_join` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_export_template_join` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_export_templates`
--

DROP TABLE IF EXISTS `sys_export_templates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_export_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `db_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据库名称',
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模板名称',
  `table_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表名称',
  `template_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模板标识',
  `template_info` text COLLATE utf8mb4_unicode_ci,
  `limit` bigint DEFAULT NULL COMMENT '导出限制',
  `order` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '排序',
  PRIMARY KEY (`id`),
  KEY `idx_sys_export_templates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_export_templates`
--

LOCK TABLES `sys_export_templates` WRITE;
/*!40000 ALTER TABLE `sys_export_templates` DISABLE KEYS */;
INSERT INTO `sys_export_templates` VALUES (1,'2025-10-30 15:16:56.781','2025-10-30 15:16:56.781',NULL,'','api','sys_apis','api','{\n\"path\":\"路径\",\n\"method\":\"方法（大写）\",\n\"description\":\"方法介绍\",\n\"api_group\":\"方法分组\"\n}',NULL,'');
/*!40000 ALTER TABLE `sys_export_templates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_ignore_apis`
--

DROP TABLE IF EXISTS `sys_ignore_apis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_ignore_apis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api路径',
  `method` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'POST' COMMENT '方法',
  PRIMARY KEY (`id`),
  KEY `idx_sys_ignore_apis_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_ignore_apis`
--

LOCK TABLES `sys_ignore_apis` WRITE;
/*!40000 ALTER TABLE `sys_ignore_apis` DISABLE KEYS */;
INSERT INTO `sys_ignore_apis` VALUES (1,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/swagger/*any','GET'),(2,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/api/freshCasbin','GET'),(3,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/uploads/file/*filepath','GET'),(4,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/health','GET'),(5,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/uploads/file/*filepath','HEAD'),(6,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/autoCode/llmAuto','POST'),(7,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/system/reloadSystem','POST'),(8,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/base/login','POST'),(9,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/base/captcha','POST'),(10,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/init/initdb','POST'),(11,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/init/checkdb','POST'),(12,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/info/getInfoDataSource','GET'),(13,'2025-10-30 15:16:44.972','2025-10-30 15:16:44.972',NULL,'/info/getInfoPublic','GET');
/*!40000 ALTER TABLE `sys_ignore_apis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_operation_records`
--

DROP TABLE IF EXISTS `sys_operation_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_operation_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ip` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求ip',
  `method` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求方法',
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求路径',
  `status` bigint DEFAULT NULL COMMENT '请求状态',
  `latency` bigint DEFAULT NULL COMMENT '延迟',
  `agent` text COLLATE utf8mb4_unicode_ci COMMENT '代理',
  `error_message` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '错误信息',
  `body` text COLLATE utf8mb4_unicode_ci COMMENT '请求Body',
  `resp` text COLLATE utf8mb4_unicode_ci COMMENT '响应Body',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  KEY `idx_sys_operation_records_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_operation_records`
--

LOCK TABLES `sys_operation_records` WRITE;
/*!40000 ALTER TABLE `sys_operation_records` DISABLE KEYS */;
INSERT INTO `sys_operation_records` VALUES (1,'2025-10-30 15:53:51.493','2025-10-30 15:53:51.493',NULL,'127.0.0.1','GET','/api/getApiGroups',200,148624500,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\"]},\"msg\":\"成功\"}',1),(2,'2025-10-30 15:56:10.104','2025-10-30 15:56:10.104',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1310839333,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":7}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(3,'2025-10-30 15:56:22.250','2025-10-30 15:56:22.250',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,461308208,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单有角色正在作为首页，不可删除\"}',1),(4,'2025-10-30 15:56:52.579','2025-10-30 15:56:52.579',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,144779792,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":5}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单存在子菜单不可删除\"}',1),(5,'2025-10-30 15:57:40.332','2025-10-30 15:57:40.332',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,440181166,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单有角色正在作为首页，不可删除\"}',1),(6,'2025-10-30 15:58:31.337','2025-10-30 15:58:31.337',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1574001750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":2}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(7,'2025-10-30 15:58:48.986','2025-10-30 15:58:48.986',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1345941917,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":8}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(8,'2025-10-30 15:58:54.070','2025-10-30 15:58:54.070',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,152727958,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":9}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单存在子菜单不可删除\"}',1),(9,'2025-10-30 16:00:15.873','2025-10-30 16:00:15.873',NULL,'127.0.0.1','POST','/menu/updateBaseMenu',200,740064167,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"CreatedAt\":\"2025-10-30T15:16:53.592+07:00\",\"UpdatedAt\":\"2025-10-30T15:16:53.592+07:00\",\"parentId\":0,\"path\":\"dashboard\",\"name\":\"dashboard\",\"hidden\":true,\"component\":\"view/dashboard/index.vue\",\"sort\":1,\"meta\":{\"activeName\":\"\",\"keepAlive\":false,\"defaultMenu\":false,\"title\":\"仪表盘\",\"icon\":\"odometer\",\"closeTab\":false,\"transitionType\":\"\"},\"authoritys\":null,\"children\":null,\"parameters\":[],\"menuBtn\":[]}','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(10,'2025-10-30 16:01:07.234','2025-10-30 16:01:07.234',NULL,'127.0.0.1','GET','/api/getApiGroups',200,145864750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\"]},\"msg\":\"成功\"}',1),(11,'2025-10-30 16:22:37.142','2025-10-30 16:22:37.142',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,702411375,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(12,'2025-10-30 16:26:13.701','2025-10-30 16:26:13.701',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,149924583,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":5}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单存在子菜单不可删除\"}',1),(13,'2025-10-30 16:35:26.550','2025-10-30 16:35:26.550',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,148997625,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":5}','{\"code\":7,\"data\":{},\"msg\":\"删除失败:此菜单存在子菜单不可删除\"}',1),(14,'2025-10-30 16:38:53.446','2025-10-30 16:38:53.446',NULL,'127.0.0.1','GET','/api/getApiGroups',200,146015834,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\"]},\"msg\":\"成功\"}',1),(15,'2025-10-30 16:48:00.038','2025-10-30 16:48:00.038',NULL,'127.0.0.1','GET','/api/getApiGroups',200,147709750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\"]},\"msg\":\"成功\"}',1),(16,'2025-10-30 16:53:27.419','2025-10-30 16:53:27.419',NULL,'127.0.0.1','GET','/api/getApiGroups',200,154340625,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\"]},\"msg\":\"成功\"}',1),(17,'2025-10-30 18:16:28.235','2025-10-30 18:16:28.235',NULL,'127.0.0.1','GET','/api/getApiGroups',200,154797042,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\"]},\"msg\":\"成功\"}',1),(18,'2025-10-30 18:27:24.266','2025-10-30 18:27:24.266',NULL,'127.0.0.1','GET','/api/getApiGroups',200,214366708,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"botList\":\"表\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"表\"]},\"msg\":\"成功\"}',1),(19,'2025-10-30 18:27:32.600','2025-10-30 18:27:32.600',NULL,'127.0.0.1','POST','/api/getApiById',200,154141958,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"id\":141}','{\"code\":0,\"data\":{\"api\":{\"ID\":141,\"CreatedAt\":\"2025-10-30T18:26:11.381+07:00\",\"UpdatedAt\":\"2025-10-30T18:26:11.381+07:00\",\"path\":\"/botList/getBotList\",\"description\":\"获取表列表\",\"apiGroup\":\"表\",\"method\":\"GET\"}},\"msg\":\"获取成功\"}',1),(20,'2025-10-30 18:27:44.175','2025-10-30 18:27:44.175',NULL,'127.0.0.1','POST','/api/updateApi',200,843486125,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":141,\"CreatedAt\":\"2025-10-30T18:26:11.381+07:00\",\"UpdatedAt\":\"2025-10-30T18:26:11.381+07:00\",\"path\":\"/botList/getBotList\",\"description\":\"获取表列表\",\"apiGroup\":\"表\",\"method\":\"GET\"}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(21,'2025-10-30 20:02:49.801','2025-10-30 20:02:49.801',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1522016583,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(22,'2025-10-30 20:02:54.390','2025-10-30 20:02:54.390',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1497919417,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(23,'2025-10-30 20:03:04.978','2025-10-30 20:03:04.978',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,682741166,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(24,'2025-10-30 20:06:58.737','2025-10-30 20:06:58.737',NULL,'127.0.0.1','GET','/api/getApiGroups',200,153486791,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"botList\":\"表\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"机器人列表\",\"表\"]},\"msg\":\"成功\"}',1),(25,'2025-10-30 20:44:35.067','2025-10-30 20:44:35.067',NULL,'127.0.0.1','POST','/menu/deleteBaseMenu',200,1659992417,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":36}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(26,'2025-10-30 20:50:52.778','2025-10-30 20:50:52.778',NULL,'127.0.0.1','GET','/api/getApiGroups',200,147132625,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"botList\":\"表\",\"bot_msg_mgr\":\"机器人消息管理\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"机器人列表\",\"表\",\"机器人消息管理\"]},\"msg\":\"成功\"}',1),(27,'2025-10-30 20:51:21.470','2025-10-30 20:51:21.470',NULL,'127.0.0.1','DELETE','/api/deleteApisByIds',200,1738307541,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ids\":[141,140,139,138]}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(28,'2025-10-30 20:51:37.664','2025-10-30 20:51:37.664',NULL,'127.0.0.1','DELETE','/api/deleteApisByIds',200,1061491417,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ids\":[137,136]}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(29,'2025-10-30 20:54:12.628','2025-10-30 20:54:12.628',NULL,'127.0.0.1','GET','/api/getApiGroups',200,147546583,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"bot_msg_mgr\":\"机器人消息管理\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"机器人消息管理\"]},\"msg\":\"成功\"}',1),(30,'2025-10-30 21:16:26.022','2025-10-30 21:16:26.022',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,691081750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(31,'2025-10-30 21:16:42.298','2025-10-30 21:16:42.298',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1488200708,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(32,'2025-10-30 21:16:44.323','2025-10-30 21:16:44.323',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,681861708,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(33,'2025-10-30 21:17:21.336','2025-10-30 21:17:21.336',NULL,'127.0.0.1','POST','/authority/createAuthority',200,1277464458,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"authorityId\":1,\"authorityName\":\"管理员\",\"parentId\":0}','{\"code\":0,\"data\":{\"authority\":{\"CreatedAt\":\"2025-10-30T21:17:20.228+07:00\",\"UpdatedAt\":\"2025-10-30T21:17:20.372+07:00\",\"DeletedAt\":null,\"authorityId\":1,\"authorityName\":\"管理员\",\"parentId\":0,\"dataAuthorityId\":null,\"children\":null,\"menus\":[{\"ID\":1,\"CreatedAt\":\"2025-10-30T21:17:20.518+07:00\",\"UpdatedAt\":\"2025-10-30T21:17:20.518+07:00\",\"parentId\":0,\"path\":\"dashboard\",\"name\":\"dashboard\",\"hidden\":false,\"component\":\"view/dashboard/index.vue\",\"sort\":1,\"meta\":{\"activeName\":\"\",\"keepAlive\":false,\"defaultMenu\":false,\"title\":\"仪表盘\",\"icon\":\"setting\",\"closeTab\":false,\"transitionType\":\"\"},\"authoritys\":null,\"children\":null,\"parameters\":null,\"menuBtn\":null}],\"defaultRouter\":\"dashboard\"}},\"msg\":\"创建成功\"}',1),(34,'2025-10-30 21:17:36.875','2025-10-30 21:17:36.875',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1487438750,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(35,'2025-10-30 21:18:31.883','2025-10-30 21:18:31.883',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,698685417,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(36,'2025-10-30 21:18:50.789','2025-10-30 21:18:50.789',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,767462416,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[8881,9528]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(37,'2025-10-30 21:18:53.473','2025-10-30 21:18:53.473',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,753805834,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[9528]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(38,'2025-10-30 21:18:55.565','2025-10-30 21:18:55.565',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,499535208,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[]}','{\"code\":7,\"data\":{},\"msg\":\"修改失败\"}',1),(39,'2025-10-30 21:19:03.307','2025-10-30 21:19:03.307',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,767751417,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[9528]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(40,'2025-10-30 21:19:13.786','2025-10-30 21:19:13.786',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,758509291,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[1]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(41,'2025-10-30 21:19:32.447','2025-10-30 21:19:32.447',NULL,'127.0.0.1','POST','/authority/deleteAuthority',200,462009708,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"authorityId\":9528}','{\"code\":7,\"data\":{},\"msg\":\"删除失败此角色有用户正在使用禁止删除\"}',1),(42,'2025-10-30 21:19:39.609','2025-10-30 21:19:39.609',NULL,'127.0.0.1','POST','/authority/deleteAuthority',200,457771875,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"authorityId\":9528}','{\"code\":7,\"data\":{},\"msg\":\"删除失败此角色有用户正在使用禁止删除\"}',1),(43,'2025-10-30 21:19:51.575','2025-10-30 21:19:51.575',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,448248125,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":2,\"authorityIds\":[]}','{\"code\":7,\"data\":{},\"msg\":\"修改失败\"}',1),(44,'2025-10-30 21:19:56.653','2025-10-30 21:19:56.653',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,782748208,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":2,\"authorityIds\":[888]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(45,'2025-10-30 21:20:02.375','2025-10-30 21:20:02.375',NULL,'127.0.0.1','GET','/api/getApiGroups',200,246232500,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"bot_msg_mgr\":\"机器人消息管理\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"机器人消息管理\"]},\"msg\":\"成功\"}',1),(46,'2025-10-30 21:20:25.280','2025-10-30 21:20:25.280',NULL,'127.0.0.1','POST','/authority/deleteAuthority',200,2621180875,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"authorityId\":9528}','{\"code\":0,\"data\":{},\"msg\":\"删除成功\"}',1),(47,'2025-10-30 21:33:29.846','2025-10-30 21:33:29.846',NULL,'127.0.0.1','GET','/api/getApiGroups',200,160303958,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"bot_mgr\":\"机器人\",\"bot_msg_mgr\":\"机器人消息管理\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"机器人消息管理\",\"机器人\"]},\"msg\":\"成功\"}',1),(48,'2025-10-30 21:33:47.796','2025-10-30 21:33:47.796',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1441853250,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(49,'2025-10-30 21:33:49.847','2025-10-30 21:33:49.847',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,688300583,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(50,'2025-10-30 21:42:15.544','2025-10-30 21:42:15.544',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1484327458,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(51,'2025-10-30 21:42:28.708','2025-10-30 21:42:28.708',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,732359875,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[1,888]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(52,'2025-10-30 21:43:09.065','2025-10-30 21:43:09.065',NULL,'127.0.0.1','PUT','/authority/updateAuthority',200,664936791,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"authorityId\":1,\"authorityName\":\"管理员\",\"parentId\":0}','{\"code\":0,\"data\":{\"authority\":{\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"authorityId\":1,\"authorityName\":\"管理员\",\"parentId\":0,\"dataAuthorityId\":null,\"children\":null,\"menus\":null,\"defaultRouter\":\"\"}},\"msg\":\"更新成功\"}',1),(53,'2025-10-30 21:43:20.325','2025-10-30 21:43:20.325',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1570422458,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(54,'2025-10-30 21:43:22.920','2025-10-30 21:43:22.920',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,931444500,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(55,'2025-10-30 21:43:37.806','2025-10-30 21:43:37.806',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,1047401125,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[888]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(56,'2025-10-30 21:43:39.963','2025-10-30 21:43:39.963',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,468287542,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[]}','{\"code\":7,\"data\":{},\"msg\":\"修改失败\"}',1),(57,'2025-10-30 21:43:41.507','2025-10-30 21:43:41.507',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,475785000,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[]}','{\"code\":7,\"data\":{},\"msg\":\"修改失败\"}',1),(58,'2025-10-30 21:43:48.863','2025-10-30 21:43:48.863',NULL,'127.0.0.1','POST','/user/setUserAuthorities',200,835782667,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"ID\":1,\"authorityIds\":[1]}','{\"code\":0,\"data\":{},\"msg\":\"修改成功\"}',1),(59,'2025-10-30 21:44:59.689','2025-10-30 21:44:59.689',NULL,'127.0.0.1','POST','/casbin/updateCasbin',200,697358083,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"更新成功\"}',1),(60,'2025-10-30 21:46:17.308','2025-10-30 21:46:17.308',NULL,'127.0.0.1','POST','/authority/setDataAuthority',200,1175883000,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"CreatedAt\":\"2025-10-30T21:17:20.228+07:00\",\"UpdatedAt\":\"2025-10-30T21:43:19.358+07:00\",\"DeletedAt\":null,\"authorityId\":1,\"authorityName\":\"管理员\",\"parentId\":0,\"dataAuthorityId\":[{\"authorityId\":1,\"authorityName\":\"管理员\"}],\"children\":[],\"menus\":null,\"defaultRouter\":\"dashboard\"}','{\"code\":0,\"data\":{},\"msg\":\"设置成功\"}',1),(61,'2025-10-30 21:46:44.723','2025-10-30 21:46:44.723',NULL,'127.0.0.1','POST','/authority/setDataAuthority',200,1355328083,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"CreatedAt\":\"2025-10-30T15:16:45.987+07:00\",\"UpdatedAt\":\"2025-10-30T21:42:14.64+07:00\",\"DeletedAt\":null,\"authorityId\":888,\"authorityName\":\"普通用户\",\"parentId\":0,\"dataAuthorityId\":[{\"authorityId\":888,\"authorityName\":\"普通用户\"},{\"authorityId\":8881,\"authorityName\":\"普通用户子角色\"}],\"children\":[{\"CreatedAt\":\"2025-10-30T15:16:45.987+07:00\",\"UpdatedAt\":\"2025-10-30T15:16:59.246+07:00\",\"DeletedAt\":null,\"authorityId\":8881,\"authorityName\":\"普通用户子角色\",\"parentId\":888,\"dataAuthorityId\":[],\"children\":[],\"menus\":null,\"defaultRouter\":\"dashboard\"}],\"menus\":null,\"defaultRouter\":\"dashboard\"}','{\"code\":0,\"data\":{},\"msg\":\"设置成功\"}',1),(62,'2025-10-30 21:46:48.941','2025-10-30 21:46:48.941',NULL,'127.0.0.1','GET','/api/getApiGroups',200,143728250,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{}','{\"code\":0,\"data\":{\"apiGroupMap\":{\"api\":\"api\",\"attachmentCategory\":\"媒体库分类\",\"authority\":\"角色\",\"authorityBtn\":\"按钮权限\",\"autoCode\":\"代码生成器历史\",\"bot_mgr\":\"机器人\",\"bot_msg_mgr\":\"机器人消息管理\",\"casbin\":\"casbin\",\"customer\":\"客户\",\"email\":\"email\",\"fileUploadAndDownload\":\"文件上传与下载\",\"info\":\"公告\",\"jwt\":\"jwt\",\"menu\":\"菜单\",\"simpleUploader\":\"断点续传(插件版)\",\"sysDictionary\":\"系统字典\",\"sysDictionaryDetail\":\"系统字典详情\",\"sysExportTemplate\":\"导出模板\",\"sysOperationRecord\":\"操作记录\",\"sysParams\":\"参数管理\",\"sysVersion\":\"版本控制\",\"system\":\"系统服务\",\"user\":\"系统用户\"},\"groups\":[\"jwt\",\"系统用户\",\"api\",\"角色\",\"casbin\",\"菜单\",\"分片上传\",\"文件上传与下载\",\"系统服务\",\"客户\",\"代码生成器\",\"模板配置\",\"代码生成器历史\",\"系统字典详情\",\"系统字典\",\"操作记录\",\"断点续传(插件版)\",\"email\",\"按钮权限\",\"导出模板\",\"公告\",\"参数管理\",\"媒体库分类\",\"版本控制\",\"机器人消息管理\",\"机器人\"]},\"msg\":\"成功\"}',1),(63,'2025-10-30 23:32:53.490','2025-10-30 23:32:53.490',NULL,'127.0.0.1','PUT','/authority/updateAuthority',200,461858666,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','{\"authorityId\":1,\"AuthorityName\":\"管理员\",\"parentId\":0,\"defaultRouter\":\"bot_mgr\"}','{\"code\":0,\"data\":{\"authority\":{\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"authorityId\":1,\"authorityName\":\"管理员\",\"parentId\":0,\"dataAuthorityId\":null,\"children\":null,\"menus\":null,\"defaultRouter\":\"bot_mgr\"}},\"msg\":\"更新成功\"}',1),(64,'2025-10-30 23:32:55.312','2025-10-30 23:32:55.312',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1496141375,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1),(65,'2025-10-30 23:32:59.807','2025-10-30 23:32:59.807',NULL,'127.0.0.1','POST','/menu/addMenuAuthority',200,1503898167,'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36','','[超出记录长度]','{\"code\":0,\"data\":{},\"msg\":\"添加成功\"}',1);
/*!40000 ALTER TABLE `sys_operation_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_params`
--

DROP TABLE IF EXISTS `sys_params`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_params` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '参数名称',
  `key` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '参数键',
  `value` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '参数值',
  `desc` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '参数说明',
  PRIMARY KEY (`id`),
  KEY `idx_sys_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_params`
--

LOCK TABLES `sys_params` WRITE;
/*!40000 ALTER TABLE `sys_params` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_params` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_authority`
--

DROP TABLE IF EXISTS `sys_user_authority`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_authority` (
  `sys_user_id` bigint unsigned NOT NULL,
  `sys_authority_authority_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`sys_user_id`,`sys_authority_authority_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_authority`
--

LOCK TABLES `sys_user_authority` WRITE;
/*!40000 ALTER TABLE `sys_user_authority` DISABLE KEYS */;
INSERT INTO `sys_user_authority` VALUES (1,1),(2,888);
/*!40000 ALTER TABLE `sys_user_authority` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_users`
--

DROP TABLE IF EXISTS `sys_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `uuid` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户UUID',
  `username` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录名',
  `password` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录密码',
  `nick_name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT '系统用户' COMMENT '用户昵称',
  `header_img` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'https://qmplusimg.henrongyi.top/gva_header.jpg' COMMENT '用户头像',
  `authority_id` bigint unsigned DEFAULT '888' COMMENT '用户角色ID',
  `phone` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户手机号',
  `email` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户邮箱',
  `enable` bigint DEFAULT '1' COMMENT '用户是否被冻结 1正常 2冻结',
  `origin_setting` text COLLATE utf8mb4_unicode_ci COMMENT '配置',
  PRIMARY KEY (`id`),
  KEY `idx_sys_users_username` (`username`),
  KEY `idx_sys_users_deleted_at` (`deleted_at`),
  KEY `idx_sys_users_uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_users`
--

LOCK TABLES `sys_users` WRITE;
/*!40000 ALTER TABLE `sys_users` DISABLE KEYS */;
INSERT INTO `sys_users` VALUES (1,'2025-10-30 15:16:54.505','2025-10-30 21:43:48.511',NULL,'52e3169d-44c1-4e1e-bf07-1fb805bad974','admin','$2a$10$fpajQeX0aXa8TPaJzQw7/utDYc5QYObGDS00KDJUHJQBCIC/eJcZG','Mr.奇淼','https://qmplusimg.henrongyi.top/gva_header.jpg',1,'17611111111','333333333@qq.com',1,NULL),(2,'2025-10-30 15:16:54.505','2025-10-30 21:19:56.352',NULL,'94d549b8-8dd1-4a99-bf49-90c00ba8c312','a303176530','$2a$10$lwgFI96wA9dDNBBakOG20OpPbbBQQllNmrriyho1JnHarSvswMhvm','用户1','https://qmplusimg.henrongyi.top/1572075907logo.png',888,'17611111111','333333333@qq.com',1,NULL);
/*!40000 ALTER TABLE `sys_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_versions`
--

DROP TABLE IF EXISTS `sys_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_versions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `version_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '版本名称',
  `version_code` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '版本号',
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '版本描述',
  `version_data` text COLLATE utf8mb4_unicode_ci COMMENT '版本数据JSON',
  PRIMARY KEY (`id`),
  KEY `idx_sys_versions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_versions`
--

LOCK TABLES `sys_versions` WRITE;
/*!40000 ALTER TABLE `sys_versions` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_versions` ENABLE KEYS */;
UNLOCK TABLES;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-10-30 23:34:13
