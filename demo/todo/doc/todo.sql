/*
 Navicat Premium Data Transfer

 Source Server         : pi
 Source Server Type    : MySQL
 Source Server Version : 50737
 Source Host           : 192.168.31.70:3306
 Source Schema         : apijson_go_todo

 Target Server Type    : MySQL
 Target Server Version : 50737
 File Encoding         : 65001

 Date: 11/11/2022 17:26:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for _access
-- ----------------------------
DROP TABLE IF EXISTS `_access`;
CREATE TABLE `_access`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `debug` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否为调试表，只允许在开发环境使用，测试和线上环境禁用',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '实际表名，例如 apijson_user',
  `alias` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '外部调用的表别名，例如 User',
  `get` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 get 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]\n用 JSON 类型不能设置默认值，反正权限对应的需求是明确的，也不需要自动转 JSONArray。\nTODO: 直接 LOGIN,CONTACT,CIRCLE,OWNER 更简单，反正是开发内部用，不需要复杂查询。',
  `head` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 head 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',
  `gets` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 gets 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',
  `heads` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]' COMMENT '允许 heads 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',
  `post` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"OWNER\", \"ADMIN\"]' COMMENT '允许 post 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',
  `put` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"OWNER\", \"ADMIN\"]' COMMENT '允许 put 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',
  `delete` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '[\"OWNER\", \"ADMIN\"]' COMMENT '允许 delete 的角色列表，例如 [\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\"]',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `detail` varchar(1000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name_UNIQUE`(`name`) USING BTREE,
  UNIQUE INDEX `alias_UNIQUE`(`alias`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '权限配置(必须)' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of _access
-- ----------------------------
INSERT INTO `_access` VALUES (2, 0, 't_user', 'User', '[\"OWNER\",\"PARTNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"UNKNOWN\",\"LOGIN\",\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '2021-07-28 22:02:41', '用户公开信息表');
INSERT INTO `_access` VALUES (4, 0, 't_todo', 'Todo', '[\"OWNER\", \"PARTNER\",\"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\",\"OWNER\", \"ADMIN\"]', '[\"OWNER\",\"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '2021-07-28 22:02:41', '代办事项表');
INSERT INTO `_access` VALUES (5, 0, '_function', 'Function', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[]', '[]', '[]', '2018-11-29 00:38:15', '框架本身需要');
INSERT INTO `_access` VALUES (6, 0, 'privacy', 'Privacy', '[\"OWNER\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '2022-10-26 10:56:15', NULL);
INSERT INTO `_access` VALUES (8, 0, 'notice', 'Notice', '[\"UNKNOWN\",\"LOGIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '2022-10-26 10:56:35', NULL);
INSERT INTO `_access` VALUES (12, 0, 'notice_inner', 'NoticeInner', '[\"LOGIN\"]', '[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '[\"OWNER\", \"ADMIN\"]', '2022-10-26 10:56:53', NULL);

-- ----------------------------
-- Table structure for _access_ext
-- ----------------------------
DROP TABLE IF EXISTS `_access_ext`;
CREATE TABLE `_access_ext`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `table` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '表名',
  `row_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '关联主键字段名,联合主键时使用,分割',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `fields_get` json NULL COMMENT 'get时字段配置',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of _access_ext
-- ----------------------------
INSERT INTO `_access_ext` VALUES (1, 't_user', 'id', '2022-10-23 23:10:19', '{\"OWNER\": {\"in\": {\"user_id\": [\"=\"], \"username\": [\"*\"], \"created_at\": [\"$%\", \"=\"]}, \"out\": {\"id\": \"\", \"user_id\": \"\", \"username\": \"\", \"created_at\": \"\"}}, \"default\": {\"in\": {\"user_id\": [\"=\"], \"username\": [\"*\"], \"created_at\": [\"$%\", \"=\"]}, \"out\": {\"id\": \"\", \"username\": \"\"}}}');
INSERT INTO `_access_ext` VALUES (2, 't_todo', 'id', '2022-10-23 23:10:09', '{\"default\": {\"in\": {\"note\": [\"*\"], \"title\": [\"*\"], \"partner\": [\"*\"], \"user_id\": [\"=\", \"$%\"], \"created_at\": [\"$%\", \"=\"]}, \"out\": {\"title\": \"\", \"user_id\": \"\", \"created_at\": \"\"}}}');

-- ----------------------------
-- Table structure for _function
-- ----------------------------
DROP TABLE IF EXISTS `_function`;
CREATE TABLE `_function`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `debug` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。',
  `userId` bigint(20) NOT NULL COMMENT '管理员用户Id',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '方法名',
  `arguments` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '参数列表，每个参数的类型都是 String。\n用 , 分割的字符串 比 [JSONArray] 更好，例如 array,item ，更直观，还方便拼接函数。',
  `demo` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '可用的示例。\nTODO 改成 call，和返回值示例 back 对应。',
  `detail` varchar(1000) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '详细描述',
  `type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'Object' COMMENT '返回值类型。TODO RemoteFunction 校验 type 和 back',
  `version` tinyint(4) NOT NULL DEFAULT 0 COMMENT '允许的最低版本号，只限于GET,HEAD外的操作方法。\nTODO 使用 requestIdList 替代 version,tag,methods',
  `tag` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '允许的标签.\nnull - 允许全部\nTODO 使用 requestIdList 替代 version,tag,methods',
  `methods` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '允许的操作方法。\nnull - 允许全部\nTODO 使用 requestIdList 替代 version,tag,methods',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `back` varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '返回值示例',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '远程函数。强制在启动时校验所有demo是否能正常运行通过' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of _function
-- ----------------------------
INSERT INTO `_function` VALUES (1, 0, 0, 'sayHello', 'name', '{\"name\": \"test\"}', '最简单的远程函数', 'Object', 0, NULL, NULL, '2021-07-28 20:04:27', NULL);
INSERT INTO `_function` VALUES (2, 0, 0, 'isUserCanPutTodo', 'todoId', '{\"todoId\": 123}', '用来判定todo的写权限。', 'Object', 0, NULL, NULL, '2021-07-28 20:04:27', NULL);
INSERT INTO `_function` VALUES (3, 0, 0, 'getNoteCountAPI', '', '{}', '计数当前登录用户的todo数，展示如何在远程函数内部操作db', 'Object', 0, NULL, NULL, '2021-07-28 20:04:27', NULL);
INSERT INTO `_function` VALUES (4, 0, 0, 'rawSQLAPI', 'id', '{\"id\": \"_DOCUMENT_ONLY_\"}', '展示如何用裸SQL操作', 'Object', 0, NULL, NULL, '2021-07-28 20:04:27', NULL);
INSERT INTO `_function` VALUES (10, 0, 0, 'countArray', 'array', '{\"array\": [1, 2, 3]}', '（框架启动自检需要）获取数组长度。没写调用键值对，会自动补全 \"result()\": \"countArray(array)\"', 'Object', 0, NULL, NULL, '2018-10-13 16:23:23', NULL);
INSERT INTO `_function` VALUES (11, 0, 0, 'isContain', 'array,value', '{\"array\": [1, 2, 3], \"value\": 2}', '（框架启动自检需要）判断是否数组包含值。', 'Object', 0, NULL, NULL, '2018-10-13 16:23:23', NULL);
INSERT INTO `_function` VALUES (12, 0, 0, 'getFromArray', 'array,position', '{\"array\": [1, 2, 3], \"result()\": \"getFromArray(array,1)\"}', '（框架启动自检需要）根据下标获取数组里的值。position 传数字时直接作为值，而不是从所在对象 request 中取值', 'Object', 0, NULL, NULL, '2018-10-13 16:30:31', NULL);
INSERT INTO `_function` VALUES (13, 0, 0, 'getFromObject', 'object,key', '{\"key\": \"id\", \"object\": {\"id\": 1}}', '（框架启动自检需要）根据键获取对象里的值。', 'Object', 0, NULL, NULL, '2018-10-13 16:30:31', NULL);

-- ----------------------------
-- Table structure for _request
-- ----------------------------
DROP TABLE IF EXISTS `_request`;
CREATE TABLE `_request`  (
  `id` int(10) NOT NULL COMMENT '唯一标识',
  `debug` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。',
  `version` tinyint(4) NOT NULL DEFAULT 1 COMMENT 'GET,HEAD可用任意结构访问任意开放内容，不需要这个字段。\n其它的操作因为写入了结构和内容，所以都需要，按照不同的version选择对应的structure。\n\n自动化版本管理：\nRequest JSON最外层可以传  “version”:Integer 。\n1.未传或 <= 0，用最新版。 “@order”:”version-“\n2.已传且 > 0，用version以上的可用版本的最低版本。 “@order”:”version+”, “version{}”:”>={version}”',
  `method` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'GETS' COMMENT '只限于GET,HEAD外的操作方法。',
  `tag` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '标签',
  `structure` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '结构。\nTODO 里面的 PUT 改为 UPDATE，避免和请求 PUT 搞混。',
  `detail` varchar(10000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '详细说明',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '请求参数校验配置(必须)。\r\n最好编辑完后删除主键，这样就是只读状态，不能随意更改。需要更改就重新加上主键。\r\n\r\n每次启动服务器时加载整个表到内存。\r\n这个表不可省略，model内注解的权限只是客户端能用的，其它可以保证即便服务端代码错误时也不会误删数据。' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of _request
-- ----------------------------
INSERT INTO `_request` VALUES (2, 0, 1, 'POST', 'api_register', '{\"User\": {\"MUST\": \"username,realname\", \"REFUSE\": \"id\", \"UNIQUE\": \"username\"}, \"Credential\": {\"MUST\": \"pwdHash\", \"UPDATE\": {\"id@\": \"User/id\"}}}', '注意tag名小写开头，则不会被默认映射到表', '2021-07-29 02:15:40');
INSERT INTO `_request` VALUES (3, 0, 1, 'PUT', 'User', '{\"REFUSE\": \"username\", \"UPDATE\": {\"@role\": \"OWNER\"}}', 'user 修改自身数据', '2021-07-29 20:49:20');
INSERT INTO `_request` VALUES (4, 0, 1, 'POST', 'Todo', '{\"MUST\": \"title\", \"UPDATE\": {\"@role\": \"OWNER\"}, \"REFUSE\": \"id,user_id\"}', '增加todo', '2021-07-29 21:18:50');
INSERT INTO `_request` VALUES (5, 0, 1, 'PUT', 'Todo', '{\"Todo\":{ \"MUST\":\"id\",\"REFUSE\": \"userId\", \"INSERT\": {\"@role\": \"OWNER\"}} }', '修改todo', '2021-07-29 22:05:57');
INSERT INTO `_request` VALUES (6, 0, 1, 'DELETE', 'Todo', '{\"MUST\": \"id\", \"REFUSE\": \"!\", \"INSERT\": {\"@role\": \"OWNER\"}}', '删除todo', '2021-07-29 22:10:32');
INSERT INTO `_request` VALUES (8, 0, 1, 'PUT', 'helper+', '{\"Todo\": {\"MUST\": \"id,helper+\", \"INSERT\": {\"@role\": \"OWNER\"}}}', '增加todo helper', '2021-07-30 05:46:34');
INSERT INTO `_request` VALUES (9, 0, 1, 'PUT', 'helper-', '{\"Todo\": {\"MUST\": \"id,helper-\", \"INSERT\": {\"@role\": \"OWNER\"}}}', '删除todo helper', '2021-07-30 05:46:34');
INSERT INTO `_request` VALUES (10, 0, 1, 'POST', 'Todo:[]', '{\"Todo[]\": [{\"MUST\": \"title\", \"REFUSE\": \"id\"}], \"UPDATE\": {\"@role\": \"OWNER\"}}', '批量增加todo', '2021-08-01 12:51:31');
INSERT INTO `_request` VALUES (11, 0, 1, 'PUT', 'Todo:[]', '{\"Todo[]\":[{ \"MUST\":\"id\",\"REFUSE\": \"userId\", \"UPDATE\": {\"checkCanPut-()\": \"isUserCanPutTodo(id)\"}}] }', '每项单独设置（现在不生效）', '2021-08-01 12:51:31');
INSERT INTO `_request` VALUES (12, 0, 1, 'PUT', 'Todo[]', '{\"Todo\":{ \"MUST\":\"id{}\",\"REFUSE\": \"userId\", \"UPDATE\": {\"checkCanPut-()\": \"isUserCanPutTodo(id)\"}} }', '指定全部改（现在不生效）', '2021-08-01 12:51:31');
INSERT INTO `_request` VALUES (13, 0, 1, 'DELETE', 'Todo[]', '{\"Todo\": {\"MUST\": \"id{}\", \"REFUSE\": \"!\", \"INSERT\": {\"@role\": \"OWNER\"}}}', '删除todo', '2021-08-01 18:35:15');

-- ----------------------------
-- Table structure for notice
-- ----------------------------
DROP TABLE IF EXISTS `notice`;
CREATE TABLE `notice`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `content` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `created_by` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of notice
-- ----------------------------
INSERT INTO `notice` VALUES (2, '公告测试', '这是第一条公告', '2022-10-26 11:09:35', NULL);

-- ----------------------------
-- Table structure for notice_inner
-- ----------------------------
DROP TABLE IF EXISTS `notice_inner`;
CREATE TABLE `notice_inner`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `content` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `created_by` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of notice_inner
-- ----------------------------
INSERT INTO `notice_inner` VALUES (2, '“三体游戏” 版本更新，停机维护通知', '本次版本更新新增若干\"主\"的世界的特性， 将在11-11 11:11:11停机重启,  请同志们注意', '2022-10-26 11:12:29', NULL);

-- ----------------------------
-- Table structure for privacy
-- ----------------------------
DROP TABLE IF EXISTS `privacy`;
CREATE TABLE `privacy`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `secret_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of privacy
-- ----------------------------

-- ----------------------------
-- Table structure for t_todo
-- ----------------------------
DROP TABLE IF EXISTS `t_todo`;
CREATE TABLE `t_todo`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `note` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `partner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '与谁一起',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 69 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of t_todo
-- ----------------------------
INSERT INTO `t_todo` VALUES (18, '10001', '找史强去唱k', '唱 真的爱你', '2022-10-24 17:55:42', NULL, '10002');
INSERT INTO `t_todo` VALUES (20, '10001', '找林云去逛街', NULL, '2022-10-24 17:56:34', NULL, '10004');
INSERT INTO `t_todo` VALUES (22, '10001', '找丁仪发呆', NULL, '2022-10-24 17:56:56', NULL, '10003');
INSERT INTO `t_todo` VALUES (24, '10002', '找丁仪发呆', NULL, '2022-10-24 17:59:47', NULL, '10003');
INSERT INTO `t_todo` VALUES (26, '10002', '找丁仪发呆', NULL, '2022-10-24 18:00:03', NULL, '10003');
INSERT INTO `t_todo` VALUES (28, '10002', '找丁仪发呆', NULL, '2022-10-24 18:00:55', NULL, '10003');
INSERT INTO `t_todo` VALUES (30, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:01', NULL, '10003');
INSERT INTO `t_todo` VALUES (32, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:02', NULL, '10003');
INSERT INTO `t_todo` VALUES (34, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:02', NULL, '10003');
INSERT INTO `t_todo` VALUES (36, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:03', NULL, '10003');
INSERT INTO `t_todo` VALUES (38, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:03', NULL, '10003');
INSERT INTO `t_todo` VALUES (40, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:04', NULL, '10003');
INSERT INTO `t_todo` VALUES (42, '10002', '找丁仪发呆', NULL, '2022-10-24 18:01:04', NULL, '10003');
INSERT INTO `t_todo` VALUES (44, '10001', '找丁仪发呆', NULL, '2022-10-24 18:10:23', NULL, '10003');
INSERT INTO `t_todo` VALUES (46, '10002', '找丁仪发呆', NULL, '2022-10-24 18:10:26', NULL, '10003');
INSERT INTO `t_todo` VALUES (48, '10002', '找丁仪发呆', NULL, '2022-10-24 18:11:27', NULL, '10003');
INSERT INTO `t_todo` VALUES (50, '10002', '找丁仪发呆', NULL, '2022-10-24 18:11:28', NULL, '10003');
INSERT INTO `t_todo` VALUES (52, '10002', '找丁仪发呆', NULL, '2022-10-24 18:11:29', NULL, '10003');
INSERT INTO `t_todo` VALUES (54, '10002', '找丁仪发呆', NULL, '2022-10-24 18:11:30', NULL, '10003');
INSERT INTO `t_todo` VALUES (56, '10002', '找丁仪发呆', NULL, '2022-10-25 17:43:13', NULL, NULL);
INSERT INTO `t_todo` VALUES (58, '10002', '给林云搬家', NULL, NULL, NULL, NULL);
INSERT INTO `t_todo` VALUES (60, '12312', '找丁仪发呆', NULL, '2022-10-26 15:28:24', NULL, NULL);
INSERT INTO `t_todo` VALUES (62, '12312', '找丁仪发呆', NULL, '2022-10-26 15:28:56', NULL, NULL);
INSERT INTO `t_todo` VALUES (64, '10002', '找丁仪发呆', NULL, '2022-10-26 15:32:04', NULL, NULL);
INSERT INTO `t_todo` VALUES (66, '10002', '找丁仪发呆', NULL, '2022-10-26 15:32:48', NULL, NULL);
INSERT INTO `t_todo` VALUES (68, '10001', '找丁仪发呆', NULL, '2022-10-26 15:32:53', NULL, NULL);

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `realname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `User_id_uindex`(`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO `t_user` VALUES (2, '10001', 'wangmiao', '汪淼', '2022-10-24 17:04:11');
INSERT INTO `t_user` VALUES (4, '10002', 'shiqiang', '史强', '2022-10-24 17:06:09');
INSERT INTO `t_user` VALUES (6, '10003', 'dingyi', '丁仪', '2022-10-24 17:06:57');
INSERT INTO `t_user` VALUES (8, '10004', 'linyun', '林云', '2022-10-24 17:07:23');

SET FOREIGN_KEY_CHECKS = 1;
